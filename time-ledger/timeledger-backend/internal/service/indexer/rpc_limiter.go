package indexer

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"
)

/*
RPC Limiter + Retry
-------------------
- 全局 RPS 限制
- 每次 RPC 最多 3 次重试
- 一旦识别为 rate limit，立即返回 ErrRateLimited
*/

var ErrRateLimited = errors.New("rpc rate limited")

type RPCLimiter struct {
	ch chan struct{}
}

// NewRPCLimiter
// rps = 每秒允许的 RPC 次数（Alchemy free 建议 2~3）
func NewRPCLimiter(rps int) *RPCLimiter {
	ch := make(chan struct{}, rps)
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			for i := 0; i < rps; i++ {
				select {
				case ch <- struct{}{}:
				default:
				}
			}
		}
	}()
	return &RPCLimiter{ch: ch}
}

func (l *RPCLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.ch:
		return nil
	}
}

// callRPCWithRetry
// ----------------
// 所有 RPC 必须通过这个函数调用
//
// 行为：
// - 先等待 limiter
// - 最多 3 次尝试
// - 指数退避：100ms / 200ms / 400ms
// - 一旦识别为 rate limit，立即返回 ErrRateLimited
func callRPCWithRetry[T any](
	ctx context.Context,
	limiter *RPCLimiter,
	rpcName string,
	chainID int64,
	block uint64,
	fn func() (T, error),
) (T, error) {

	var zero T

	backoff := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
	}

	for i := 0; i < len(backoff); i++ {

		if err := limiter.Wait(ctx); err != nil {
			return zero, err
		}

		start := time.Now()
		res, err := fn()
		cost := time.Since(start)

		// 成功
		if err == nil {
			log.Printf(
				"[rpc.ok] chain_id=%d rpc=%s block=%d cost_ms=%d attempt=%d",
				chainID,
				rpcName,
				block,
				cost.Milliseconds(),
				i+1,
			)
			return res, nil
		}

		// 被限流：立刻终止
		if isRateLimitErr(err) {
			log.Printf(
				"[rpc.rate_limited] chain_id=%d rpc=%s block=%d cost_ms=%d attempt=%d err=%v",
				chainID,
				rpcName,
				block,
				cost.Milliseconds(),
				i+1,
				err,
			)
			return zero, ErrRateLimited
		}

		// 普通错误：重试
		log.Printf(
			"[rpc.retry] chain_id=%d rpc=%s block=%d cost_ms=%d attempt=%d backoff_ms=%d err=%v",
			chainID,
			rpcName,
			block,
			cost.Milliseconds(),
			i+1,
			backoff[i].Milliseconds(),
			err,
		)

		time.Sleep(backoff[i])
	}

	return zero, errors.New("rpc retry exhausted")
}

// 判断是否为 RPC 限流错误
func isRateLimitErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "429") ||
		strings.Contains(msg, "rate limit") ||
		strings.Contains(msg, "too many requests")
}
