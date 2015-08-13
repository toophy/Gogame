package thread

import (
	"bytes"
	"github.com/toophy/Gogame/event"
)

// 事件 : 线程投递的日志
type Event_thread_log struct {
	event.EventNormal
	Data bytes.Buffer
}

// 事件执行
func (this *Event_thread_log) Exec(home interface{}) bool {
	home.(*LogThread).Add_log(this.Data)
	return true
}

// 事件 : 线程投递的日志
type Event_flush_log struct {
	event.EventNormal
}

// 事件执行
func (this *Event_flush_log) Exec(home interface{}) bool {
	home.(*LogThread).Flush_log()

	evt := &Event_flush_log{}
	evt.Init("", 300)
	home.(*LogThread).PostEvent(evt)

	return true
}
