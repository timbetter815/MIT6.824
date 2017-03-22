6.824 2017 Lecture 2: Infrastructure: RPC and threads
# 6.824 2017课程2：基础：RPC及threads

Most commonly-asked question: Why Go?
  6.824 used to use C++
    students spent time fixing bugs unrelated to distributed systems
      e.g., they freed objects that were still in use
  Go allows you to concentrate on distributed systems problems
    good support for concurrency
    good support for RPC
    garbage-collected (no use after freeing problems)
    type safe
  We like programming in Go
    simple language to learn
    a fashionable language when we redid the 6.824 labs
  After the tutorial, use https://golang.org/doc/effective_go.html
  Russ Cox will give a guest lecture 3/16


## 经常被问到的问题：为什么选择Go？
1. 6.824之前使用c++
    1. 学生们花了好大一部分时间在修复bug、而非分布式系统本身
    2. 例如：使用已经被释放的对象

2. Go允许大家集中精力在分布式系统本身的那些问题上
    1. 很好的支持并发
    2. 很好的支持RPC
    3. 自动垃圾回收
    4. 类型安全

3. 我们更喜欢使用Go编程
    1. 语言学习起来简单
    2. 使用更加流行的语言来完成6.824实验
    3. 在本教程中，使用[Effective Go](https://golang.org/doc/effective_go.html)（也可以在本项目MIT6.824/docs下找到中文版文档：Effective Go中文版.pdf）


Threads
  threads are a useful structuring tool
  Go calls them goroutines; everyone else calls them threads
  they can be tricky
### Threads
+ 线程是有用的结构化工具
+ Go将至称为goroutines，其他语言叫threads
+ 它们相对比较复杂


Why threads?
  Allow you to exploit concurrency, which shows up naturally in distributed systems
  I/O concurrency:
    While waiting for a response from another server, process next request
  Multicore:
    Threads run in parallel on several cores
### 为什么使用Threads
1. 允许大家开发并发程序，这些自然而然的在分布式系统中会出现
2. I/O并发：
    1. 当在等待其他server的response时，可以处理下一个请求
3. 多核：
    1. 线程并行运行在多个核上

Thread = "thread of execution"
  threads allow one program to (logically) execute many things at once
  the threads share memory
  each thread includes some per-thread state:
    program counter, registers, stack
### Thread = "thread of execution"
1. 线程允许一个程序在同一时间执行多个事情
2. 多个线程共享内存
3. 每个线程包括一些自己线程的状态：
    1. 程序计数器
    2. 寄存器
    3. 栈空间

How many threads in a program?
  As many as "useful" for the application
  Go encourages you to create many threads
    Typically more threads than cores
    The Go runtime schedules them on the available cores
  Go threads are not free, but you should think as if they are
    Creating a thread is more expensive than a method call
### 一个程序多少个线程合适？
1. 一个程序中有用线程越多越好
2. Go鼓励大家创建更多的线程
    1. 通常线程数量远远大于核数
    2. Go运行时环境调度这些线程到可用的核心上执行
3. Go创建一个线程也不完全是无开销的，因此需要考虑创建线程比一次方法调用的开销更大
    
Threading challenges:
  sharing data 
     one thread reads data that another thread is changing?
     can lead to race conditions
     -> don't share or coordinate sharing (e.g., w. mutexes)
  coordination between threads
    e.g. wait for all Map threads to finish
    can lead to deadlocks (typically easier to notice than races)
    -> use Go channels or WaitGroup
  granularity of concurrencyq
     coarse-grained -> simple, but little concurrency/parallelism
     fine-grained -> more concurrency, more races and deadlocks
### 线程挑战
1. 共享数据
    1. 一个线程读取数据，另一个线程一定会感知该数据的变化？
    2. 可能会导致竞态条件
    3. 不共享或者协调共享（eg：互斥信号量）
2. 线程之间的协调
    1. 等待所有map线程完成
    2. 导致死锁（通常更简单使用通知相对于竞态）

Crawler exercise: two challenges
  Arrange for I/O concurrency
    While fetching an URL, work on another URL
  Fetch each URL *once*
    want to avoid wasting network bandwidth
    want to be nice to remote servers
    => Need to some way of keeping track which URLs visited 
  [Work on different URLs on different cores]
    less agree that is important
### 爬虫联系:两个挑战
1. I/O的并发
    1. 当获取一个url，而该url又work在另一个线程
2. 获取每一个URL一次
    1. 阻止网络带宽的浪费
    2. 对远程服务端更友好
    3. 需要一些方式记录URL的访问轨迹
3. [不同的url在不同的核上]


Solutions [see handout: crawler.go]
  Eliminate depth --- use fetched instead
  Sequential one: pass fetched map to recursive calls
    Doesn't overlap I/O when fetcher takes a long time
    Doesn't exploit multiple cores
  Solution with a go routines and shared fetched map
    Create a thread for each url
      What happens we don't pass u?  (race)
    Why locks?  (remove them and every appears to works!)
      What can go wrong without locks?
        Checking and marking of url isn't atomic
	So it can happen that we fetch the same url twice.
	  T1 checks fetched[url], T2 checks fetched[url]
	  Both see that url hasn't been fetched
	  Both return false and both fetch, which is wrong
      This called a *race condition*
        The bug shows up only with some thread interleavings
	Hard to find, hard to reason about
      Go can detect races for you (go run -race crawler.go)
      Note that checking and marking of visiting must be atomic
    How can we can decide we are done with processing a page?
      waitGroup
  Solution with channels
    Channels: general-purse mechanism to coordinate threads
      bounded-buffer of messages
      several threads can send and receive on channel
        (Go runtime uses locks internally)
    Sending or receiving can block
      when channel is full
      when channel is empty
    Dispatch every url fetch through a master thread
      No races for fetched map, because it isn't shared!
### 解决办法[处理方式：crawler.go]
1. 消除深度--使用fetched代替
2. 序列化：传递fetched map递推调用
    1. 当获取数据花费比较长时间时，不同时处理
    2. 不支持多核
3. 通过goroutine及共享fetched map解决
    1. 为每一个url创建一个线程
    2. 为什么使用lock？（移除出现在work中的它们）
        1. 没有锁将会出现什么问题？
        2. 不会自动检查和标记url
4. 可能对于同一个url获取两次
    1. T1及T2都检查fetched[url]。
    2. 同时观察到url没有被获取
    3. 同时返回false，同时fetch，当出现错误
    4. 这个被称之为竞态条件
        1. 这个bug出现在一些线程之间交错
5. 难以发现，难以推导
    1. Go能够发现竞态通过例如：go run -race crawler.go
    2. url访问的检查和标记必须是自动的
6. 我们如何决定处理页面工作已经完成？
    1. waitGroup
7. 使用channels解决
    1. 通道：广泛采用的机制处理线程协作
        1. 有界的缓存区通道
        2. 多个线程能够对通道发送及接受
        3. （Go runtime内部使用锁）
    2. 发送和接受会阻塞
        1. channel满了
        2. channel空了
    3. 分发每一个url的获取通过主线程
        1. 没有竞态从fetched map获取，因此它们并没有共享


What is the best solution?
  All concurrent ones are substantial more difficult than serial one
  Some Go designers argue avoid shared-memory
    I.e. Use channels only
  My lab solutions use many concurrency features
    locks when sharing is natural
      e.g., several server threads sharing a map
    channels for coordination between threads
      e.g., producer/consumer-style concurrency
  Use Go's race detector:
    https://golang.org/doc/articles/race_detector.html
    go test -race mypkg
### 什么是最好的解决办法
  1. 并发比序列化更复杂
  2. 一些Go的设计者们极力反对共享内存
    1. 例如：只建议仅仅使用channels
  3. 我们的实验将使用多种并发特性
    1. 当共享是自然的则使用locks
        1. 例如： 多个服务器server共享map
    2. 多个线程协作使用channels
        1. 例如：生产者消费者模式并发
  4. 使用Go的竞态检测：
     [Golang 数据竞态检测(官方文档中英文翻译)更多详情，点击前往 ](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/go_race_detector.md?dir=0&filepath=LEC2_+RPCAndThreads%2Fgo_race_detector.md&oid=650fa0ca8013b5fd9a420889a598da567798128b&sha=e1d3e2bbbc17a3e94b5d1af9ac20f1b841148d2e)

     [英文原始链接，点击前往 ](https://golang.org/doc/articles/race_detector.html)


Remote Procedure Call (RPC)
  a key piece of distributed system machinery; all the labs use RPC
  goal: easy-to-program network communication
    hides most details of client/server communication
    client call is much like ordinary procedure call
    server handlers are much like ordinary procedures
  RPC is widely used!
### 远程过程调用
> 1.分布式系统的关键;所有试验使用RPC
> 2.目标：易于编程的网络通信
>   1. 隐藏client/server之间的大部份细节
>   2. client的调用类似普通本地调用
>   3. server处理非常类似普通过程调用
> 3. RPC被广泛的使用

```
RPC ideally makes net communication look just like fn call:
  Client:
    z = fn(x, y)
  Server:
    fn(x, y) {
      compute
      return z
    }
  RPC aims for this level of transparency
```

> RPC旨在使得网络通信像普通本地功能调用一样：
```
  Client:
    z = fn(x, y)
  Server:
    fn(x, y) {
      compute
      return z
    }
```
> RPC的目标是这个水平的透明度

Go example: kv.go  [see handout: kv.go]
  Client "dials" server and invokes Call()
    Call similar to a regular function call
  Server handles each request in a separate thread
    Concurrency!  Thus, locks for keyvalue.
> Go示例：kv.go  [see handout: kv.go]
> 1. Client连通server及调用Call()
>   1. Call类似于一个常规的函数调用
> 2. Server处理每一个request为每一个分布的线程
>   1. 并发！因此对keyvalue加锁lock

RPC message diagram:
  Client             Server
    request--->
       <---response

Software structure
  client app         handlers
    stubs           dispatcher
   RPC lib           RPC lib
     net  ------------ net
 
A few details:
  Which server function (handler) to call?
    In Go specified in Call()
  Marshalling: format data into packets
    Tricky for arrays, pointers, objects, &c
    Go's RPC library is pretty powerful!
    some things you cannot pass: e.g., channels, functions
  Binding: how does client know who to talk to?
    Maybe client supplies server host name
    Maybe a name service maps service names to best server host
### 更多详细：
> 1. 哪一个server的函数将被调用？
>   1. Go中指定为Call()函数
> 2. 序列化：格式数据转换成包
>   1. 棘手的数组、指针、对象、&c
>   2. Go的RPC库设计非常优秀
>   3. 一些东西你无法传递：例如：通道、函数等
> 3. 绑定：client如何知道与谁通信？
>   1. client可以通过server的hostname
>   2. 可能是service名称映射所有服务中的最好的host

RPC problem: what to do about failures?
  e.g. lost packet, broken network, slow server, crashed server
###　RPC问题：如何处理失败？
> 例如：丢包，网络故障、服务端处理速度慢，服务端crash奔溃

What does a failure look like to the client RPC library?
  Client never sees a response from the server
  Client does *not* know if the server saw the request!
    Maybe server/net failed just before sending reply
  [diagram of lost reply]
### 失败对于client意味着什么？
> 1. client没收接受到来自server的回应
> 2. client是否server已经收到了请求request
>   1. 有可能是server/net在发送reply回应之前发送失败

Simplest scheme: "at least once" behavior
  RPC library waits for response for a while
  If none arrives, re-send the request
  Do this a few times
  Still no response -- return an error to the application
### 最简单的模式：“至少一次”行为
> 1. RPC库等待response回应一段时间
> 2. 如果没有收到返回，则重新发送请求
> 3. 重复上述多次
> 4. 如果仍然没有回应--则直接返回error给应用程序处理

Q: is "at least once" easy for applications to cope with?
#### Q：“至少一次”对于应用程序是否容易处理？

Simple problem w/ at least once:
  client sends "deduct $10 from bank account"
#### 一个简单的“至少一次”问题：client发送“从银行账户减少10美元”

Q: what can go wrong with this client program?
  Put("k", 10) -- an RPC to set key's value in a DB server
  Put("k", 20) -- client then does a 2nd Put to same key
  [diagram, timeout, re-send, original arrives very late]
#### Q: 这个客户端程序会出现什么错误？
> 1. Put("k", 10) --一个RPC调用在数据库服务器中设置键值对。
> 2. Put("k",20) -- 客户端对同一个键设置其他值。

Q: is at-least-once ever OK?
  yes: if it's OK to repeat operations, e.g. read-only op
  yes: if application has its own plan for coping w/ duplicates
    which you will need for Lab 1
#### Q：“至少一次”一直能很好的工作？
> 1. 是的：它比较适用重复操作，例如只读操作
> 2. 是的：如果应用程序自己计划复制多个写副本（例如实验lab1）

Better RPC behavior: "at most once"
  idea: server RPC code detects duplicate requests
    returns previous reply instead of re-running handler
  Q: how to detect a duplicate request?
  client includes unique ID (XID) with each request
    uses same XID for re-send
  server:
    if seen[xid]:
      r = old[xid]
    else
      r = handler()
      old[xid] = r
      seen[xid] = true
### 更好的RPC行为：最多一次”
> 1. 想法：server RPC代码处理重复的请求，返回之前的回应而不是重新处理
> 2. Q：如何处理重复的请求？
> 3. client拥有一个唯一的id对应每一个request请求，通过唯一的id来处理重新发送
> 4. server: 如果已经计算过直接从old表中查找之前记录返回，否则直接计算，并保存此次记录。

some at-most-once complexities
  this will come up in labs 3 and on
  how to ensure XID is unique?
    big random number?
    combine unique client ID (ip address?) with sequence #?
  server must eventually discard info about old RPCs
    when is discard safe?
    idea:
      unique client IDs
      per-client RPC sequence numbers
      client includes "seen all replies <= X" with every RPC
      much like TCP sequence #s and acks
    or only allow client one outstanding RPC at a time
      arrival of seq+1 allows server to discard all <= seq
    or client agrees to keep retrying for < 5 minutes
      server discards after 5+ minutes
  how to handle dup req while original is still executing?
    server doesn't know reply yet; don't want to run twice
    idea: "pending" flag per executing RPC; wait or ignore
### 关于"最多一次"的复杂些
> 1. 这些将出现在lab3，以及如何确保XID的唯一性？
>   1. 大的随机数
>   1. 结合唯一的client id（例如ip地址）？
> 2. server必须能够有能力评估及丢弃老的RPC请求信息
>   1. 如何保证丢弃是安全的？
>   1. 一些想法：
>       1. 唯一的client id
>       1. 每个client RPC序列号
>       1. client的每一个RPC请求包含"seen all replies <=X"
>       1. 类似tcp中的seq和ack
>   1. 或者每次只允许一个RPC调用，到达的是seq+1，那么忽略其他小于seq
>   1. 客户端最多可以尝试5分钟，服务器会忽略大于5分钟之后的请求。
> 2. 当原来的请求还在执行，怎么样处理相同seq的请求？
>   1. 服务器不想运行两次，也不想回复。
>   1. 想法：给每个执行的RPC，pending标识；等待或者忽略。


What if an at-most-once server crashes and re-starts?
  if at-most-once duplicate info in memory, server will forget
    and accept duplicate requests after re-start
  maybe it should write the duplicate info to disk?
  maybe replica server should also replicate duplicate info?
### 如果“最多一次”的server奔溃或者重启？
> 1. 如果服务器将副本信息保存在内存中，服务器会忘记请求，同时在重启之后接受相同的请求。
> 2. 也许，你应该将副本信息保存到磁盘？
> 3. 也许，副本服务器应该继续保存副本信息？

What about "exactly once"?
  at-most-once plus unbounded retries plus fault-tolerant service
  Lab 3
### 关于“恰好一次”？
> 最多一次+无限重试+容错服务(Lab 3)

Go RPC is "at-most-once"
  open TCP connection
  write request to TCP connection
  TCP may retransmit, but server's TCP will filter out duplicates
  no retry in Go code (i.e. will NOT create 2nd TCP connection)
  Go RPC code returns an error if it doesn't get a reply
    perhaps after a timeout (from TCP)
    perhaps server didn't see request
    perhaps server processed request but server/net failed before reply came back
### Go RPC 采用“最多一次”模型
> 1. 打开TCP链接
> 1. 写入请求到TCP链接通道
> 1. TCP可能会重传，但是服务端的TCP能够过滤重复请求
> 1. Go 代码没有重试机制（例如：不会创建第二个TCP连接）
> 1. Go RPC返回一个错误error，如果没有获取到回复
>   1. 也许是超时之后
>   1. 也许是服务器没有接收到请求
>   1. 也许服务器处理完了请求，但是server/net在返回回复的时候发生故障


Go RPC's at-most-once isn't enough for Lab 1
  it only applies to a single RPC call
  if worker doesn't respond, the master re-send to it to another worker
    but original worker may have not failed, and is working on it too
  Go RPC can't detect this kind of duplicate
    No problem in lab 1, which handles at application level
    Lab 3 will explicitly detect duplicates
### Go RPC的“最多一次”不适应于Lab1
> 1. 它只能够适用于单个RPC调用
> 2. 如果worker没有回应，那么master将会重新发送job到其他worker,但是原来的worker并没有失败，且仍然在执行同一个job
> 3. Go RPC 不能处理如下重复
>   1. lab1中没有问题，lab1处理这些在应用程序级别
>   2. lab3将需要明确的处理重复