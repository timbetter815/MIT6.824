## 第3次课程学习顺序及相应课程资源:
> 1. [l-gfs-short-zh_cn.md ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/l-gfs-short-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2Fl-gfs-short-zh_cn.md&oid=c9b89fca5b1d85f0e7f48a3c50e6eff7ec0334f3&sha=ed481bdbb6a41443aa2ce04390c44b9653ef0a78)
> 2. [Lecture 3: GFS论文 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/gfs-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2Fgfs-zh_cn.md&oid=4e69900cf9e548c526bb4981a97ee55dd0dbd55f&sha=ed481bdbb6a41443aa2ce04390c44b9653ef0a78)
　　　　[原始资源：GFS论文英文文档 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/raw/dev/LEC3_GFSAndRaft/gfs.pdf)
> 3. [gfs-faq-zh_cn.md ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/gfs-faq-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2Fgfs-faq-zh_cn.md&oid=9c784174eae8abb6ec83a166a1794ec232a45717&sha=ed481bdbb6a41443aa2ce04390c44b9653ef0a78)
> 4. [raft-zh_cn.md ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/raft-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2Fraft-zh_cn.md&oid=665ca438613092ec8d9080c58184b360f840c5fd&sha=2ccc943ea154b0b5783772823671da262a7d48be)
　　　　[原始资源：raft论文英文文档 ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/raw/dev/LEC3_GFSAndRaft/raft.pdf)
> 5. [6.824_Lab_2_Raft-zh_cn.md ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC3_GFSAndRaft/6.824_Lab_2_Raft-zh_cn.md?dir=0&filepath=LEC3_GFSAndRaft%2F6.824_Lab_2_Raft-zh_cn.md&oid=381f13961fa57427ec0f746a20825cf1fa63da1e&sha=ed481bdbb6a41443aa2ce04390c44b9653ef0a78)

本次课程主要学习了GFS论文，及raft扩展论文及实现Raft实验。



## 补充：Raft论文中关于“拜占庭将军问题”：
在分布式理论中，经常看到"在非拜占庭错误情况下，算法是有效的..."或者说"...可以容忍非拜占庭失效"。
那么什么是拜占庭将军问题呢？

> 在分布式计算上，不同的计算机通过讯息交换，尝试达成共识；
> 但有时候，系统上的协调计算机（Coordinator / Commander）或成员计算机 （Member / Lieutanent）
> 可能因系统错误并交换错的讯息，导致影响最终的系统一致性。
> 即：在系统不同的机器之间会传递错误的消息，这种情况即为拜占庭问题。
> 这与“网络分割、机器崩溃...”是不同的。(例如：由于系统自身的内存或者cpu出现内部故障，或者黑客攻击，伪造消息情况等，即：拜占庭将军叛变。)
> 比如Raft协议不能容忍拜占庭问题，但是能够在非拜占庭错误情况下，
> 有网络延迟、分区、丢包、冗余和乱序等错误情况出现时，都可以保证其操作的正确性。


---------------------------------------------------------------------------------

### 如果有对应翻译文件，则链接到对应翻译地址，否则链接到mit网站地址，也可以直接在本项目对应课程本地查看全部原始资源！
PS：其中文件名相同，以.md结尾的为对应的中英文翻译markdown文件。




>*联系方式：*
>博客网址：[https://my.oschina.net/tantexian/blog ![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/click.jpg)](https://my.oschina.net/tantexian/blog)
>邮箱：tantexian@qq.com