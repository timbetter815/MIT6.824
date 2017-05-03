package mytest

import (
	"sync/atomic"
	"fmt"
	"testing"
	"time"
	"math/rand"
	"sync"
	"context"
)

func Test_cancel1(t *testing.T) {
	for index := 0; index < 10; index++ {

		ok := broadcastVoteWithContextCancel()
		if ok {
			fmt.Printf("%v get vote success.\n", index)
		} else {
			fmt.Printf("%v get vote fail.\n", index)
		}
	}
}

func broadcastVoteWithContextCancel() (bool) {
	var successNum int32 = 0
	successCh := make(chan bool, MaxServer)
	onceCancel := sync.Once{}
	ctx, cancel := context.WithCancel(context.TODO())
	for index := 0; index < MaxServer; index++ {
		go func(server int, ctx context.Context, cancel context.CancelFunc) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			fmt.Printf("exeNum == %v\n", atomic.AddInt32(&exeNum, 1))
			ok := getRequestVoteWithContext(server)
			if ok {
				if atomic.AddInt32(&successNum, 1) > MaxServer/2 {
					onceCancel.Do(func() {
						cancel()
					})
					successCh <- true
				} else {
					successCh <- false
				}
			} else {
				successCh <- false
			}
		}(index, ctx, cancel)
	}

	for index := 0; index < MaxServer; index++ {
		result := <-successCh
		if result {
			return true
		}
	}
	return false
}

func getRequestVoteWithContext(server int) bool {
	if server == MaxServer-2 {
		time.Sleep(1 * time.Second)
	}
	// 模拟1/5的几率失败
	intn := rand.Intn(5)
	if intn == 0 {
		return false
	}
	return true
}
