package mytest

import (
	"fmt"
	"testing"
)

func myDefer() {
	fmt.Println("I,m myDefer...")
}

/**
 * 定义多少次defer则会调用多少次defer定义语句或者函数
 *
 * @author tantexian<my.oschina.net/tantexian>
 * @since 2017/3/5
 * @params
 */
func TestMultiDefer(t *testing.T) {
	defer myDefer()
	fmt.Println("I,m main executing...")
	defer myDefer()
}
