package mytest

import "testing"

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