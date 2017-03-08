package mytest

import (
	"testing"
	"fmt"
)

/**
 * @author tantexian(https://my.oschina.net/tantexian/blog)
 * @since 2017/3/7
 * @params
 */
func TestSwitchNoBreake(t *testing.T) {
	str := "1"
	switch str {
	case "1":
		fmt.Println("case 1")
	case "2":
		fmt.Println("case 2")
	case "3":
		fmt.Println("case 3")
	default:
		fmt.Println("no switch!!!")
	}

}

func TestSwitchWithBreake(t *testing.T) {
	str := "1"
	switch str {
	case "1":
		fmt.Println("case 1")
		break
	case "2":
		fmt.Println("case 2")
		break
	case "3":
		fmt.Println("case 3")
		break
	default:
		fmt.Println("no switch!!!")
	}

	// 被延期执行的函数，它的参数（包括接收者，如果函数是一个方法）是在defer执行的时候被求值的，而不
	// 是在调用执行的时候。这样除了不用担心变量随着函数的执行值会改变，这还意味着单个被延期执行的调
	// 用点可以延期多个函数执行。因此输出结果为：4 3 2 1 0
	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i)
	}

	a := [...]string {1: "no error", 2: "Eio", 3: "invalid argument"}
	for _, val := range a {
		fmt.Println(val)
	}
	fmt.Println("-----------------------")

	s := []string {1: "no error", 2: "Eio", 3: "invalid argument"}
	for _, val := range s {
		fmt.Println(val)
	}
	fmt.Println("-----------------------")


	m := map[int]string{1: "no error", 2: "Eio", 3: "invalid argument"}
	for _, val := range m {
		fmt.Println(val)
	}
	fmt.Println("-----------------------")


}