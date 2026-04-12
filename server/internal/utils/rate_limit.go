package utils

import (
	"sync"
	"time"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	rate      int           // 单位时间内允许的请求数
	per       time.Duration // 时间窗口
	ips       map[string][]time.Time // IP地址及其请求时间
	mu        sync.Mutex
	cleanupInterval time.Duration
}

// NewRateLimiter 创建一个新的速率限制器
func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		rate:      rate,
		per:       per,
		ips:       make(map[string][]time.Time),
		cleanupInterval: time.Minute * 5,
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// 清理过期的请求记录
	rl.cleanupIP(ip, now)

	// 检查是否超过速率限制
	if len(rl.ips[ip]) >= rl.rate {
		return false
	}

	// 记录新的请求
	rl.ips[ip] = append(rl.ips[ip], now)
	return true
}

// cleanupIP 清理指定IP的过期请求记录
func (rl *RateLimiter) cleanupIP(ip string, now time.Time) {
	var valid []time.Time
	cutoff := now.Add(-rl.per)

	for _, t := range rl.ips[ip] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) == 0 {
		delete(rl.ips, ip)
	} else {
		rl.ips[ip] = valid
	}
}

// cleanup 定期清理所有过期的请求记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.per)

		for ip, times := range rl.ips {
			var valid []time.Time
			for _, t := range times {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}

			if len(valid) == 0 {
				delete(rl.ips, ip)
			} else {
				rl.ips[ip] = valid
			}
		}
		rl.mu.Unlock()
	}
}
