#### 本项目主要针对自己对MIT6.824课程学习和翻译的记录，也方便后续其他学习者作为参考。
有些内容也会记录在博客网址：[https://my.oschina.net/tantexian/blog ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://my.oschina.net/tantexian/blog)

### [MIT6.824网址 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](http://nil.csail.mit.edu/6.824/2017/)

6.824: Distributed Systems

Spring 2017
TR1-2:30, room 54-100

![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/start-index1.png?dir=0&filepath=resources%2Fstatic%2Fimg%2Fstart-index1.png&oid=d46b06da44b6e9b7a4842f25b40d9af2f4eb5b46&sha=b8814ea7f159eb661b26c1bb6f8449c9dbc2874a)


6.824是麻省理工学院的、关于分布式系统的编程实验课程。主要包括容错、复制、一致性等内容，总共包括23次课，4次试验及2次考试。

左上角链接补充解释：
* Information：课程的介绍。
* Schedule：课程表安排，安排了第1课到23课的排课日期及每次课程教授的内容及课程内容介绍和考试（包括提前需要准备的基础知识及布置的作业）。
* Submissions：课程答题入口，注册自己的账号。
* Labs: 1 2 3 4：表示本次课程需要完成的四个实验。
* Questions：问答。
* Past Exams：往年考试试卷及答案。
* 2000 Web Site - 2016 Web Site：2000年-2016年以往6.824课程网址。

### 学习入门:
该课程的学习入口：[Schedule入口 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](http://nil.csail.mit.edu/6.824/2017/schedule.html)开始学习LEC1-LEC23以及课程学习资料。
![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/scheduler-2017.png?dir=0&filepath=resources%2Fstatic%2Fimg%2Fscheduler-2017.png&oid=9c38ebc90be76b60603943bec0b9769289e6da86&sha=dec89cbe68c8c8eb7bd1a5020da1ec3bfe3a534d)




### 入门开始（必读）：

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


---
## 第一次课程学习顺序及相应课程资源:
> 1. [LEC 1: Introduction ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC1_introduction/101.md?dir=0&filepath=LEC1_introduction%2F101.md&oid=b0d1831741f38a5f100a721b7e3cd9a69d709822&sha=b78bf6a66ea41b865876cf6fe0f065602c5e4eb7)
> 2. [Preparation: Read MapReduce (2004) ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](http://nil.csail.mit.edu/6.824/2017/papers/mapreduce.pdf)
　　　　[原始资源：MapReduce论文英文文档 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/raw/dev/LEC1_introduction/mapreduce.pdf)
> 3. [Assigned: Lab 1: MapReduce ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](http://nil.csail.mit.edu/6.824/2017/labs/lab-1.html)

本次课程学习了MapReduce的论文及简单版本mapreduce相关实现。



## 第2次课程学习顺序及相应课程资源:
> 1. [Lecture 2: Infrastructure: RPC and threads ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/l-rpc.md?dir=0&filepath=LEC2_+RPCAndThreads%2Fl-rpc.md&oid=d9882eaf0b502618b7791a710b17f64e39821aea&sha=3433911775cc3fa59246fa90478b30f322927c3f)
> 2. [Crawler ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/crawler.go?dir=0&filepath=LEC2_+RPCAndThreads%2Fcrawler.go&oid=ff0cd66a881cd48c13635a5be87fdc69d9936330&sha=3433911775cc3fa59246fa90478b30f322927c3f)
> 3. [K/V ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/kv.go?dir=0&filepath=LEC2_+RPCAndThreads%2Fkv.go&oid=3dfee994be6ddcd7c5df368cd9015cc76fba646a&sha=3433911775cc3fa59246fa90478b30f322927c3f)
> 4. [阅读A Tour of Go ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://tour.golang.org/welcome/1)
> 5. [go-faq.txt ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/go-faq.md?dir=0&filepath=LEC2_+RPCAndThreads%2Fgo-faq.md&oid=069edc9588eea4d31410d791ac68ab8d4cc97fa4&sha=6c3f736d5b76322c2492a4275f745137cc277b97)
> 6. [扩展资源：Effective Go中文版.pdf ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/raw/dev/docs/Effective%20Go%E4%B8%AD%E6%96%87%E7%89%88.pdf)
> 7. [扩展资源：Golang内存模型(官方文档中英文翻译）![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/go_mem_model.md?dir=0&filepath=LEC2_+RPCAndThreads%2Fgo_mem_model.md&oid=1d6554a60b8a59d0a6be126356e04dfc405bf6d0&sha=80f489a4f27c4b551794d7fba60a4f67d6568515)
> 7. [扩展资源：Golang数据竞态检测(官方文档中英文翻译) ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC2_%20RPCAndThreads/go_race_detector.md?dir=0&filepath=LEC2_+RPCAndThreads%2Fgo_race_detector.md&oid=04409064747458409c4606735ad42a2b465707f9&sha=80f489a4f27c4b551794d7fba60a4f67d6568515)

本次课程主要学习了Go语言的基础知识、并发编程、RPC示例。以及“at least one”，“at most one”分布式系统设计要素及重难点。



## 第3次课程学习顺序及相应课程资源:
> 1. [l-gfs-short-zh_cn.md ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/l-gfs-short-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2Fl-gfs-short-zh_cn.md&oid=c9b89fca5b1d85f0e7f48a3c50e6eff7ec0334f3&sha=ed481bdbb6a41443aa2ce04390c44b9653ef0a78)
> 2. [Lecture 3: GFS论文 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/gfs-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2Fgfs-zh_cn.md&oid=4e69900cf9e548c526bb4981a97ee55dd0dbd55f&sha=ed481bdbb6a41443aa2ce04390c44b9653ef0a78)
　　　　[原始资源：GFS论文英文文档 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/raw/dev/LEC3_GFSAndRaft/gfs.pdf)
> 3. [gfs-faq-zh_cn.md ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/gfs-faq-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2Fgfs-faq-zh_cn.md&oid=9c784174eae8abb6ec83a166a1794ec232a45717&sha=ed481bdbb6a41443aa2ce04390c44b9653ef0a78)
> 4. [raft-zh_cn.md ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/raft-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2Fraft-zh_cn.md&oid=665ca438613092ec8d9080c58184b360f840c5fd&sha=2ccc943ea154b0b5783772823671da262a7d48be)
　　　　[原始资源：raft论文英文文档 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/raw/dev/LEC3_GFSAndRaft/raft.pdf)
> 5. [6.824_Lab_2_Raft-zh_cn.md ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/6.824_Lab_2_Raft-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2F6.824_Lab_2_Raft-zh_cn.md&oid=381f13961fa57427ec0f746a20825cf1fa63da1e&sha=ed481bdbb6a41443aa2ce04390c44b9653ef0a78)

本次课程主要学习了GFS论文，及raft扩展论文及实现Raft实验。





---
## 如果有对应翻译文件，则链接到对应翻译地址，否则链接到mit网站地址，也可以直接在本项目对应课程目录下本地查看全部资源！
## 翻译文档全部以zh_cn.md结尾。

#### GOPATH环境配置：
为了能够让6.824代码中能够正确导入对应的源码路径，因此本示例中将源码路径：C:\Users\ASUS\Desktop\dev\GOPATH\src\git.oschina.net\tantexian\MIT6.824\6.824-golabs-2017添加到GOPATH中


<br></br>
----------
## 关于代码编写规范：

### 如何开始？
* [下载最新版代码](https://git.oschina.net/tantexian/MIT6.824)
* [`使用指南及相关文档`]：所有目录下对应README.md为入口文档。


----------


### 开发规范`必读`
* 源文件使用Unix换行、UTF-8文件编码
* 请在git clone命令之前执行`git config --global core.autocrlf false`，确保本地代码使用Unix换行格式
* 请在非主干分支上开发，禁止提交本地未测试运行通过代码到线上分支
* 每次提交及之前(正常来说需要先pull --rebase,解决冲突)，对代码进行修改必须有相对应的解释说明
* 正常组内开发人员提交代码，需要经过经过审核后方可提交（且需要有统一格式注释，参照注释类型3）
  


### 注释规范
* 对于注释，请遵照以下规范：
* 注释类型1（适用于结构体或者包名注释）、

```
// 方法对象名 xxx
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/20 or v1.0.0 
```

* 注释类型2（适用于功能确定的单行注释）、

```
// 由于是顺序消息，因此只能选择一个queue生产和消费消息
```

* 注释类型3（适用于修改它人代码注释）、

```
// xxx  Modify: tantexian, <my.oschina.net/tantexian> 2017/3/20 or v1.0.0 
// xxx  Add: tantexian, <my.oschina.net/tantexian> 2017/3/20 or v1.0.0 
```
  
* 关于TODO、FIXME、XXX注释规范（后续再加上）、

```
// TODO: + 说明：xxx Author: tantexian, <my.oschina.net/tantexian>  Since: 2017/3/20 or v1.0.0 
```
如果代码中有TODO该标识，说明在标识处有功能代码待编写，待实现的功能在说明中会简略说明。

```
// FIXME: + 说明：xxx Author: tantexian, <my.oschina.net/tantexian>  Since: 2017/3/20 or v1.0.0 
```
如果代码中有FIXME该标识，说明标识处代码需要修正，甚至代码是错误的，不能工作，需要修复，如何修正会在说明中简略说明。

```
// XXX: + 说明：xxx Author: tantexian, <my.oschina.net/tantexian>  Since: 2017/3/20 or v1.0.0 
```
如果代码中有XXX该标识，说明标识处代码虽然实现了功能，但是实现的方法有待商榷，希望将来能改进，要改进的地方会在说明中简略说明。



### 开发IDE
* 开发工具不做统一规定（Idea、Eclipse都可以），建议使用Idea
* 建议使用最新版格式Idea，附下载地址：http://pan.baidu.com/s/1slMkXY1
* 附Idea属性格式注释文件下载地址：http://pan.baidu.com/s/1hrU3IgW（其中java版本Idea使用zz或者zzz命令来生成注释，golang使用gg或者ggg）

----------


>*联系方式：*
>博客网址：[https://my.oschina.net/tantexian/blog ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://my.oschina.net/tantexian/blog)
>邮箱：tantexian@qq.com