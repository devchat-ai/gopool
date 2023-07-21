package gopool

import (
    "sync"
    "testing"
    "time"

    "github.com/daniel-hutao/spinlock"
)

func TestGoPoolWithMutex(t *testing.T) {
    pool := NewGoPool(100, WithLock(new(sync.Mutex)))
    for i := 0; i < 1000; i++ {
        pool.AddTask(func() {
            time.Sleep(10 * time.Millisecond)
        })
    }
    pool.Release()
}

func TestGoPoolWithSpinLock(t *testing.T) {
    pool := NewGoPool(100, WithLock(new(spinlock.SpinLock)))
    for i := 0; i < 1000; i++ {
        pool.AddTask(func() {
            time.Sleep(10 * time.Millisecond)
        })
    }
    pool.Release()
}

func BenchmarkGoPoolWithMutex(b *testing.B) {
    var wg sync.WaitGroup
    var taskNum = int(1e6)
    pool := NewGoPool(5e4, WithLock(new(sync.Mutex)))

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wg.Add(taskNum)
        for num := 0; num < taskNum; num++ {
            pool.AddTask(func() {
                time.Sleep(10 * time.Millisecond)
                wg.Done()
            })
        }
    }
    wg.Wait()
	b.StopTimer()
    pool.Release()
}

func BenchmarkGoPoolWithSpinLock(b *testing.B) {
    var wg sync.WaitGroup
    var taskNum = int(1e6)
    pool := NewGoPool(5e4, WithLock(new(spinlock.SpinLock)))

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        wg.Add(taskNum)
        for num := 0; num < taskNum; num++ {
            pool.AddTask(func() {
                time.Sleep(10 * time.Millisecond)
                wg.Done()
            })
        }
    }
    wg.Wait()
	b.StopTimer()
    pool.Release()
}

func BenchmarkGoroutines(b *testing.B) {
    var wg sync.WaitGroup
    var taskNum = int(1e6)

    for i := 0; i < b.N; i++ {
        wg.Add(taskNum)
        for num := 0; num < taskNum; num++ {
            go func() {
                time.Sleep(10 * time.Millisecond)
                wg.Done()
            }()
        }
    }
}
