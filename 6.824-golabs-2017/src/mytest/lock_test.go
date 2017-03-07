package mytest

import (
	"sync"
	"fmt"
	"testing"
	"time"
)

/**
 * @author tantexian(https://my.oschina.net/tantexian/blog)
 * @since 2017/3/6
 * @params
 */
type  MyLock struct {
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

