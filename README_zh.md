<div align="center">
</br>

<img src="./logo/gopool-logo-350.png" width="120">

# GoPool

[![PRs welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](https://makeapullrequest.com)
[![build and test](https://github.com/devchat-ai/gopool/workflows/CI/badge.svg)](https://github.com/devchat-ai/gopool/actions)
[![go report](https://goreportcard.com/badge/github.com/devchat-ai/gopool?style=flat)](https://goreportcard.com/report/github.com/devchat-ai/gopool)
[![release](https://img.shields.io/github/release/devchat-ai/gopool.svg)](https://github.com/devchat-ai/gopool/releases/)

| [English](README.md) | 中文 |
| --- | --- |

</div>

欢迎来到 GoPool，这是**一个95%的代码由GPT生成的项目**。你可以在[pro.devchat.ai](https://pro.devchat.ai)找到相应的 Commit 和 Prompt 列表。

GoPool 是一个用 Golang 实现的**高性能**、**功能丰富**、**简单易用**的工作池库。它会管理和回收一组 goroutine 来并发完成任务，从而提高你的应用程序的效率和性能。

## 性能测试

这个表格展示了三个 Go 库 GoPool、[ants](https://github.com/panjf2000/ants) 和 [pond](https://github.com/alitto/pond)的性能测试结果。表格包括每个库处理 100 万个任务所需的时间和内存消耗（以 MB 为单位）。

|     项目    | 处理一百万任务耗时 (s) | 内存消耗 (MB) |
|----------------|:----------------------------:|:-----------------------:|
| GoPool         | 1.13                         | 1.23                    |
| [ants](https://github.com/panjf2000/ants) | 1.43 | 9.49                 |
| [pond](https://github.com/alitto/pond)    | 3.51 | 1.88                 |

你可以通过运行下列命令在你的机器上测试 GoPool、ants 和 pond 的性能：

```bash
$ go test -benchmem -run=^$ -bench ^BenchmarkGoPoolWithMutex$ github.com/devchat-ai/gopool
$ go test -benchmem -run=^$ -bench ^BenchmarkAnts$ github.com/devchat-ai/gopool
$ go test -benchmem -run=^$ -bench ^BenchmarkPond$ github.com/devchat-ai/gopool
```

在我的电脑上，性能测试的结果如下：

- GoPool

```bash
$ go test -benchmem -run=^$ -bench ^BenchmarkGoPoolWithMutex$ github.com/devchat-ai/gopool
goos: darwin
goarch: arm64
pkg: github.com/devchat-ai/gopool
BenchmarkGoPoolWithMutex-10    	       1	1131753125 ns/op	1966192 B/op	 13609 allocs/op
PASS
ok  	github.com/devchat-ai/gopool	1.5085s
```

- ants

```bash
$ go test -benchmem -run=^$ -bench ^BenchmarkAnts$ github.com/devchat-ai/gopool
goos: darwin
goarch: arm64
pkg: github.com/devchat-ai/gopool
BenchmarkAnts-10    	       1	1425282750 ns/op	 9952656 B/op	   74068 allocs/op
PASS
ok  	github.com/devchat-ai/gopool	1.730s
```

- pond

```bash
$ go test -benchmem -run=^$ -bench ^BenchmarkPond$ github.com/devchat-ai/gopool
goos: darwin
goarch: arm64
pkg: github.com/devchat-ai/gopool
BenchmarkPond-10    	       1	3512323792 ns/op	 1288984 B/op	   11106 allocs/op
PASS
ok  	github.com/devchat-ai/gopool	3.946s
```

## 特性

<div align="center">
<img src="./logo/gopool.png" width="750">
</div>

- [x] **任务队列**：GoPool 使用一个线程安全的任务队列来存储等待处理的任务。多个工作器可以同时从这个队列中获取任务。任务队列的大小可配置。

- [x] **并发控制**：GoPool 可以控制并发任务的数量，防止系统过载。

- [x] **动态工作器调整**：GoPool 可以根据任务数量和系统负载动态调整工作器的数量。

- [x] **优雅关闭**：GoPool 可以优雅地关闭。当没有更多的任务或收到关闭信号时，它会停止接受新的任务，并等待所有进行中的任务完成后再关闭。

- [x] **任务错误处理**：GoPool 可以处理任务执行过程中出现的错误。

- [x] **任务超时处理**：GoPool 可以处理任务执行超时。如果一个任务在指定的超时期限内没有完成，该任务被认为失败，返回一个超时错误。

- [x] **任务结果获取**：GoPool 提供了一种获取任务结果的方式。

- [x] **任务重试**：GoPool 为失败的任务提供了重试机制。

- [x] **锁定制**：GoPool 支持不同类型的锁。你可以使用内置的`sync.Mutex`或自定义锁，如`spinlock.SpinLock`。

- [ ] **任务优先级**：GoPool 支持任务优先级。优先级更高的任务会被优先处理。

## 安装

要安装GoPool，使用`go get`：

```bash
go get -u github.com/devchat-ai/gopool
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
        pool.AddTask(func() (interface{}, error){
            time.Sleep(10 * time.Millisecond)
            return nil, nil
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
        pool.AddTask(func() (interface{}, error){
            time.Sleep(10 * time.Millisecond)
            return nil, nil
        })
    }
    pool.Wait()
}
```

## 配置任务队列大小

GoPool 使用一个线程安全的任务队列来存储等待处理的任务。多个工作器可以同时从这个队列中获取任务。任务队列的大小可配置。可以通过在创建池时设置 `WithQueueSize` 选项来配置任务队列的大小。

这是一个如何配置 GoPool 任务队列大小的示例：

```go
package main

import (
    "time"

    "github.com/devchat-ai/gopool"
)

func main() {
    pool := gopool.NewGoPool(100, gopool.WithTaskQueueSize(5000))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error){
            time.Sleep(10 * time.Millisecond)
            return nil, nil
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
        pool.AddTask(func() (interface{}, error){
            time.Sleep(10 * time.Millisecond)
            return nil, nil
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

## 任务重试

GoPool 支持任务重试。如果任务失败，可以重试指定的次数。可以通过在创建池时设置 `WithRetryCount` 选项来启用此功能。

以下是如何使用带有任务重试的 GoPool 的示例：

```go
package main

import (
    "errors"
    "fmt"

    "github.com/devchat-ai/gopool"
)

func main() {
    pool := gopool.NewGoPool(100, gopool.WithRetryCount(3))
    defer pool.Release()

    for i := 0; i < 1000; i++ {
        pool.AddTask(func() (interface{}, error) {
            return nil, errors.New("task error")
        })
    }
    pool.Wait()
}
```

在这个示例中，如果任务失败，它将重试最多3次。
