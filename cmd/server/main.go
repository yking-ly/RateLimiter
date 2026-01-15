package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ratelimiter/pkg/ratelimit"
	"sync"
	"time"
)

var (
	limiter     RateLimiter
	limiterType string
	mu          sync.Mutex
)

type RateLimiter interface {
	Allow() bool
}

type ConfigRequest struct {
	Type       string  `json:"type"`
	Rate       float64 `json:"rate"`
	Capacity   float64 `json:"capacity"`
	Limit      int     `json:"limit"`
	WindowSecs float64 `json:"windowSecs"`
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/api/configure", handleConfigure)
	http.HandleFunc("/api/request", handleRequest)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleConfigure(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	limiterType = config.Type

	switch config.Type {
	case "token_bucket":
		limiter = ratelimit.NewTokenBucket(config.Rate, config.Capacity)
	case "leaky_bucket":
		limiter = ratelimit.NewLeakyBucket(config.Rate, config.Capacity)
	case "fixed_window":
		limiter = ratelimit.NewFixedWindowLimiter(config.Limit, time.Duration(config.WindowSecs*float64(time.Second)))
	case "sliding_log":
		limiter = ratelimit.NewSlidingWindowLogLimiter(config.Limit, time.Duration(config.WindowSecs*float64(time.Second)))
	case "sliding_counter":
		limiter = ratelimit.NewSlidingWindowCounterLimiter(config.Limit, time.Duration(config.WindowSecs*float64(time.Second)))
	case "concurrent":
		cl := ratelimit.NewConcurrentLimiter(config.Limit)
		limiter = &AutoReleasingConcurrentLimiter{cl: cl, workDuration: 500 * time.Millisecond}
	default:
		http.Error(w, "Unknown limiter type", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "configured", "message": fmt.Sprintf("Switched to %s", limiterType)})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	l := limiter
	mu.Unlock()

	if l == nil {
		http.Error(w, "Limiter not configured", http.StatusBadRequest)
		return
	}

	allowed := l.Allow()

	resp := map[string]interface{}{
		"allowed": allowed,
		"ts":      time.Now().UnixMilli(),
	}
	json.NewEncoder(w).Encode(resp)
}

type AutoReleasingConcurrentLimiter struct {
	cl           *ratelimit.ConcurrentLimiter
	workDuration time.Duration
}

func (a *AutoReleasingConcurrentLimiter) Allow() bool {
	if a.cl.Acquire() {
		go func() {
			time.Sleep(a.workDuration)
			a.cl.Release()
		}()
		return true
	}
	return false
}
