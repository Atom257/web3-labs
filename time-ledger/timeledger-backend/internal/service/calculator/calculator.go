package calculator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
	"github.com/Atom257/web3-labs/timeledger-backend/internal/repository"
)

type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

func New(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

// StartHourly：每小时跑一次
func (s *Service) StartHourly(ctx context.Context) {
	// 对齐到整点（时间可修改）
	next := time.Now().UTC().Truncate(time.Hour).Add(time.Hour)
	timer := time.NewTimer(time.Until(next))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
	}

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		now := time.Now().UTC().Truncate(time.Hour)
		if err := s.RunOnce(ctx, now); err != nil {
			log.Printf("[calculator] run failed: %v", err)
		} else {
			log.Printf("[calculator] run ok at %s", now.Format(time.RFC3339))
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// RunOnce：对所有 DB 中的 Active 合约，计算积分并落库
func (s *Service) RunOnce(ctx context.Context, now time.Time) error {

	// 从数据库获取任务
	contracts, err := repository.GetActiveContracts(ctx, s.db)
	if err != nil {
		return fmt.Errorf("load active contracts failed: %w", err)
	}

	for _, c := range contracts {
		// 以 safe block 的 block_time 作为积分上界
		safeT, err := s.safeBlockTime(ctx, c.ChainID, c.Address)
		if err != nil {
			log.Printf("[ERROR] get safeBlockTime failed chain=%d contract=%s: %v", c.ChainID, c.Address, err)
			continue
		}

		// 传入 SysContract 对象，以便后续获取分表名
		if err := s.runContract(ctx, c, safeT); err != nil {
			log.Printf("[ERROR] runContract failed chain=%d contract=%s: %v", c.ChainID, c.Address, err)
		}
	}
	return nil
}

func (s *Service) runContract(ctx context.Context, contract models.SysContract, now time.Time) error {
	chainID := contract.ChainID
	addr := contract.Address

	// 1) 找到本合约下“所有已出现过的用户”（从 user_balance）
	var accounts []string
	if err := s.db.WithContext(ctx).
		Model(&models.UserBalance{}).
		Where("chain_id=? AND contract_address=?", chainID, addr).
		Pluck("account", &accounts).Error; err != nil {
		return fmt.Errorf("load accounts failed chain=%d contract=%s: %w", chainID, addr, err)
	}

	if len(accounts) == 0 {
		return nil
	}

	// 2) 确保每个 account 都有 user_point（没有则初始化）
	for _, acct := range accounts {
		initT, err := s.firstSeenTime(ctx, chainID, addr, acct, now)
		if err != nil {
			return fmt.Errorf("load first seen time failed acct=%s: %w", acct, err)
		}

		up := models.UserPoint{
			ChainID:         chainID,
			ContractAddress: addr,
			Account:         acct,
			TotalPoints:     "0",
			LastCalcTime:    initT,
			UpdatedAt:       time.Now().UTC(),
		}

		if err := s.db.WithContext(ctx).
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&up).Error; err != nil {
			return fmt.Errorf("ensure user_point failed acct=%s: %w", acct, err)
		}
	}

	// 获取动态表名 (例如 user_point_log_1)
	logTableName := contract.GetLogTableName()

	// 3) 对每个账户补算
	for _, acct := range accounts {
		if err := s.calcOneAccount(ctx, chainID, addr, acct, now, logTableName); err != nil {
			log.Printf("[ERROR] calcOneAccount failed: chain=%d contract=%s account=%s err=%v", chainID, addr, acct, err)
			continue
		}
	}

	return nil
}

func (s *Service) calcOneAccount(
	ctx context.Context,
	chainID int64,
	contract, account string,
	now time.Time,
	logTableName string, // 【新增参数】动态表名
) error {

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1) 锁住 user_point
		var up models.UserPoint
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(
				"chain_id=? AND contract_address=? AND account=?",
				chainID, contract, account,
			).
			First(&up).Error; err != nil {
			return err
		}

		t0 := up.LastCalcTime.UTC()
		t1 := now.UTC()

		if !t0.Before(t1) {
			return nil
		}

		// 2) 计算积分（调用 compute.go 中的函数，逻辑不变）
		pd, err := ComputePointsDelta(
			ctx,
			tx,
			chainID,
			contract,
			account,
			t0,
			t1,
		)
		if err != nil {
			return fmt.Errorf("compute points failed acct=%s: %w", account, err)
		}

		deltaPoints := pd.Total
		if deltaPoints.IsNegative() {
			deltaPoints = decimal.Zero
		}

		nowUTC := time.Now().UTC()

		// 3) 0分也推进时间
		if deltaPoints.IsZero() {
			return tx.Model(&models.UserPoint{}).
				Where("id=?", up.ID).
				Updates(map[string]any{
					"last_calc_time": t1,
					"updated_at":     nowUTC,
				}).Error
		}

		// 4) 写入积分派生事实（user_point_log）
		for _, seg := range pd.Segments {
			if seg.Points.IsZero() {
				continue
			}

			pl := models.UserPointLog{
				ChainID:         chainID,
				ContractAddress: contract,
				Account:         account,
				FromTime:        seg.FromTime,
				ToTime:          seg.ToTime,
				Balance:         seg.Balance.String(),
				Points:          seg.Points.String(),
				RateNumerator:   seg.RateNumerator,
				RateDenominator: seg.RateDenominator,
				CreatedAt:       nowUTC,
			}

			// 关键！使用 Table(logTableName) 写入分表
			if err := tx.Table(logTableName).
				Clauses(clause.OnConflict{DoNothing: true}).
				Create(&pl).Error; err != nil {
				return err
			}
		}

		// 5) 更新积分快照（user_point 总表不变）
		total, err := decimal.NewFromString(up.TotalPoints)
		if err != nil {
			return fmt.Errorf("invalid total_points in db: %w", err)
		}
		total = total.Add(deltaPoints)

		return tx.Model(&models.UserPoint{}).
			Where("id=?", up.ID).
			Updates(map[string]any{
				"total_points":   total.String(),
				"last_calc_time": t1,
				"updated_at":     nowUTC,
			}).Error
	})
}

var ErrNotFound = gorm.ErrRecordNotFound

func isNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func (s *Service) firstSeenTime(
	ctx context.Context,
	chainID int64,
	contract string,
	account string,
	fallback time.Time,
) (time.Time, error) {

	type row struct {
		T time.Time
	}

	var r row
	if err := s.db.WithContext(ctx).
		Model(&models.BalanceLog{}).
		Select("MIN(block_time) AS t").
		Where("chain_id=? AND contract_address=? AND account=?",
			chainID, contract, account,
		).
		Scan(&r).Error; err != nil {
		return time.Time{}, err
	}

	if r.T.IsZero() {
		return fallback.UTC(), nil
	}
	return r.T.UTC(), nil
}

func (s *Service) safeBlockTime(
	ctx context.Context,
	chainID int64,
	contract string,
) (time.Time, error) {

	type row struct {
		LastBlockTime time.Time
	}

	var r row
	err := s.db.WithContext(ctx).
		Model(&models.BlockCursor{}).
		Select("last_block_time").
		Where("chain_id=? AND contract_address=?", chainID, contract).
		Scan(&r).Error

	if err != nil {
		return time.Time{}, err
	}

	if r.LastBlockTime.IsZero() {
		return time.Time{}, fmt.Errorf("no safe block_time found")
	}

	return r.LastBlockTime.UTC(), nil
}
