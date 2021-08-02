package ThreadPool

import (
	"runtime"

	Task "github.com/skkim-01/task"
)

type Pool struct {
	TaskQueue         chan *Task.Task    // Thread Pool Task Queue
	MapWorkerThreads  map[string]*Worker // workers
	_maxTaskQueueSize int32              // Task channel Size
	_activeWorkerSize int32              // Active Worker Size
	_maxWorkerSize    int32              // CPU Core * 4
	_isRunning        bool
}

// NewPool : Create new thread pool
func NewPool() *Pool {
	runtime.GOMAXPROCS(runtime.NumCPU())
	maxWorker := int32(runtime.NumCPU() * 4)

	// Create Pool Object
	p := Pool{
		TaskQueue:         make(chan *Task.Task, maxWorker*4),
		MapWorkerThreads:  make(map[string]*Worker),
		_maxTaskQueueSize: maxWorker * 4,
		_activeWorkerSize: 0,
		_maxWorkerSize:    maxWorker,
		_isRunning:        false,
	}

	// Worker start
	p._initWorkers()
	p._isRunning = true

	return &p
}
