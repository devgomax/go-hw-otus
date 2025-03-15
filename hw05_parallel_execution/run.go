package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task, min(len(tasks), n))
	var errorsCount int32
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&wg, jobs, &errorsCount, int32(m))
	}

	err := loadFrom(tasks, jobs, &errorsCount, int32(m))

	wg.Wait()

	return err
}

func loadFrom(tasks []Task, jobs chan<- Task, errorsCount *int32, limit int32) error {
	defer close(jobs)

	for _, task := range tasks {
		if limit > 0 && atomic.LoadInt32(errorsCount) >= limit {
			return ErrErrorsLimitExceeded
		}
		jobs <- task
	}

	return nil
}

func worker(wg *sync.WaitGroup, jobs <-chan Task, errorsCount *int32, limit int32) {
	defer wg.Done()

	for task := range jobs {
		if limit > 0 && atomic.LoadInt32(errorsCount) >= limit {
			return
		}

		err := task()
		if err != nil {
			atomic.AddInt32(errorsCount, 1)
		}
	}
}
