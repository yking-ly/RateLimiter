package ratelimit

import (
	"sync"
	"time"
)

type SlidingWindowCounterLimiter struct {
	limit                float64
	windowSize           time.Duration
	previousWindowCounts int
	currentWindowCounts  int
	currentWindowStart   time.Time
	mu                   sync.Mutex
}

func NewSlidingWindowCounterLimiter(limit int, windowSize time.Duration) *SlidingWindowCounterLimiter {
	return &SlidingWindowCounterLimiter{
		limit:                float64(limit),
		windowSize:           windowSize,
		currentWindowStart:   time.Now(),
		previousWindowCounts: 0,
		currentWindowCounts:  0,
	}
}

func (sw *SlidingWindowCounterLimiter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()

	timeSinceStart := now.Sub(sw.currentWindowStart)

	if timeSinceStart >= sw.windowSize {
		windowsPassed := int(timeSinceStart / sw.windowSize)

		if windowsPassed == 1 {
			sw.previousWindowCounts = sw.currentWindowCounts
			sw.currentWindowCounts = 0
			sw.currentWindowStart = sw.currentWindowStart.Add(sw.windowSize)
		} else {
			sw.previousWindowCounts = 0
			sw.currentWindowCounts = 0
			sw.currentWindowStart = now
		}

		timeSinceStart = now.Sub(sw.currentWindowStart)
	}

	weightPrev := float64(sw.windowSize-timeSinceStart) / float64(sw.windowSize)

	estimatedCount := float64(sw.previousWindowCounts)*weightPrev + float64(sw.currentWindowCounts)

	if estimatedCount < sw.limit {
		sw.currentWindowCounts++
		return true
	}

	return false
}
