package ratelimit

import (
	"math"
	"sync"
	"time"
)

type LeakyBucket struct {
	rate     float64
	capacity float64
	water    float64
	lastLeak time.Time
	mu       sync.Mutex
}

func NewLeakyBucket(rate float64, capacity float64) *LeakyBucket {
	return &LeakyBucket{
		rate:     rate,
		capacity: capacity,
		water:    0,
		lastLeak: time.Now(),
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.leak()

	if lb.water < lb.capacity {
		lb.water++
		return true
	}

	return false
}

func (lb *LeakyBucket) leak() {
	now := time.Now()
	elapsed := now.Sub(lb.lastLeak).Seconds()

	leakedAmount := elapsed * lb.rate

	if leakedAmount > 0 {
		lb.water = math.Max(0, lb.water-leakedAmount)
		lb.lastLeak = now
	}
}
