package mapreduce

import (
	"fmt"
	"sync"
)

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	// map阶段
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	// reduce阶段
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	/**
	 * @author tantexian<my.oschina.net/tantexian>
	 * @since 2017/3/4
	 * @params
	 */
	var wg sync.WaitGroup

	// 如果为map阶段，则循环获取mapfile文件，如果为reduce阶段，则循环nreduce次数
	for i := 0; i < ntasks; i++ {
		wg.Add(1)

		go func(index int) {
			// 函数退出时候执行
			defer wg.Done()
			STARTRPC:
			// 获取空闲worker
			idleWorker := <-registerChan

			// 根据net/rpc[Go官方库RPC开发指南:https://my.oschina.net/tantexian/blog/851914]
			// 可知 此处rpc调用Worker结构体的DoTask方法
			success := call(idleWorker, "Worker.DoTask", DoTaskArgs{jobName, mapFiles[index], phase, index, n_other}, nil)
			if success == true {
				// 执行成功则将worker继续放入到空闲worker池中
				go func() {// 因为任何对通道发送信息及获取信息都是阻塞的，因此需要使用goroutine
					registerChan <- idleWorker
				}()
			} else {
				fmt.Printf("Master: RPC %s DoTask error\n", idleWorker)
				// 执行失败则重新执行
				goto STARTRPC
			}
		}(i)
	}

	// 一直等到wg为0
	wg.Wait()
	fmt.Printf("Schedule: %v phase done\n", phase)
}
