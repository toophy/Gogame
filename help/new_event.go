package help

import (
	"errors"
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
	Remove() error
	Cancel() error
	IsHeader() bool
	LeaveTimer() IEvent // 返回下一个事件
	LeaveObj() IEvent   // 返回下一个事件
	Time_push(header IEvent)
	Obj_push(header IEvent)
	get_time_pre() IEvent
	get_time_next() IEvent
	get_obj_pre() IEvent
	get_obj_next() IEvent
	get_name_pre() IEvent
	get_name_next() IEvent
	set_time_pre(e IEvent)
	set_time_next(e IEvent)
	set_obj_pre(e IEvent)
	set_obj_next(e IEvent)
	set_name_pre(e IEvent)
	set_name_next(e IEvent)
}

type Event struct {
	time_pre   IEvent        // 定时器使用 : 前一个事件
	time_next  IEvent        // 定时器使用 : 后一个事件
	obj_pre    IEvent        // 关联对象使用 : 前一个事件
	obj_next   IEvent        // 关联对象使用 : 后一个事件
	name_pre   IEvent        // 别名使用 : 前一个事件
	name_next  IEvent        // 别名使用 : 后一个事件
	start_time time.Duration // 开启定时器
	iterate    in32          // 迭代次数(0:表示不迭代,>0:迭代并每次+1)
	interval   time.Duration // 循环间隔时间
	name       string        // 事件别名, 特殊时刻使用, 一般不使用
	header     bool          // 事件头
}

func (this *Event) Name() string {
	return this.name
}

func (this *Event) Remove() error {
	if this.IsHeader() {
		return errors.New("[W] 不能Remove事件头")
	}

	if len(this.name) > 0 {
		this.LeaveName()
	}

	this.LeaveTimer()
	this.LeaveObj()
}

func (this *Event) Cancel() error {
	if this.IsHeader() {
		return errors.New("[W] 不能Cancel事件头")
	}

	if len(this.name) > 0 {
		this.LeaveName()
	}

	this.LeaveTimer()
	this.LeaveObj()
}

func (this *Event) IsHeader() bool {
	return this.header
}

func (this *Event) LeaveTimer() IEvent {
	if this.IsHeader() {
		return nil
	}

	this.get_time_pre().set_time_next(this.get_time_next())
	this.get_time_next().set_time_pre(this.get_time_pre())
	this.set_time_next(this)
	this.set_time_pre(this)
}

func (this *Event) LeaveObj() IEvent {
	if this.IsHeader() {
		return nil
	}

	this.get_obj_pre().set_obj_next(this.get_obj_next())
	this.get_obj_next().set_obj_pre(this.get_obj_pre())
	this.set_obj_next(this)
	this.set_obj_pre(this)
}

func (this *Event) Timer_push(header IEvent) {
	if this.IsHeader() {
		return
	}

	old_pre := header.get_time_pre()
	header.set_time_pre(this)
	this.set_time_next(header)
	this.set_time_pre(old_pre)
	old_pre.set_time_next(this)
}

func (this *Event) Obj_push(header IEvent) {
	if this.IsHeader() {
		return
	}

	old_pre := header.get_obj_pre()
	header.set_obj_pre(this)
	this.set_obj_next(header)
	this.set_obj_pre(old_pre)
	old_pre.set_obj_next(this)
}

func (this *Event) set_time_pre(e IEvent) {
	this.time_pre = e
}

func (this *Event) set_time_next(e IEvent) {
	this.time_next = e
}

func (this *Event) set_obj_pre(e IEvent) {
	this.obj_pre = e
}

func (this *Event) set_obj_next(e IEvent) {
	this.obj_next = e
}

func (this *Event) get_time_pre() IEvent {
	return this.time_pre
}

func (this *Event) get_time_next() IEvent {
	return this.time_next
}

func (this *Event) get_obj_pre() IEvent {
	return this.obj_pre
}

func (this *Event) get_obj_next() IEvent {
	return this.obj_next
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
func (this *EvnetPool) Event_push(task IEvent) false {
	// 有同名事件存在
	if len(task.Name()) > 0 {
		if _, ok := this.Names[task.Name()]; ok {
			return false
		}
		// 链接到Names

	}

	//task.Start() -> 确定 是那个 心跳开始时间
	key := this.Make_timer(task.Start())
	if v, ok := this.Timer[key]; !ok {
		header := new(Event)
		header.header = true
		this.Timer[key] = header
	}
	// 链接到Timer
}

// 生成定时器关键字
func (this *EvnetPool) Make_timer(t time.Duration) time.Duration {
	return t
}

// 删除任务
func (this *EvnetPool) Event_remove(name string) error {
	if _, ok := this.Names[name]; ok {
		return this.Names[name].Remove()
	}
	return errors.New("[W] 没有找到事件 : " + name)
}

// 取消任务
func (this *EvnetPool) Event_cancel(name string) (err error) {
	if _, ok := this.Names[name]; ok {
		return this.Names[name].Cancel()
	}
	return errors.New("[W] 没有找到事件 : " + name)
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
		e.Remove()
	} else {
		// 断开时间链, 重新投递一次
		e.SetStart(e.Start() + e.Interval())
		this.Event_push(e)
	}
	return
}
