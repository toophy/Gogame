package event

import (
	"github.com/toophy/Gogame/help"
)

type EventObj struct {
	event_header help.ListNode // 所有关联事件,便于卸载
}

// 必须调用
func (this *EventObj) InitEventHeader() {
	this.event_header.Init("", 100)
}

// 压入定时器事件
func (this *EventObj) PostEvent(e *help.ListNode) bool {
	if e != nil {
		pre := this.event_header.Pre

		this.event_header.Pre = e
		e.Next = this.event_header
		e.Pre = pre
		pre.Next = e

		return true
	}
	return false
}

// 获取事件列表头
func (this *EventObj) GetEventHeader() *help.ListNode {
	return &this.event_header
}
