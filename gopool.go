package gopool

type Task func()

type GoPool struct {
    TaskQueue  chan Task
    MaxWorkers int
    Workers    []*Worker
}

func NewGoPool(maxWorkers int) *GoPool {
    pool := &GoPool{
        TaskQueue:  make(chan Task),
        MaxWorkers: maxWorkers,
        Workers:    make([]*Worker, maxWorkers),
    }
    for i := 0; i < maxWorkers; i++ {
        worker := newWorker(pool.TaskQueue)
        pool.Workers[i] = worker
        worker.start()
    }
    return pool
}

func (p *GoPool) AddTask(task Task) {
    p.TaskQueue <- task
}

func (p *GoPool) Release() {
    close(p.TaskQueue)
    for _, worker := range p.Workers {
        <-worker.TaskQueue
    }
}
