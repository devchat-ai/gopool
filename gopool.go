package gopool

type Task func()

type GoPool struct {
    TaskQueue  chan Task
    MaxWorkers int
    Workers    []*Worker
}

func NewGoPool(maxWorkers int) *GoPool {
    return &GoPool{
        TaskQueue:  make(chan Task),
        MaxWorkers: maxWorkers,
        Workers:    make([]*Worker, maxWorkers),
    }
}

func (p *GoPool) AddTask(task Task) {
    // Implementation here
}

func (p *GoPool) Release() {
    // Implementation here
}
