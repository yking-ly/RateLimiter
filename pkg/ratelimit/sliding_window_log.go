package ratelimit

import (
	"sync"
	"time"
)

type SlidingWindowLogLimiter struct {
	limit      int
	windowSize time.Duration
	logs       []time.Time
	mu         sync.Mutex
}

func NewSlidingWindowLogLimiter(limit int, windowSize time.Duration) *SlidingWindowLogLimiter {
	return &SlidingWindowLogLimiter{
		limit:      limit,
		windowSize: windowSize,
		logs:       make([]time.Time, 0),
	}
}

func (sw *SlidingWindowLogLimiter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	boundary := now.Add(-sw.windowSize)

	newLogs := sw.logs[:0]
	for _, t := range sw.logs {
		if t.After(boundary) {
			newLogs = append(newLogs, t)
		}
	}
	sw.logs = newLogs

	if len(sw.logs) < sw.limit {
		sw.logs = append(sw.logs, now)
		return true
	}

	return false
}
