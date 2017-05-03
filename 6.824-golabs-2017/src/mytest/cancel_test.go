package mytest

import (
	"sync/atomic"
	"fmt"
	"testing"
	"time"
	"math/rand"
	"sync"
)

var exeNum int32 = 0

func Test_cancel(t *testing.T) {
	for index := 0; index < 10; index++ {
		ok := broadcastVoteWithCancel()
		if ok {
			fmt.Printf("%v get vote success.\n", index)
		} else {
			fmt.Printf("%v get vote fail.\n", index)
		}
	}
}

func wasCancelled(cancelCh chan struct{}) bool {
	select {
	case <-cancelCh:
		return true
	default:
		return false
	}
}

func doCancel(cancelCh chan struct{}) {
	close(cancelCh)
}

func broadcastVoteWithCancel() (bool) {
	cancelChan := make(chan struct{})
	var successNum int32 = 0
	successCh := make(chan bool, MaxServer)
	onceCancel := sync.Once{}
	for index := 0; index < MaxServer; index++ {
		go func(server int) {
			if wasCancelled(cancelChan) {
				fmt.Printf("now[%v]: wasCancelled.\n", time.Now())
				return
			} else {
				fmt.Printf("exeNum == %v\n", atomic.AddInt32(&exeNum, 1))
				ok := getRequestVote1(server)
				if ok {
					if atomic.AddInt32(&successNum, 1) > MaxServer/2 {
						onceCancel.Do(func() {
							doCancel(cancelChan)
						})
						successCh <- true
					} else {
						successCh <- false
					}
				} else {
					successCh <- false
				}
			}
		}(index)
	}

	for index := 0; index < MaxServer; index++ {
		result := <-successCh
		if result {
			return true
		}
	}
	return false
}

func getRequestVote1(server int) bool {
	if server == MaxServer-2 {
		time.Sleep(1 * time.Second)
	}
	// 模拟1/5的几率失败
	intn := rand.Intn(2)
	if intn == 0 {
		return false
	}
	return true
}
