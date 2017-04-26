package mytest

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	"math/rand"
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

func TestTime3(t *testing.T) {
	fmt.Printf("start time == %v\n", time.Now())

	success := make(chan bool)

	for index := 0; index < 10; index++ {
		time.Sleep(1 * time.Second)
		go func() {
			if rand.Int()%5 == 0 {
				//success <- true
			}
		}()

	}

	select {
	case y := <-success:
		fmt.Printf("success %v\n", y)
	case x := <-time.After(2 * time.Second):
		fmt.Printf("[%v] timeout == %v\n", x, time.Now()) // 立即输出
	}

	fmt.Printf("end time == %v\n", time.Now()) // 立即输出
}

func TestTime4(t *testing.T) {
	timeout := 8
	fmt.Printf("Start: %v\n", time.Now())
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	go func() {
		for {
			fmt.Printf("%v\n", time.Now())
			time.Sleep(1 * time.Second)
		}
	}()
	select {
	case <-timer.C:
		fmt.Printf("timeout: %v\n", time.Now())
	}
	fmt.Printf("Finished: %v\n", time.Now())
}

func TestTime5(t *testing.T) {
	timeStr := time.Now().String()
	fmt.Printf("%v\n", timeStr)
	time, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", timeStr)
	if err == nil {
		fmt.Printf("time == %v\n", time)
	} else{
		fmt.Printf("%v\n",err)
	}

}
