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

# 接下来开始lab1：
### 首先下载源代码[git仓库地址](git://g.csail.mit.edu/6.824-golabs-2017 6.824)，且配置好该源代码src路径到GOPATH中。

---
#### 设置linux环境下gopath（windows类似）：
vim ~.bash_profile
>      export GOROOT=/home/go
>      export PATH=$GOROOT/bin:$PATH
>      export GOPATH=/home/gopath:/home/gopath/src/git.oschina.net/tantexian/MIT6.824/6.824-golabs-2017
>      alias cdgo='cd /home/gopath/src'
>      alias cdmit='cd /home/gopath/src/git.oschina.net/tantexian/MIT6.824/6.824-golabs-2017/src'
---

## Part 1：Map/Reduce输入及输出
1. 从测试代码入手：~/6.824-golabs-2017/src/mapreduce/test_test.go
2. test_test.go为测试用例代码，其中覆盖了很多场景用例，TestSequentialSingle()测试方法入手，分析跟踪代码。
3. 更多具体代码逻辑解析，直接在源代码中进行注释，请查看对应源代码注释。

### 本次需要完成代码：
1. 将输出分割成map任务的函数和为reduce任务收集全部输入的函数。
2. 实现common_map.go里面的doMap函数和common_reduce.go里面的doReduce函数。

### 测试：
>      cd "$GOPATH/src/mapreduce"
>      go test -run Sequential

### 分析（更多详细请参考对应[git地址:](https://git.oschina.net/tantexian/MIT6.824/)中的代码）：
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


## Part 2：Single-worker单词统计
1. 从代码main/wc.go入手
2. 在linux上执行go run wc.go master sequential pg-*.txt
3. 完成代码之后，将得到统计完成所有word出现次数的输出文件


### 本次需要完成代码：
1. main/wc.go#mapF()（实现可以参考part1中的mapFunc）
2. main/wc.go#reduceF()（实现可以参考part1中的reduceFunc）

### 测试：
>      cd "$GOPATH/src/main"
>      go run wc.go master sequential pg-*.txt


### 分析（更多详细请参考对应[git地址:](https://git.oschina.net/tantexian/MIT6.824/)中的代码）：
1. main/wc.go#mapF()（实现可以参考part1中的mapFunc）：根据传入的文件内容，将之分割为一个个单词，然后范围key为单词，value为单词出现次数的keyval数组
2. main/wc.go#reduceF()：根据传入进来的keyval数组，返回当前word出现的次数
3. 其中输入文件数量为map执行任务数量，其中nreduce为执行reduce任务数量也即输出中间文件个数


## Part 3：分布式mapreduce任务
1. 在这部分实验中，你将会完成另外一个版本的MapReduce,将工作分散到一系列的工作线程，为了充分利用多核的作用。
虽然工作不是分布在真正的多机上面，不过你的实现可以使用RPC和Channel来模拟一个真正的分布式实现。
2. mapreduce/master.go，主要管理mapreduce任务；mapreduce/common_rpc.go处理rpc
3. 本次实验即实现mapreduce/schedule.go#schedule()函数。master将会调用两次schedule()（Map时期和reduce时期）
4. schedule()的主要目的是：分发任务到可用的works。通常情况下，任务数远远多于worker线程，因此master将会每次分配一个任务给这些works。
schedule()必须等到所有的任务完成，才能够返回。
6. schedule()通过读取registerChan参数来感知workers，channel为每个workers产生一个字符串（包含channel的address）。
一些workers在schedule()被调用前可能已经存在，有些则是在schedule()运行过程中启动。所有workers都会出现在registerChan，
schedule()将会调度使用所有的workers，包括那些在它启动之后加入的。
7. schedule()通过发送Worker.DoTask RPC命令给worker让它执行任务。mapreduce/common_rpc.go#DoTaskArgs定义了RPC参数。
RPC参数file只用于map任务，代表map中哪一个被读文件名。schedule()在mapFiles中能够找到这些文件名字。
8. 通过mapreduce/common_rpc.go#call()方法发送RPC命令给worker。第一个参数为worker的address（从registerChan中获得）
第二个参数为："Worker.DoTask"，第三个参数为DoTaskArgs结构体，最后一个参数为nil。
9. part3只需要修改schedule.go文件中代码，其他代码不允许修改提交。


### 本次需要完成代码：
修改完成schedule.go中的代码


### 测试：
>      cd "$GOPATH/src/mapreduce"
>      go test -run TestBasic


### 分析（更多详细请参考对应[git地址:](https://git.oschina.net/tantexian/MIT6.824/)中的代码）：
1. mapreduce#schedule()函数将会被调用两次，第一次为map阶段，第二次为reduce阶段
2. 如果为map阶段，则执行map任务数量为map输入文件个数，如果为reduce阶段，则reduce任务数为该nreduce值
（其中输入文件数量为map执行任务数量，其中nreduce为执行reduce任务数量也即输出中间文件个数）
3. 由于map或者reduce阶段，都需要等所有任务完成才能返回，因此可以使用sync.WaitGroup等待所有任务完成再返回
4. 使用for循环，执行ntask次任务，每次任务使用RPC调用，call()远程调用worker.go#DoTask方法
5. 如果远程调用DoTask方法返回成功，则将当期worker放置到空闲worker通道中，供下一次复用
如果返回失败，则continue，重新获取一个新的空闲worker来执行当前任务



## Part 4：处理worker执行任务失败情况
1. 本部分将需要处理worker失败的情况，由于worker为无状态的，因此mapreduce处理这种情况相对比较容易
2. 如果master通过RPC分配任务给worker失败（或者超时），那么master需要将该任务重新分配给一个其他的worker来处理
3. RPC失败，并不总是意味着该worker不能执行该任务，有可能是RPC的reply丢失、或者master的RPC超时（但是worker仍然在执行任务）
因此有可能两个worker执行着相同的任务，产生相同的输出。
PS：本次不需要处理master失败的情况，由于master是有状态的因此FT容错处理会相对困难，在后续课程实验中将会来解决该问题。


### 本次需要完成代码：
修改完成schedule.go中的代码

### 测试：
>      cd "$GOPATH/src/mapreduce"
>      go test -run Failure


### 分析（更多详细请参考对应[git地址:](https://git.oschina.net/tantexian/MIT6.824/)中的代码）：
1. mapreduce#schedule()函数中，获取对于的远程RPC调用DoTask函数的返回值。
2. 如果DoTask返回值为执行成功，则说明该次任务被当前worker处理成功执行完成，那么则将该worker回收到当前空闲workers中，供下一次任务调度执行使用
3. 如果DoTask返回值为执行失败，则当前该次任务需要继续被分配重新分给一个新的worker去执行
PS：此部分其他在Part 3中已经完成。



## Part 5：倒排索引（可选的额外加分）
1. 本次实验将创建map及reduce来实现倒排序功能。倒排索引广泛运用在计算科学领域，特别是资料搜索。
广义地说，一个反向索引就是map,保存感兴趣的基础数据索引，来指向数据的原始位置。
例如：在上下文搜索中，map即保存关键字指向所在文档的索引集合。
2. 本次实验需要修改main/ii.go中的mapF及reduceF，让之能够产生一个倒排索引。
3. 运行ii.go将输出元组列表，像如下格式：
```
$ go run ii.go master sequential pg-*.txt
$ head -n5 mrtmp.iiseq
A: 16 pg-being_ernest.txt,pg-dorian_gray.txt,pg-dracula.txt,pg-emma.txt,pg-frankenstein.txt,pg-great_expectations.txt,pg-grimm.txt,pg-huckleberry_finn.txt,pg-les_miserables.txt,pg-metamorphosis.txt,pg-moby_dick.txt,pg-sherlock_holmes.txt,pg-tale_of_two_cities.txt,pg-tom_sawyer.txt,pg-ulysses.txt,pg-war_and_peace.txt
ABC: 2 pg-les_miserables.txt,pg-war_and_peace.txt
ABOUT: 2 pg-moby_dick.txt,pg-tom_sawyer.txt
ABRAHAM: 1 pg-dracula.txt
ABSOLUTE: 1 pg-les_miserables.txt
```


### 本次需要完成代码：
main#ii.go#mapF()/reduceF()

### 测试：
>      cd "$GOPATH/src/main"
>      go run ii.go master sequential pg-*.txt


### 分析（更多详细请参考对应[git地址:](https://git.oschina.net/tantexian/MIT6.824/)中的代码）：
1. mapF()根据输入文件名filename及内容，将内容分割为一个个单词，然后返回当前key=word，val=filename的值
2. reduceF(),将同一个key对应的多个filename数组进行计算得出出现次数，
其次将所有出现该单词的所有文件的文件名查找出来（相同文件名需要去重）
PS：去重使用map来实现，由于底层使用hash算法复杂度为O(1),数组具有n个元素，因此整体复杂度为O(n)


## 执行所有测试：
>      cd "$GOPATH/src/main"
>      bash ./test-mr.sh












