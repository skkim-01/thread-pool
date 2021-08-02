package main

import (
	"fmt"
	"time"

	sig "github.com/skkim-01/signal-waiter"
	Task "github.com/skkim-01/task"

	ThreadPool "github.com/skkim-01/thread-pool/src"
)

// thread pool test
func main() {
	fmt.Println("[log] start")

	go poolThread()

	sig.Wait()
}

func _inv(p *ThreadPool.Pool) {
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			t1 := NewLocalTask(i)
			p.AddTask(t1)
		} else {
			t2 := NewLocalTask2(i)
			p.AddTask(t2)
		}
	}
	fmt.Println("[Log] Done: Insert tasks")
}

func poolThread() {
	p := ThreadPool.NewPool()

	p.Stat()

	//p.RetriveWorkers()

	go _inv(p)

	time.Sleep(time.Second * 5)

	bForce := true
	fmt.Println("[dbg] p.close start:", bForce)
	p.Close(bForce)
	fmt.Println("[dbg] p.close done:", bForce)

	// close
	sig.Close()
}

/////////////////////////////////////////////
// Sample task 1: running in thread pool
/////////////////////////////////////////////

type SampleTask_1 struct{}

func (t *SampleTask_1) Invoke(count int) (int, string, error) {
	time.Sleep(time.Second * 2)
	fmt.Println("[*] SampleTask_1.Invoke(): count:", count)
	return count, "Invoke", nil
}

var SampleTask_1_CALLBACK Task.TaskCallback = func(cbResult []interface{}) {
	// Print callback result
	fmt.Print("### LOG ### SampleTask_1_CALLBACK.Result:")
	fmt.Println(cbResult...)
}

func NewLocalTask(i int) *Task.Task {
	t := Task.NewTask(&SampleTask_1{}, "Invoke", SampleTask_1_CALLBACK, i)
	return t
}

/////////////////////////////////////////////
// Sample task 2: running in thread pool
/////////////////////////////////////////////

type SampleTask_2 struct{}

var SampleTask_2_CALLBACK Task.TaskCallback = func(cbResult []interface{}) {
	// Print callback result
	fmt.Print("### LOG ### SampleTask_2_CALLBACK.Result:")
	fmt.Println(cbResult...)
}

func (t *SampleTask_2) Execute(count int) (int, string, error) {
	time.Sleep(time.Second * 2)
	fmt.Println("[*] SampleTask_2.Execute(): count:", count)
	return count, "Execute", nil
}

func NewLocalTask2(i int) *Task.Task {
	t := Task.NewTask(&SampleTask_2{}, "Execute", SampleTask_2_CALLBACK, i)
	return t
}
