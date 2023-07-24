package gopool

import "sync"

// Option represents an option for the pool.
type Option func(*goPool)

// WithLock sets the lock for the pool.
func WithLock(lock sync.Locker) Option {
    return func(p *goPool) {
        p.lock = lock
        p.cond = sync.NewCond(p.lock)
    }
}

// WithMinWorkers sets the minimum number of workers for the pool.
func WithMinWorkers(minWorkers int) Option {
    return func(p *goPool) {
        p.minWorkers = minWorkers
    }
}
