package mytest

import (
	"testing"
	"fmt"
	"io"
)

type ByteSize float64

const (
	_ = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
)

func Test1(t *testing.T) {
	fmt.Printf("KB=%v, MB=%v, GB=%v, TB=%v", KB, MB, GB, TB)
	io.Writer()
}
