package mytest

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestTime1(t *testing.T) {
	fmt.Printf("start time == %v\n", time.Now())
	time.After(time.Duration(3))
	fmt.Printf("end time == %v\n", time.Now()) // 立即输出
}

// TestTime2 用于测试定时心跳任务
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/29
func TestTime2(t *testing.T) {
	heartbeatChan := make(chan string, 10)
	fmt.Printf("start time == %v\n", time.Now())

	// 启动goroutine每隔1秒钟发送一次心跳到heartbeatChan  Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/29
	go func() {
		for i := 0; i < 10; i++ {
			str := strconv.Itoa(i) + ":" + time.Now().String()
			heartbeatChan <- str
			time.Sleep(1 * time.Second)
		}
	}()

	// 启动goroutine接收心跳，两秒未收到心跳则认为超时  Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/29
	go func() {
		for {
			select {
			case msg := <-heartbeatChan:
				fmt.Printf("heartbeat msg : %v time == %v\n", msg, time.Now())
			case timeCh := <-time.After(2 * time.Second):
				fmt.Printf("#### Timeout timeCh:%v time == %v\n", timeCh, time.Now())
			}
		}
	}()

	time.Sleep(60 * time.Second)
	fmt.Printf("main is exit: %v\n", time.Now())
}
