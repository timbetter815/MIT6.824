Q: When should we use Go channels as opposed to sync.Mutex?
## Q: 什么时候我们应该使用Go channels而不是sync.Mutex?

A: There is no clear answer to this question. If there's information
that multiple goroutines need to use, and those goroutines don't
specifically interact with each other, I tend to store the information
in a shared data structure protected by locks. I use channels for
producer-consumer interactions between goroutines, and when one
goroutine needs to wait for another goroutine to do something.
> 1. A: 没有明确的答案回答这个问题。如果有多个了goroutine需要使用共同的信息,
> 且这些goroutine之间没有相互作用,我倾向于使用locks来保护这些共享数据结构。
> 2. 通常使用channel来实现生产者及消费者之间的信息交互，以及适用于一个goroutine
> 需要等待另外一个goroutine去完成某些操作。

Q: Why does 6.824 use Go?
### 为什么6.824使用Go？

A: Go is a good fit for distibuted systems programming: it supports
concurrency well, it has a convenient RPC library, and it is
garbage-collected. There are other languages that might work as well
for 6.824, but Go was particularly fashionable when we were creating
this set of labs. 6.824 used C++ before Go; C++ worked pretty well,
but its lack of garbage collection made threaded code particularly
bug-prone.
> 1. A: Go非常适合分布式系统开发：对并发支持良好，拥有高效的RPC库，具有垃圾回收机制。
> 2. 当然也可能有其他语言适合于6.824，但是Go语言比较流行。
> 3. 6.824之前使用c++语言，c++表现的更好，但是缺少自动垃圾回收是的多线程编程更易于出现bug。

Q: How can I wait for goroutines to finish?
### 如何等待goroutine执行完毕？

A: Try sync.WaitGroup. Or you could create a channel, have each
goroutine send something on the channel when it finishes, and have the
main goroutine wait for the appropriate number of channel messages.
> 1. A:使用sync.WaitGroup。或者创建一个通道，当goroutine完成时候，发送完成信号给该通道，
> main主线程等待合适数量的channel信息。

Q: Why is map not thread-safe?
### Q：为什么map不是线程安全的？

A: That is a bit of surprise to everyone who starts using Go. My guess
is the designers didn't want programs to pay for locking for a map
that isn't shared by multiple goroutines.
> 1. A：开始接触go的人对此会感到有点奇怪。我的猜想是设置者并不希望程序为了map而花费较大开销，
> 因此map不能用于多线程之间共享。

Q: Why does Go have pointers, instead of just object references?
### 为什么Go拥有指针，而不是仅仅支持对象引用？

A: The designers of Go intended it to be used for systems programming,
where the programmer often wants control of how things are stored and
passed around. For example, Java passes integers, references, etc. by
value, which means that a callee cannot change the caller's value of
the integer, reference, etc. In Go, a caller can pass a pointer to a
variable to the callee, and the callee can use the pointer to modify
the caller's variable.
> 1. A:Go的设计者倾向于go为系统编程，程序经常需要控制存储及传递对象。
> 2. 例如：java传递integers、引用等等。如果通过值传递方式，意味着被调用者不能改变调用者传过来的(integer, reference等等)值。
> 3. 在Go语言中，调用者能够传递一个变量的指针给被调用者，这样被调用者能够修改调用者传递过来的值。