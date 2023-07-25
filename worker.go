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
                var result interface{}
                var err error

                if pool.timeout > 0 {
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
                        if err != nil && pool.errorCallback != nil {
                            pool.errorCallback(err)
                        } else if pool.resultCallback != nil {
                            pool.resultCallback(result)
                        }
                    case <-ctx.Done():
                        // The context timed out, the task took too long
                        if pool.errorCallback != nil {
                            pool.errorCallback(fmt.Errorf("Task timed out"))
                        }
                    }
                } else {
                    // If timeout is not set or is zero, just run the task
                    result, err = t()
                    if err != nil && pool.errorCallback != nil {
                        pool.errorCallback(err)
                    } else if pool.resultCallback != nil {
                        pool.resultCallback(result)
                    }
                }
            }
            pool.pushWorker(workerIndex)
        }
    }()
}
