/**
 * Golang内存模型happen-before示例代码
 *
 * @author tantexian, <my.oschina.net/tantexian>
 * @since 2017/4/19
 * @params
 */

package mytest

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

//--------------------------------示例1代码分割线----------------------------------
var a string

func f1() {
	print(a)
}

// 启动goroutine线程happens before goroutine的执行
// 由于线程没有使用同步或者通道，因此有可能主线程在goroutine线程之前或者之后完成
// 所以"hello, world"有可能在Test_hello1函数返回之后才会执行打印，所以可能没有输出
func Test_hello1(t *testing.T) {
	a = "hello, world"
	go f1()
}

//--------------------------------示例2代码分割线----------------------------------
// goroutine的退出不保证一定会先于其他事件发生。
func Test_hello2(t *testing.T) {
	go func() { a = "hello" }()
	print(a) // 由于没有使用同步相关事件，因此有可能a的赋值对其他线程是不可见的
}

//--------------------------------示例3代码分割线----------------------------------
var c3 = make(chan int, 10)
var a3 string

func f3() {
	a3 = "hello, world"
	c3 <- 0 // 通道的发送
}

// Golang 内存模型规则1：
// Channel communication 是用来解决多个goroutine之间的同步的主要方法：
// 每个对特定channel的send操作 通常都会对应于其它goroutine中的receive操作，
// 那么通道的发送happens before与之对应的通道接收完成（此处可以是缓存通道）

// f函数中：对a的赋值happens before对于c通道的发送操作，
// 而对通道c的发送操作happens before对于c通道的接收完成，
// 而对通道c的接收完成happens before对应a的打印,
// 由于happens before的传递原则：a的赋值happens before print(a)发生，
// 因此能确保打印“hello, world”。
func Test_main3(t *testing.T) {
	go f3()
	<-c3 // 通道的接收
	print(a3)
}

//--------------------------------示例4代码分割线----------------------------------
var c4 = make(chan int, 10)

//var c4 = make(chan int)
var a4 string

// Golang 内存模型规则2：
// 通道的关闭 happens before 接收到返回值0（由于通道关闭将会发送0到该通道）

// 由于a4的赋值 happens before close(c4)
// close(c4) happens before 通道的接收(<-c4)
// 通道的接收(<-c4) happens before print(a4)
// 因此对a4的赋值 happens before print(a4)，所以能确保打印出a4的值
func f4() {
	a4 = "hello, world"
	// c4 <- 0 // 通道的发送
	close(c4)
}

// 因此能确保打印“hello, world”。
func Test_main4(t *testing.T) {
	go f4()
	<-c4 // 通道的接收
	print(a4)
}

//--------------------------------示例5代码分割线----------------------------------
var c5 = make(chan int)
var a5 string

func f5() {
	a5 = "hello, world"
	<-c5
}

// Golang 内存模型规则3：
// 无缓冲的通道的接收操作 happens before 发送操作结束之前（即发送的值立马会传递到接收，等接收完成，发送才会结束）

// 结合规则1和规则3 happen-before 规则 : 同一个无缓冲通道：send开始 -> receive开始 -> receive结束 -> send结束

// 对a5的赋值 happens before c5通道的接收（<-c5），
// c5通道的接收（<-c5） happens before 对c5通道的发送（c5 <- 0）结束。[即：如果c5发送结束，那么肯定c5接收完成了。]
// 对c5通道的发送（c5 <- 0）结束  happens before print(a5)，因此保证一定能打印出a5值
func Test_main5(t *testing.T) {
	go f5()
	c5 <- 0
	print(a5)
}

//--------------------------------示例6代码分割线----------------------------------
var c6 = make(chan int, 1) // 代表为缓存chan, 拥有一个缓存区间
var a6 string

func f6() {
	a6 = "hello, world"
	<-c6
}

// Golang 内存模型规则5（规则3的缓存通道规则）：
// The kth receive on a channel with capacity C happens before the k+Cth send from that channel completes.
// 第k次对通道的接收（通道的缓冲大小为C） happen-before 第k+C次的send的完成。

// 如果通道为缓存通道，则不能保证示例5中的happens-before原则，因此不保证一定打印a6的内容
// 如果发送使用两次c6 <- 0 c6 <- 0 操作将缓存区塞满，然后能够达到示例6的效果
// 总结：对于无缓冲通道或者缓存通道，当塞满的通道的（包括缓存区），那么当前位置会被阻塞，直到有接受者取走一条信号
func Test_main6(t *testing.T) {
	go f6()
	c6 <- 0 // 如果使用两次c6 <- 0 c6 <- 0 操作将缓存区塞满，然后能够达到示例6的效果
	print(a6)
}

//--------------------------------示例7代码分割线----------------------------------
var limitChan = make(chan int, 3)
var workTasks = 10
var work = make([]func(), 0, workTasks)

func intWork() {
	for i := 0; i < workTasks; i++ {
		tmpI := i // 解决闭包问题，定义临时变量tmp
		work = append(work, func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("[%v]: Task%v\n", time.Now().String(), tmpI)
		})
	}
}

// 对于缓冲通道limitChan容量为3，由于每次新开启一个goroutine
// 因此很快send三次塞满了limitChan，
// 根据Golang 内存模型规则5，第1次对通道limitChan的receive happen-before 与第1+3次对通道的send的完成。
// 所以第4次对通道的send完成之前，必须等到第1次对通道的接收完成，
// 因此实现了限制同时最大limit个并发执行
func Test_main7(t *testing.T) {
	// 初始化work
	intWork()

	wg := sync.WaitGroup{}
	wg.Add(workTasks)

	// 遍历work
	for _, w := range work {
		go func(w func()) {
			limitChan <- 1
			w()
			<-limitChan
			wg.Done()
		}(w)
	}

	wg.Wait()
	// select {}
}

//--------------------------------示例8代码分割线----------------------------------
var lock sync.Mutex
var a8 string

func f8() {
	a8 = "hello, world lock"
	lock.Unlock()
}

// Golang 内存模型规则6：
// 对于任意 sync.Mutex 或 sync.RWMutex 变量l。
// 如果 n < m ，那么第n次 l.Unlock() happens before 第m次l.Lock()调用返回。

// 同一个互斥变量的解锁 happen-before 与该变量的下一次加锁（先发生解锁，再发生加锁）
// 对同一个变量lock的解锁n次和加锁操作m次，
// 如果n<m,即解锁次数1小于加锁次数2，因此解锁操作先于加锁操作完成
// 即f8函数的unlock先于第二个lock完成，而第二个lock又先于print完成，因此保证a1的赋值能被打印出来
func Test_main8(t *testing.T) {
	lock.Lock()
	go f8()
	lock.Lock()
	print(a8)
}
