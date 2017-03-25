6.824 2017 Lecture 3: GFS Case Study
# 6.824 2017 第三课：GFS案例研究

The Google File System
Sanjay Ghemawat, Howard Gobioff, and Shun-Tak Leung
SOSP 2003
## Google文件系统

Why are we reading this paper?
  the file system for map/reduce
  main themes of 6.824 show up in this paper
    trading consistency for simplicity and performance
    motivation for subsequent designs
  good systems paper -- details from apps all the way to network
    performance, fault-tolerance, consistency
  influential
    many other systems use GFS (e.g., bigtable)
    HDFS (Hadoop's Distributed File Systems) based on GFS
## 为什么阅读该论文？
> 1. 文件系统为map/reduce做支撑
> 2. 6.824主题出现在此篇论文中
>   1. 简单高效的交易的一致性
>   2. 为了后续设计
> 3. 优秀的系统论文--app与网络连接所有方式的细节
>   1. 性能，容错，一致性
> 4. 有影响力的
>   1. 许多其他系统使用GFS（例如：bigtable）
>   2. HDFS（Hadoop分布式系统）基于GFS

What is consistency?
  A correctness condition
  Important when data is replicated and concurrently accessed by applications 
    if an application performs a write, what will a later read observe?
      what if the read is from a different application?
  Weak consistency
    read() may return stale data  --- not the result of the most recent write
  Strong consistency
    read() always returns the data from the most recent write()
  General trade-off:
    strong consistency is nice for application writers
    strong consistency is bad for performance
  Many correctness conditions (often called consistency models)
    Today first peak; will show up in almost every paper we read this term
## 什么是一致性？
> 1. 正确性条件
> 2. 最重要的：当数据多个副本被应用程序并发访问
>   1. 如果应用程序完成写操作，那么之后的读将会观察到什么值？如果读是从其他的应用程序发起的呢？
> 3. 弱一致性
>   1. read()可能返回过期的数据--并不是最近写的数据
> 4. 强一致性
>   1. read()一直能够返回最新写的数据
> 5. 常用的权衡：
>   1. 强一致性对应用程序的写更加友好
>   1. 强一致性对性能有影响
> 5. 许多正确性条件（通常称为一致性模型）
>   1. 频率出现很高；几乎将出现在这学期我们读的每一篇论文


"Ideal" consistency model
  A replicated files behaves like as a non-replicated file system
    picture: many clients on the same machine accessing files on a single disk
  If one application writes, later reads will observe that write
  What if two application concurrently write to the same file
    In file systems often undefined  --- file may have some mixed content
  What if two application concurrently write to the same directory
    One goes first, the other goes second
## 理想的一致性模型
> 1. 副本文件的行为能够像没有副本的文件系统一样
>   1. 想象：许多客户端在同一个机器访问单个磁盘的文件
> 2. 如果应用程序写完成，那么之后的读将能够看到这个写的变化
> 3. 如果两个应用程序并发的写相同的文件
>   1. 文件系统中通常没有定义--文件可能出现混合的内容
> 4. 如果两个应用程序并发的写相同的目录
>   1. 一个先,其他的第二

Challenges to achieving ideal consistency
  Concurrency
  Machine failures
  Network partitions
## 达到理想一致性的挑战
> 1. 并发
> 2. 机器故障
> 3. 网络分区

Why are these challenges difficult to overcome:
  Requires communication between clients and servers
    May cost performance
  Protocols can become complex --- see next week
    Difficult to implement system correctly
  Many systems in 6.824 don't provide ideal
    GFS is one example
## 为什么这些挑战难以克服：
> 1. 客户端和服务器端需要通信
>   1. 可能会损失性能
> 1. 协议会变得更加复杂--详见下周课程
>   1. 系统的容错性难以实现
> 1. 6.824的许多系统不提供理想一致性
>   1. GFS就是不提供理想一致性的例子


Central challenge in GFS:
  With so many machines failures are common
    assume a machine fails once per year
    w/ 1000 machines, ~3 will fail per day.
  High-performance: many concurrent readers and writers
    Map/Reduce jobs read and store final result in GFS
    Note: *not* the temporary, intermediate files
  Use network efficiently
  These challenges difficult combine with "ideal" consistency
## GFS中面临的核心挑战
> 1. 如此多机器，因此故障的出现成为常态
>   1. 假设每个机器一年出现一个故障
>   1. 如果有1000台机器，那么一天则会有3台出现故障
> 1. 高性能：高并发的读和写
>   1. Map/Reduce任务读及存储最终的结果到GFS中
>   1. 注意：不是临时的的，中间文件
> 1. 使用高效网络
> 1. 这些挑战困难结合“理想”的一致性

High-level design
  Directories, files, names, open/read/write
    But not POSIX
  100s of Linux chunk servers with disks
    store 64MB chunks (an ordinary Linux file for each chunk)
    each chunk replicated on three servers
    Q: why 3x replication?
    Q: Besides availability of data, what does 3x replication give us?
       load balancing for reads to hot files
       affinity
    Q: why not just store one copy of each file on a RAID'd disk?
       RAID isn't commodity
       Want fault-tolerance for whole machine; not just storage device
    Q: why are the chunks so big?
  GFS master server knows directory hierarchy
    for dir, what files are in it
    for file, knows chunk servers for each 64 MB
    master keeps state in memory
      64 bytes of metadata per each chunk
    master has private recoverable database for metadata
      master can recovery quickly from power failure
    shadow masters that lag a little behind master
      can be promoted to master
## 高层设计
> 1. 目录，文件，名字，打开/读/写
>   1. 但是不符合POSIX
> 1. 上百的Linux块服务器与磁盘块
>   1. 存储64MB块（一个普通的Linux文件对应每个块）
>   1. 每个块副本复制在三个服务器上
>   1. Q: 为什么三个副本？
>   1. Q: 除了数据的可用性，三份副本还有什么作用？
>       1. 为热点文件的读负载均衡
>       1. 类同
>   1. 为什么不仅仅为每一个文件存储一个拷贝在RAID磁盘上？
>       1. RIAD不是日常商品
>       1. 想要对整个机器进行容错，不仅仅针对存储
>   1. 为什么块文件如何大？
> 1. GFS主服务器知道目录层次结构
>   1. 对于目录，知道里面存了哪些文件
>   1. 对于文件，知道块服务对于每个块的64MB文件
>   1. master保存状态在内存中（每一个块保存64字节元数据）
>   1. master可以自主恢复数据库元数据（master能够快速从掉电故障中恢复）
>   1. shadow masters 一点点落后于master（能够被提升为master）


Basic file operations:
  Client read:
    send file name and offset to master
    master replies with set of servers that have that chunk
      response includes version # of chunk
      clients cache that information
    ask nearest chunk server
      checks version #
      if version # is wrong, re-contact master
  Client append
    ask master where to store
      maybe master chooses a new set of chunk servers if crossing 64 MB
      master responds with chunk servers and version #
        one chunk server is primary
    Clients pushes data to replicas
      Replicas form a chain
      Chain respects network topology
      Allows fast replication
    Client contacts primary when data is on all chunk servers
      primary assigns sequence number
      primary applies change locally
      primary forwards request to replicas
      primary responds to client after receiving acks from all replicas
    If one replica doesn't respond, client retries
      After contacting master
  Master can appoint new master if master doesn't refresh lease
  Master replicates chunks if number replicas drop below some number
  Master rebalances replicas
## 基本文件操作
> 1. 客户端读：
>   1. 发送文件的名字及偏移量到master
>   1. master回复拥有chunk的服务器组
>       1. 回复包括chunk的版本号
>       1. 客户端缓存信息
>   1. 询问最近的chunk服务器
>       1. 检查版本
>       1. 如果版本错误，重连master
> 1. 客户端追加
>   1. 询问master存储到哪儿
>      1. master可能选择一系列新的chunk服务器如果当前chunk服务器超过64MB
>      1. master回应对应的chunks服务器及版本(其中包含一个主chunk server)
>   1. 客户端推送数据到副本
>      1. 副本构成一个链
>      1. 链涉及网络拓扑
>      1. 允许快速复制
>   1. 客户端联系主chunk server，当数据是所有块服务器上时
>      1. 主chunk server分配序列号
>      1. 主chunk server应用于本地变化
>      1. 主chunk server将请求转发到副本
>      1. 主chunk server响应客户端，当收到从所有副本的ack时
>   1. 如果一个副本没有回应，那么客户端将重试
>      1. 然后联系master
> 1. Master能选举一个新的master，当在规定租赁时期没有更新
> 1. Master复制chunk块，如果副本数量低于一些数字
> 1. Master重新均衡副本


Does GFS achieve "ideal" consistency?
  Two cases: directories and files
  Directories: yes, but...
    Yes: strong consistency (only one copy)
    But:
      master may go down and GFS is unavailable
        shadow master can serve read-only operations, which may return stale data
        Q: Why not write operations?
	  split-brain syndrome (see next lecture)
  Files: not always
    Mutations with atomic appends
      A file can have duplicate entries and holes
        if primary fails to contact a replica, the primary reports an error to client
        client retries and primary picks a new offset
	record can be duplicated at two offsets
	while other replicas may have a hole at one offset
    An "unlucky" client can read stale data for short period of time
      A failed mutation leaves chunks inconsistent
        The primary chunk server updated chunk
        But then failed and the replicas are out of date
      A client may read an not-up-to-date chunk
      When client refreshes lease it will learn about new version #
    Mutations without atomic append
      data of several clients maybe intermingled
      concurrent writes on non-replicated Unix can also result in a strange outcome
      if you are, use atomic append or a temporary file and atomically rename
## GFS是否达到了理想的一致性？
> 1. 两个关注点：目录和文件
> 1. 目录：是的，但是...
>   1. 是的：强一致性（仅仅一份copy）
>   1. 但是：master可能会宕掉（导致整个GFS不可用）
>       1. shadow master能够提供只读操作服务，可以返回过期的数据
>       1. Q: 为什么不能提供写操作？
>   1. 出现脑裂（详见下一课）
> 1. 文件：不总是
>   1. 原子追加变化
>       1. 一个文件可以有重复的条目和漏洞
>           1. 如果主副本连接其他副本失败，主副本报告error给client
>           1. 客户端重试及主副本选择一个新offset
>   1. 记录能够被副本在两个offset上
>   1. 而其他副本可能有一个hole在一个偏移量
>   1. 一个不幸的客户端在一个短暂时期可能会读到过期的数据
>       1. 主chunk server更新chunk
>       1. 但是这时出现失败且副本过期
>   1. 客户端可能读到一个没有按时更新的chunk
>   1. 当客户端重新刷新租约将获取一个新的版本
> 1. 非原子追加变化
>   1. 多个客户端的数据可能会混合
>   1. 并发写在没有副本的Unix系统也会出现奇怪的结果
>   1. 如果是你，请使用原子追加或者中间文件且原子重命名



Authors claims weak consistency is not a big problems for apps    
  Most file updates are append-only updates
    Application can use UID in append records to detect duplicates
    Application may just read less data (but not stale data)
  Application can use temporary files and atomic rename
## 作者声称弱一致性对应用程序不是一个大问题
> 1. 大部分文件更新仅仅是追加更新
>   1. 应用程序能够在追加记录后使用UID处理多份记录
>   1. 应用程序可以读取更少的数据（但是不是过期的）
> 1. 应用程序能够使用中间文件及原子重命名操作


Performance (Figure 3)
  huge aggregate throughput for read (3 copies, striping)
    125 MB/sec in aggregate
    Close to saturating network
  writes to different files lower than possible maximum
    authors blame their network stack
    it causes delays in propagating chunks from one replica to next
  concurrent appends to single file
    limited by the server that stores last chunk
## 高性能（图表3）
> 1. 大集群吞吐量读（3份，条带化）
>   1. 集群支持125MB/s
>   1. 接近饱和的网络
> 1. 写入不同的文件低于可能的最大值
>   1. 作者责怪他们网络堆栈
>   1. 它导致延迟传播从一个副本块到下一个
> 1. 并发追加到一个文件
>   1. 受限于服务器的最后一个chunk


Summary
  case study of performance, fault-tolerance, consistency
    specialized for MapReduce applications
  what works well in GFS?
    huge sequential reads and writes
    appends
    huge throughput (3 copies, striping)
    fault tolerance of data (3 copies)
  what less well in GFS?
    fault-tolerance of master
    small files (master a bottleneck)
    clients may see stale data
    appends maybe duplicated
 ## 总结
 > 1. 高性能，容错，一致性的案例研究
 >   1. 特别是MapReduce应用
 > 1. GFS哪些表现优秀？
 >   1. 大量序列化读写
 >   1. 追加
 >   1. 大吐吞量（3份，条带化）
 >   1. 数据容错（）3份
 > 1. GFS哪些表现不好？
 >   1. master的容错
 >   1. 小文件（master是一个瓶颈）
 >   1. 客户端可能看到过期数据
 >   1. 追加可能重复

References
## 引用
  http://queue.acm.org/detail.cfm?id=1594206  (discussion of gfs evolution)
  http://highscalability.com/blog/2010/9/11/googles-colossus-makes-search-real-time-by-dumping-mapreduce.html
