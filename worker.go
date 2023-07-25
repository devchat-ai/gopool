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

func (w *worker) executeTaskWithTimeout(t task, pool *goPool) (result interface{}, err error) {
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), pool.timeout)
    defer cancel()

    // Create a channel to receive the result of the task
    done := make(chan struct{})

    // Run the task in a separate goroutine
    go func() {
        result, err = t()
        close(done)
    }()

    // Wait for the task to finish or for the context to timeout
    select {
    case <-done:
        // The task finished successfully
        return result, err
    case <-ctx.Done():
        // The context timed out, the task took too long
        return nil, fmt.Errorf("Task timed out")
    }
}

func (w *worker) executeTaskWithoutTimeout(t task, pool *goPool) (result interface{}, err error) {
    // If timeout is not set or is zero, just run the task
    return t()
}

func (w *worker) handleResult(result interface{}, err error, pool *goPool) {
    if err != nil && pool.errorCallback != nil {
        pool.errorCallback(err)
    } else if pool.resultCallback != nil {
        pool.resultCallback(result)
    }
}
