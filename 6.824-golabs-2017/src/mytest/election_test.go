package mytest

import (
	"sync/atomic"
	"fmt"
	"testing"
	"time"
	"math/rand"
)

const MaxServer = 10

func Test_election(t *testing.T) {
	for index := 0; index < 10; index++ {
		ok, i := broadcastVote(index)
		if ok {
			fmt.Printf("%v get vote success.\n", i)
		} else {
			fmt.Printf("%v get vote fail.\n", i)
		}
	}
}

func broadcastVote(i int) (bool, int) {
	var successNum int32 = 0
	successCh := make(chan bool, MaxServer)
	for index := 0; index < MaxServer; index++ {
		go func(server int) {
			ok := getRequestVote(server)
			if ok {
				addInt32 := atomic.AddInt32(&successNum, 1)
				if addInt32 > MaxServer/2 {
					successCh <- true
				} else {
					successCh <- false
				}
			} else {
				successCh <- false
			}
		}(index)
	}

	for index := 0; index < MaxServer; index++ {
		ch := <-successCh
		if ch {
			return true, i
		}
	}
	return false, i
}

func getRequestVote(server int) bool {
	if server == MaxServer/2+1 {
		time.Sleep(1 * time.Second)
	}
	// 模拟一半的几率失败
	intn := rand.Intn(2)
	if intn == 0 {
		return false
	}
	return true
}
