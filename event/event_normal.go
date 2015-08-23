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
	GetNodeTimer() *help.ListNode  // 获取定时器节点
	GetNodeObj() *help.ListNode    // 获取对象节点
	SetNodeTimer(e *help.ListNode) // 设置定时器节点
	SetNodeObj(e *help.ListNode)   // 设置对象节点
	// 定时器链表
	GetPreTimer() IEvent   // 获取前一个定时器事件
	GetNextTimer() IEvent  // 获取下一个定时器事件
	SetPreTimer(e IEvent)  // 设置前一个定时器事件
	SetNextTimer(e IEvent) // 设置下一个定时器事件
	// 对象链表
	GetPreObj() IEvent   // 获取前一个对象事件
	GetNextObj() IEvent  // 获取下一个对象事件
	SetPreObj(e IEvent)  // 设置前一个对象事件
	SetNextObj(e IEvent) // 设置下一个对象事件
	// 触发时间
	GetTouchTime() uint64            // 获取定时器触发时间戳
	SetTouchTime(t uint64)           // 设置定时器时间戳
	SetDelayTime(d uint64, c uint64) // 设置定时器相对时间, c是当前时间戳
	// 打印自己
	PrintSelf() // 打印自己
	PrintList() // 打印列表(只对EventHeader有用)
}

// 普通事件
type EventNormal struct {
	name       string         // 名称
	node_timer *help.ListNode // 定时器节点
	node_obj   *help.ListNode // 对象节点
	touch_time uint64         // 定时器触发时间戳
}

func (this *EventNormal) Init(name string, t uint64) {
	this.name = name
	this.touch_time = t
}

func (this *EventNormal) GetNodeTimer() *help.ListNode {
	return this.node_timer
}

func (this *EventNormal) GetNodeObj() *help.ListNode {
	return this.node_obj
}

func (this *EventNormal) SetNodeTimer(e *help.ListNode) {
	this.node_timer = e
}

func (this *EventNormal) SetNodeObj(e *help.ListNode) {
	this.node_obj = e
}

func (this *EventNormal) Exec(home interface{}) bool {
	println("Normal Exec")
	return true
}

func (this *EventNormal) GetName() string {
	return this.name
}

func (this *EventNormal) GetPreTimer() *help.ListNode {
	if this.node_timer {
		return this.node_timer.Pre
	}
	return nil
}

func (this *EventNormal) GetNextTimer() *help.ListNode {
	if this.node_timer {
		return this.node_timer.Next
	}
	return nil
}

func (this *EventNormal) SetPreTimer(e *help.ListNode) {
	if this.node_timer {
		this.node_timer.Pre = e
	}
}

func (this *EventNormal) SetNextTimer(e *help.ListNode) {
	if this.node_timer {
		this.node_timer.Next = e
	}
}

func (this *EventNormal) GetPreObj() *help.ListNode {
	if this.node_obj {
		return this.node_obj.Pre
	}
	return nil
}

func (this *EventNormal) GetNextObj() *help.ListNode {
	if this.node_obj {
		return this.node_obj.Next
	}
}

func (this *EventNormal) SetPreObj(e *help.ListNode) {
	if this.node_obj {
		this.node_obj.Pre = e
	}
}

func (this *EventNormal) SetNextObj(e *help.ListNode) {
	if this.node_obj {
		this.node_obj.Next = e
	}
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
