package gopool

type Worker struct {
    TaskQueue chan Task
}

func newWorker(taskQueue chan Task) *Worker {
    return &Worker{
        TaskQueue: taskQueue,
    }
}

func (w *Worker) start() {
    go func() {
        for task := range w.TaskQueue {
            if task != nil {
                task()
            }
        }
    }()
}
