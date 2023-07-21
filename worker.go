package gopool

type Worker struct {
    TaskQueue chan Task
}

func NewWorker(taskQueue chan Task) *Worker {
    return &Worker{
        TaskQueue: taskQueue,
    }
}

func (w *Worker) Start() {
    // Implementation here
}
