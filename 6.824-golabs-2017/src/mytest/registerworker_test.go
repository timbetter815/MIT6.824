package mytest

import (
	"fmt"
	"testing"
	"sync"
	"strconv"
	"math/rand"
	"time"
)

type Task struct {
	index            int
	handleWorkerName string
	status           string
}

func TestRegisterWorker(t *testing.T) {
	allWorkers := make(chan string)
	go registerWorker(allWorkers, 10)
	// var idleWorker string 此处不能定义为全局变量，否则会有并发问题
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()
			START:
			idleWorker := <-allWorkers

			num := rand.Intn(100)
			if (num % 10 == 0) {
				fmt.Printf("---- index==%v is fail, idleWorker==%v \n", index, idleWorker)
				goto START

			} else {
				fmt.Printf("index==%v is success, idleWorker==%v \n", index, idleWorker)
				// allWorkers <- idleWorker 所有针对chan的操作默认都是阻塞的，因此需要启动goroutine
				go func() {
					// 模拟执行成功
					allWorkers <- idleWorker
				}()
			}
		}(i)

	}
	wg.Wait()
	fmt.Println("All works task done!!!")
}

func registerWorker(allWorkers chan string, workNum int) {
	for i := 0; i < workNum; i++ {
		if (i % 5 == 0) {
			time.Sleep(time.Second)
		}
		go func(index int) {
			allWorkers <- ("worker" + strconv.Itoa(index))
		}(i)
	}
}