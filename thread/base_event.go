package thread

import (
	"github.com/toophy/Gogame/event"
)

// 事件 : 线程关闭
type Event_close_thread struct {
	event.EventNormal
	Master IThread
}

// 事件执行
func (this *Event_close_thread) Exec() bool {
	if this.Master != nil {
		this.Master.pre_close_thread()
		return true
	}

	println("没找到线程")
	return true
}
