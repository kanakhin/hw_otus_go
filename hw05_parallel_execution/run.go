package hw05parallelexecution

import (
	"errors"
	"slices"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Counter struct {
	counter int
	mu      sync.Mutex
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	inProgressCount := Counter{}
	errorCount := Counter{}
	wg := &sync.WaitGroup{}

	for len(tasks) > 0 && (errorCount.Get() < m || m <= 0) {
		if inProgressCount.Get() >= n {
			continue
		}

		task := tasks[0]
		tasks = slices.Delete(tasks, 0, 1)

		wg.Add(1)

		inProgressCount.Inc()

		go func(task Task) {
			defer func() {
				inProgressCount.Dec()

				wg.Done()
			}()

			err := task()
			if err != nil {
				errorCount.Inc()
			}
		}(task)
	}

	wg.Wait()

	if m > 0 && errorCount.Get() >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func (c *Counter) Inc() {
	c.mu.Lock()
	c.counter++
	c.mu.Unlock()
}

func (c *Counter) Dec() {
	c.mu.Lock()
	c.counter--
	c.mu.Unlock()
}

func (c *Counter) Get() int {
	defer c.mu.Unlock()

	c.mu.Lock()

	return c.counter
}
