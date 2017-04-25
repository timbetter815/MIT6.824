package mytest

import (
	"testing"
	"time"
	"sync/atomic"
	"log"
	"runtime"
	"fmt"
	"sync"
)

var voteNum int32 = 0

const Max = 100

func Test_NoClose(t *testing.T) {
	once := sync.Once{}
	for index1 := 0; index1 < 100; index1++ {
		fmt.Printf("-------------new for index1[%v]\n", index1)
		var ci = make(chan int)
		// var finishedCh = make(chan bool)
		for index2 := 0; index2 < Max; index2++ {
			go func(i1 int, i2 int) {
				success := sendRequestVote(i2)
				if success {
					atomic.AddInt32(&voteNum, 1)
				}

				if voteNum > Max/2 {
					select {
					case ci <- i2:
					default:
						fmt.Println("channel is full !")
						return
					}
					log.Printf("[Quorum-END]: finishNum[%v] > Max[%v]\n", voteNum, Max)
				}
			}(index1, index2)

			once.Do(func() {
				go func() {
					fmt.Printf("once.Do print NumGoroutine\n")
					for {
						fmt.Printf("runtime.NumGoroutine() == %v\n", runtime.NumGoroutine())
						time.Sleep(1000 * time.Millisecond)
					}
				}()
			})
		}
		select {
		case <-ci:
			log.Printf("Has receive finishNum %v\n", ci)
		case <-time.After(3 * time.Second):
			log.Printf("Timeout %v", time.Now())
		}
		time.Sleep(1000 * time.Millisecond)
	}

	time.Sleep(100 * time.Second)
}

func sendRequestVote(server int) bool {
	if server%10 == 0 {
		log.Printf("server[%v] dont send vote.\n", server)
		time.Sleep(3 * time.Second)
		return false
	}
	return true
}
