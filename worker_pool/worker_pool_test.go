package worker_pool

import (
	"context"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1100*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		WorkerPool()
		close(done)
	}()

	select {
	case <-ctx.Done():
		t.Fatal("test timed out")
	case <-done:
		// Test completed successfully
	}
}
