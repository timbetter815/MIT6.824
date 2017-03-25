GFS FAQ
# GFS问题
Q: Why is atomic record append at-least-once, rather than exactly
once?
## Q:为什么原子追加使用“至少一次”，而不是“恰好一次”？

It is difficult to make the append exactly once, because a primary
would then need to keep state to perform duplicate detection. That
state must be replicated across servers so that if the primary fails,
this information isn't lost. You will implement exactly once in lab
3, but with more complicated protocols that GFS uses.
> 1. 比较难以实现“恰好一次”追加，因为主副本希望保持状态进行重复检测。
> 那些状态必须跨机器复制，因此如果主副本失败发生故障，这些信息也不会丢失。
> 2. 你将会使用“恰好一次”实现lab3，但是需要使用比GFS更复杂的协议。

Q: How does an application know what sections of a chunk consist of
padding and duplicate records?
## Q: 应用程序怎么知道chunk块包括什么及重复的记录？

A: To detect padding, applications can put a predictable magic number
at the start of a valid record, or include a checksum that will likely
only be valid if the record is valid. The application can detect
duplicates by including unique IDs in records. Then, if it reads a
record that has the same ID as an earlier record, it knows that they
are duplicates of each other. GFS provides a library for applications
that handles these cases.
> 1. A:为了检测填充,应用程序可以把一个可预测的魔法数字设置在有效的记录的开始出,
> 或者包括一个校验和（只有记录是有效的才有效）。
> 2. 应用程序能够检测重复通过包括一个唯一的ID在记录中。因此如果读到一个相同的id，
> 那么则知道是重复数据。
> 3. GFS为应用程序提供基础库来解决这些问题。

Q: How can clients find their data given that atomic record append
writes it at an unpredictable offset in the file?
## Q:客户端如何找到它们那些原子记录追加写在文件中不可预知的offset的数据？

A: Append (and GFS in general) is mostly intended for applications
that read entire files. Such applications will look for every record
(see the previous question), so they don't need to know the record
locations in advance. For example, the file might contain the set of
link URLs encountered by a set of concurrent web crawlers. The
file offset of any given URL doesn't matter much; readers just want to
be able to read the entire set of URLs.
> 1. A：追加(GFS)主要是用于应用程序读取整个文件。因此应用程序将查找每一个记录（查看之前的问题），
> 这样它们不需要预先知道记录位置。例如，文件可能包含一系列并行爬虫web产生的url链接；
> 对于给定的任何URL，文件offset关系不大；读取程序只需要能够读取整个urls即可。

Q: The paper mentions reference counts -- what are they?
## Q：论文中提及的引用计数--它们到底是什么？

A: They are part of the implementation of copy-on-write for snapshots.
When GFS creates a snapshot, it doesn't copy the chunks, but instead
increases the reference counter of each chunk. This makes creating a
snapshot inexpensive. If a client writes a chunk and the master
notices the reference count is greater than one, the master first
makes a copy so that the client can update the copy (instead of the
chunk that is part of the snapshot). You can view this as delaying the
copy until it is absolutely necessary. The hope is that not all chunks
will be modified and one can avoid making some copies.
> 1. A:它们部分实现了写时拷贝快照。当GFS创建一个快照，它们不会复制chunks，
> 但是会增加每一个chunk的引用计数。这样使得创建一个快照成本更低。
> 2. 如果客户端对chunk进行写入操作且master通知引用计数超过一，
> 那么master首先创建一个复制副本使得客户端能够更新副本（而不是chunk的快照的一部分）
> 你能够观察这些将会延时复制直到真的必要。
> 3. 意愿：并不是所有的块将被修改,可以避免副本复制。

Q: If an application uses the standard POSIX file APIs, would it need
to be modified in order to use GFS?
## Q: 如果一个应用程序使用标准的POSIX文件API，那么它需要为了使用GFS而修改吗？

A: Yes, but GFS isn't intended for existing applications. It is
designed for newly-written applications, such as MapReduce programs.
> 1. 是的，但是GFS并不适应与已经存在的应用。它是为新编写的应用程序设计的，例如MapReduce程序。

Q: How does GFS determine the location of the nearest replica?
## Q: GFS如何确定本地最近的副本？

A: The paper hints that GFS does this based on the IP addresses of the
servers storing the available replicas. In 2003, Google must have
assigned IP addresses in such a way that if two IP addresses are close
to each other in IP address space, then they are also close together
in the machine room.
> 1. 论文暗示GFS根据服务器IP地址存储可用的副本。在2003年，google已经采用这种方式分配ip地址
> 即：如果两个ip地址空间比较相近，那么它们也仍然在机器空间比较相近。

Q: Does Google still use GFS?
## Google仍然在使用GFS吗？

A: GFS is still in use by Google and is the backend of other storage
systems such as BigTable. GFS's design has doubtless been adjusted
over the years since workloads have become larger and technology has
changed, but I don't know the details. HDFS is a public-domain clone
of GFS's design, which is used by many companies.
> A：GFS仍然在google被使用且是其他系统后端，例如bigtable。
> GFS的设计毫无疑问在这些年进行过调整，随着工作负载变得越来越大以及技术的改变，
> 但是我也不清楚具体的细节。HDFS是一个开源的模仿GFS的设计，目前被用于许多公司。

Q: Won't the master be a performance bottleneck?
## Q：master将成为性能瓶颈？

A: It certainly has that potential, and the GFS designers took trouble
to avoid this problem. For example, the master keeps its state in
memory so that it can respond quickly. The evaluation indicates that
for large file/reads (the workload GFS is targeting), the master is
not a bottleneck. For small file operations or directory operations,
the master can keep up (see 6.2.4).
> 1. 当然有这种可能性，GFS设计者们带来麻烦为了阻止这个问题。
> 例如：master保持状态在内存中以至于它们能够快速响应。评估表明在大文件读（就算GFS的负载接近目标值）master不是瓶颈。
> 2. 对于小文件操作或者目录操作，master也能够跟上（详见6.2.4）

Q: How acceptable is it that GFS trades correctness for performance
and simplicity?
## Q: 如何接受GFS交易正确性，关于性能及简单性？

A: This a recurring theme in distributed systems. Strong consistency
usually requires protocols that are complex and require chit-chat
between machines (as we will see in the next few lectures). By
exploiting ways that specific application classes can tolerate relaxed
consistency, one can design systems that have good performance and
sufficient consistency. For example, GFS optimizes for MapReduce
applications, which need high read performance for large files and are
OK with having holes in files, records showing up several times, and
inconsistent reads. On the other hand, GFS would not be good for
storing account balances at a bank.
> 1. 这在分布式系统一个反复出现的主题.强一致性通常需要复杂的协议以及需要保存对话在两台机器之间（我们将会在接下来几课中见到）。
> 2. 某些特定场景下的程序能够容忍一定程度的不一致性，那么则可以设计出高性能及足够一致性的系统。
> 3. 例如：GFS针对mapreduce应用进行优化，这些应用需要针对大文件高的读性能，且运行文件中有holes。
> 另一方面，GFS不适用于存储银行账户余额。

Q: What if the master fails?
## 如果master发生故障？

A: There are replica masters with a full copy of the master state; an
unspecified mechanism switches to one of the replicas if the current
master fails (Section 5.1.3). It's possible a human has to intervene
to designate the new master. At any rate, there is almost certainly a
single point of failure lurking here that would in theory prevent
automatic recovery from master failure. We will see in later lectures
how you could make a fault-tolerant master using Raft.
> 1. A:存在副本master保存了master全部状态；一个未详细说明的机制能够在master故障时候将副本master切换为master（章节5.1.3）。
> 2. 也有可能人工进行干预，指定新的master。
> 3. 无论如何,几乎肯定会有一个单点故障理论在这里预防自动恢复，当master故障时。
> 4. 我们将在后续课程中使用Raft来设计一个容错master。
