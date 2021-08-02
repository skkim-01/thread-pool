package ThreadPool

import Task "github.com/skkim-01/task"

// Worker structure definition
type Worker struct {
	WorkerID    string          // Worker Unique ID
	TaskQueue   chan *Task.Task // Task channel: Queue
	IsTasking   int32           // Tasking flag: atomic int
	StopChannel chan bool       // Thread stop channel
}

// NewWorker : Create new thread pool worker
func NewWorker(taskQueue chan *Task.Task) *Worker {
	w := Worker{
		WorkerID:    _generateUuid(),
		TaskQueue:   taskQueue,
		StopChannel: make(chan bool),
		IsTasking:   0,
	}
	go w.WorkRoutine()
	return &w
}
