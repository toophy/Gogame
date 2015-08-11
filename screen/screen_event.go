package screen

import (
	"fmt"
	"github.com/toophy/Gogame/event"
)

// 事件 : 场景心跳
type Event_heart_beat struct {
	event.EventNormal
	Screen_ *Screen
}

// 事件执行
func (this *Event_heart_beat) Exec() bool {

	this.Screen_.Tolua_heart_beat()

	evt := &Event_heart_beat{Screen_: this.Screen_}
	evt.Init("", 3000)
	this.Screen_.PostEvent(evt)
	this.Screen_.Get_thread().PostEvent(evt)

	return true
}

// 事件 : 线程问好
type Event_thread_hello struct {
	event.EventNormal
	SrcThread int32
	DstThread int32
	Chat      string
}

// 事件执行
func (this *Event_thread_hello) Exec() bool {

	fmt.Printf("%d %s", this.SrcThread, this.Chat)

	evt := &Event_thread_hello{SrcThread: this.DstThread, DstThread: this.SrcThread, Chat: "wo 看看你"}
	evt.Init("", 3000)
	//this.

	return true
}
