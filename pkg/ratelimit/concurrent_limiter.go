package ratelimit

import "sync"

type ConcurrentLimiter struct {
	limit   int
	current int
	mu      sync.Mutex
}

func NewConcurrentLimiter(limit int) *ConcurrentLimiter {
	return &ConcurrentLimiter{
		limit:   limit,
		current: 0,
	}
}

func (c *ConcurrentLimiter) Acquire() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.current < c.limit {
		c.current++
		return true
	}
	return false
}

func (c *ConcurrentLimiter) Release() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.current > 0 {
		c.current--
	}
}

func (c *ConcurrentLimiter) Current() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.current
}
