package thread

import (
	"fmt"
	"testing"
	"time"
)

type _task struct {
	EventNormal
}

func (t *_task) Exec(home interface{}) bool {
	fmt.Printf("Task %d Executed.\n", t.Id())
	return true
}

type _taskx struct {
	EventNormal
	name string
}

func (t *_taskx) Exec(home interface{}) bool {
	fmt.Printf("Taskx(%s) %d Executed.\n", t.name, t.Id())
	return true
}

func Test(t *testing.T) {
	mythread := &Thread{}
	err := mythread.Init_thread(1, 100)
	if err != nil {
		panic(err.Error())
	}

	go mythread.Run_thread()

	t1 := &_task{}
	t1.Init("wowo")
	t1.SetTouchTime(100)
	mythread.PostEvent(t1)

	t2 := &_taskx{}
	t2.Init("wowo2")
	t2.SetTouchTime(200)
	mythread.PostEvent(t2)

	time.Sleep(20 * time.Second)
}
