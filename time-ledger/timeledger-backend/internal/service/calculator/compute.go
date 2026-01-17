package calculator

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/Atom257/web3-labs/timeledger-backend/internal/models"
)

type ratePoint struct {
	At              time.Time
	RateNumerator   int64
	RateDenominator int64
	Rate            decimal.Decimal // 仅用于计算 = num/den
}

type balPoint struct {
	At           time.Time
	BlockNumber  int64
	LogIndex     int
	BalanceAfter *big.Int
}

type PointSegment struct {
	FromTime        time.Time
	ToTime          time.Time
	Balance         decimal.Decimal
	RateNumerator   int64
	RateDenominator int64
	Points          decimal.Decimal
}

type PointsDelta struct {
	Total    decimal.Decimal
	Segments []PointSegment
}

// ComputePointsDelta：计算 [t0, t1) 内应增加的积分（不落库）
//
// 核心：把时间线切成小段：
// - 余额变化点（balance_log.block_time）
// - rate 变化点（point_rate.effective_time）
// 对每段做： balance * rate * (durationSeconds / 3600)
func ComputePointsDelta(
	ctx context.Context,
	db *gorm.DB,
	chainID int64,
	contract, account string,
	t0, t1 time.Time,
) (PointsDelta, error) {

	// 1) 取 t0 时刻起始余额：最后一条 block_time < t0 的 balance_after
	startBal, err := loadBalanceAt(ctx, db, chainID, contract, account, t0)
	if err != nil {
		return PointsDelta{}, err
	}

	// 2) 取 [t0, t1) 内所有余额变化点（按 block_time, block_number, log_index 排序）
	bals, err := loadBalanceChanges(ctx, db, chainID, contract, account, t0, t1)
	if err != nil {
		return PointsDelta{}, err
	}

	// 3) 取 rate 时间线：需要 t0 生效的 rate + (t0, t1) 内的变更
	rates, err := loadRateTimeline(ctx, db, chainID, contract, t0, t1)
	if err != nil {
		return PointsDelta{}, err
	}
	if len(rates) == 0 {
		return PointsDelta{}, fmt.Errorf("no point_rate found for chain=%d contract=%s", chainID, contract)
	}

	// 4) 合并所有切割点（t0, t1 + rate/balance change times）
	// 我们用“指针推进”而不是构建全集表，更稳更快。
	curTime := t0.UTC()
	curBal := startBal

	// 余额指针
	bi := 0
	// rate 指针：rates 按 At 升序
	ri := rateIndexAtOrBefore(rates, curTime)

	// 当前生效的 rate（包含 num/den + decimal）
	curRatePoint := rates[ri]
	curRate := curRatePoint.Rate

	total := decimal.Zero

	segments := make([]PointSegment, 0)

	for curTime.Before(t1) {
		nextTime := t1

		// 下一个余额变化时间
		if bi < len(bals) {
			bt := bals[bi].At
			if bt.After(curTime) && bt.Before(nextTime) {
				nextTime = bt
			} else if bt.Equal(curTime) {
				// 同时刻多条变化（同 block_time），按 block/log 顺序依次应用
				curBal = bals[bi].BalanceAfter
				bi++
				continue
			}
		}

		// 下一个 rate 变化时间
		if ri+1 < len(rates) {
			rt := rates[ri+1].At
			if rt.After(curTime) && rt.Before(nextTime) {
				nextTime = rt
			} else if rt.Equal(curTime) {
				ri++
				curRatePoint = rates[ri]
				curRate = curRatePoint.Rate
				continue
			}
		}

		// 计算 [curTime, nextTime) 的积分
		if nextTime.After(curTime) {
			seconds := int64(nextTime.Sub(curTime).Seconds())

			if seconds > 0 && curBal.Sign() > 0 && curRate.GreaterThan(decimal.Zero) {
				segPoints := pointsForSegment(curBal, curRate, seconds)
				// 记录积分变化的 log
				segments = append(segments, PointSegment{
					FromTime:        curTime,
					ToTime:          nextTime,
					Balance:         decimal.NewFromBigInt(curBal, -18),
					RateNumerator:   curRatePoint.RateNumerator,
					RateDenominator: curRatePoint.RateDenominator,
					Points:          segPoints,
				})

				total = total.Add(segPoints)
			}
			curTime = nextTime

		}

		// 应用在 nextTime 发生的变化（下一轮循环会处理 Equal 情况）
	}

	return PointsDelta{
		Total:    total,
		Segments: segments,
	}, nil

}

func pointsForSegment(balanceWei *big.Int, rate decimal.Decimal, seconds int64) decimal.Decimal {
	// balance * rate * seconds / 3600
	balDec := decimal.NewFromBigInt(balanceWei, -18)
	sec := decimal.NewFromInt(seconds)
	return balDec.Mul(rate).Mul(sec).Div(decimal.NewFromInt(3600))
}

/* ---------------- DB loads ---------------- */

func loadBalanceAt(
	ctx context.Context,
	db *gorm.DB,
	chainID int64,
	contract, account string,
	t time.Time,
) (*big.Int, error) {

	// 找最后一条 block_time < t 的 balance_after
	type row struct {
		BalanceAfter string
	}
	var r row
	err := db.WithContext(ctx).
		Model(&models.BalanceLog{}).
		Select("balance_after").
		Where("chain_id=? AND contract_address=? AND account=? AND block_time < ?",
			chainID, contract, account, t,
		).
		Order("block_time DESC, block_number DESC, log_index DESC").
		Limit(1).
		Scan(&r).Error
	if err != nil {
		return nil, err
	}

	if r.BalanceAfter == "" {
		return big.NewInt(0), nil
	}

	bi, ok := new(big.Int).SetString(r.BalanceAfter, 10)
	if !ok {
		return nil, fmt.Errorf("invalid balance_after in db: %s", r.BalanceAfter)
	}
	return bi, nil
}

func loadBalanceChanges(
	ctx context.Context,
	db *gorm.DB,
	chainID int64,
	contract, account string,
	t0, t1 time.Time,
) ([]balPoint, error) {

	type row struct {
		BlockTime    time.Time
		BlockNumber  int64
		LogIndex     int
		BalanceAfter string
	}

	var rows []row
	if err := db.WithContext(ctx).
		Model(&models.BalanceLog{}).
		Select("block_time, block_number, log_index, balance_after").
		Where("chain_id=? AND contract_address=? AND account=? AND block_time >= ? AND block_time < ?",
			chainID, contract, account, t0, t1,
		).
		Order("block_time ASC, block_number ASC, log_index ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]balPoint, 0, len(rows))
	for _, r := range rows {
		bi, ok := new(big.Int).SetString(r.BalanceAfter, 10)
		if !ok {
			return nil, fmt.Errorf("invalid balance_after: %s", r.BalanceAfter)
		}
		out = append(out, balPoint{
			At:           r.BlockTime.UTC(),
			BlockNumber:  r.BlockNumber,
			LogIndex:     r.LogIndex,
			BalanceAfter: bi,
		})
	}

	// 已经按 SQL 排序了；这里再稳定一下也行
	sort.SliceStable(out, func(i, j int) bool {
		if !out[i].At.Equal(out[j].At) {
			return out[i].At.Before(out[j].At)
		}
		if out[i].BlockNumber != out[j].BlockNumber {
			return out[i].BlockNumber < out[j].BlockNumber
		}
		return out[i].LogIndex < out[j].LogIndex
	})

	return out, nil
}

func loadRateTimeline(
	ctx context.Context,
	db *gorm.DB,
	chainID int64,
	contract string,
	t0, t1 time.Time,
) ([]ratePoint, error) {

	// 取 t0 时刻生效的最后一条 rate（effective_time <= t0）
	type row struct {
		EffectiveTime   time.Time
		RateNumerator   int64
		RateDenominator int64
	}

	var base row
	if err := db.WithContext(ctx).
		Model(&models.PointRate{}).
		Select("effective_time, rate_numerator, rate_denominator").
		Where("chain_id=? AND contract_address=? AND effective_time <= ?",
			chainID, contract, t0,
		).
		Order("effective_time DESC").
		Limit(1).
		Scan(&base).Error; err != nil {
		return nil, err
	}

	if base.RateDenominator == 0 {
		// 没有 <= t0 的规则，尝试找第一条 > t0 的规则（否则就是没有任何规则）
		// 你们要求必须有默认规则，所以这里报错更好。
		return nil, fmt.Errorf("no base rate at/before t0=%s", t0.Format(time.RFC3339))
	}

	out := []ratePoint{{
		At:              base.EffectiveTime.UTC(),
		RateNumerator:   base.RateNumerator,
		RateDenominator: base.RateDenominator,
		Rate:            decimal.NewFromInt(base.RateNumerator).Div(decimal.NewFromInt(base.RateDenominator)),
	}}

	// 取 (t0, t1) 内的变更
	var changes []row
	if err := db.WithContext(ctx).
		Model(&models.PointRate{}).
		Select("effective_time, rate_numerator, rate_denominator").
		Where("chain_id=? AND contract_address=? AND effective_time > ? AND effective_time < ?",
			chainID, contract, t0, t1,
		).
		Order("effective_time ASC").
		Find(&changes).Error; err != nil {
		return nil, err
	}

	for _, c := range changes {
		out = append(out, ratePoint{
			At:              c.EffectiveTime.UTC(),
			RateNumerator:   c.RateNumerator,
			RateDenominator: c.RateDenominator,
			Rate:            decimal.NewFromInt(c.RateNumerator).Div(decimal.NewFromInt(c.RateDenominator)),
		})
	}

	// 确保有序
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].At.Before(out[j].At)
	})

	return out, nil
}

/* ---------------- rate helpers ---------------- */

func rateIndexAtOrBefore(rates []ratePoint, t time.Time) int {
	// rates 按 At 升序
	idx := 0
	for i := 0; i < len(rates); i++ {
		if rates[i].At.After(t) {
			break
		}
		idx = i
	}
	return idx
}

func rateAt(rates []ratePoint, t time.Time) decimal.Decimal {
	return rates[rateIndexAtOrBefore(rates, t)].Rate
}
