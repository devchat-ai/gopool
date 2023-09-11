package gopool

import (
	"sync"
	"testing"
	"time"

	"github.com/alitto/pond"
	"github.com/daniel-hutao/spinlock"
	"github.com/panjf2000/ants/v2"
)

const (
	PoolSize = 1e4
	TaskNum  = 1e6
)

func BenchmarkGoPool(b *testing.B) {
	pool := NewGoPool(PoolSize)
	defer pool.Release()

	taskFunc := func() (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return nil, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for num := 0; num < TaskNum; num++ {
			pool.AddTask(taskFunc)
		}
		pool.Wait()
	}
	b.StopTimer()
}

func BenchmarkGoPoolWithSpinLock(b *testing.B) {
	pool := NewGoPool(PoolSize, WithLock(new(spinlock.SpinLock)))
	defer pool.Release()

	taskFunc := func() (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return nil, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for num := 0; num < TaskNum; num++ {
			pool.AddTask(taskFunc)
		}
		pool.Wait()
	}
	b.StopTimer()
}

func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup
	var taskNum = int(TaskNum)

	b.ResetTimer()
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

func BenchmarkPond(b *testing.B) {
	pool := pond.New(PoolSize, 0, pond.MinWorkers(PoolSize))

	taskFunc := func() {
		time.Sleep(10 * time.Millisecond)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < TaskNum; i++ {
			pool.Submit(taskFunc)
		}
		pool.StopAndWait()
	}
	b.StopTimer()
}

func BenchmarkAnts(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := ants.NewPool(PoolSize)
	defer p.Release()

	taskFunc := func() {
		time.Sleep(10 * time.Millisecond)
		wg.Done()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < TaskNum; i++ {
			wg.Add(1)
			_ = p.Submit(taskFunc)
		}
		wg.Wait()
	}
	b.StopTimer()
}
