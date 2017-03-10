package mytest

import (
	"testing"
	"fmt"
	"sync"
	"time"
	"math/rand"
)

var MaxOutstanding = 10
var clientRequests chan *Request
var quit chan bool
var once sync.Once

func initChan() {
	quit = make(chan bool)
	clientRequests = make(chan *Request, MaxOutstanding)
}

type Request struct {
	args       []int
	f          func([]int) int
	resultChan chan int
}

func sum(a []int) (s int) {
	for _, v := range a {
		s += v
	}
	return
}

func handle(queue chan *Request) {
	fmt.Printf("handle queue==%v\n", queue)
	for req := range queue {
		fmt.Printf("handle req==%v\n", req)
		req.resultChan <- req.f(req.args)
	}
}

func Serve(clientRequests chan *Request, quit chan bool) {
	// Start handlers
	for i := 0; i < MaxOutstanding; i++ {
		go handle(clientRequests)
	}
	<-quit // Wait to be told to exit.
}

func sendRequestToServer() {
	arg1 := rand.Intn(10)
	arg2 := rand.Intn(10)
	arg3 := rand.Intn(10)
	request := &Request{[]int{arg1, arg2, arg3}, sum, make(chan int)}
	fmt.Printf("[%v]: client:[{arg1:%v}{arg2:%v}{arg3:%v}]\n", time.Now(), arg1, arg2, arg3)
	// Send request
	clientRequests <- request
	// Wait for response.
	fmt.Printf("client:[{arg1:%v}{arg2:%v}{arg3:%v}]answer: %d\n", arg1, arg2, arg3, <-request.resultChan)

}

func TestClient(t *testing.T) {
	once.Do(initChan)
	// 每次间隔1000毫秒产生一次计算请求
	for i := 0; i < MaxOutstanding; i++ {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		go sendRequestToServer()
	}
}

func TestStartServer(t *testing.T) {
	once.Do(initChan)
	go Serve(clientRequests, quit)
	for {
		time.Sleep(time.Second)
		fmt.Printf("[%v]: I'm server, and runing... \n", time.Now())
	}
	//quit <- true
}