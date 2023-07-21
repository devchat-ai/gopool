package gopool

import (
    "sync"
)

type Task func()

type GoPool struct {
    Workers    []*Worker
    MaxWorkers int
    workerStack []int
    taskQueue chan Task
    mutex sync.Mutex
    cond *sync.Cond
}

func NewGoPool(maxWorkers int) *GoPool {
    pool := &GoPool{
        MaxWorkers: maxWorkers,
        Workers:    make([]*Worker, maxWorkers),
        workerStack: make([]int, maxWorkers),
        taskQueue: make(chan Task, 1e6),
    }
    pool.cond = sync.NewCond(&pool.mutex)
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
    p.cond.L.Lock()
    for len(p.workerStack) != p.MaxWorkers {
        p.cond.Wait()
    }
    p.cond.L.Unlock()
    for _, worker := range p.Workers {
        close(worker.TaskQueue)
    }
    p.Workers = nil
    p.workerStack = nil
}

func (p *GoPool) popWorker() int {
    p.mutex.Lock()
    workerIndex := p.workerStack[len(p.workerStack)-1]
    p.workerStack = p.workerStack[:len(p.workerStack)-1]
    p.mutex.Unlock()
    return workerIndex
}

func (p *GoPool) pushWorker(workerIndex int) {
    p.mutex.Lock()
    p.workerStack = append(p.workerStack, workerIndex)
    p.mutex.Unlock()
    p.cond.Signal()
}

func (p *GoPool) dispatch() {
    for task := range p.taskQueue {
        p.cond.L.Lock()
        for len(p.workerStack) == 0 {
            p.cond.Wait()
        }
        p.cond.L.Unlock()
        workerIndex := p.popWorker()
        p.Workers[workerIndex].TaskQueue <- task
    }
}
