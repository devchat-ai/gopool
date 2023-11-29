package gopool_test

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/daniel-hutao/spinlock"
	"github.com/devchat-ai/gopool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gopool", func() {
	Describe("With Mutex", func() {
		It("should work correctly", func() {
			pool := gopool.NewGoPool(100, gopool.WithLock(new(sync.Mutex)))
			defer pool.Release()
			for i := 0; i < 1000; i++ {
				pool.AddTask(func() (interface{}, error) {
					time.Sleep(10 * time.Millisecond)
					return nil, nil
				})
			}
			pool.Wait()
		})
	})

	Describe("With SpinLock", func() {
		It("should work correctly", func() {
			pool := gopool.NewGoPool(100, gopool.WithLock(new(spinlock.SpinLock)))
			defer pool.Release()
			for i := 0; i < 1000; i++ {
				pool.AddTask(func() (interface{}, error) {
					time.Sleep(10 * time.Millisecond)
					return nil, nil
				})
			}
			pool.Wait()
		})
	})

	Describe("With Error", func() {
		It("should work correctly", func() {
			var errTaskError = errors.New("task error")
			pool := gopool.NewGoPool(100, gopool.WithErrorCallback(func(err error) {
				Expect(err).To(Equal(errTaskError))
			}))
			defer pool.Release()

			for i := 0; i < 1000; i++ {
				pool.AddTask(func() (interface{}, error) {
					return nil, errTaskError
				})
			}
			pool.Wait()
		})
	})

	Describe("With Result", func() {
		It("should work correctly", func() {
			var expectedResult = "task result"
			pool := gopool.NewGoPool(100, gopool.WithResultCallback(func(result interface{}) {
				Expect(result).To(Equal(expectedResult))
			}))
			defer pool.Release()

			for i := 0; i < 1000; i++ {
				pool.AddTask(func() (interface{}, error) {
					return expectedResult, nil
				})
			}
			pool.Wait()
		})
	})

	Describe("With Retry", func() {
		It("should work correctly", func() {
			var retryCount = int32(3)
			var taskError = errors.New("task error")
			var taskRunCount int32 = 0

			pool := gopool.NewGoPool(100, gopool.WithRetryCount(int(retryCount)))
			defer pool.Release()

			pool.AddTask(func() (interface{}, error) {
				atomic.AddInt32(&taskRunCount, 1)
				if taskRunCount <= retryCount {
					return nil, taskError
				}
				return nil, nil
			})

			pool.Wait()

			Expect(atomic.LoadInt32(&taskRunCount)).To(Equal(retryCount + 1))
		})
	})

	Describe("With Timeout", func() {
		It("should work correctly", func() {
			var taskRun int32

			pool := gopool.NewGoPool(100, gopool.WithTimeout(100*time.Millisecond), gopool.WithErrorCallback(func(err error) {
				Expect(err.Error()).To(Equal("task timed out"))
				atomic.StoreInt32(&taskRun, 1)
			}))
			defer pool.Release()

			pool.AddTask(func() (interface{}, error) {
				time.Sleep(200 * time.Millisecond)
				return nil, nil
			})

			pool.Wait()

			Expect(atomic.LoadInt32(&taskRun)).To(Equal(int32(1)))
		})
	})

	Describe("With Timeout and Result Hook", func() {
		It("should work correctly", func() {
			var taskRun int32

			pool := gopool.NewGoPool(100, gopool.WithTimeout(100*time.Millisecond), gopool.WithErrorCallback(func(err error) {
				if err != nil {
					Expect(err.Error()).To(Equal("task timed out"))
					atomic.StoreInt32(&taskRun, 1)
				}
			}), gopool.WithResultCallback(func(i interface{}) {
				println("result callback",i)
				if v, ok := i.(int); ok {
					Expect(v).To(Equal(1))
				}
			}))
			defer pool.Release()

			pool.AddTask(func() (interface{}, error) {
				time.Sleep(200 * time.Millisecond)
				return nil, nil
			})

			pool.AddTask(func() (interface{}, error) {
				return 1, nil
			})

			pool.Wait()

			Expect(atomic.LoadInt32(&taskRun)).To(Equal(int32(1)))
		})
	})

	Describe("With MinWorkers", func() {
		It("should work correctly", func() {
			var minWorkers = 50

			pool := gopool.NewGoPool(100, gopool.WithMinWorkers(minWorkers))
			defer pool.Release()

			Expect(pool.GetWorkerCount()).To(Equal(minWorkers))
		})
	})

	Describe("With TaskQueueSize", func() {
		It("should work correctly", func() {
			size := 5000
			pool := gopool.NewGoPool(100, gopool.WithTaskQueueSize(size))
			defer pool.Release()

			Expect(pool.GetTaskQueueSize()).To(Equal(size))
		})
	})
})
