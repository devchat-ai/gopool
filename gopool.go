package gopool

type Task func()

type GoPool struct {
    Workers    []*Worker
    MaxWorkers int
    workerStack []int
}

func NewGoPool(maxWorkers int) *GoPool {
    pool := &GoPool{
        MaxWorkers: maxWorkers,
        Workers:    make([]*Worker, maxWorkers),
        workerStack: make([]int, maxWorkers),
    }
    for i := 0; i < maxWorkers; i++ {
        worker := newWorker()
        pool.Workers[i] = worker
        pool.workerStack[i] = i
        worker.start(pool, i)
    }
    return pool
}

func (p *GoPool) AddTask(task Task) {
    workerIndex := p.popWorker()
    p.Workers[workerIndex].TaskQueue <- task
}

func (p *GoPool) Release() {
    for _, worker := range p.Workers {
        close(worker.TaskQueue)
    }
}

func (p *GoPool) popWorker() int {
    workerIndex := p.workerStack[len(p.workerStack)-1]
    p.workerStack = p.workerStack[:len(p.workerStack)-1]
    return workerIndex
}

func (p *GoPool) pushWorker(workerIndex int) {
    p.workerStack = append(p.workerStack, workerIndex)
}
