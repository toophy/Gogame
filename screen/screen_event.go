package screen

import (
	"fmt"
	"github.com/toophy/Gogame/event"
	"github.com/toophy/Gogame/jiekou"
)

// 事件 : 场景心跳
type Event_heart_beat struct {
	event.EventNormal
	Screen_ *Screen
}

// 事件执行
func (this *Event_heart_beat) Exec(home interface{}) bool {

	this.Screen_.Tolua_heart_beat()

	//
	evt_hello := &Event_thread_hello{SrcThread: this.Screen_.Get_thread().Get_thread_id(), Chat: "wo 看看你", Replay: false}
	if this.Screen_.Get_thread().Get_thread_id() == 1 {
		evt_hello.DstThread = 2
	} else if this.Screen_.Get_thread().Get_thread_id() == 2 {
		evt_hello.DstThread = 1
	}
	evt_hello.Init("", 3000)
	home.(jiekou.IScreenThread).PostThreadMsg(evt_hello.DstThread, evt_hello)

	//
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
	Replay    bool
}

// 事件执行
func (this *Event_thread_hello) Exec(home interface{}) bool {

	fmt.Printf("%d %s", this.SrcThread, this.Chat)

	if !this.Replay {
		println("need replay")
		evt := &Event_thread_hello{SrcThread: this.DstThread, DstThread: this.SrcThread, Chat: "wo 也来看看你", Replay: true}
		evt.Init("", 3000)
		home.(jiekou.IScreenThread).PostThreadMsg(evt.DstThread, evt)
	}

	return true
}
