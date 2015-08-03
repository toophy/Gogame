package screen

import (
	"github.com/toophy/Gogame/help"
	//"time"
	"fmt"
)

// 事件 : 场景心跳
type Event_heart_beat struct {
	help.Task
	Screen_ *Screen
}

// 事件执行
func (t *Event_heart_beat) Exec() error {

	t.Screen_.Tolua_heart_beat()
	fmt.Println("heart")

	// n := time.Duration(time.Now().UnixNano())
	// t.Screen_.thread.Task_push(&Event_heart_beat{
	// 	Task: help.Task{
	// 		Id_:       t.Id_ + 1,
	// 		Start_:    n + 1*time.Second,
	// 		Interval_: time.Second,
	// 		Iterate_:  1*time.Second,
	// 	},
	// 	Screen_: t.Screen_,
	// })
	return nil
}
