package function

import (
	"sync"
	"time"
)

type TokenBucketLimiter struct {
	Capacity     int           // 桶容量
	FillRate     time.Duration // 添加令牌速率，如每10ms加1个令牌
	tokens       int           // 当前令牌数
	lastFillTime time.Time     // 上次添加令牌时间
	mutex        sync.Mutex
}

// NewTokenBucketLimiter 创建新的令牌桶限流器
func NewTokenBucketLimiter(capacity int, fillRate time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		Capacity:     capacity,
		FillRate:     fillRate,
		tokens:       capacity,
		lastFillTime: time.Now(),
	}
}

// Allow 检查是否允许请求通过，返回布尔值
func (l *TokenBucketLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	// 计算应该添加的令牌数
	elapsed := now.Sub(l.lastFillTime)
	newTokens := int(elapsed / l.FillRate)

	if newTokens > 0 {
		l.tokens += newTokens
		if l.tokens > l.Capacity {
			l.tokens = l.Capacity
		}
		l.lastFillTime = now
	}

	if l.tokens <= 0 {
		return false
	}
	l.tokens--
	return true
}
