package ratelimiter

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var tooManyRequestsError = errors.New("too many requests")

type RateLimiter struct {
	limit                 int
	currentRequestCounter int
	expiresAt             time.Time
	lock                  sync.Mutex
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		limit:                 limit,
		currentRequestCounter: 0,
		expiresAt:             time.Now().Add(time.Millisecond * 1000),
	}
}

func (rl *RateLimiter) Process(currentFunc func() error) error {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	if time.Now().After(rl.expiresAt) {

		rl.currentRequestCounter = 0
		rl.expiresAt = time.Now().Add(time.Millisecond * 1000)
	}
	if rl.currentRequestCounter >= rl.limit {
		return tooManyRequestsError
	}

	rl.currentRequestCounter += 1
	err := currentFunc()
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}

func TestRateLimiter(requestLimitPerSecond int, port string) {
	rl := NewRateLimiter(requestLimitPerSecond)
	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		err := rl.Process(func() error {
			return nil
		})
		if errors.Is(err, tooManyRequestsError) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	})
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
