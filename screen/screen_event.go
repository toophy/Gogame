package screen

import (
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
