package mytest

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

/**
 * @author tantexian(https://my.oschina.net/tantexian/blog)
 * @since 2017/3/6
 * @params
 */
type MyLock struct {
	a    int
	lock *sync.Mutex
}

func (myLock *MyLock) GlobalLock_A() {
	myLock.lock.Lock()
	fmt.Println("START-A: I'm testA....")
	time.Sleep(time.Millisecond * 100)
	fmt.Println("END-A: I'm testA....")
	myLock.lock.Unlock()
}

func (myLock *MyLock) GlobalLock_B() {
	myLock.lock.Lock()
	fmt.Println("START-B: I'm testB....")
	time.Sleep(time.Millisecond * 200)
	fmt.Println("END-B: I'm testB....")
	myLock.lock.Unlock()
}

// 此次使用的是全局锁，所有全局方法竞争一把锁
func Test_GlobalLock(t *testing.T) {
	lock := sync.Mutex{}
	myLock1 := MyLock{1, &lock}
	go myLock1.GlobalLock_A()
	go myLock1.GlobalLock_B()

	myLock2 := MyLock{1, &lock}
	go myLock1.GlobalLock_A()
	go myLock2.GlobalLock_B()

	time.Sleep(time.Second * 10)
}

func (myLock *MyLock) OBJLock_A() {
	lock := sync.Mutex{}
	lock.Lock()
	fmt.Println("START-A: I'm testA....")
	time.Sleep(time.Millisecond * 100)
	fmt.Println("END-A: I'm testA....")
	lock.Unlock()
}

func (myLock *MyLock) OBJLock_B() {
	lock := sync.Mutex{}
	lock.Lock()
	fmt.Println("START-B: I'm testB....")
	time.Sleep(time.Millisecond * 200)
	fmt.Println("END-B: I'm testB....")
	lock.Unlock()
}

// 各方法使用各自对象的锁
func Test_OBJLock(t *testing.T) {
	lock := sync.Mutex{}
	myLock1 := MyLock{1, &lock}
	go myLock1.OBJLock_A()
	go myLock1.OBJLock_A()
	go myLock1.OBJLock_A()
	go myLock1.OBJLock_B()

	myLock2 := MyLock{1, &lock}
	go myLock1.OBJLock_A()
	go myLock2.OBJLock_B()

	time.Sleep(time.Second * 10)
}

var rwlock sync.RWMutex
var str3 string

func fun3() {
	time.Sleep(100 * time.Millisecond)
	str3 = "hello, world"
	rwlock.Unlock()

}

func fun4() {
	time.Sleep(100 * time.Millisecond)
	str3 = "hello, world2"
	rwlock.Lock()

}

func TestRWMutex(t *testing.T) {
	rwlock.Lock()
	go fun3()
	go fun4()
	str3 = "hello, world3"
	rwlock.RLock() // 这是第n(1)次出现RLock，应该要先于第n+1(2)次lock的发生，即应该先于fun4中的lock之前发生。
	print(str3)
}

var l sync.Mutex
var a1 string

// 由于unlock在a1赋值之前，因此主线程可能看不到a1赋值的值
func fun1() {
	l.Unlock()
	time.Sleep(100 * time.Millisecond)
	a1 = "hello, world"

}

func TestMutex1(t *testing.T) {
	l.Lock()
	go fun1()
	l.Lock()
	print(a1)
}

func f2() {
	a1 = "hello, world"
	l.Unlock()

}

// 对同一个变量l的解锁n次和加锁操作m次，如果n<m,即解锁次数1小于加锁次数2，因此解锁操作先于加锁操作完成
// 即f函数的unlock先于第二个lock完成，而第二个lock又先于print完成，因此保证a1的赋值能被打印出来
func TestMutex2(t *testing.T) {
	l.Lock()
	go f2()
	l.Lock() // 如果此处语句屏蔽，则a1打印为空
	print(a1)
	time.Sleep(time.Second)
}

var a111 string

func f() {
	println(a111)
}

func hello() {
	a111 = "hello, world"
	go f()
}

func hello1() {
	go func() {
		a111 = "hello"
	}()
	print(a111)
	time.Sleep(time.Second)
}

func TestGoroutine(t *testing.T) {
	hello1()
}

var limit = make(chan int, 5)

func printHello(i int) {
	fmt.Println(strconv.Itoa(i))
}

func TestLimit(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			limit <- 1
			go printHello(index)
			<-limit
		}(i)
	}
	wg.Wait()
	// select {}
}
