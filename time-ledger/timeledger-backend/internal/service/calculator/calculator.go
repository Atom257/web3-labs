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
)

type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

func New(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

// StartHourly：每小时跑一次（在 main.go 里启动 goroutine 调用）
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

// / RunOnce：对所有链/合约，计算 [last_calc_time, safe_block_time) 的积分并落库
func (s *Service) RunOnce(ctx context.Context, now time.Time) error {
	for _, chain := range s.cfg.Chains {
		for _, c := range chain.Contracts {

			//	以 safe block 的 block_time 作为积分上界
			safeT, err := s.safeBlockTime(ctx, chain.ChainID, c.Address)
			if err != nil {
				return err
			}

			//	用 safeT（而不是 now）进行积分计算
			if err := s.runContract(ctx, chain.ChainID, c.Address, safeT); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) runContract(ctx context.Context, chainID int64, contract string, now time.Time) error {
	// 1) 找到本合约下“所有已出现过的用户”（从 user_balance）
	var accounts []string
	if err := s.db.WithContext(ctx).
		Model(&models.UserBalance{}).
		Where("chain_id=? AND contract_address=?", chainID, contract).
		Pluck("account", &accounts).Error; err != nil {
		return fmt.Errorf("load accounts failed chain=%d contract=%s: %w", chainID, contract, err)
	}

	if len(accounts) == 0 {
		return nil
	}

	// 2) 确保每个 account 都有 user_point（没有则初始化：total_points=0, last_calc_time=now）
	//    注意：你也可以初始化为“第一次出现 balance 的 block_time”，这里先用 now，简单且一致。
	for _, acct := range accounts {

		initT, err := s.firstSeenTime(ctx, chainID, contract, acct, now)
		if err != nil {
			return fmt.Errorf("load first seen time failed acct=%s: %w", acct, err)
		}

		up := models.UserPoint{
			ChainID:         chainID,
			ContractAddress: contract,
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

	// 3) 对每个账户补算（从 last_calc_time 开始补到 now）
	for _, acct := range accounts {
		if err := s.calcOneAccount(ctx, chainID, contract, acct, now); err != nil {
			// 记录错误但不中断循环
			log.Printf("[ERROR] calcOneAccount failed: chain=%d contract=%s account=%s err=%v", chainID, contract, acct, err)
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
) error {

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1) 锁住 user_point，避免并发重复计算
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

		// 2) 计算积分（唯一真实来源：balance_log）
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
			// 防御：理论上不会出现
			deltaPoints = decimal.Zero
		}

		nowUTC := time.Now().UTC()

		// 3) 如果本段积分为 0，也推进 last_calc_time，避免重复扫描
		if deltaPoints.IsZero() {
			return tx.Model(&models.UserPoint{}).
				Where("id=?", up.ID).
				Updates(map[string]any{
					"last_calc_time": t1,
					"updated_at":     nowUTC,
				}).Error
		}

		// 4) 写入积分派生事实（user_point_log）：按稳定区间拆分
		for _, seg := range pd.Segments {

			// 没有积分的不记录
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

			if err := tx.
				Clauses(clause.OnConflict{DoNothing: true}).
				Create(&pl).Error; err != nil {
				return err
			}
		}

		// 5) 更新积分快照（user_point）
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

// 如果你觉得 record not found 太吵：可在 gorm logger 里调级别。
// 这里不在业务层 suppress。
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
	// 找该账户在 balance_log 中最早的 block_time
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
		// 理论上不会：account 来源于 user_balance，但做防御
		return fallback.UTC(), nil
	}

	return r.T.UTC(), nil
}

// 查 safe block 时间
func (s *Service) safeBlockTime(
	ctx context.Context,
	chainID int64,
	contract string,
) (time.Time, error) {

	type row struct {
		LastBlockTime time.Time
	}

	var r row
	// 查 block_cursor 表的last_block_time
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
