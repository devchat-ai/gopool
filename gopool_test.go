package gopool

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/daniel-hutao/spinlock"
)

func TestGoPoolWithMutex(t *testing.T) {
	pool := NewGoPool(100, WithLock(new(sync.Mutex)))
	defer pool.Release()
	for i := 0; i < 1000; i++ {
		pool.AddTask(func() (interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			return nil, nil
		})
	}
	pool.Wait()
}

func TestGoPoolWithSpinLock(t *testing.T) {
	pool := NewGoPool(100, WithLock(new(spinlock.SpinLock)))
	defer pool.Release()
	for i := 0; i < 1000; i++ {
		pool.AddTask(func() (interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			return nil, nil
		})
	}
	pool.Wait()
}

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

func TestGoPoolWithError(t *testing.T) {
	var errTaskError = errors.New("task error")
	pool := NewGoPool(100, WithErrorCallback(func(err error) {
		if err != errTaskError {
			t.Errorf("Expected error %v, but got %v", errTaskError, err)
		}
	}))
	defer pool.Release()

	for i := 0; i < 1000; i++ {
		pool.AddTask(func() (interface{}, error) {
			return nil, errTaskError
		})
	}
	pool.Wait()
}

func TestGoPoolWithResult(t *testing.T) {
	var expectedResult = "task result"
	pool := NewGoPool(100, WithResultCallback(func(result interface{}) {
		if result != expectedResult {
			t.Errorf("Expected result %v, but got %v", expectedResult, result)
		}
	}))
	defer pool.Release()

	for i := 0; i < 1000; i++ {
		pool.AddTask(func() (interface{}, error) {
			return expectedResult, nil
		})
	}
	pool.Wait()
}
