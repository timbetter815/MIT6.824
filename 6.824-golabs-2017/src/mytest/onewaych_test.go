package mytest

import (
	"fmt"
	"testing"
)

func Test_chan_reveice(t *testing.T) {
	ci := make(chan int)
	go just_receive(ci)
	ci <- 1
}

func Test_chan_send(t *testing.T) {
	ci := make(chan int)
	go just_send(ci)
	fmt.Printf("%v\n", <-ci)
}

func just_receive(receive <-chan int) {
	fmt.Printf("%v\n", <-receive)
}

func just_send(send chan<- int) {
	send <- 888
}

/*func Test_chan_error(t *testing.T) {
	ci := make(chan int)
	go just_receive_error(ci)
	fmt.Printf("%v\n", <-ci)
}

func just_receive_error(receive <-chan int) {
	receive <- 666
}*/
