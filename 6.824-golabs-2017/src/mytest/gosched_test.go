package mytest

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"
)

var sched Sched
var GOMAXPROCS = 100
var nTasks = 10000
var doneTask = make([]int, 0)

// G 为goroutine接口
type G interface {
	Run()
}

// WorkTask 为可以运行的goroutine
type WorkTask struct {
	name string
	id   int
}

func (workTask *WorkTask) Run() {
	fmt.Printf("%v %v I,m start workTask...\n", workTask.id, workTask.name)
	doneTask = append(doneTask, workTask.id)
	workTime := time.Millisecond * time.Duration(rand.Intn(1000))
	time.Sleep(workTime)
	fmt.Printf("%v %v spend %v Millisecond\n", workTask.id, workTask.name, workTime)
}

// 调度函数用来将就绪的GQueue列表中的goroutine调度给m运行
type Sched struct {
	GQueue []G
	lock   sync.Mutex
}

func (sched *Sched) add(goroutineTask G) {
	sched.GQueue = append(sched.GQueue, goroutineTask)
}

func allTaskDone() bool {
	fmt.Println(len(doneTask))
	if nTasks == len(doneTask) {
		return true
	}
	return false
}

// 代表machine机器线程，能进行任务调度运行
func M() {
	for {
		sched.lock.Lock()
		if len(sched.GQueue) > 0 {
			g := sched.GQueue[0]
			sched.GQueue = sched.GQueue[1:]
			sched.lock.Unlock()
			g.Run()
		} else {
			sched.lock.Unlock()
		}
	}
}

func ProduceTaskToSched() {
	for i := 0; i < nTasks; i++ {
		workTask := new(WorkTask)
		workTask.id = i
		workTask.name = "goroutine-" + strconv.Itoa(i)
		sched.add(workTask)
		workTime := time.Millisecond * time.Duration(rand.Intn(100))
		time.Sleep(workTime)
	}
}

/**
 * 模拟go的线程调度模型：其中G、M分别对应上述数据结构，P对应为GOMAXPROCS
 * @author tantexian
 * @since 2017/3/16
 * @params
 */
func TestSched(t *testing.T) {
	go ProduceTaskToSched()
	// GOMAXPROCS代表机器核数，即机器最大同时执行的线程数量
	for i := 0; i < GOMAXPROCS; i++ {
		go M()
	}
	for {
		if allTaskDone() {
			fmt.Println("All task has done.")
			break
		}
		time.Sleep(time.Second)
	}

	sort.Ints(doneTask)
	fmt.Println(doneTask)

}
