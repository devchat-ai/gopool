package gopool

import (
	"sync"
	"testing"
	"time"

	"github.com/daniel-hutao/spinlock"
)

func BenchmarkGoPoolWithMutex(b *testing.B) {
	var wg sync.WaitGroup
	var taskNum = int(1e6)
	pool := NewGoPool(1e4, WithLock(new(sync.Mutex)))
	defer pool.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(taskNum)
		for num := 0; num < taskNum; num++ {
			pool.AddTask(func() (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
				return nil, nil
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkGoPoolWithSpinLock(b *testing.B) {
	var wg sync.WaitGroup
	var taskNum = int(1e6)
	pool := NewGoPool(1e4, WithLock(new(spinlock.SpinLock)))
	defer pool.Release()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(taskNum)
		for num := 0; num < taskNum; num++ {
			pool.AddTask(func() (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
				return nil, nil
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}

func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup
	var taskNum = int(1e6)

	for i := 0; i < b.N; i++ {
		wg.Add(taskNum)
		for num := 0; num < taskNum; num++ {
			go func() (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
				return nil, nil
			}()
		}
		wg.Wait()
	}
}
