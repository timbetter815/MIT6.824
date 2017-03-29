package mytest

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
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

	//time.Sleep(3 * time.Second)
	//close(ci)
	go func() {
		for i := range ci {
			fmt.Println(i)
		}
	}()
	go put2Ch(ci, 4)
	go put2Ch(ci, 5)
	go put2Ch(ci, 6)
	time.Sleep(time.Second)
	go put2Ch(ci, 7)
	go put2Ch(ci, 8)
	go put2Ch(ci, 9)
	time.Sleep(time.Second)
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

/**
 * @author tantexian(https://my.oschina.net/tantexian/blog)
 * @since 2017/3/13
 * 输出值为：0，0 | 1，-2 |
 * 对于不同实例pos或者neg，不会共享sum。但是对于同一个实例的多次调用会共享共同变量sum
 */
func Test3(t *testing.T) {
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}
}

// 但是对于同一个实例的多次调用会共享共同变量sum
func Test4(t *testing.T) {
	pos := adder()
	for i := 0; i < 10; i++ {
		fmt.Println(pos(i))
	}
}

// 不同实例则不会共享共同变量sum
func Test5(t *testing.T) {
	adder1 := adder()
	adder2 := adder()
	adder3 := adder()
	a1 := adder1(10)
	a11 := adder1(888)
	a2 := adder2(100)
	a3 := adder3(100)
	fmt.Printf("a1=%v, a11=%v, a2=%v, a3=%v", a1, a11, a2, a3)
}

func Test6(t *testing.T) {
	s := make([]int, 2)
	p1 := &s[0]
	ptr := unsafe.Pointer(&s[0])
	fmt.Printf("s=%v, p1=%v, ptr=%v\n", s, p1, ptr)
}

type X struct {
	a int
}

func Test7(t *testing.T) {
	x := &X{}
	fmt.Printf("x.a == %v", x.a)
//	x.a = nil
	/*if x.a == nil {

	}*/
}
