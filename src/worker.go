package ThreadPool

import (
	"crypto/rand"
	"fmt"
	"sync/atomic"
)

// Close: Stop thread pool worker
func (w *Worker) Close() {
	// TaskChannel & TaskEventChannel will be closed in thread pool object
	w.StopChannel <- true
	<-w.StopChannel
	defer close(w.StopChannel)
	w.WorkerID = ""
	atomic.StoreInt32(&w.IsTasking, 0)
}

// WorkRoutine: Loop thread for running task
func (w *Worker) WorkRoutine() {
_SELECT_SECTION:
	select {
	case <-w.StopChannel:
		w.StopChannel <- true
		goto _EXIT
	case task := <-w.TaskQueue:
		atomic.AddInt32(&w.IsTasking, 1)
		task.InvokeTask()
		atomic.AddInt32(&w.IsTasking, -1)
		goto _SELECT_SECTION
	}
_EXIT:
}

// IsRunning : Get worker's running state
func (w *Worker) IsRunning() int32 {
	return atomic.LoadInt32(&w.IsTasking)
}

func _generateUuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
