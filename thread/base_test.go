package thread

import (
	"fmt"
	"testing"
	"time"
)

type _task struct {
	start    time.Duration
	iterate  int
	interval time.Duration
	id       int
}

func (t *_task) Start() time.Duration {
	return t.start
}

func (t *_task) SetStart(tm time.Duration) {
	t.start = tm
}

func (t *_task) Interval() time.Duration {
	return 0
}
func (t *_task) Iterate() int {
	return 0
}

func (t *_task) Id() interface{} {
	return t.id
}
func (t *_task) Exec() error {
	fmt.Printf("Task %d Executed.\n", t.Id())
	return nil
}
func (t *_task) Cancel() error {
	return nil
}

//////////////////

type _taskx struct {
	start    time.Duration
	iterate  int
	interval time.Duration
	id       int
	name     string
}

func (t *_taskx) Start() time.Duration {
	return t.start
}

func (t *_taskx) SetStart(tm time.Duration) {
	t.start = tm
}

func (t *_taskx) Interval() time.Duration {
	return 0
}
func (t *_taskx) Iterate() int {
	return 0
}

func (t *_taskx) Id() interface{} {
	return t.id
}
func (t *_taskx) Exec() error {
	fmt.Printf("Taskx(%s) %d Executed.\n", t.name, t.Id())
	return nil
}
func (t *_taskx) Cancel() error {
	return nil
}

func Test(t *testing.T) {
	mythread := &Thread{}
	mythread.Init_thread(1, 100)

	go mythread.Run_thread()

	n := time.Duration(time.Now().UnixNano())
	mythread.Task_push(&_task{
		id:       1,
		start:    n + time.Second,
		interval: time.Second,
		iterate:  0,
	})
	mythread.Task_push(&_task{
		id:       2,
		start:    n + time.Second,
		interval: time.Second,
		iterate:  0,
	})
	mythread.Task_push(&_task{
		id:       3,
		start:    n + 3*time.Second,
		interval: time.Second,
		iterate:  0,
	})
	mythread.Task_push(&_taskx{
		id:       4,
		start:    n + 3*time.Second,
		interval: time.Second,
		iterate:  0,
		name:     "你妹啊",
	})

	time.Sleep(15 * time.Second)
}
