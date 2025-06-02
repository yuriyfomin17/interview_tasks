package worker_pool

import (
	"fmt"
	"sync"
	"time"
)

func WorkerPool() {
	numWorkers := 3
	jobsArr := []string{"a", "b", "c"}
	jobs := make(chan string)
	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				// imitate some work
				time.Sleep(1 * time.Second)
				fmt.Println("current job is", job)
			}
		}()
	}
	go func() {
		for _, jobValue := range jobsArr {
			jobs <- jobValue
		}
		close(jobs)
	}()
	wg.Wait()
}
