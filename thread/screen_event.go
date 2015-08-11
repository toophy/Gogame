package thread

import (
	"github.com/toophy/Gogame/event"
)

// 事件 : 场景增/删
type Event_open_screen struct {
	event.EventNormal
	Screen_oid_    int32
	Screen_name_   string
	Screen_thread_ *ScreenThread
	Open           bool
}

// 事件执行
func (this *Event_open_screen) Exec(home interface{}) bool {
	if this.Open {
		if this.Screen_thread_.Add_screen(this.Screen_name_, this.Screen_oid_) {
			println("打开场景成功")
			return true
		}
		println("打开场景失败")
		return true
	}

	if this.Screen_thread_.Del_screen(this.Screen_oid_) {
		println("关闭场景成功")
		return true
	}
	println("关闭场景失败")
	return true
}
