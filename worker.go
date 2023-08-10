package gopool

import (
	"context"
	"fmt"
)

// worker represents a worker in the pool.
type worker struct {
	taskQueue chan task
}

func newWorker() *worker {
	return &worker{
		taskQueue: make(chan task, 1),
	}
}

// start starts the worker in a separate goroutine.
// The worker will run tasks from its taskQueue until the taskQueue is closed.
// For the length of the taskQueue is 1, the worker will be pushed back to the pool after executing 1 task.
func (w *worker) start(pool *goPool, workerIndex int) {
	go func() {
		for t := range w.taskQueue {
			if t != nil {
				result, err := w.executeTask(t, pool)
				w.handleResult(result, err, pool)
			}
			pool.pushWorker(workerIndex)
		}
	}()
}

// executeTask executes a task and returns the result and error.
// If the task fails, it will be retried according to the retryCount of the pool.
func (w *worker) executeTask(t task, pool *goPool) (result interface{}, err error) {
	for i := 0; i <= pool.retryCount; i++ {
		if pool.timeout > 0 {
			result, err = w.executeTaskWithTimeout(t, pool)
		} else {
			result, err = w.executeTaskWithoutTimeout(t, pool)
		}
		if err == nil || i == pool.retryCount {
			return result, err
		}
	}
	return
}

// executeTaskWithTimeout executes a task with a timeout and returns the result and error.
func (w *worker) executeTaskWithTimeout(t task, pool *goPool) (result interface{}, err error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), pool.timeout)
	defer cancel()

	// Create a channel to receive the result of the task
	resultChan := make(chan interface{})
	errChan := make(chan error)

	// Run the task in a separate goroutine
	go func() {
		res, err := t()
		select {
		case resultChan <- res:
		case errChan <- err:
		case <-ctx.Done():
			// The context was cancelled, stop the task
			return
		}
	}()

	// Wait for the task to finish or for the context to timeout
	select {
	case result = <-resultChan:
		err = <-errChan
		// The task finished successfully
		return result, err
	case <-ctx.Done():
		// The context timed out, the task took too long
		return nil, fmt.Errorf("task timed out")
	}
}

func (w *worker) executeTaskWithoutTimeout(t task, pool *goPool) (result interface{}, err error) {
	// If timeout is not set or is zero, just run the task
	return t()
}

// handleResult handles the result of a task.
func (w *worker) handleResult(result interface{}, err error, pool *goPool) {
	if err != nil && pool.errorCallback != nil {
		pool.errorCallback(err)
	} else if pool.resultCallback != nil {
		pool.resultCallback(result)
	}
}
