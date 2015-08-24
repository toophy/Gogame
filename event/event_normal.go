package event

import (
	"github.com/toophy/Gogame/help"
)

// 事件接口
type IEvent interface {
	Init(name string, t uint64) // 初始化(name可以为空, t是触发时间)
	Exec(home interface{}) bool // 执行
	GetName() string            // 获取别名
	//
	GetNodeTimer() *help.ListNode // 获取定时器节点
	GetNodeObj() *help.ListNode   // 获取对象节点
	// 定时器链表
	GetPreTimer() *help.ListNode   // 获取前一个定时器事件
	GetNextTimer() *help.ListNode  // 获取下一个定时器事件
	SetPreTimer(e *help.ListNode)  // 设置前一个定时器事件
	SetNextTimer(e *help.ListNode) // 设置下一个定时器事件
	// 对象链表
	GetPreObj() *help.ListNode   // 获取前一个对象事件
	GetNextObj() *help.ListNode  // 获取下一个对象事件
	SetPreObj(e *help.ListNode)  // 设置前一个对象事件
	SetNextObj(e *help.ListNode) // 设置下一个对象事件
	// 触发时间
	GetTouchTime() uint64            // 获取定时器触发时间戳
	SetTouchTime(t uint64)           // 设置定时器时间戳
	SetDelayTime(d uint64, c uint64) // 设置定时器相对时间, c是当前时间戳
	// 打印自己
	PrintSelf() // 打印自己
}

// 普通事件
type EventNormal struct {
	name       string        // 名称
	node_timer help.ListNode // 定时器节点
	node_obj   help.ListNode // 对象节点
	touch_time uint64        // 定时器触发时间戳
}

func (this *EventNormal) Init(name string, t uint64) {
	this.name = name
	this.touch_time = t
}

func (this *EventNormal) GetNodeTimer() *help.ListNode {
	return &this.node_timer
}

func (this *EventNormal) GetNodeObj() *help.ListNode {
	return &this.node_obj
}

func (this *EventNormal) Exec(home interface{}) bool {
	println("Normal Exec")
	return true
}

func (this *EventNormal) GetName() string {
	return this.name
}

func (this *EventNormal) GetPreTimer() *help.ListNode {
	return this.node_timer.Pre
}

func (this *EventNormal) GetNextTimer() *help.ListNode {
	return this.node_timer.Next
}

func (this *EventNormal) SetPreTimer(e *help.ListNode) {
	this.node_timer.Pre = e
}

func (this *EventNormal) SetNextTimer(e *help.ListNode) {
	this.node_timer.Next = e
}

func (this *EventNormal) GetPreObj() *help.ListNode {
	return this.node_obj.Pre
}

func (this *EventNormal) GetNextObj() *help.ListNode {
	return this.node_obj.Next
}

func (this *EventNormal) SetPreObj(e *help.ListNode) {
	this.node_obj.Pre = e
}

func (this *EventNormal) SetNextObj(e *help.ListNode) {
	this.node_obj.Next = e
}

func (this *EventNormal) GetTouchTime() uint64 {
	return this.touch_time
}

func (this *EventNormal) SetTouchTime(t uint64) {
	this.touch_time = t
}

func (this *EventNormal) SetDelayTime(d uint64, c uint64) {
	this.touch_time = c + d
}

func (this *EventNormal) PrintSelf() {
	println("  {E} Is normal")
}
