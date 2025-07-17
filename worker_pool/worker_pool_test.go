package worker_pool

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWorkerPool_ShouldExecuteWithinLimits(t *testing.T) {
	workerPool := NewWorkerPool(3)
	timeNow := time.Now()
	workerPool.SubmitTasks([]int{1, 2, 3})
	fmt.Println("time elapsed:", time.Now().Sub(timeNow))
	require.True(t, time.Since(timeNow) <= 2*time.Second)
}
