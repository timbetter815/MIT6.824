####本项目主要针对自己对MIT6.824课程学习和翻译的记录，也方便后续其他学习者作为参考。
有些内容也会记录在博客网址：[https://my.oschina.net/tantexian/blog](https://my.oschina.net/tantexian/blog)

###[MIT6.824网址](http://nil.csail.mit.edu/6.824/2017/)

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
该课程的学习入口：[Schedule入口](http://nil.csail.mit.edu/6.824/2017/schedule.html)开始学习LEC1-LEC23以及课程学习资料。
![](https://git.oschina.net/tantexian/MIT6.824/raw/dev/resources/static/img/scheduler-2017.png?dir=0&filepath=resources%2Fstatic%2Fimg%2Fscheduler-2017.png&oid=9c38ebc90be76b60603943bec0b9769289e6da86&sha=dec89cbe68c8c8eb7bd1a5020da1ec3bfe3a534d)


### 整体课程介绍：
>[请点击：英文介绍及中文翻译](https://git.oschina.net/tantexian/MIT6.824/blob/dev/LEC1/101.md?dir=0&filepath=LEC1%2F101.md&oid=10b8b103df840aac63d52672faa063bf9a7d1453&sha=2fa7e192e9229a20cae87e5011212091cb95a1bb)

----------

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
* 注释类型1、

```
/**
 * 顺序消息的生产者（顺序消息的消费者与普通消费者一致）
 *
 * @author xxx
 * @since 2016/6/27
 * @params ctx：Context类型上下文
 */
```

* 注释类型2、

```
// 由于是顺序消息，因此只能选择一个queue生产和消费消息
```

* 注释类型3、

```
// xxx 2016/7/11 Add by xxx Or // xxx 2016/7/11 Edit by xxx
```

* 关于TODO、FIXME、XXX注释规范、

```
// TODO: + 说明：
```
如果代码中有TODO该标识，说明在标识处有功能代码待编写，待实现的功能在说明中会简略说明。

```
// FIXME: + 说明：
```
如果代码中有FIXME该标识，说明标识处代码需要修正，甚至代码是错误的，不能工作，需要修复，如何修正会在说明中简略说明。

```
// XXX: + 说明：
```
如果代码中有XXX该标识，说明标识处代码虽然实现了功能，但是实现的方法有待商榷，希望将来能改进，要改进的地方会在说明中简略说明。



### 开发IDE
* 开发工具不做统一规定（Idea、Eclipse都可以），建议使用Idea按照golang插件
* 建议使用最新版格式Idea，附下载地址：http://pan.baidu.com/s/1slMkXY1
* 附Idea属性格式注释文件下载地址：http://pan.baidu.com/s/1hrU3IgW（其中java版本Idea使用zz或者zzz命令来生成注释，golang使用gg或者ggg）

----------


>*联系方式：*
>博客网址：[https://my.oschina.net/tantexian/blog](https://my.oschina.net/tantexian/blog)
>邮箱：tantexian@qq.com