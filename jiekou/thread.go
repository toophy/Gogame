package jiekou

import (
	"github.com/toophy/Gogame/event"
	lua "github.com/toophy/gopher-lua"
)

const (
	Tid_master = iota
	Tid_screen_1
	Tid_screen_2
	Tid_screen_3
	Tid_screen_4
	Tid_screen_5
	Tid_screen_6
	Tid_screen_7
	Tid_screen_8
	Tid_screen_9
	Tid_net_1
	Tid_net_2
	Tid_net_3
	Tid_db_1
	Tid_db_2
	Tid_db_3
	Tid_log
	Tid_last
)

const (
	Evt_gap_time  = 16     // 心跳时间(毫秒)
	Evt_gap_bit   = 4      // 心跳时间对应得移位(快速运算使用)
	Evt_lay1_time = 160000 // 第一层事件池最大支持时间(毫秒)
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
	LogDebug(f string, v ...interface{})          // 线程日志 : 调试[D]级别日志
	LogInfo(f string, v ...interface{})           // 线程日志 : 信息[I]级别日志
	LogWarn(f string, v ...interface{})           // 线程日志 : 警告[W]级别日志
	LogError(f string, v ...interface{})          // 线程日志 : 错误[E]级别日志
	LogFatal(f string, v ...interface{})          // 线程日志 : 致命[F]级别日志
}
