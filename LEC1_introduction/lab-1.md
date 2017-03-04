### 6.824 Lab 1: MapReduce
(PS：后续只会对文档重要部分进行翻译)

原文地址：[链接](http://nil.csail.mit.edu/6.824/2017/labs/lab-1.html)

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

You should not need to modify any other files,
but reading them might be useful in order to understand
how the other methods fit into the overall architecture of the system.

----------

> #### 6. master发送关闭RPC给对应的workers，然后关闭自己的RPC服务。
### 在本次课程接下来的练习中：
>      1. 你将不得不需要编写或者修改doMap, doReduce，schedule方法
>      2. 上述分别在common_map.go, common_reduce.go, and schedule.go文件中。
>      3. 你需要编写map 和 reduce方法（在../main/wc.go文件中）


>#### 你不需要更改其他文件，但是阅读它们可与帮你理解其他的方法适应系统的总体架构。

----------

## 接下来开始lab1：
### 首先下载源代码[git仓库地址](git://g.csail.mit.edu/6.824-golabs-2017 6.824)，且配置好该源代码src路径到GOPATH中。

---
##### 设置linux环境下gopath（windows类似）：
vim ~.bash_profile
>      export GOROOT=/home/go
>      export PATH=$GOROOT/bin:$PATH
>      export GOPATH=/home/gopath:/home/gopath/src/git.oschina.net/tantexian/MIT6.824/6.824-golabs-2017
>      alias cdgo='cd /home/gopath/src'
>      alias cdmit='cd /home/gopath/src/git.oschina.net/tantexian/MIT6.824/6.824-golabs-2017/src'
---

### Part 1：Map/Reduce输入及输出
1. 从测试代码入手：~/6.824-golabs-2017/src/mapreduce/test_test.go
2. test_test.go为测试用例代码，其中覆盖了很多场景用例，TestSequentialSingle()测试方法入手，分析跟踪代码。
3. 更多具体代码逻辑解析，直接在源代码中进行注释，请查看对应源代码注释。

本次需要完成代码：
1. 将输出分割成map任务的函数和为reduce任务收集全部输入的函数。
2. 实现common_map.go里面的doMap函数和common_reduce.go里面的doReduce函数。

测试：
>      cd "$GOPATH/src/mapreduce"
>      go test -run Sequential

分析（更多详细请参考对应[git地址:](https://git.oschina.net/tantexian/MIT6.824/)中的代码）：
>      common_map.go#domap():（本次test中mapF（）进行了单词分割统计）
1. 根据input file使用utilio将该文件的内容全部读取出来保存到contents
2. 调用mapF将文件内容contents分割为一个个单词，然后返回一个keyval数组（其中key为该单词，value为单词出现的次数（默认都有出现1次））
3. 根据nreduce的参数值，创建对应的文件数量的中间输出文件
4. 为了保证相同的key对应的keyval值保存到同一个文件中（方便后续doReduce任务统计相同key出现次数）
5. 步骤4中使用对相同的key值使用hash取模%nreduce来达到相同key保存到同一个文件中

>      common_reduce.go#doReduce():（本次test中reduceF（）直接返回空）
1. 将nreduce个中间文件（domap步骤生成），内容读取出来
2. 将1中读取内容统计到一个keyval的map中，其中map的key为单词word，value为对应出现的次数数组
3. 将keyval的map中的每个项（一次处理一个key），调用reduceF（）函数处理
4. 其中reduceF（）函数将具有相同key的数组进行计算，计算出当前key的单词出现过多少次
5. 将4中对每一个单词出现多少次的计算结果保存到一个输出文件中


### Part 2：Single-worker单词统计
1. 从代码main/wc.go入手
2. 在linux上执行go run wc.go master sequential pg-*.txt
3. 完成代码之后，将得到统计完成所有word出现次数的输出文件


本次需要完成代码：
1. main/wc.go#mapF()（实现可以参考part1中的mapFunc）
2. main/wc.go#reduceF()（实现可以参考part1中的reduceFunc）

测试：
>      cd "$GOPATH/src/main"
>      go run wc.go master sequential pg-*.txt


分析（更多详细请参考对应[git地址:](https://git.oschina.net/tantexian/MIT6.824/)中的代码）：
1. main/wc.go#mapF()（实现可以参考part1中的mapFunc）：根据传入的文件内容，将之分割为一个个单词，然后范围key为单词，value为单词出现次数的keyval数组
2. main/wc.go#reduceF()：根据传入进来的keyval数组，返回当前word出现的次数