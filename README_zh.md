
# GoPool

欢迎来到 GoPool，这是**一个95%的代码由GPT生成的项目**。你可以在[pro.devchat.ai](https://pro.devchat.ai)找到相应的 Commit 和 Prompt 列表。

GoPool 是一个用 Golang 实现的高性能、功能丰富、易于使用的工作池库。它会管理和回收一组 goroutine 来并发完成任务，从而提高你的应用程序的效率和性能。

## 特性

- [x] **任务队列**：GoPool 使用一个线程安全的任务队列来存储等待处理的任务。多个工作器可以同时从这个队列中获取任务。

- [x] **并发控制**：GoPool 可以控制并发任务的数量，防止系统过载。

- [x] **动态工作器调整**：GoPool 可以根据任务数量和系统负载动态调整工作器的数量。

- [x] **优雅的关闭**：GoPool 可以优雅地关闭。当没有更多的任务或收到关闭信号时，它会停止接受新的任务，并等待所有进行中的任务完成后再关闭。

- [x] **任务错误处理**：GoPool 可以处理任务执行过程中出现的错误。

- [x] **任务超时处理**：GoPool 可以处理任务执行超时。如果一个任务在指定的超时期限内没有完成，该任务被认为失败，返回一个超时错误。

- [x] **任务结果获取**：GoPool 提供了一种获取任务结果的方式。

- [x] **任务重试**：GoPool 为失败的任务提供了重试机制。

- [x] **锁定定制**：GoPool 支持不同类型的锁。你可以使用内置的`sync.Mutex`或自定义锁，如`spinlock.SpinLock`。

- [ ] **任务优先级**：GoPool 支持任务优先级。优先级更高的任务会被优先处理。

## 安装

要安装GoPool，使用`go get`：

```bash
go get github.com/devchat-ai/gopool@v0.2.0
```

## 使用

这是一个如何使用带有 `sync.Mutex` 的GoPool 的简单示例：

```go
package main

import (
    "sync"
    "time"

    "github.com/devchat-ai/gopool"
)

func main() {
    pool := gopool.NewGoPool(100)
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() {
            time.Sleep(10 * time.Millisecond)
        })
    }
    pool.Wait()
}
```

这是如何使用带有 `spinlock.SpinLock` 的 GoPool 的示例：

```go
package main

import (
    "time"

    "github.com/daniel-hutao/spinlock"
    "github.com/devchat-ai/gopool"
)

func main() {
    pool := gopool.NewGoPool(100, gopool.WithLock(new(spinlock.SpinLock)))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() {
            time.Sleep(10 * time.Millisecond)
        })
    }
    pool.Wait()
}
```

## 动态工作器调整

GoPool 支持动态工作器调整。这意味着池中的工作器数量可以根据队列中的任务数量增加或减少。可以通过在创建池时设置 MinWorkers 选项来启用此功能。

这是如何使用动态工作器调整的 GoPool 的示例：

```go
package main

import (
    "time"

    "github.com/devchat-ai/gopool"
)

func main() {
    pool := gopool.NewGoPool(100, gopool.WithMinWorkers(50))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() {
            time.Sleep(10 * time.Millisecond)
        })
    }
    pool.Wait()
}
```

在这个示例中，池开始时有50个工作器。如果队列中的任务数量超过(MaxWorkers - MinWorkers) / 2 + MinWorkers，池将添加更多的工作器。如果队列中的任务数量少于 MinWorkers，池将移除一些工作器。

## 任务超时处理

GoPool支持任务超时。如果一个任务花费的时间超过指定的超时时间，它将被取消。可以通过在创建池时设置 `WithTimeout` 选项来启用此功能。

这是如何使用任务超时的 GoPool 的示例：

```go
package main

import (
    "time"

    "github.com/devchat-ai/gopool"
)

func main() {
    pool := gopool.NewGoPool(100, gopool.WithTimeout(1*time.Second))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error) {
            time.Sleep(2 * time.Second)
            return nil, nil
        })
    }
    pool.Wait()
}
```

在这个示例中，如果任务花费的时间超过1秒，任务将被取消。

## 任务错误处理

GoPool 支持任务错误处理。如果一个任务返回一个错误，错误回调函数将被调用。可以通过在创建池时设置 `WithErrorCallback` 选项来启用此功能。

这是如何使用错误处理的 GoPool 的示例：

```go
package main

import (
    "errors"
    "fmt"

    "github.com/devchat-ai/gopool"
)

func main() {
    pool := gopool.NewGoPool(100, gopool.WithErrorCallback(func(err error) {
        fmt.Println("Task error:", err)
    }))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error) {
            return nil, errors.New("task error")
        })
    }
    pool.Wait()
}
```

在这个示例中，如果一个任务返回一个错误，错误将被打印到控制台。

## 任务结果获取

GoPool 支持任务结果获取。如果一个任务返回一个结果，结果回调函数将被调用。可以通过在创建池时设置 `WithResultCallback` 选项来启用此功能。

这是如何使用任务结果获取的 GoPool 的示例：

```go
package main

import (
    "fmt"

    "github.com/devchat-ai/gopool"
)

func main() {
    pool := gopool.NewGoPool(100, gopool.WithResultCallback(func(result interface{}) {
        fmt.Println("Task result:", result)
    }))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error) {
            return "task result", nil
        })
    }
    pool.Wait()
}
```

在这个示例中，如果一个任务返回一个结果，结果将被打印到控制台。

## 性能测试

我们进行了几个性能测试来评估 GoPool 的效率和性能。以下是结果：

- **TestGoPoolWithMutex**：

```bash
$ go test -benchmem -run=^$ -bench ^BenchmarkGoPoolWithMutex$ github.com/devchat-ai/gopool

goos: darwin
goarch: arm64
pkg: github.com/devchat-ai/gopool
=== RUN   BenchmarkGoPoolWithMutex
BenchmarkGoPoolWithMutex
BenchmarkGoPoolWithMutex-10            2         803105167 ns/op        17416408 B/op    1017209 allocs/op
PASS
ok      github.com/devchat-ai/gopool    2.586s
```

- **TestGoPoolWithSpinLock**：

```bash
$ go test -benchmem -run=^$ -bench ^BenchmarkGoPoolWithSpinLock$ github.com/devchat-ai/gopool

goos: darwin
goarch: arm64
pkg: github.com/devchat-ai/gopool
=== RUN   BenchmarkGoPoolWithSpinLock
BenchmarkGoPoolWithSpinLock
BenchmarkGoPoolWithSpinLock-10                 2         662952562 ns/op        17327176 B/op    1016087 allocs/op
PASS
ok      github.com/devchat-ai/gopool    2.322s
```

- **BenchmarkGoroutines**： 

```bash
$ go test -benchmem -run=^$ -bench ^BenchmarkGoroutines$ github.com/devchat-ai/gopool

goos: darwin
goarch: arm64
pkg: github.com/devchat-ai/gopool
=== RUN   BenchmarkGoroutines
BenchmarkGoroutines
BenchmarkGoroutines-10                 3         371622847 ns/op        96642458 B/op    2005219 allocs/op
PASS
ok      github.com/devchat-ai/gopool    2.410s
```

请注意，实际性能可能会根据具体的使用情况和系统环境而变化。
