package ratelimit

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	limit       int
	windowSize  time.Duration
	count       int
	windowStart time.Time
	mu          sync.Mutex
}

func NewFixedWindowLimiter(limit int, windowSize time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:       limit,
		windowSize:  windowSize,
		windowStart: time.Now(),
		count:       0,
	}
}

func (fw *FixedWindowLimiter) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	if now.Sub(fw.windowStart) >= fw.windowSize {
		fw.windowStart = now
		fw.count = 0
	}

	if fw.count < fw.limit {
		fw.count++
		return true
	}

	return false
}
