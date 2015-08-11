package jiekou

import (
	"github.com/toophy/Gogame/event"
	lua "github.com/toophy/gopher-lua"
)

// 场景线程接口
type IScreenThread interface {
	Get_thread_id() int32                         // 获取线程ID
	Get_thread_name() string                      // 获取线程名称
	GetLuaState() *lua.LState                     // !!!只能获取, 不准许保存指针, 获取LState
	PostEvent(a event.IEvent) bool                // 投递定时器事件
	PostThreadMsg(tid int32, a event.IEvent) bool // 投递线程间消息
	GetEvent(name string) event.IEvent            // 通过别名获取事件
	RemoveEvent(e event.IEvent) bool              // 删除事件, 只能操作本线程事件
	RemoveEventList(header event.IEvent)          // 删除一整个事件列表
	PopTimer(e event.IEvent)                      // 从线程事件中弹出指定事件, 只能操作本线程事件
	PopObj(e event.IEvent)                        // 从关联对象中弹出指定事件, 只能操作本线程事件
}
