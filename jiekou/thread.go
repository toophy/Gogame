package jiekou

import (
	lua "github.com/toophy/gopher-lua"
	"time"
)

// 线程任务接口
type ITask interface {
	// Start time
	Start() time.Duration
	SetStart(time.Duration)
	// interval of test executing, effective with Iterate() returns none-zero
	Interval() time.Duration
	// repeating times, 0 means don't repeat
	Iterate() int
	Id() interface{}
	Exec() error
	Cancel() error
}

// 场景线程接口
type IScreenThread interface {
	Get_thread_id() int32                   // 获取线程ID
	Get_thread_name() string                // 获取线程名称
	Task_push(task ITask)                   // 任务推送
	Task_remove(id interface{})             // 任务删除
	Task_cancel(id interface{}) (err error) // 任务取消
	GetLuaState() *lua.LState               // !!!只能获取, 不准许保存指针, 获取LState
}
