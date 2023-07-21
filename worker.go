package gopool

type Worker struct {
    TaskQueue chan Task
}

func newWorker() *Worker {
    return &Worker{
        TaskQueue: make(chan Task, 1),
    }
}

func (w *Worker) start(pool *GoPool, workerIndex int) {
    go func() {
        for task := range w.TaskQueue {
            if task != nil {
                task()
            }
            pool.pushWorker(workerIndex)
        }
    }()
}
