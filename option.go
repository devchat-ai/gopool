package gopool

import (
	"sync"
	"time"
)

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

// WithTimeout sets the timeout for the pool.
func WithTimeout(timeout time.Duration) Option {
	return func(p *goPool) {
		p.timeout = timeout
	}
}

// WithResultCallback sets the result callback for the pool.
func WithResultCallback(callback func(interface{})) Option {
	return func(p *goPool) {
		p.resultCallback = callback
	}
}

// WithErrorCallback sets the error callback for the pool.
func WithErrorCallback(callback func(error)) Option {
	return func(p *goPool) {
		p.errorCallback = callback
	}
}

// WithRetryCount sets the retry count for the pool.
func WithRetryCount(retryCount int) Option {
	return func(p *goPool) {
		p.retryCount = retryCount
	}
}

// WithTaskQueueSize sets the size of the task queue for the pool.
func WithTaskQueueSize(size int) Option {
	return func(p *goPool) {
		p.taskQueueSize = size
	}
}
