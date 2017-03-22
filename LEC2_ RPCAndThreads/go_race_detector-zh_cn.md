Data Race Detector
# Golang 数据竞态检测(官方文档中英文翻译)
[Golang 数据竞态检测(官方文档中英文翻译)更多详情，点击前往 ](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/go_race_detector.md?dir=0&filepath=LEC2_+RPCAndThreads%2Fgo_race_detector.md&oid=650fa0ca8013b5fd9a420889a598da567798128b&sha=e1d3e2bbbc17a3e94b5d1af9ac20f1b841148d2e)

[英文原始链接，点击前往 ](https://golang.org/doc/articles/race_detector.html)

Introduction
Usage
Report Format
Options
Excluding Tests
How To Use
Typical Data Races
Race on loop counter
Accidentally shared variable
Unprotected global variable
Primitive unprotected variable
Supported Systems
Runtime Overhead
> 1. 介绍
> 1. 使用
> 1. 报告格式
> 1. 选项
> 1. 排除测试
> 1. 如何使用
> 1. 典型的数据竞态
>   1. 循环计数器的竞态示例
>   1. 意外的变量共享
>   1. 未保护的全局变量
>   1. 未保护的原始变量
> 1. 支持的系统
> 1. 运行时开销


Introduction
Data races are among the most common and hardest to debug types of bugs in concurrent systems.
A data race occurs when two goroutines access the same variable concurrently
and at least one of the accesses is a write. See the The Go Memory Model for details.
## 介绍
> 在并发系统中数据竞态是非常常见且很难调试的bug。数据竞态发生在两个goroutines并发的访问同一个变量，
> 且至少有一个访问时写。更新详情请参考Go内存模型[点击前往 ](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/go_mem_model.md?dir=0&filepath=LEC2_+RPCAndThreads%2Fgo_mem_model.md&oid=1d6554a60b8a59d0a6be126356e04dfc405bf6d0&sha=63f49486f3284c7d23cdabf8f38422c0e6b3004e)


Here is an example of a data race that can lead to crashes and memory corruption:
> 这里是数据竞态的一个例子,可以导致崩溃和内存损坏

func main() {
	c := make(chan bool)
	m := make(map[string]string)
	go func() {
		m["1"] = "a" // First conflicting access.
		c <- true
	}()
	m["2"] = "b" // Second conflicting access.
	<-c
	for k, v := range m {
		fmt.Println(k, v)
	}
}

Usage
## 使用

To help diagnose such bugs, Go includes a built-in data race detector. To use it, add the -race flag to the go command:
$ go test -race mypkg    // to test the package
$ go run -race mysrc.go  // to run the source file
$ go build -race mycmd   // to build the command
$ go install -race mypkg // to install the package
> 为了诊断这些bug，go包括一个内建额数据竞态检测工具，通过增加-race标志给go的命令行：
```
$ go test -race mypkg    // 测试包
$ go run -race mysrc.go  // 运行整个源文件
$ go build -race mycmd   // build命令
$ go install -race mypkg // 安装包
```

Report Format
# 报告格式
When the race detector finds a data race in the program, it prints a report.
The report contains stack traces for conflicting accesses,
as well as stacks where the involved goroutines were created. Here is an example:
> 如果race detector发现程序中有数据竞态，那么将会打印该竞态报告。
> 报告将会包含冲突访问的stack的路径，以及与stack相关的goroutine。
> 下面为报告示例：

```
WARNING: DATA RACE
Read by goroutine 185:
  net.(*pollServer).AddFD()
      src/net/fd_unix.go:89 +0x398
  net.(*pollServer).WaitWrite()
      src/net/fd_unix.go:247 +0x45
  net.(*netFD).Write()
      src/net/fd_unix.go:540 +0x4d4
  net.(*conn).Write()
      src/net/net.go:129 +0x101
  net.func·060()
      src/net/timeout_test.go:603 +0xaf

Previous write by goroutine 184:
  net.setWriteDeadline()
      src/net/sockopt_posix.go:135 +0xdf
  net.setDeadline()
      src/net/sockopt_posix.go:144 +0x9c
  net.(*conn).SetDeadline()
      src/net/net.go:161 +0xe3
  net.func·061()
      src/net/timeout_test.go:616 +0x3ed

Goroutine 185 (running) created at:
  net.func·061()
      src/net/timeout_test.go:609 +0x288

Goroutine 184 (running) created at:
  net.TestProlongTimeout()
      src/net/timeout_test.go:618 +0x298
  testing.tRunner()
      src/testing/testing.go:301 +0xe8
 ```

Options
## 选项
The GORACE environment variable sets race detector options. The format is:
GORACE="option1=val1 option2=val2"
> GORACE环境变量设置race detector。格式如下：
> GORACE="option1=val1 option2=val2"

The options are:
log_path (default stderr): The race detector writes its report to a file named log_path.pid.
The special names stdout and stderr cause reports to be written to standard output and standard error, respectively.
exitcode (default 66): The exit status to use when exiting after a detected race.
strip_path_prefix (default ""): Strip this prefix from all reported file paths, to make reports more concise.
history_size (default 1): The per-goroutine memory access history is 32K * 2**history_size elements.
Increasing this value can avoid a "failed to restore the stack" error in reports, at the cost of increased memory usage.
halt_on_error (default 0): Controls whether the program exits after reporting first data race.
```
选项如下：
1. log_path：默认为stderr，报告写入文件配置项
2. exitcode： 默认66，如果存在惊天，竞态检测后的设置退出状态。
3. strip_path_prefix： 默认为空“”，去掉带前缀的所有报告文件路径,使报道更简洁
4. history_size ：默认为1，per-goroutine内存访问历史是32 k * 2 * * history_size元素。
增加这个值可以避免“未能恢复堆栈”错误报告,内存使用成本增加。
5. halt_on_error：默认为0，控制在报告完第一个数据竞态之后是否退出程序。
```

Example:
$ GORACE="log_path=/tmp/race/report strip_path_prefix=/my/go/sources/" go test -race
> 竞态检测参数示例：
> GORACE="log_path=/tmp/race/report strip_path_prefix=/my/go/sources/" go test -race


Excluding Tests
## 排除测试

When you build with -race flag, the go command defines additional build tag race.
You can use the tag to exclude some code and tests when running the race detector.
> 当使用-race标志时，go的命令定义了额外的build竞态标记。你可以使用这个标记来执行代码及测试用例当你运行竞态检测时。
> 示例如下：

Some examples:
```
// +build !race

package foo

// The test contains a data race. See issue 123.
func TestFoo(t *testing.T) {
	// ...
}

// The test fails under the race detector due to timeouts.
func TestBar(t *testing.T) {
	// ...
}

// The test takes too long under the race detector.
func TestBaz(t *testing.T) {
	// ...
}
```


How To Use
## 如何使用

To start, run your tests using the race detector (go test -race).
The race detector only finds races that happen at runtime,
so it can't find races in code paths that are not executed.
If your tests have incomplete coverage,
you may find more races by running a binary built with -race under a realistic workload.
> 最开始，使用go test -race运行所有测试用例来检测竞态。
> 竞态检测仅仅能够发现运行时races.因此找不到没有执行的代码路径上是否有竞态存在。
> 如果你的测试没有完全覆盖，你可能发现更多的竞态通过运行二进制文件（通过-race构建），在真实环境中。

Typical Data Races
## 典型的数据竞态

Here are some typical data races. All of them can be detected with the race detector.
> 下面为一些典型的数据竞态示例。它们都能通过race detector检测出来。

Race on loop counter
### 循环计数器的竞态示例

```
func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i) // Not the 'i' you are looking for.
			wg.Done()
		}()
	}
	wg.Wait()
}
```

The variable i in the function literal is the same variable used by the loop,
so the read in the goroutine races with the loop increment. (This program typically prints 55555, not 01234.)
The program can be fixed by making a copy of the variable:
> 字面函数中的变量i，与循环中使用的变量i是相同的变量，因此i在变量中的读，与循环中的i变量增加形成了竞态。
> 程序将会打印5555，而不是01234.
> 程序能通过赋值一次变量而修复，如下例：
```
func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(j int) {
			fmt.Println(j) // Good. Read local copy of the loop counter.
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```


Accidentally shared variable
### 意外的变量共享

// ParallelWrite writes data to file1 and file2, returns the errors.
> 并行写入数据到file1及file2中，返回错误。
```
func ParallelWrite(data []byte) chan error {
	res := make(chan error, 2)
	f1, err := os.Create("file1")
	if err != nil {
		res <- err
	} else {
		go func() {
			// This err is shared with the main goroutine,此处的err使用了主线程中的err，因此发生竞态
			// so the write races with the write below.
			_, err = f1.Write(data)
			res <- err
			f1.Close()
		}()
	}
	f2, err := os.Create("file2") // The second conflicting write to err.
	if err != nil {
		res <- err
	} else {
		go func() {
			_, err = f2.Write(data)
			res <- err
			f2.Close()
		}()
	}
	return res
}
```

The fix is to introduce new variables in the goroutines (note the use of :=):
> bug修复为，在goroutine线程中，使用全新变量err即可
			...
			_, err := f1.Write(data)
			...
			_, err := f2.Write(data)
			...

Unprotected global variable
## 未保护的全局变量

If the following code is called from several goroutines,
it leads to races on the service map. Concurrent reads and writes of the same map are not safe:
> 如果下述代码被多个goroutine调用，将导致service map竞态.对同一个map进行并发读写是不安全的：

var service map[string]net.Addr

```
func RegisterService(name string, addr net.Addr) {
	service[name] = addr
}

func LookupService(name string) net.Addr {
	return service[name]
}
```

To make the code safe, protect the accesses with a mutex:
> 为了保证代码安全，使用互斥锁保护变量：

```
var (
	service   map[string]net.Addr
	serviceMu sync.Mutex
)

func RegisterService(name string, addr net.Addr) {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	service[name] = addr
}

func LookupService(name string) net.Addr {
	serviceMu.Lock()
	defer serviceMu.Unlock()
	return service[name]
}
```

Primitive unprotected variable
## 为保护的原始变量

Data races can happen on variables of primitive types as well (bool, int, int64, etc.), as in this example:
> 数据竞态可能会发生在原始类型变量上（bool，int，int64等），如下所示：

```
type Watchdog struct{ last int64 }

func (w *Watchdog) KeepAlive() {
	w.last = time.Now().UnixNano() // First conflicting access.
}

func (w *Watchdog) Start() {
	go func() {
		for {
			time.Sleep(time.Second)
			// Second conflicting access.
			if w.last < time.Now().Add(-10*time.Second).UnixNano() {
				fmt.Println("No keepalives for 10 seconds. Dying.")
				os.Exit(1)
			}
		}
	}()
}
```

Even such "innocent" data races can lead to hard-to-debug problems caused by non-atomicity of the memory accesses,
interference with compiler optimizations, or reordering issues accessing processor memory .
> 即使是这样的“无辜的”数据竞态也会导致很难debug发现问题(由于non-atomicity引起的内存访问),
> 由于编译器的优化，或者重排序问题访问处理器内存。

A typical fix for this race is to use a channel or a mutex.
To preserve the lock-free behavior, one can also use the sync/atomic package.
> 一个典型的修复该竞态问题是使用channel或者mutex。保持锁定的行为，也仍然可以使用sync/atomic包。

```
type Watchdog struct{ last int64 }

func (w *Watchdog) KeepAlive() {
	atomic.StoreInt64(&w.last, time.Now().UnixNano())
}

func (w *Watchdog) Start() {
	go func() {
		for {
			time.Sleep(time.Second)
			if atomic.LoadInt64(&w.last) < time.Now().Add(-10*time.Second).UnixNano() {
				fmt.Println("No keepalives for 10 seconds. Dying.")
				os.Exit(1)
			}
		}
	}()
}
```

Supported Systems
## 支持系统

The race detector runs on darwin/amd64, freebsd/amd64, linux/amd64, and windows/amd64.
> 竞态检测可以运行在darwin/amd64, freebsd/amd64, linux/amd64, and windows/amd64平台上。

Runtime Overhead
## 运行时开销

The cost of race detection varies by program,
but for a typical program,
memory usage may increase by 5-10x and execution time by 2-20x.
> 竞态检测的成本开销因项目而异，但是对于一个常规的项目，内存的使用可能会增加5%~10%，执行时间增加2%~20%。



