package ThreadPool

import (
	"sync/atomic"
	"time"
)

// TODO ?
// 1. Increase Worker: [condition: ActivWokerSize / MaxWorkerSize > 80%]
// 2. Decrease Worker: [condition: ActivWokerSize / MaxWorkerSize < 40%]

// _initWorkers: Starting workers at pool start: cpu*4/2
func (p *Pool) _initWorkers() {
	var nActiveWorker int = int(float64(p._maxWorkerSize) / 2)
	p._increaseWorker(nActiveWorker)
}

// _closeForce: Kill thread-pool
func (p *Pool) _closeForce() {
	// Stop add task
	p._isRunning = false
	defer close(p.TaskQueue)

	// Flush task channel
	for {
		if 0 >= len(p.TaskQueue) {
			break
		}
		_ = <-p.TaskQueue
	}

	// Stop workers
	poolMutex.Lock()
	for _, v := range p.MapWorkerThreads {
		v.Close()
	}
	p.MapWorkerThreads = nil
	p._maxTaskQueueSize = 0
	p._activeWorkerSize = 0
	poolMutex.Unlock()
}

// _closeForce: Stop thread-pool when task will be done
func (p *Pool) _closeWait() {
	// Stop add task
	p._isRunning = false
	defer close(p.TaskQueue)

	// Wait for flush queue
	for 0 >= len(p.TaskQueue) {
		time.Sleep(time.Second * 1)
	}

	// Stop workers
	poolMutex.Lock()
	for _, v := range p.MapWorkerThreads {
		v.Close()
	}
	p.MapWorkerThreads = nil
	p._maxTaskQueueSize = 0
	p._activeWorkerSize = 0
	poolMutex.Unlock()
}

// _increaseWorker: Create worker at runtime (currently not use)
func (p *Pool) _increaseWorker(inc int) {
	for i := 0; i < inc; i++ {
		w := NewWorker(p.TaskQueue)
		atomic.AddInt32(&p._activeWorkerSize, 1)
		p.MapWorkerThreads[w.WorkerID] = w
	}
}

// _decreaseWorker: Kill worker at runtime (currently not use)
func (p *Pool) _decreaseWorker(dec int) {
	nCurrentCnt := 0

	for _, v := range p.MapWorkerThreads {
		isRunnsing := atomic.LoadInt32(&v.IsTasking)
		if 0 == isRunnsing {
			v.Close()
			atomic.AddInt32(&p._activeWorkerSize, -1)
			nCurrentCnt = nCurrentCnt + 1
		}
	}
}
