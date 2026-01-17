package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/config"
	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
)

// PointRateRepository 定义积分规则的存取接口
// 只负责数据库读写，不包含业务逻辑
type PointRateRepository interface {
	// 获取某个时间点生效的积分规则
	GetEffectiveRate(
		ctx context.Context,
		chainID int64,
		contract string,
		at time.Time,
	) (*models.PointRate, error)

	// 获取某个时间区间内的积分规则变更
	ListBetween(
		ctx context.Context,
		chainID int64,
		contract string,
		start, end time.Time,
	) ([]models.PointRate, error)

	// 获取某合约下所有积分规则（按生效时间升序）
	ListAll(
		ctx context.Context,
		chainID int64,
		contract string,
	) ([]models.PointRate, error)

	// 新增一条积分规则
	Create(
		ctx context.Context,
		rate *models.PointRate,
	) error
}

// pointRateRepo 是基于 GORM 的实现
type pointRateRepo struct {
	db *gorm.DB
}

// NewPointRateRepo 创建 PointRateRepository 实例
func NewPointRateRepo(db *gorm.DB) PointRateRepository {
	return &pointRateRepo{db: db}
}

// GetEffectiveRate 获取 at 时间点生效的积分规则
// 规则：effective_time <= at 的最后一条
func (r *pointRateRepo) GetEffectiveRate(
	ctx context.Context,
	chainID int64,
	contract string,
	at time.Time,
) (*models.PointRate, error) {

	var pr models.PointRate
	err := r.db.WithContext(ctx).
		Where(
			"chain_id = ? AND contract_address = ? AND effective_time <= ?",
			chainID, contract, at,
		).
		Order("effective_time DESC").
		Limit(1).
		First(&pr).Error

	if err != nil {
		return nil, err
	}

	return &pr, nil
}

// ListBetween 返回 (start, end) 区间内的积分规则变更
// 不包含 start，包含 end 之前的所有规则
func (r *pointRateRepo) ListBetween(
	ctx context.Context,
	chainID int64,
	contract string,
	start, end time.Time,
) ([]models.PointRate, error) {

	var rates []models.PointRate
	err := r.db.WithContext(ctx).
		Where(
			"chain_id = ? AND contract_address = ? AND effective_time > ? AND effective_time < ?",
			chainID, contract, start, end,
		).
		Order("effective_time ASC").
		Find(&rates).Error

	if err != nil {
		return nil, err
	}

	return rates, nil
}

// ListAll 获取某合约下所有积分规则
func (r *pointRateRepo) ListAll(
	ctx context.Context,
	chainID int64,
	contract string,
) ([]models.PointRate, error) {

	var rates []models.PointRate
	err := r.db.WithContext(ctx).
		Where(
			"chain_id = ? AND contract_address = ?",
			chainID, contract,
		).
		Order("effective_time ASC").
		Find(&rates).Error

	if err != nil {
		return nil, err
	}

	return rates, nil
}

// Create 新增一条积分规则
// 注意：不做任何冲突检测或业务校验，由上层 service 保证
func (r *pointRateRepo) Create(
	ctx context.Context,
	rate *models.PointRate,
) error {

	return r.db.WithContext(ctx).Create(rate).Error
}

// EnsureDefaultPointRate 确保每个合约至少有一条积分规则
// 如果已存在规则，则什么也不做
func EnsureDefaultPointRate(
	ctx context.Context,
	db *gorm.DB,
	cfg *config.Config,
) error {

	now := time.Now().UTC()

	for _, chain := range cfg.Chains {
		for _, c := range chain.Contracts {

			rates := []models.PointRate{
				// 初始积分规则：5%
				{
					ChainID:         chain.ChainID,
					ContractAddress: c.Address,
					RateNumerator:   5,
					RateDenominator: 100,
					EffectiveTime:   time.Unix(0, 0).UTC(),
					CreatedAt:       now,
				},
				// 升级积分规则：8%
				{
					ChainID:         chain.ChainID,
					ContractAddress: c.Address,
					RateNumerator:   8,
					RateDenominator: 100,
					EffectiveTime:   time.Date(2026, 1, 15, 3, 1, 0, 0, time.UTC),
					CreatedAt:       now,
				},
			}

			for _, r := range rates {
				if err := db.WithContext(ctx).
					Clauses(clause.OnConflict{DoNothing: true}).
					Create(&r).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}
