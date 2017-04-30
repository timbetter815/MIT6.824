package mytest

import (
	"testing"
	"fmt"
)

func Test_rangech(t *testing.T) {
	max := 10
	rangeChan := make(chan bool, max)

	for index := 0; index < max; index++ {
		go func(int) {
			if index%5 == 0 {
				rangeChan <- true
			} else {
				rangeChan <- false
			}
		}(index)
	}

	/*for index := 0; index < max; index++ {
		fmt.Printf("c == %v\n", <-rangeChan)

	}*/

	// FIXME:Why deadlock?
	for c := range rangeChan {
		fmt.Printf("rangeChan == %v\n", rangeChan)
		fmt.Printf("c == %v\n", c)
	}
}
