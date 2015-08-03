package help

import (
	"time"
)

// 线程任务
type Task struct {
	Start_    time.Duration
	Iterate_  int
	Interval_ time.Duration
	Id_       int
}

func (t *Task) Start() time.Duration {
	return t.Start_
}

func (t *Task) SetStart(tm time.Duration) {
	t.Start_ = tm
}

func (t *Task) Interval() time.Duration {
	return 0
}
func (t *Task) Iterate() int {
	return 0
}

func (t *Task) Id() interface{} {
	return t.Id_
}

func (t *Task) Exec() error {
	return nil
}

func (t *Task) Cancel() error {
	return nil
}
