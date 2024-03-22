package ratelimiter

import (
	"sync"
	"time"
)

type CallerFunction func(x CallerParam) ReturnResult
type CallerParam interface{}
type ReturnResult interface{}

type RateLimiter struct {
	startTime      time.Time
	limit          time.Duration
	initialWindow  time.Duration
	taskCount      int
	mutex          sync.Mutex
	resetThreshold int
	waitCh         chan struct{} // Channel to control the number of waiting goroutines
}

func NewRateLimiter(limit, initialWindow time.Duration, resetThreshold, maxWaitingGoroutines int) *RateLimiter {
	// Initialize the channel with the specified buffer size and fill it to capacity
	waitCh := make(chan struct{}, maxWaitingGoroutines)
	for i := 0; i < maxWaitingGoroutines; i++ {
		waitCh <- struct{}{}
	}
	return &RateLimiter{
		startTime:      time.Now(),
		limit:          limit,
		initialWindow:  initialWindow,
		resetThreshold: resetThreshold,
		waitCh:         waitCh,
	}
}

func (r *RateLimiter) Execute(task CallerFunction, param CallerParam) ReturnResult {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Reset task count if elapsed time exceeds initial window
	if time.Since(r.startTime) >= r.initialWindow {
		r.taskCount = 0
		r.startTime = time.Now()
	}

	// If within initial window and task count is below threshold, execute immediately
	if r.taskCount < r.resetThreshold {
		r.taskCount++
	} else {
		// If outside initial window or task count exceeds threshold, rate-limit
		elapsed := time.Since(r.startTime)
		if elapsed < r.initialWindow {
			time.Sleep(r.initialWindow - elapsed)
		}
	}

	// Retrieve token from the channel, effectively blocking if buffer is empty
	<-r.waitCh

	// Execute task
	result := task(param)

	// Release token back to the channel
	r.waitCh <- struct{}{}
	return result
}
