package mytest

import (
	"fmt"
	"sort"
	"sync"
	"testing"
)

type MutexLock struct {
	ids  []int
	lock sync.Mutex
}

func (mutexLock *MutexLock) add(i int) {
	mutexLock.ids = append(mutexLock.ids, i)
}

type MutexLockPtr struct {
	ids  []int
	lock *sync.Mutex
}

func (mutexLockPtr *MutexLockPtr) add(i int) {
	mutexLockPtr.ids = append(mutexLockPtr.ids, i)
}

var mutexLock = new(MutexLock)
var mutexLockPtr = new(MutexLockPtr)

func TestMutexLock(t *testing.T) {
	for i := 0; i < 100; i++ {
		mutexLock.lock.Lock()
		mutexLock.add(i)
		mutexLock.lock.Unlock()
	}
	ids := mutexLock.ids
	sort.Ints(ids)
	fmt.Println(ids)
}

func TestMutexLockPtr(t *testing.T) {
	mutexLockPtr.lock = new(sync.Mutex) // 如何为*sync.mutex则需要显示初始化
	for i := 0; i < 100; i++ {
		mutexLockPtr.lock.Lock()
		mutexLockPtr.add(i)
		mutexLockPtr.lock.Unlock()
	}
	ids := mutexLockPtr.ids
	sort.Ints(ids)
	fmt.Println(ids)
}
