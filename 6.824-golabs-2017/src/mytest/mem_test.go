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
// A send on a channel happens before the corresponding receive from that channel completes.
// Channel communication 是用来解决多个goroutine之间的同步的主要方法：
// 每个对特定channel的send操作 通常都会对应于其它goroutine中的receive操作，
// 那么通道的发送happens before与之对应的通道接收完成（此处可以是缓存通道和非缓冲通道）

// f函数中：对a的赋值happens before对于c通道的发送操作，
// 而对通道c的发送操作happens before对于c通道的接收完成，
// 而对通道c的接收完成happens before对应a的打印,
// 由于happens before的传递原则：a的赋值happens before print(a)发生，
// 因此能确保打印“hello, world”。
func Test_main3(t *testing.T) {
	go f3()
	<-c3 // 通道的接收（没有收到值，则会一直阻塞直到有值）
	print(a3)
}

//--------------------------------示例4代码分割线----------------------------------
var c4 = make(chan int, 10)

//var c4 = make(chan int)
var a4 string

// Golang 内存模型规则2：
// The closing of a channel happens before a receive that returns a zero value because the channel is closed.
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
	// 如果使用两次c6 <- 0 c6 <- 0 操作将缓存区塞满，然后能够达到示例6的效果
	// 缓冲区send塞满了，则必须阻塞等待receive取走，才能解除阻塞
	c6 <- 0
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

// Golang 内存模型规则7：
// 对于读写锁任何调用变量l的RLock，
// 第n次调用l.RLock必定发生在第n次l.Unlock之后,
// 与之相对应的l.RUnlock必定发生在第n+1次l.Lock之前。

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

//--------------------------------示例9代码分割线----------------------------------
var rwLock sync.RWMutex
var a9 string

func f9() {
	a9 = "hello, world rwLock a9  "
	rwLock.Unlock() // 第n次l.Unlock happen-before 第n次调用l.RLock
}

func f91() {
	a9 = "hello, world rwLock a9 change again !!!\n"
	rwLock.RUnlock() // 与之相对应的n次l.RUnlock happen-before 第n+1次l.Lock
}

// Golang 内存模型规则7：
// For any call to l.RLock on a sync.RWMutex variable l,
// there is an n such that the l.RLock happens (returns) after the n'th call to l.Unlock
// and the matching l.RUnlock happens before the n+1'th call to l.Lock.
// 对于读写锁任何调用变量l的RLock，
// 第n次调用l.RLock必定发生在第n次l.Unlock之后,
// 与之相对应的l.RUnlock必定发生在第n+1次l.Lock之前。

// 对于读写锁任何调用变量l,第n次Unlock happen-before 第n次RLock（Lock()、Unlock()必须成对出现）
// 即对读写锁加读写锁Lock()，必须等到该读写锁Unlock()解锁完，才能进行加读锁RLock()

// 总结：对于读写锁，每次加锁必须等待对应解锁操作完成，才能进行新的加锁
func Test_main9(t *testing.T) {
	for i := 0; i < 10; i++ {
		rwLock.Lock()
		go f9()
		rwLock.RLock()
		print(a9)

		go f91()
		rwLock.Lock()
		print(a9)
		rwLock.Unlock()
	}
}

//--------------------------------示例10代码分割线----------------------------------
var a10 string
var once10 sync.Once

func setup10() {
	a10 = "hello, world once\n"
}

func doprint() {
	// 由于第一次goroutine通过once.Do调用setup10()的返回 happen-before 后续goroutine 的调用返回
	// 因此能保证后续goroutine执行时候，第一次goroutine肯定已经执行完毕setup10，即，对a10的赋值
	once10.Do(setup10)
	print("a10 == " + a10)
}

// Golang 内存模型规则8：
// A single call of f() from once.Do(f) happens (returns) before any call of once.Do(f) returns.
// 通过once.Do(f)第一次调用f的返回 happen-before 任何其它的once.Do(f)返回

// sync包提供了一个安全的机制保证多个goroutine的初始化只被执行一次。
// 多个线程都会调用执行once.Do(f)函数，但是f函数只会执行一次。

// 总结：通过once.Do(f)对f函数的多次调用，永远能保证，第n次调用时候，第一次的f函数已经执行完毕，并返回
func Test_main_twoprint(t *testing.T) {
	go doprint()
	go doprint()
	go doprint()
	time.Sleep(1 * time.Second) // 此处用于阻止主线程过快结束，而其他goroutine还没来得及执行打印
}

//--------------------------------无效错误的惯例11代码分割线----------------------------------
var a11, b11 int

func f11() {
	a11 = 1
	b11 = 2
}

func g11() {
	print(b11)
	print(a11)
	print(" | ")

}

// 有可能打印出2和0（即b已经赋值为2，但是a还未被赋值）
func Test_main11(t *testing.T) {
	for i := 0; i < 100; i++ {
		go f11()
		g11()
	}
}

//--------------------------------无效错误的惯例12代码分割线----------------------------------
var a12 string
var done bool
var once12 sync.Once

func setup12() {
	a12 = "hello, world\n"
	done = true
}

func doprint12() {
	if !done {
		once12.Do(setup12)
	}
	print(a12)
}

// 这里不能保证，在doprint函数里面，能够观察到变量done的值，就暗示着能够观察到对a的赋值变化。
// 这里可能会错误的打印空，而不是“hello, world”。
func Test_main_twoprint12(t *testing.T) {
	go doprint12()
	go doprint12()
	time.Sleep(1 * time.Second) // 此处用于阻止主线程过快结束，而其他goroutine还没来得及执行打印
}

//--------------------------------无效错误的惯例13代码分割线----------------------------------
var a13 string
var done13 bool

func setup13() {
	a13 = "hello, world"
	done13 = true
}

// 上述示例所示，在主函数main中，观察到done值为true并不能保证一定能看到对a的赋值变化,
// 因此本示例程序也有可能打印空值。 更糟糕的是，setup函数对done=true的赋值，
// 并不一定能被main函数观察到，因此这两个线程之间并没有使用任何的同步事件。
// 因此for循环在主函数中并不能保证一定会循环结束。
func Test_main13(t *testing.T) {
	go setup13()
	for !done13 {
	}
	print(a13)
}

//--------------------------------无效错误的惯例14代码分割线----------------------------------
type T14 struct {
	msg14 string
}

var g14 *T14

func setup14() {
	t14 := new(T14)
	t14.msg14 = "hello, world msg14"
	g14 = t14
}

// 即使主函数观察到g!=null且退出了循环，
// 也不能保证main能够观察到g.msg的赋值内容。
func Test_main14(t *testing.T) {
	go setup14()
	for g14 == nil {
	}
	print(g14.msg14)
}

// 总结：无效错误的惯例一般，一般都是由于没有采用正确的同步导致，
// 因此，都可以采用同样的解决办法：使用明确的同步
