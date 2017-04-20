package mytest

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_closure(t *testing.T) {
	var voteNum int32 = 0
	var wg = sync.WaitGroup{}
	var len int32 = 100

	wg.Add(int(len))
	for index := int32(0); index < len; index++ {
		go func() {
			num := rand.Intn(3)
			if num%3 == 0 || num%2 == 0 {
				atomic.AddInt32(&voteNum, 1)
			}

			if voteNum > len/2 {
				fmt.Printf("[vote success]: %v\n", voteNum)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("[result]: %v", voteNum)
}
