package gopool

import (
	"sync"
	"testing"
	"time"
)

func TestGoPool(t *testing.T) {
	pool := NewGoPool(100)
	for i := 0; i < 1000; i++ {
		pool.AddTask(func() {
			time.Sleep(10 * time.Millisecond)
		})
	}
	pool.Release()
}

func BenchmarkGoPool(b *testing.B) {
	var wg sync.WaitGroup
	var taskNum = int(1e6)
	pool := NewGoPool(1e4)

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
	// pool.Release()
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
