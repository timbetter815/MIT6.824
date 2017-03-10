package mytest

import (
	"testing"
	"fmt"
	"time"
)

type T struct {
	msg string
}

var g *T

func setup() {
	t := new(T)
	t.msg = "hello, world"
	g = t
}

func TestTmp(t *testing.T) {
	go setup()
	for g == nil {
	}
	print(g.msg)
}

func put2Ch(ci chan int, val int) {
	ci <- val
}

func TestRangeChan(t *testing.T) {
	ci := make(chan int, 4)
	go put2Ch(ci, 1)
	go put2Ch(ci, 2)
	go put2Ch(ci, 3)

	time.Sleep(3 * time.Second)
	//close(ci)
	for i := range ci {
		fmt.Println(i)
	}
	time.Sleep(3 * time.Second)

}