package mytest

import "testing"

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

// Channel communication 主要的方法来解决多个goroutine之间的同步。
// 每发送一个特定的信道匹配相应的接收通道，通常在不同的goroutine中。
// 通道的发送happens before与之对应的通道接收完成（此处可以是缓存通道）
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

// 无缓冲的通道的接收操作 happens before 发送操作结束之前（即发送的值立马会传递到接收，等接收完成，发送才会结束）
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

// The kth receive on a channel with capacity C happens before the k+Cth send from that channel completes.
// 第k次对通道的接收
// 如果通道为缓存通道，则不能保证示例5中的happens-before原则，因此不保证一定打印a6的内容
// 如果发送使用两次c6 <- 0 c6 <- 0 操作将缓存区塞满，然后能够达到示例6的效果
// 总结：对于无缓冲通道或者缓存通道，当塞满的通道的（包括缓存区），那么当前位置会被阻塞，直到有接受者取走一条信号
func Test_main6(t *testing.T) {
	go f6()
	c6 <- 0 // 如果使用两次c6 <- 0 c6 <- 0 操作将缓存区塞满，然后能够达到示例6的效果
	print(a6)
}

//--------------------------------示例7代码分割线----------------------------------
