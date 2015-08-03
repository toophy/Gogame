package help

import (
	"time"
)

// 线程任务接口
type IEvent interface {
	Start() time.Duration
	SetStart(time.Duration)
	Interval() time.Duration
	Iterate() int
	Name() string
	Exec() error
	Cancel() error
	IsHeader() bool
}

type Event struct {
	time_pre   *Event        // 定时器使用 : 前一个事件
	time_next  *Event        // 定时器使用 : 后一个事件
	obj_pre    *Event        // 关联对象使用 : 前一个事件
	obj_next   *Event        // 关联对象使用 : 后一个事件
	start_time time.Duration // 开启定时器
	iterate    in32          // 迭代次数(0:表示不迭代,>0:迭代并每次+1)
	interval   time.Duration // 循环间隔时间
	name       string        // 事件别名, 特殊时刻使用, 一般不使用
	header     bool          // 事件头
}

type EvnetPool struct {
	Timer      map[time.Duration]IEvent // 这里的关键字是心跳时间, 也就是一个区间的开始时间,
	Names      map[string]IEvent        // 别名事件
	Start_time time.Duration            // 事件系统开启时间
	HeartTime  time.Duration            // 心跳时间
}

// 初始化事件
func (this *EvnetPool) Init(heart time.Duration) {
	this.Timer = make(map[time.Duration]IEvent)
	this.Names = make(map[string]IEvent)
	this.Start_time = time.Now()
	this.HeartTime = heart
}

// 投递事件
func (this *EvnetPool) Event_push(task IEvent) {
}

// 删除任务
func (this *EvnetPool) Event_remove(id interface{}) {
	delete(this.tasks, id)
}

// 取消任务
func (this *EvnetPool) Event_cancel(name string) (err error) {
	err = this.tasks[id].Cancel()
	delete(this.tasks, id)
	return
}

// 执行一个任务
func (this *EvnetPool) event_exec(e IEvent) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	err = e.Exec()
	if e.Iterate() == 0 {
		delete(this.tasks, e.Id())
	} else {
		e.SetStart(e.Start() + e.Interval())
		this.Event_push(e)
	}
	return
}
