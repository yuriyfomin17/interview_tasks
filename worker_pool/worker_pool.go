package worker_pool

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	numWorkers int
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{numWorkers}
}

func (wp *WorkerPool) SubmitTasks(tasks []int) {
	jobs := make(chan int, len(tasks))
	outCh := make(chan int, len(tasks))
	wg := &sync.WaitGroup{}
	for workerId := 0; workerId < wp.numWorkers; workerId++ {

		go wp.worker(jobs, outCh, wg, workerId)
	}
	go func() {
		wg.Wait()
		close(outCh)
	}()
	go func() {
		for _, task := range tasks {
			jobs <- task
		}
		close(jobs)
	}()
	for outVal := range outCh {
		fmt.Printf("received outVal: %d\n", outVal)
	}
}

func (wp *WorkerPool) worker(jobs chan int, outCh chan int, wg *sync.WaitGroup, workerId int) {
	defer wg.Done()
	wg.Add(1)
	for job := range jobs {
		fmt.Printf("worker %d received job: %d\n", workerId, job)
		time.Sleep(1 * time.Second)
		outCh <- job
	}
}
