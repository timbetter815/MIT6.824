The Go Memory Model
# Golang内存模型(官方文档中英文翻译）
[点击：更多详情请前往博客地址](https://my.oschina.net/tantexian)
[点击：或者git地址](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/go_mem.md?dir=0&filepath=LEC2_+RPCAndThreads%2Fgo_mem.md&oid=5fca03003d6e62f098716d7bfbd9fb7be91ec3d1&sha=4d19b0880fadc8c2a63696ec736cac5eafa99c03)
[点击: Go Memory Model 示例代码 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/6.824-golabs-2017/src/mytest/mem_test.go?dir=0&filepath=6.824-golabs-2017%2Fsrc%2Fmytest%2Fmem_test.go&oid=180dc8be3f38c1996293098eb8a55e6bdbcd51b9&sha=473914894c0eae52cd3904b4d110d3ff4be9cd36)

Version of May 31, 2014

Introduction
Advice
Happens Before
Synchronization
Initialization
Goroutine creation
Goroutine destruction
Channel communication
Locks
Once
Incorrect synchronization
> 1. 介绍
> 2. 建议
> 3. Happens Before原则
> 4. 同步
> 5. 初始化
> 6. Goroutine创建
> 7. Goroutine销毁
> 8. 通道通信
> 9. 锁
> 10. 一次
> 11. 不正确的同步

Introduction
## 介绍
The Go memory model specifies the conditions under which reads of a variable in one goroutine
can be guaranteed to observe values produced by writes to the same variable in a different goroutine.
### Go内存模型指出一个goroutine在读取一个变量时候，能够观察到其他线程对相同变量写的结果值。

Advice
## 建议

Programs that modify data being simultaneously accessed by multiple goroutines must serialize such access.

To serialize access, protect the data with channel operations or other synchronization primitives such as those in the sync and sync/atomic packages.

If you must read the rest of this document to understand the behavior of your program, you are being too clever.

Don't be clever.
> 程序中多个goroutine同时修改相同数据，必须使之能够序列化访问。
> 为了序列化访问：保护数据安全可以使用channel操作或者其他同步原语例如sync及sync/atomic包。
> 如果你必须阅读本文档的其余部分来理解您的程序的行为,这样将是很聪明的做法。当然也不要太聪明。


Happens Before
## Happens Before原则

Within a single goroutine, reads and writes must behave as if they executed in the order specified by the program.
That is, compilers and processors may reorder the reads and writes executed within a single goroutine
only when the reordering does not change the behavior within that goroutine as defined by the language specification.
Because of this reordering, the execution order observed by one goroutine may differ from the order perceived by another.
For example, if one goroutine executes a = 1; b = 2;, another might observe the updated value of b before the updated value of a.
> 在单个goroutine,读写必须表现得好像他们指定的程序的顺序执行一致。
> 也就是说,编译器和处理器可能会对读写进行重新排序，在一个执行goroutine中，
> 只有当重新排序并不会改变goroutine在语言规范定义行为时，才会进行重排序。
> 因为重新排序之后,执行顺序在一个goroutine可能不同于另一个goroutine中观察会保持一致。
> 举个例子,如果一个goroutine执行a = 1,b = 2; 另一个可以观察b的值更新，在a的值更新之前。

To specify the requirements of reads and writes, we define happens before,
a partial order on the execution of memory operations in a Go program.
If event e1 happens before event e2, then we say that e2 happens after e1.
Also, if e1 does not happen before e2 and does not happen after e2,
then we say that e1 and e2 happen concurrently.
> 为了详细说明Go程序的读写在内存中的偏序关系，我们定义了happen-before原则。
> 如果时间e1发生在事件e2之前，那么则e2发生在e1之后。
> 如果e1既不发生在e2之前也不发生在它之后，那么则说明e1和e2同时发生。

Within a single goroutine, the happens-before order is the order expressed by the program.
> 在一个单goroutine中，happen-before的顺序与程序中定义的顺序保持一致。

A read r of a variable v is allowed to observe a write w to v if both of the following hold:
1、r does not happen before w.
2、There is no other write w' to v that happens after w but before r.
> 如果满足下述两个规则：那么对变量v的读操作r能够观察到对变量v的写操作w：
> 1. r没有发生在w之前
> 2. 没有其他针对变量v的写操作发生在w操作之后，r之前

To guarantee that a read r of a variable v observes a particular write w to v,
ensure that w is the only write r is allowed to observe.
That is, r is guaranteed to observe w if both of the following hold:
1、w happens before r.
2、Any other write to the shared variable v either happens before w or after r.
> 为了保证对变量v的读操作能够观察到对变量v的特定的写操作w的值，确保w是唯一的写操作，r的变化允许被观察到。
> 那么，如果满足下述两个原则，r将能够观察到w
> 1. w发生在r之前
> 2. 任何对共享变量v的写操作发生在写操作w之前或者读操作r之后

This pair of conditions is stronger than the first pair;
it requires that there are no other writes happening concurrently with w or r.
> 这对条件强于第一条；它需要没有其他任何的写操作与读操作w及写操作r同时发生。

Within a single goroutine, there is no concurrency, so the two definitions are equivalent:
a read r observes the value written by the most recent write w to v.
When multiple goroutines access a shared variable v,
they must use synchronization events to establish happens-before conditions that ensure reads observe the desired writes.
> 在单一的goroutine中，由于不可能同时发生，因此两个定义也是等价的：
> 针对变量v的r操作能够观察到最近对改变了的写操作
> 当多个goroutine同时访问一个共享变量v时，必须使用sync同步事件建立happen-before条件来确保读操作观察到期望的写操作的变化。

The initialization of variable v with the zero value for v's type behaves as a write in the memory model.
> 变量v的初始化以零值作为默认值的类型的行为作为一个写在内存模型。

Reads and writes of values larger than a single machine word
behave as multiple machine-word-sized operations in an unspecified order.
> 如果读取和写入的值大于一个机器字,那么多个多机器字节变量的操作发生在一个未指明的顺序。

Synchronization
## 同步

Initialization
## 初始化

Program initialization runs in a single goroutine,
but that goroutine may create other goroutines, which run concurrently.
> 程序初始化运行在单个goroutine上，但是goroutine可能创建其他goroutine，并发运行。

If a package p imports package q, the completion of q's init functions
happens before the start of any of p's.
> 如果一个包p导入了一个包q,q的初始化函数的结束先于任何p的开始完成。

The start of the function main.main happens after all init functions have finished.
> 函数main.mian在所有init初始化函数的完成后发生。

Goroutine creation
## Goroutine创建

The go statement that starts a new goroutine happens before the goroutine's execution begins.
> 用于启动goroutine的go语句在goroutine之前运行。

For example, in this program:
```
var a string

func f() {
	print(a)
}

func hello() {
	a = "hello, world"
	go f()
}
```
calling hello will print "hello, world" at some point in the future (perhaps after hello has returned).
> 调用hello函数，会在某个时刻打印“hello, world”（有可能是在hello函数返回之后）。


Goroutine destruction
## Goroutine销毁

The exit of a goroutine is not guaranteed to happen before any event in the program.
> goroutine的退出不保证一定会先于其他事件发生。

For example, in this program:
```
var a string

func hello() {
	go func() { a = "hello" }()
	print(a) // 由于没有使用同步相关事件，因此有可能a的赋值对其他线程是不可见的
}
```
the assignment to a is not followed by any synchronization event,
so it is not guaranteed to be observed by any other goroutine.
In fact, an aggressive compiler might delete the entire go statement.
> 对a的分配赋值没有跟随任何同步事件之后, 所以不能保证会被其他goroutine观察到。
事实上,一个优秀的编译器可能会删除整个go声明语句。

If the effects of a goroutine must be observed by another goroutine,
use a synchronization mechanism such as a lock or channel communication to establish a relative ordering.
> 如果想让一个goroutine的执行操作影响必须对另外一个goroutine可见，
可以使用同步机制例如lock或者channel通信来确定相关顺序。

Channel communication
## 通道通信

Channel communication is the main method of synchronization between goroutines.
Each send on a particular channel is matched to a corresponding receive from that channel,
usually in a different goroutine.
> 通道通信是主要的方法来解决多个goroutine之间的同步。
每发送一个特定的信道匹配相应的接收通道，通常在不同的goroutine中。


A send on a channel happens before the corresponding receive from that channel completes.
> 一个通道的发送发生与相应的通道接收之前完成

This program:
````
var c = make(chan int, 10)
var a string

func f() {
	a = "hello, world"
	c <- 0
}

func main() {
	go f()
	<-c
	print(a)
}
````
is guaranteed to print "hello, world".
The write to a happens before the send on c,
which happens before the corresponding receive on c completes,
which happens before the print.
> 确保一定能够打印出“hello，world”。
对a的写先于对c通道信息的发送，c通道的发送一定会先于c通道接受的完成，因此这些都先于print完成。

The closing of a channel happens before a receive that returns a zero value
because the channel is closed.
> 关闭一个通道发生之前收到返回零值 因为通道是关闭的

In the previous example, replacing c <- 0 with close(c) yields a program
with the same guaranteed behavior.
> 在上一个示例中,用close(c)取代c < - 0,同样能够保证相同的行为

A receive from an unbuffered channel happens before the send on
that channel completes.
> 从无缓冲channel中接收操作 发生在该channel的发送操作结束前

This program (as above, but with the send and receive statements swapped
and using an unbuffered channel):
````
var c = make(chan int)
var a string

func f() {
	a = "hello, world"
	<-c
}
func main() {
	go f()
	c <- 0
	print(a)
}
````
is also guaranteed to print "hello, world".
The write to a happens before the receive on c,
which happens before the corresponding send on c completes,
which happens before the print.
> 任然能够保证一定会输出“hello，world”
对a的写操作先于通道c的信号接收，这也发生在相应的发送完成之前，当然也先于print的发生。

If the channel were buffered (e.g., c = make(chan int, 1))
then the program would not be guaranteed to print "hello, world".
(It might print the empty string, crash, or do something else.)
> 如果是缓存通道，那么程序则不一定能保证会打印“hello，world”。（可能打印空字符串或者crash、或者其他未知情况）

The kth receive on a channel with capacity C
happens before the k+Cth send from that channel completes.
> kth收到一个通道先于k+Cth发送通道完成.

This rule generalizes the previous rule to buffered channels.
It allows a counting semaphore to be modeled by a buffered channel:
the number of items in the channel corresponds to the number of active uses,
the capacity of the channel corresponds to the maximum number of simultaneous uses,
sending an item acquires the semaphore, and receiving an item releases the semaphore.
This is a common idiom for limiting concurrency.
> 这些定义了缓存channel通道前序规则。
允许一个计数信号量通过缓存通道来模拟：
通道缓存数量与活跃用户保存一致，通道的最大容量与并发同时执行的最大人数保持一致，发送一个item获取信号量及接受一个item释放信号量。这个一个惯用的方式来限制并发。

This program starts a goroutine for every entry in the work list,
but the goroutines coordinate using the limit channel to
ensure that at most three are running work functions at a time.
> 示例：程序为工作列表中的每一个条目启动一个goroutine，但是各个goroutine之间使用受限制的缓存通道来确保最多只能有三个运行work在同时执行。

````
var limit = make(chan int, 3)

func main() {
	for _, w := range work {
		go func(w func()) {
			limit <- 1
			w()
			<-limit
		}(w)
	}
	select{}
}
````
## Locks

The sync package implements two lock data types, sync.Mutex and sync.RWMutex.
> sync包实现了两种类型的锁：互斥锁及读写锁

For any sync.Mutex or sync.RWMutex variable l and n < m, call n of l.Unlock()
happens before call m of l.Lock() returns.
> 对于任意 sync.Mutex 或 sync.RWMutex 变量l。 如果 n < m ，那么第n次 l.Unlock() 调用在第 m次 l.Lock()调用返回前发生。

This program:
````
// 对同一个变量l的解锁n次和加锁操作m次，如果n<m,即解锁次数1小于加锁次数2，因此解锁操作先于加锁操作完成
// 即f函数的unlock先于第二个lock完成，而第二个lock又先于print完成，因此保证a1的赋值能被打印出来
var l sync.Mutex
var a string

func f() {
	a = "hello, world"
	l.Unlock()
}

func main() {
	l.Lock()
	go f()
	l.Lock()
	print(a)
}

````
is guaranteed to print "hello, world". The first call to l.Unlock() (in f)
happens before the second call to l.Lock() (in main) returns, which happens before the print.
> 可以确保输出“hello, world”结果。因为，第一次 l.Unlock() 调用（在f函数中）在第二次 l.Lock() 调用（在main 函数中）返回之前发生，
也就是在 print 函数调用之前发生。

For any call to l.RLock on a sync.RWMutex variable l,
there is an n such that the l.RLock happens (returns) after call n to l.Unlock
and the matching l.RUnlock happens before call n+1 to l.Lock.
> 对于任何调用变量l的RLock，第n次调用l.RLock必定发生在第n次l.Unlock之后,
与之相对应的l.RUnlock必定发生在第n+1次l.Lock之前。

## Once
## 一次执行

The sync package provides a safe mechanism for initialization in the presence of multiple goroutines
through the use of the Once type. Multiple threads can execute once.Do(f) for a particular f,
but only one will run f(), and the other calls block until f() has returned.
> sync包提供了一个安全的机制保证多个goroutine的初始化只被执行一次。多个线程都会调用执行once.Do(f)函数，但是f函数只会执行一次。

A single call of f() from once.Do(f) happens (returns) before any call of once.Do(f) returns.
> 通过once.Do(f)单次调用f()的返回先于其他调用once.Do(f)的返回。

In this program:
````
var a string
var once sync.Once

func setup() {
	a = "hello, world"
}

func doprint() {
	once.Do(setup)
	print(a)
}

func twoprint() {
	go doprint()
	go doprint()
}
````
calling twoprint causes "hello, world" to be printed twice.
The first call to doprint runs setup once.
> twoprint函数将会两次打印输出，通过once.Do(setup)调用的setup只会有一次输出。

Incorrect synchronization
## 错误的同步

Note that a read r may observe the value written by a write w that happens concurrently with r.
Even if this occurs,it does not imply that reads happening after r will observe writes that happened before w.
> [错误认为:] 如果读操作和写操作同时发生，读操作r可以观察到写操作w的变化。
即使发生这种情况,它也并不意味着即使读操作发生在写操作r之后就能观察到写入值writes就发生在w之前。

In this program:
````
var a, b int

func f() {
	a = 1
	b = 2
}

func g() {
	print(b)
	print(a)
}

func main() {
	go f()
	g()
}
````
it can happen that g prints 2 and then 0.
> 有可能打印出2和0（即b已经赋值为2，但是a还未被赋值）

This fact invalidates a few common idioms.
> 下述示例为一些无效错误的惯例。

Double-checked locking is an attempt to avoid the overhead of synchronization.
For example, the twoprint program might be incorrectly written as:
> 双检锁是为了避免同步的开销，例如,twoprint程序可能会错误地写成：

```
var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func doprint() {
	if !done {
		once.Do(setup)
	}
	print(a)
}

func twoprint() {
	go doprint()
	go doprint()
}
```

but there is no guarantee that, in doprint, observing the write to done implies observing the write to a.
This version can (incorrectly) print an empty string instead of "hello, world".
> 这里不能保证，在doprint函数里面，能够观察到变量done的值，就暗示着能够观察到对a的赋值变化。
这里可能会错误的打印空，而不是“hello, world”。

Another incorrect idiom is busy waiting for a value, as in:
> 下述为另外一个错误示例（一直忙等）

````
var a string
var done bool

func setup() {
	a = "hello, world"
	done = true
}

func main() {
	go setup()
	for !done {
	}
	print(a)
}
````
As before, there is no guarantee that, in main,
observing the write to done implies observing the write to a,
so this program could print an empty string too.
Worse, there is no guarantee that the write to done will ever be observed by main,
since there are no synchronization events between the two threads.
The loop in main is not guaranteed to finish.
> 上述示例所示，在主函数main中，观察到done值为true并不能保证一定能看到对a的赋值变化, 因此本示例程序也有可能打印空值。
> 更糟糕的是，setup函数对done=true的赋值，并不一定能被main函数观察到，因此这两个线程之间并没有使用任何的同步事件。
> 因此for循环在主函数中并不能保证一定会循环结束。

There are subtler variants on this theme, such as this program.
> 下述为一些更为难以发现的变体(错误示例)

````
type T struct {
	msg string
}

var g *T

func setup() {
	t := new(T)
	t.msg = "hello, world"
	g = t
}

func main() {
	go setup()
	for g == nil {
	}
	print(g.msg)
}
````
Even if main observes g != nil and exits its loop,
there is no guarantee that it will observe the initialized value for g.msg.
> 即使主函数观察到g!=null且退出了循环，也不能保证main能够观察到g.msg的赋值内容。

In all these examples, the solution is the same: use explicit synchronization.
> 对于所有的示例，都可以采用同样的解决办法：使用明确的同步。