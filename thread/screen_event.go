package thread

import (
	"errors"
	//"fmt"
)

type Event_open_screen struct {
	Task
	Screen_oid_    int32
	Screen_name_   string
	Screen_thread_ *Screen
	Open           bool
}

func (t *Event_open_screen) Exec() error {
	if t.Open {
		if t.Screen_thread_.Add_screen(t.Screen_name_, t.Screen_oid_) {
			return nil
		}
		return errors.New("打开场景失败")
	}

	if t.Screen_thread_.Del_screen(t.Screen_oid_) {
		return nil
	}
	return errors.New("关闭场景失败")
}
