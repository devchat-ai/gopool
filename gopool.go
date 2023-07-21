package gopool

import "time"

type Task func()

type GoPool struct {
    Workers    []*Worker
    MaxWorkers int
    workerStack []int
    taskQueue chan Task
}

func NewGoPool(maxWorkers int) *GoPool {
    pool := &GoPool{
        MaxWorkers: maxWorkers,
        Workers:    make([]*Worker, maxWorkers),
        workerStack: make([]int, maxWorkers),
        taskQueue: make(chan Task, 1e6),
    }
    for i := 0; i < maxWorkers; i++ {
        worker := newWorker()
        pool.Workers[i] = worker
        pool.workerStack[i] = i
        worker.start(pool, i)
    }
    go pool.dispatch()
    return pool
}

func (p *GoPool) AddTask(task Task) {
    p.taskQueue <- task
}

func (p *GoPool) Release() {
    close(p.taskQueue)
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

func (p *GoPool) dispatch() {
    for task := range p.taskQueue {
        for len(p.workerStack) == 0 {
            time.Sleep(time.Millisecond)
        }
        workerIndex := p.popWorker()
        p.Workers[workerIndex].TaskQueue <- task
    }
}
