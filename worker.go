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
                if pool.timeout > 0 {
                    // Create a context with timeout
                    ctx, cancel := context.WithTimeout(context.Background(), pool.timeout)
                    defer cancel()

                    // Create a channel to receive the result of the task
                    done := make(chan struct{})

                    // Run the task in a separate goroutine
                    go func() {
                        t()
                        close(done)
                    }()

                    // Wait for the task to finish or for the context to timeout
                    select {
                    case <-done:
                        // The task finished successfully
                    case <-ctx.Done():
                        // The context timed out, the task took too long
                        fmt.Println("Task timed out")
                    }
                } else {
                    // If timeout is not set or is zero, just run the task
                    t()
                }
            }
            pool.pushWorker(workerIndex)
        }
    }()
}
