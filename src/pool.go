package ThreadPool

import (
	"errors"
	"fmt"
	"sync"

	Task "github.com/skkim-01/task"
)

var poolMutex sync.Mutex

// Close: stop thread pool.
func (p *Pool) Close(bForce bool) {
	if false == bForce {
		p._closeWait()
	} else {
		// force close will be flush task queue
		p._closeForce()
	}

}

// AddTask
func (p *Pool) AddTask(task *Task.Task) error {
	if false == p._isRunning {
		return errors.New("Thread Pool is closed")
	}

	// Task queue checker: performance?
	//---- Is it necessary to thread-pool's task watch threads?
	poolMutex.Lock()
	queueRate := int(float64(len(p.TaskQueue)) / float64(p._maxTaskQueueSize) * 100)
	if queueRate > 80 {
		p.IncreaseTaskQueueSize()
		fmt.Println("p._max", p._maxTaskQueueSize)
	}
	poolMutex.Unlock()

	p.TaskQueue <- task
	return nil
}

// IncreaseTaskQueueSize: Increase task queue size in runtime: size*2
func (p *Pool) IncreaseTaskQueueSize() {
	//poolMutex.Lock()
	tmpQueue := make(chan *Task.Task, p._maxTaskQueueSize*2)
	for len(p.TaskQueue) != 0 {
		t := <-p.TaskQueue
		tmpQueue <- t
	}

	p.TaskQueue = make(chan *Task.Task, p._maxTaskQueueSize*2)
	p.TaskQueue = tmpQueue

	for _, v := range p.MapWorkerThreads {
		v.TaskQueue = p.TaskQueue
	}
	p._maxTaskQueueSize = p._maxTaskQueueSize * 2
	//poolMutex.Unlock()
}

// for debug
func (p *Pool) TaskCount() {
	fmt.Println(len(p.TaskQueue), cap(p.TaskQueue))
}

// for debug
func (p *Pool) Stat() {
	fmt.Println("Thread Pool Stat------------------------")
	fmt.Println("MaxTaskQueueSize:", p._maxTaskQueueSize)
	fmt.Println("_currentTaskSize:", len(p.TaskQueue))
	fmt.Println("_activeWorkerSize:", p._activeWorkerSize)
	fmt.Println("_maxWorkerSize:", p._maxWorkerSize)
	fmt.Println("----------------------------------------")
}

// for debug
func (p *Pool) RetriveWorkers() {
	count := 0
	fmt.Println("Retrive Workers-------------------------")
	for k, v := range p.MapWorkerThreads {
		fmt.Println("count:", count, " id:", k, " state:", v.IsTasking)
		count++
	}
	fmt.Println("----------------------------------------")
}
