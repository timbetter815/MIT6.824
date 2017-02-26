### 6.824 Lab 1: MapReduce
(PS：后续只会对文档重要部分进行翻译)

[原文地址](http://nil.csail.mit.edu/6.824/2017/labs/lab-1.html)

简介

在这实验中你们将会通过建立MapReduce的库方法去学习Golang编程语言，同时学习分布式系统怎么进行容错。
在第一部分你将会编写简单的MapReduce程序。在第二部分你将会编写master，master往workers分配任务，
处理workers的错误。库的接口和容错的方式跟之前MapReduce论文中描述的很类似。

软件

你们将会使用Golang语言实现这个以及其他的全部实验。Golang的官方网站包含了很多教程资料，也许你们会想去查看资料。

我们将会向你们提供MapReduce的部分实现，包括分布式和非分布式的操作。你将通过git(一个版本管理系统)获取到初始的实验代码。可与参考它的帮助手册学习如何使用git，或者如果你已经熟悉其他版本控制软件，你也许会发现this CS-oriented overview of git useful.

课程的[git仓库地址](git://g.csail.mit.edu/6.824-golabs-2017 6.824)

Preamble: Getting familiar with the source


The mapreduce package provides a simple Map/Reduce library.
Applications should normally call Distributed() [located in master.go] to start a job,
but may instead call Sequential() [also in master.go] to get a sequential execution for debugging.

> ### 序言：熟悉源代码
> #### mapreduce源码包提供了一个简单的Map/Reduce库。应用程序通常调用mapreduce包下面master.go文件的Distributed()启动一个job任务。
> #### 也可以调用Sequential()为调试启动一个顺序执行job。


The code executes a job as follows:
> #### 一个job任务代码的执行如下：

1. The application provides a number of input files, a map function, a reduce function,
and the number of reduce tasks (nReduce).
> 1. 应用程序提供一系列输入文件，一个map函数、一个reduce函数，以及一系列reduce任务。

2. A master is created with this knowledge.
It starts an RPC server (see master_rpc.go),
and waits for workers to register (using the RPC call Register() [defined in master.go]).
As tasks become available (in steps 4 and 5), schedule() [schedule.go] decides how to assign those tasks to workers,
and how to handle worker failures.
> 2. master基于这些知识而建立。它启动一个RPC服务（master_rpc.go中），等待workers（使用RPC 调用Register()[在master.go中]）来注册。

3. The master considers each input file to be one map task,
and calls doMap() [common_map.go] at least once for each map task.
It does so either directly (when using Sequential()) or by issuing the DoTask RPC to a worker [worker.go].
Each call to doMap() reads the appropriate file, calls the map function on that file's contents,
and writes the resulting key/value pairs to nReduce intermediate files.
doMap() hashes each key to pick the intermediate file and thus the reduce task that will process the key.
There will be nMap x nReduce files after all map tasks are done.
Each file name contains a prefix, the map task number, and the reduce task number.
If there are two map tasks and three reduce tasks, the map tasks will create these six intermediate files:
mrtmp.xxx-0-0
mrtmp.xxx-0-1
mrtmp.xxx-0-2
mrtmp.xxx-1-0
mrtmp.xxx-1-1
mrtmp.xxx-1-2
Each worker must be able to read files written by any other worker, as well as the input files. Real deployments use distributed storage systems such as GFS to allow this access even though workers run on different machines. In this lab you'll run all the workers on the same machine, and use the local file system.
> 3. master将每一个输入文件作为一个map任务，为每一个map任务调用doMap()[common_map.go]至少一次
任务直接运行(当调用Sequential方法)或者通过触发worker上的DoTask RPC调用(worker.go)。
每一个doMap会调用合适的文件，在文件内容上调用map函数，然后将对应的key-val对结果写入nReduce中间文件。
doMap()哈希每一个key到对应的中间文件，接着reduce将处理这些key。
在全部的map任务执行完成之后这里将会有nMap x nReduce文件。每一个文件包含map任务编号及reduce任务编号前缀。
假若有两个map任务和三个reduce任务，map任务将产生6个中间文件：
mrtmp.xxx-0-0
mrtmp.xxx-0-1
mrtmp.xxx-0-2
mrtmp.xxx-1-0
mrtmp.xxx-1-1
mrtmp.xxx-1-2


4. The master next calls doReduce() [common_reduce.go] at least once for each reduce task.
As with doMap(), it does so either directly or through a worker.
The doReduce() for reduce task r collects the r'th intermediate file from each map task,
and calls the reduce function for each key that appears in those files. The reduce tasks produce nReduce result files.
> 4. master的下一步工作就是为reduce任务至少调用一次doReduce函数(common_reduce.go)，
和doMap函数一样，它可以直接运行或者通过worker运行。
doReduce方法收集每一个map产生的reduce文件，然后在这些文件之上调用reduce函数。这些reduce任务产生nReduce个结果文件。

5. The master calls mr.merge() [master_splitmerge.go],
which merges all the nReduce files produced by the previous step into a single output.
> 5. master调用mr.merge()函数[master_splitmerge.go]，该函数将之前步骤产生的nReduce个文件合并成一个单独的文件输出。

6. The master sends a Shutdown RPC to each of its workers, and then shuts down its own RPC server.
Over the course of the following exercises,
you will have to write/modify doMap, doReduce, and schedule yourself.
These are located in common_map.go, common_reduce.go, and schedule.go respectively.
You will also have to write the map and reduce functions in ../main/wc.go.

----------
**********

> 6. master发送关闭RPC给对应的workers，然后关闭自己的RPC服务。
## 在本次课程接下来的练习中：
>> 1.你将不得不需要编写或者修改doMap, doReduce，schedule方法
>> 2. 上述分别在common_map.go, common_reduce.go, and schedule.go文件中。
>> 3. 你需要编写map 和 reduce方法（在../main/wc.go文件中）

You should not need to modify any other files,
but reading them might be useful in order to understand
how the other methods fit into the overall architecture of the system.
> 你不需要更改其他文件，但是阅读它们可与帮你理解其他的方法适应系统的总体架构。

**********
----------

