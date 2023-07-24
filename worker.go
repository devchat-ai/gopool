package gopool

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
                t()
            }
            pool.pushWorker(workerIndex)
        }
    }()
}
