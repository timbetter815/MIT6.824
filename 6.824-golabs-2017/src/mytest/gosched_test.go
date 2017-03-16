package mytest

import (
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

var sched = new(Sched)
var GOMAXPROCS = 20
var nTasks = 100
var doneTask = make([]int, 0)
var doneTaskLock sync.Mutex

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
	doneTaskLock.Lock()
	doneTask = append(doneTask, workTask.id)
	doneTaskLock.Unlock()
	workTime := time.Microsecond * time.Duration(rand.Intn(10))
	time.Sleep(workTime)
	fmt.Printf("%v %v spend %v Millisecond\n", workTask.id, workTask.name, workTime)
}

// 调度函数用来将就绪的GQueue列表中的goroutine调度给m运行
type Sched struct {
	GQueue []G
	lock   sync.Mutex
}

func (sched *Sched) add(goroutineTask G) {
	sched.lock.Lock()
	sched.GQueue = append(sched.GQueue, goroutineTask)
	sched.lock.Unlock()
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
	fmt.Printf("#### %v GroutineId=%v\n", time.Now(), GroutineId())
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

func GroutineId() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

/**
 * 模拟go的线程调度模型：其中G、M分别对应上述数据结构，P对应为GOMAXPROCS
 * 其中goroutine的工作模式概括：同时启动GOMAXPROCS个M同时并行不停地从Sched中的就绪GQueue列表中获取groutine执行。
 * @author tantexian
 * @since 2017/3/16
 * @params
 */
func TestSched(t *testing.T) {
	fmt.Printf("#### [%v]: TestSched GroutineId=%v\n", time.Now(), GroutineId())
	go ProduceTaskToSched()
	// GOMAXPROCS代表机器核数，即机器最大同时执行的线程数量
	for i := 0; i < GOMAXPROCS; i++ {
		go M()
	}
	for {
		doneTaskLock.Lock()
		if allTaskDone() {
			fmt.Println("All task has done.")
			break
		}
		doneTaskLock.Unlock()
		fmt.Printf("[%v]: NumGoroutine=%v NumCPU=%v\n", time.Now(), runtime.NumGoroutine(), runtime.NumCPU())
		time.Sleep(time.Second)
	}

	sort.Ints(doneTask)
	fmt.Println(doneTask)

}
