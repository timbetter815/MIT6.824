6.824 2017 Lecture 1: Introduction
>## 6.824 2017 第一课: 介绍


6.824: Distributed Systems Engineering
>## 6.824: 分布式系统工程


What is a distributed system?
  multiple cooperating computers
  big databases, P2P file sharing, MapReduce, DNS, &c
  lots of critical infrastructure is distributed!
>### 什么是分布式系统？
>* 多台机器协作
>* 大数据、P2P文件共享、MapReduce、DNS等
>* 很多关键基础设施是分布式的。


Why distributed?
  to connect physically separate entities
  to achieve security via isolation
  to tolerate faults via replication
  to scale up throughput via parallel CPUs/mem/disk/net

But:
  complex: many concurrent parts
  must cope with partial failure
  tricky to realize performance potential


>### 为什么需要分布式系统？
>* 为了连接分离的物理实体
>* 为了通过隔离实现安全性
>* 为了通过复制实现容错
>* 为了通过并行CPUs/mem/disk/net实现横向扩容
>### 然后：
>* 并行性：多个并发的部分
>* 必须处理部分失败
>* 难道达到性能最大值


Why take this course?
  interesting -- hard problems, powerful solutions
  used by real systems -- driven by the rise of big Web sites
  active research area -- lots of progress + big unsolved problems
  hands-on -- you'll build serious systems in the labs

>### 为什么选择这门课程？
>* 兴趣 -- 难题，比较好的解决方案
>* 被实际系统使用 -- 被高并发高性能网站驱动
>* 活跃的研究领域 -- 快速进步及大量未解决的问题
>* 动手 -- 你将通过实验建立一个重要的系统


COURSE STRUCTURE
>### 课程结构
http://pdos.csail.mit.edu/6.824


Course staff:
>### 课程工作人员
  Robert Morris, lecturer
  Frans Kaashoek, lecturer
  Lara Araujo, TA
  Anish Athalye, TA
  Srivatsa Bhat, TA
  Daniel Ziegler, TA

Course components:
  lectures
  readings
  two exams
  labs
  final project (optional)
>### 课程组成
>* 课程
>* 阅读
>* 两次考试测验
>* 实验
>* 最终项目（可选）

Lectures:
  big ideas, paper discussion, and labs
>### 课程
>* 好的想法，论文讨论及实验

Readings:
  research papers, some classic, some new
  the papers illustrate key ideas and important details
  many lectures focus on the papers
  please read papers before class!
  each paper has a short question for you to answer
  and you must send us a question you have about the paper
  submit question&answer by 10pm the night before
>### 阅读
>* 请课前阅读研究论文，及经典论文的关键设计意图。
>* 每篇论文都会又一个较短的问题需要回答，你必须在晚上10点钱提前该问题的答案

Mid-term exam in class; final exam during finals week
>### 中期考试和期末考试（在最后一周）


Lab goals:
  deeper understanding of some important techniques
  experience with distributed programming
  first lab is due a week from Friday
  one per week after that for a while

>### 实验目标
>* 深入理解关键技术
>* 掌握分布式编程经验
>* 第一个实验从第一周周五开始的一周时间

Lab 1: MapReduce
Lab 2: replication for fault-tolerance using Raft
Lab 3: fault-tolerant key/value store
Lab 4: sharded key/value store
>* 实验1：MapReduce
>* 实验2：使用Raft技术支撑复制来实现容错
>* 实验3：key-val容错系统
>* 实验4：共享key-val存储系统


Optional final project at the end, in groups of 2 or 3.
  The final project substitutes for Lab 4.
  You think of a project and clear it with us.
>* 可选择的最终项目，可以分成2组或者3组来完成
>* 最终的项目替代实验4
>* 你可以自己设想一个项目和我们一起来弄明白它


Lab grades depend on how many test cases you pass
  we give you the tests, so you know whether you'll do well
  careful: if it often passes, but sometimes fails,
    chances are it will fail when we run it
>* 实验学分取决于你有多少测试用例通过
>* 我们会给到测试用例来验证是否你完成的很好
>* 注意：如果你的测试用例经常是通过的，只是偶尔失败，在我们运行时则很可能失败

Debugging the labs can be time-consuming
  start early
  come to TA office hours
  ask questions on Piazza
>* 调试实验可能会比较耗时

MAIN TOPICS

This is a course about infrastructure, to be used by applications.
  About abstractions that hide distribution from applications.
  Three big kinds of abstraction:
    Storage.
    Communication.
    Computation.
  [diagram: users, application servers, storage servers]

>### 课程主题：这是一个基础课程，适用于应用程序。隐藏系统分布式复杂性进行抽象。抽象为如下三个方面：
>* 存储
>* 通信
>* 计算
>* [diagram: 用户，应用服务器，存储服务器]


A couple of topics come up repeatedly.
>### 系统中下述的主题有可能会同时出现两三个。

Topic: implementation
  RPC, threads, concurrency control.

>### 主题：实现
>* RPC, threads, 并发控制.


Topic: performance
  The dream: scalable throughput.
    Nx servers -> Nx total throughput via parallel CPU, disk, net.
    So handling more load only requires buying more computers.
  Scaling gets harder as N grows:
    Load im-balance, stragglers.
    Non-parallelizable code: initialization, interaction.
    Bottlenecks from shared resources, e.g. network.

>### 主题：性能
>* 理想：可伸缩的吞吐性能，通过购买更多的机器，利用多核CPU、磁盘、网络等来实现
>* 扩展性变得困难随着N倍的增长：负载均衡、慢的掉队节点；无法并行化的代码；资源共享瓶颈、网络等

Topic: fault tolerance
  1000s of servers, complex net -> always something broken
  We'd like to hide these failures from the application.
  We often want:
    Availability -- app can keep using its data despite failures
    Durability -- app's data will come back to life when failures are repaired
  Big idea: replicated servers.
    If one server crashes, client can proceed using the other(s).

>### 主题：容错
>* 上千的服务器，复杂的网络->这样总会出现一些问题；我们希望隐藏这些出现错误的问题从应用上。
>#### 我们通常希望：
>* 高可用--尽管数据出现失败app仍然可用
>* 持续可用--当失败被修复app的数据能够恢复到正常
>#### 更好的主意：备份服务器
>* 如果一个服务器出现问题，客户端能够继续使用其他正常备份服务器


Topic: consistency
  General-purpose infrastructure needs well-defined behavior.
    E.g. "Get(k) yields the value from the most recent Put(k,v)."
  Achieving good behavior is hard!
    "Replica" servers are hard to keep identical.
    Clients may crash midway through multi-step update.
    Servers crash at awkward moments, e.g. after executing but before replying.
    Network may make live servers look dead; risk of "split brain".
  Consistency and performance are enemies.
    Consistency requires communication, e.g. to get latest Put().
    "Strong consistency" often leads to slow systems.
    High performance often imposes "weak consistency" on applications.
  People have pursued many design points in this spectrum.

> ### 主题：一致性
* 通用的基础设施需求定义良好的行为。 例如： Get(k) 获取到的值应该是最近的 Put(k,v)设置的。
* 实现良好的行为是很困难的!
>    1. 各个备份服务器之间很难保持完全相同
>    2. 客户端可能在中途奔溃
>    3. 网络故障可能使得正常服务器无法连接；出现“脑裂”风险
* 一致性和高性能是相对的
>    1. 致性需要通信，例如获取最后的Put()操作内容
>    2. 强一致性通常会导致系统变慢
>    3. 高性能一般会导致系统的“弱一致性”
* 人民已经研究了很多设计方法处理这些问题


CASE STUDY: MapReduce
> ## 用例学习：MapReduce

Let's talk about MapReduce (MR) as a case study
  MR is a good illustration of 6.824's main topics
  and is the focus of Lab 1
> 让我们将MR作为一个案例进行讨论。 MR是课程6.284主题的一个很好的例子，也是实验1的主要关注点。

MapReduce overview
  context: multi-hour computations on multi-terabyte data-sets
    e.g. analysis of graph structure of crawled web pages
    only practical with 1000s of computers
    often not developed by distributed systems experts
    distribution can be very painful, e.g. coping with failure
  overall goal: non-specialist programmers can easily split
    data processing over many servers with reasonable efficiency.
  programmer defines Map and Reduce functions
    sequential code; often fairly simple
  MR runs the functions on 1000s of machines with huge inputs
    and hides details of distribution

> ## MapReduce概要
* 概述: 几个小时处理完TB的数据集 例如：仅仅使用1000台左右机器处理分析爬行网页的结构，假若这些系统不是由分布式系统开发专家开发的这就会非常痛苦，例如复制错误。
* 总体目标: 非专业程序员可以轻松的在合理的效率下解决的巨大的数据处理问题。程序员定义Map函数和Reduce函数、顺序代码一般都比较简单。 MR在成千的机器上面运行处理大量的数据输入，隐藏全部分布式的细节。


Abstract view of MapReduce
  input is divided into M files
  Input1 -> Map -> a,1 b,1 c,1
  Input2 -> Map ->     b,1
  Input3 -> Map -> a,1     c,1
                    |   |   |
                        |   -> Reduce -> c,2
                        -----> Reduce -> b,2
  MR calls Map() for each input file, produces set of k2,v2
    "intermediate" data
    each Map() call is a "task"
  MR gathers all intermediate v2's for a given k2,
    and passes them to a Reduce call
  final output is set of <k2,v3> pairs from Reduce()
    stored in R output files

> ## MapReduce的抽象试图
> ### 输入会被分割到M文件中，
  Input1 -> Map -> a,1 b,1 c,1
  Input2 -> Map ->     b,1
  Input3 -> Map -> a,1     c,1
                    |   |   |
                        |   -> Reduce -> c,2
                        -----> Reduce -> b,2

> ### MR调用Map()函数处理每一个输入文件，产生k2/v2中间数据，每一次Map()调用即一个任务
> ### MR收集所有的中间数v2为给定的k2,并将这些值传递给Reduce()调用，存储在R输出文件


Example: word count
  input is thousands of text files
  Map(k, v)
    split v into words
    for each word w
      emit(w, "1")
  Reduce(k, v)
    emit(len(v))

> ## 举例：单词统计
## 输入为成千上万的文本文件

MapReduce hides many painful details:
  starting s/w on servers
  tracking which tasks are done
  data movement
  recovering from failures

> ## MapReduce隐藏了许多令人痛苦的细节
* starting s/w on servers
* 追踪任务完成
* 数据移动
* 从失败中恢复

MapReduce scales well:
  N computers gets you Nx throughput.
    Assuming M and R are >= N (i.e. lots of input files and output keys).
    Maps()s can run in parallel, since they don't interact.
    Same for Reduce()s.
  So you can get more throughput by buying more computers.
    Rather than special-purpose efficient parallelizations of each application.
    Computers are cheaper than programmers!

> ## MapReduce具有很好的扩展性：
* N台机器能够获取N倍的吞吐性能。假若M、R >= N,且有众多的输入文件和输出keys。那么Maps()s和Reduce()s能够并行执行，因为它们之间相互没有影响。
* 你能够通过购买更多的机器来获取更大的吞吐量，而不需要针对每一个应用程序使用特别的并行编程技术。购买电脑比并行编程成本更低。

What will likely limit the performance?
  We care since that's the thing to optimize.
  CPU? memory? disk? network?
  In 2004 authors were limited by "network cross-section bandwidth".
    [diagram: servers, tree of network switches]
    Note all data goes over network, during Map->Reduce shuffle.
    Paper's root switch: 100 to 200 gigabits/second
    1800 machines, so 55 megabits/second/machine.
    Small, e.g. much less than disk or RAM speed.
  So they cared about minimizing movement of data over the network.
    (Datacenter networks are much faster today.)

> ## 哪些因素可能限制性能的提升？
* 我们关心的就是我们需要优化的。
* CPU? memory? disk? network?
* 2004年被网络带宽限制
>    1. [diagram: 服务器, 树形网络交换机]
>    2. 在Map->Reduce模型出现之前，所有的数据都需要通过网络
>    3. Paper's root switch: 100-200 GB/s
>    4. 1800台机器，因此将有55M/s数据流过机器
>    5. 相对于磁盘和RAM的速度，网络的速度太小了
* 因此大家更加关心如何才能使得数据在网络上尽可能小的移动（当然数据中心网络速度在今天也比之更快了）

More details (paper's Figure 1):
  master: gives tasks to workers; remembers where intermediate output is
  M Map tasks, R Reduce tasks
  input stored in GFS, 3 copies of each Map input file
  all computers run both GFS and MR workers
  many more input tasks than workers
  master gives a Map task to each worker
    hands out new tasks as old ones finish
  Map worker hashes intermediate keys into R partitions, on local disk
  no Reduce calls until all Maps are finished
  master tells Reducers to fetch intermediate data partitions from Map workers
  Reduce workers write final output to GFS (one file per Reduce task)

> ## 更多详情（论文图标1）
* 主节点：分发任务给工作节点；记住中间值输出位置
* M个Map任务，R个Reduce任务
* 输入存储在GFS、每个输入文件保存三个备份
* 所有的机器运行着GFS和MR工作线程
* 更多数量的输入任务相比较于工作线程
* 当之前老的任务处理完成，master节点分配心得Map任务给对于的worker
* Map worker将输出中间key散列映射到对应的本地磁盘的R个分区上
* 当maps任务完成，Reduce任务开始执行
* master节点告知Reducers如何从数据分区位置获取中间值
* Reduce线程最终的输出到GFS

How does detailed design reduce effect of slow network?
  Map input is read from GFS replica on local disk, not over network.
  Intermediate data goes over network just once.
    Map worker writes to local disk, not GFS.
  Intermediate data partitioned into files holding many keys.
    Big network transfers are more efficient.

> ## 有哪些reduce详细设计用来应对网络慢的问题？
* Map的输入从GFS本地磁盘备份读取，而不用通过网络
* 中间值数据只会经过网络传输一次
* Map线程结果写入到本地磁盘，而不是GFS
* 中间数据通过key被划分到多个文件
* 大的网络传输将更有效


How do they get good load balance?
  Critical to scaling -- bad for N-1 servers to wait for 1 to finish.
  But some tasks likely take longer than others.
  Solution: many more tasks than workers.
    Master hands out new tasks to workers who finish previous tasks.
    So no task is so big it dominates completion time (hopefully).
    So faster servers do more work than slower ones, finish abt the same time.

> ## 如何获取比较好的负载均衡？
* 扩容的关键--比较坏的情况是N-1个服务器等待其他服务器完成
* 有些任务处理时间长于其他任务
* 解决办法：
>    1. 任务比处理线程多
>    2. 主节点将新的任务分发给其他处理完前序任务的线程
>    3. 任务都应该分解到能够在可预期的时间内处理完成的大小
>    4. 处理越快的节点完成更多的任务

What about fault tolerance?
  I.e. what if a server crashes during a MR job?
  Hiding failures is a huge part of ease of programming!
  Why not re-start the whole job from the beginning?
  MR re-runs just the failed Map()s and Reduce()s.
    MR requires them to be pure functions:
      they don't keep state across calls,
      they don't read or write files other than expected MR inputs/outputs,
      there's no hidden communication among tasks.
    So re-execution yields the same output.
  The requirement for pure functions is a major limitation of
    MR compared to other parallel programming schemes.
    But it's critical to MR's simplicity.

> ## 关于容错？
* 例如：在处理MR任务时，服务器宕机？
* 隐藏错误是轻松编程的一个巨大部分
* 为什么不从一开始重启整个job？
* MR只会重新运行失败的Map()s或者Reduce()s
* MR需要具备如下理论功能：
>     1. 不要操持状态跨越调用时
>     2. 不能够读写非预期的MR输入输出
>     3. 任务之间没有隐藏的通信
* 因此重新执行会产生同样的输出
* MR与其他并行模型的比较是这些理论和功能需求的主要因素
* 但是这些对于MR的基础是非常关键的


Details of worker crash recovery:
  * Map worker crashes:
    master sees worker no longer responds to pings
    crashed worker's intermediate Map output is lost
      but is likely needed by every Reduce task!
    master re-runs, spreads tasks over other GFS replicas of input.
    some Reduce workers may already have read failed worker's intermediate data.
      here we depend on functional and deterministic Map()!
    master need not re-run Map if Reduces have fetched all intermediate data
      though then a Reduce crash would then force re-execution of failed Map
  * Reduce worker crashes.
    finshed tasks are OK -- stored in GFS, with replicas.
    master re-starts worker's unfinished tasks on other workers.
  * Reduce worker crashes in the middle of writing its output.
    GFS has atomic rename that prevents output from being visible until complete.
    so it's safe for the master to re-run the Reduce tasks somewhere else.

> ## worker的故障恢复详情：
* Map worker崩溃：
>     1. 主节点发现worker不再回应pings,奔溃的worker行程的中间Map输出将被丢失，但是这些确实Reduce任务所必须的
>     2. 主节点重新运行，从其他GFS备份数据中分配任务
>     3. 一些Reduce worker还是有可能读取失败的worker产生的中间数据，我们需要依靠确定的，功能性的Map（）。
>     4. 如果Reduces已经获取所有的中间数据，那么则不需要重新运行Map，除非此时Reduce奔溃，那么将强制重新执行失败的Mqp
* Reduce worker崩溃：
>     1. 任务完成ok--存储到GFS，包括备份副本
>     2. 主节点在其他节点重启没有完成的任务
* Reduce worker在写输出过程中崩溃：
>     1. GFS会自动重命名输出，保持其不可见状态，直到全部Reduce完成
>     2. 因此master在其他地方再次运行Reduce worker将会是安全的。


Other failures/problems:
  * What if the master gives two workers the same Map() task?
    perhaps the master incorrectly thinks one worker died.
    it will tell Reduce workers about only one of them.
  * What if the master gives two workers the same Reduce() task?
    they will both try to write the same output file on GFS!
    atomic GFS rename prevents mixing; one complete file will be visible.
  * What if a single worker is very slow -- a "straggler"?
    perhaps due to flakey hardware.
    master starts a second copy of last few tasks.
  * What if a worker computes incorrect output, due to broken h/w or s/w?
    too bad! MR assumes "fail-stop" CPUs and software.
  * What if the master crashes?

> ## 其他失败/问题
* 如果主节点分发相同Map任务给两个工作worker？
> 1. 有可能master认为其中一个worker已经挂掉
> 2. 它将告诉Reduce workers其中的一个
* 如果主节点分发相同Reduce任务给两个工作worker？
> 1. 它们将同时尝试输出相同的输出文件到GFS
> 2. GFS原子重命名将会阻止这种混淆; 确保只有一个完成文件可见
* 如果一个worker处理速度非常慢--一种“掉队现象”
> 1. 也许是由于硬件出现的莫名其妙原因
> 2. 主节点为心得任务开启第二个数据复制
* 如果一个worker计算一个错误的输出，在硬件或者软件故障期间
> 1. 太糟糕了！MR假设是建立在"fail-stop"的cpu和软件基础之上。
* 如果master奔溃？


For what applications *doesn't* MapReduce work well?
  Not everything fits the map/shuffle/reduce pattern.
  Small data, since overheads are high. E.g. not web site back-end.
  Small updates to big data, e.g. add a few documents to a big index
  Unpredictable reads (neither Map nor Reduce can choose input)
  Multiple shuffles, e.g. page-rank (can use multiple MR but not very efficient)
  More flexible systems allow these, but more complex model.

> ## 关于那些不能很好执行MapReduce的应用？
* 不是所有的场景适应于map/shuffle/reduce模型。
* 小的数据，因为间接成本太高, 比如非网站后端
* 大数据中的小更新，比如添加很少的文件到大的索引
* 不可预知的读(Map 和 Reduce都不能选择输入)
* Multiple shuffles, 比如网页排名 (能够使用多个MR，但是不是特别有效)
* 多数灵活的系统允许MR，但是使用更加复杂的模型


Conclusion
  MapReduce single-handedly made big cluster computation popular.
  - Not the most efficient or flexible.
  + Scales well.
  + Easy to program -- failures and data movement are hidden.
  These were good trade-offs in practice.
  We'll see some more advanced successors later in the course.
  Have fun with the lab!

> ## 总结
>* MapReduce是的大数据集群计算更加流行
> 1. \- 而不是更加有效或者复杂
> 2. \+ 扩展性更好
> 3. \+ 更加易于编程--系统失败及数据移动被隐藏
>* 这个一个非常好的权衡在实际使用中
>* 我们将看到更多的高级成功者在以后的过程中
>* 快乐的完成实验！