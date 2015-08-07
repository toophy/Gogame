package event

// import (
// 	"errors"
// )

// // 线程任务接口
// type IEvent interface {
// 	GetStart() int64
// 	SetStart(int64)
// 	Interval() int64
// 	Iterate() int
// 	Name() string
// 	Exec() error
// 	Remove() error
// 	Cancel() error
// 	IsHeader() bool
// 	SetHeader()
// 	LeaveTimer() IEvent // 返回下一个事件
// 	LeaveObj() IEvent   // 返回下一个事件
// 	Time_push(header IEvent)
// 	Obj_push(header IEvent)
// 	get_time_pre() IEvent
// 	get_time_next() IEvent
// 	get_obj_pre() IEvent
// 	get_obj_next() IEvent
// 	get_name_pre() IEvent
// 	get_name_next() IEvent
// 	set_time_pre(e IEvent)
// 	set_time_next(e IEvent)
// 	set_obj_pre(e IEvent)
// 	set_obj_next(e IEvent)
// 	set_name_pre(e IEvent)
// 	set_name_next(e IEvent)
// }

// type Event struct {
// 	time_pre   IEvent // 定时器使用 : 前一个事件
// 	time_next  IEvent // 定时器使用 : 后一个事件
// 	obj_pre    IEvent // 关联对象使用 : 前一个事件
// 	obj_next   IEvent // 关联对象使用 : 后一个事件
// 	name       string // 事件别名, 特殊时刻使用, 一般不使用
// 	start_time int64  // 开启定时器
// 	interval   int64  // 循环间隔时间
// 	iterate    int32  // 迭代次数(0:表示不迭代,>0:迭代并每次+1)
// 	header     bool   // 事件头
// }

// func (this *Event) Name() string {
// 	return this.name
// }

// func (this *Event) Remove() error {
// 	if this.IsHeader() {
// 		return errors.New("[W] 不能Remove事件头")
// 	}

// 	if len(this.name) > 0 {
// 		this.LeaveName()
// 	}

// 	this.LeaveTimer()
// 	this.LeaveObj()
// }

// func (this *Event) Cancel() error {
// 	if this.IsHeader() {
// 		return errors.New("[W] 不能Cancel事件头")
// 	}

// 	if len(this.name) > 0 {
// 		this.LeaveName()
// 	}

// 	this.LeaveTimer()
// 	this.LeaveObj()
// }

// func (this *Event) IsHeader() bool {
// 	return this.header
// }

// func (this *Event) SetHeader() {
// 	this.header = true
// }

// func (this *Event) LeaveTimer() IEvent {
// 	if this.IsHeader() {
// 		return nil
// 	}

// 	this.get_time_pre().set_time_next(this.get_time_next())
// 	this.get_time_next().set_time_pre(this.get_time_pre())
// 	this.set_time_next(this)
// 	this.set_time_pre(this)
// }

// func (this *Event) LeaveObj() IEvent {
// 	if this.IsHeader() {
// 		return nil
// 	}

// 	this.get_obj_pre().set_obj_next(this.get_obj_next())
// 	this.get_obj_next().set_obj_pre(this.get_obj_pre())
// 	this.set_obj_next(this)
// 	this.set_obj_pre(this)
// }

// func (this *Event) Timer_push(header IEvent) {
// 	if this.IsHeader() {
// 		return
// 	}

// 	old_pre := header.get_time_pre()
// 	header.set_time_pre(this)
// 	this.set_time_next(header)
// 	this.set_time_pre(old_pre)
// 	old_pre.set_time_next(this)
// }

// func (this *Event) Obj_push(header IEvent) {
// 	if this.IsHeader() {
// 		return
// 	}

// 	old_pre := header.get_obj_pre()
// 	header.set_obj_pre(this)
// 	this.set_obj_next(header)
// 	this.set_obj_pre(old_pre)
// 	old_pre.set_obj_next(this)
// }

// func (this *Event) set_time_pre(e IEvent) {
// 	this.time_pre = e
// }

// func (this *Event) set_time_next(e IEvent) {
// 	this.time_next = e
// }

// func (this *Event) set_obj_pre(e IEvent) {
// 	this.obj_pre = e
// }

// func (this *Event) set_obj_next(e IEvent) {
// 	this.obj_next = e
// }

// func (this *Event) get_time_pre() IEvent {
// 	return this.time_pre
// }

// func (this *Event) get_time_next() IEvent {
// 	return this.time_next
// }

// func (this *Event) get_obj_pre() IEvent {
// 	return this.obj_pre
// }

// func (this *Event) get_obj_next() IEvent {
// 	return this.obj_next
// }
