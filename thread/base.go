package thread

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/toophy/Gogame/event"
	"github.com/toophy/Gogame/jiekou"
	"strings"
	"time"
)

// 线程接口
type IThread interface {
	Init_thread(IThread, int32, string, int64, uint64) error // 初始化线程
	Run_thread()                                             // 运行线程
	Get_thread_id() int32                                    // 获取线程ID
	Get_thread_name() string                                 // 获取线程名称
	pre_close_thread()                                       // -- 只允许thread调用 : 预备关闭线程
	on_first_run()                                           // -- 只允许thread调用 : 首次运行(在 on_run 前面)
	on_run()                                                 // -- 只允许thread调用 : 线程运行部分
	on_end()                                                 // -- 只允许thread调用 : 线程结束回调
	PostEvent(a event.IEvent) bool                           // 投递定时器事件
	GetEvent(name string) event.IEvent                       // 通过别名获取事件
	RemoveEvent(e event.IEvent) bool                         // 删除事件, 只能操作本线程事件
	PopTimer(e event.IEvent)                                 // 从线程事件中弹出指定事件, 只能操作本线程事件
	PopObj(e event.IEvent)                                   // 从关联对象中弹出指定事件, 只能操作本线程事件
	LogDebug(f string, v ...interface{})                     // 线程日志 : 调试[D]级别日志
	LogInfo(f string, v ...interface{})                      // 线程日志 : 信息[I]级别日志
	LogWarn(f string, v ...interface{})                      // 线程日志 : 警告[W]级别日志
	LogError(f string, v ...interface{})                     // 线程日志 : 错误[E]级别日志
	LogFatal(f string, v ...interface{})                     // 线程日志 : 致命[F]级别日志
}

// 线程基本功能
type Thread struct {
	id               int32                         // Id号
	name             string                        // 线程名称
	heart_time       int64                         // 心跳时间(毫秒)
	start_time       int64                         // 线程开启时间戳
	last_time        int64                         // 最近一次线程运行时间戳
	heart_rate       float64                       // 本次心跳比率
	pre_stop         bool                          // 预备停止
	self             IThread                       // 自己, 初始化之后, 不要操作
	first_run        bool                          // 线程首次运行
	evt_lay1         []event.IEvent                // 第一层事件池
	evt_lay2         map[uint64]event.IEvent       // 第二层事件池
	evt_names        map[string]event.IEvent       // 别名
	evt_lay1Size     uint64                        // 第一层池容量
	evt_lay1Cursor   uint64                        // 第一层游标
	evt_lastRunCount uint64                        // 最近一次运行次数
	evt_currRunCount uint64                        // 当前运行次数
	evt_threadMsg    [jiekou.Tid_last]event.IEvent // 保存将要发给其他线程的事件(消息)
	evt_recvMsg      event.EventHeader             // 接收线程间消息
	log_Buffer       bytes.Buffer                  // 线程日志缓冲
	log_TimeString   string                        // 时间格式(精确到秒2015.08.13 16:33:00)
}

// 初始化线程(必须调用)
// usage : Init_thread(Tid_master, "主线程", 100)
func (this *Thread) Init_thread(self IThread, id int32, name string, heart_time int64, lay1_time uint64) error {
	if id < jiekou.Tid_master || id >= jiekou.Tid_last {
		return errors.New("[E] 线程ID超出范围 [Tid_master,Tid_last]")
	}
	if self == nil {
		return errors.New("[E] 线程自身指针不能为nil")
	}

	if lay1_time < jiekou.Evt_gap_time || lay1_time > jiekou.Evt_lay1_time {
		return errors.New("[E] 第一层支持16毫秒到160000毫秒")
	}

	if len(this.evt_names) > 0 {
		return errors.New("[E] EventHome 已经初始化过")
	}

	this.id = id
	this.name = name
	this.heart_time = heart_time * int64(time.Millisecond)
	this.start_time = time.Now().UnixNano()
	this.last_time = this.start_time
	this.heart_rate = 1.0
	this.self = self
	this.first_run = true

	// 初始化事件池
	this.evt_lay1Size = lay1_time >> jiekou.Evt_gap_bit
	this.evt_lay1Cursor = 0
	this.evt_currRunCount = 1
	this.evt_lastRunCount = this.evt_currRunCount

	this.evt_lay1 = make([]event.IEvent, this.evt_lay1Size)
	this.evt_lay2 = make(map[uint64]event.IEvent, 0)
	this.evt_names = make(map[string]event.IEvent, 0)

	for i := uint64(0); i < this.evt_lay1Size; i++ {
		this.evt_lay1[i] = new(event.EventHeader)
		this.evt_lay1[i].Init("", 100)
	}

	for i := 0; i < jiekou.Tid_last; i++ {
		this.evt_threadMsg[i] = new(event.EventHeader)
		this.evt_threadMsg[i].Init("", 100)
	}

	this.evt_recvMsg.Init("", 100)

	this.log_TimeString = time.Now().Format("2006-01-02 15:04:05")

	return nil
}

// 运行线程
func (this *Thread) Run_thread() {
	// 计算心跳误差值, 决定心跳滴答(小数), heart_time, last_time, heart_rate
	// 处理线程间接收消息, 分配到水表定时器
	// 执行水表定时器
	go func() {
		GetMaster().Add_run_thread(this.self)

		this.start_time = time.Now().UnixNano()
		this.last_time = this.start_time
		next_time := time.Duration(this.heart_time)
		run_time := int64(0)

		this.self.on_first_run()

		for {

			time.Sleep(next_time)

			this.log_TimeString = time.Now().Format("01-02 15:04:05")

			this.last_time = time.Now().UnixNano()
			this.runThreadMsg()
			this.runEvents()
			this.runOnce()
			this.self.on_run()

			this.sendThreadMsg()

			// 计算下一次运行的时间
			run_time = time.Now().UnixNano() - this.last_time
			if run_time >= this.heart_time {
				run_time = this.heart_time - 10*1000*1000
			} else if run_time < 0 {
				run_time = 0
			}

			next_time = time.Duration(this.heart_time - run_time)

			if this.pre_stop {
				// 是否有需要释放的对象?
				this.self.on_end()
				break
			}
		}

		GetMaster().Release_run_thread(this.self)
	}()
}

// 运行一次(核心流程)
func (this *Thread) runOnce() {
	// 计算心跳误差值, 决定心跳滴答(小数), heart_time, last_time, heart_rate
	// 处理线程间接收消息, 分配到水表定时器
	// 执行水表定时器
}

// 返回线程编号
func (this *Thread) Get_thread_id() int32 {
	return this.id
}

// 返回线程名称
func (this *Thread) Get_thread_name() string {
	return this.name
}

// 预备关闭线程
func (this *Thread) pre_close_thread() {
	this.pre_stop = true
}

// 投递定时器事件
func (this *Thread) PostEvent(a event.IEvent) bool {
	check_name := len(a.GetName()) > 0
	if check_name {
		if _, ok := this.evt_names[a.GetName()]; ok {
			return false
		}
	}

	if a.GetTouchTime() < 0 {
		return false
	}

	// 计算放在那一层
	pos := (a.GetTouchTime() + jiekou.Evt_gap_time - 1) >> jiekou.Evt_gap_bit
	if pos < 0 {
		pos = 1
	}

	var header event.IEvent

	if pos < this.evt_lay1Size {
		new_pos := this.evt_lay1Cursor + pos
		if new_pos >= this.evt_lay1Size {
			new_pos = new_pos - this.evt_lay1Size
		}
		pos = new_pos
		header = this.evt_lay1[pos]
	} else {
		if _, ok := this.evt_lay2[pos]; !ok {
			this.evt_lay2[pos] = new(event.EventHeader)
			this.evt_lay2[pos].Init("", 100)
		}
		header = this.evt_lay2[pos]
	}

	old_pre := header.GetPreTimer()
	header.SetPreTimer(a)
	a.SetNextTimer(header)
	a.SetPreTimer(old_pre)
	old_pre.SetNextTimer(a)

	if check_name {
		this.evt_names[a.GetName()] = a
	}

	return true
}

// 投递线程间消息
func (this *Thread) PostThreadMsg(tid int32, a event.IEvent) bool {
	if tid == this.Get_thread_id() {
		fmt.Printf("[W] %d post msg failed\n", tid)
		return false
	}
	if tid >= jiekou.Tid_master && tid < jiekou.Tid_last {
		header := this.evt_threadMsg[tid]
		old_pre := header.GetPreTimer()
		header.SetPreTimer(a)
		a.SetNextTimer(header)
		a.SetPreTimer(old_pre)
		old_pre.SetNextTimer(a)

		return true
	}
	fmt.Printf("[W] %d post msg failed2\n", tid)
	return false
}

// 通过别名获取事件
func (this *Thread) GetEvent(name string) event.IEvent {
	if _, ok := this.evt_names[name]; ok {
		return this.evt_names[name]
	}
	return nil
}

// 删除事件, 只能操作本线程事件
func (this *Thread) RemoveEvent(e event.IEvent) bool {
	if !e.IsHeader() {
		if len(e.GetName()) > 0 {
			delete(this.evt_names, e.GetName())
		}

		this.PopTimer(e)
		this.PopObj(e)

		return true
	}
	return false
}

// 删除事件, 只能操作本线程事件
func (this *Thread) RemoveEventList(header event.IEvent) {
	if header.IsHeader() {
		for {
			// 每次得到链表第一个事件(非)
			e := header.GetNextObj()
			if e.IsHeader() {
				break
			}
			this.RemoveEvent(e)
		}
	}
}

// 从线程事件中弹出指定事件, 只能操作本线程事件
func (this *Thread) PopTimer(e event.IEvent) {
	if !e.IsHeader() {
		e.GetPreTimer().SetNextTimer(e.GetNextTimer())
		e.GetNextTimer().SetPreTimer(e.GetPreTimer())
		e.SetNextTimer(nil)
		e.SetPreTimer(nil)
	}
}

// 从关联对象中弹出指定事件, 只能操作本线程事件
func (this *Thread) PopObj(e event.IEvent) {
	if !e.IsHeader() {
		e.GetPreObj().SetNextObj(e.GetNextObj())
		e.GetNextObj().SetPreObj(e.GetPreObj())
		e.SetNextObj(nil)
		e.SetPreObj(nil)
	}
}

// 接收并处理线程间消息
func (this *Thread) runThreadMsg() {

	header := event.EventHeader{}
	header.Init("", 100)

	G_thread_msg_pool.GetMsg(this.Get_thread_id(), &header) // &this.evt_recvMsg)

	for {
		// 每次得到链表第一个事件(非)
		evt := header.GetNextTimer()
		if evt.IsHeader() {
			break
		}

		// 执行事件, 删除这个事件
		evt.Exec(this.self)
		this.PopTimer(evt)
	}
}

// 发送消息间消息
func (this *Thread) sendThreadMsg() {

	// 发送日志到日志线程
	if this.log_Buffer.Len() > 0 {
		evt := &Event_thread_log{}
		evt.Init("", 100)
		evt.Data = this.log_Buffer
		fmt.Print(this.log_Buffer.String())
		this.PostThreadMsg(jiekou.Tid_log, evt)
		this.log_Buffer.Reset()
	}

	for i := int32(jiekou.Tid_master); i < jiekou.Tid_last; i++ {
		if !this.evt_threadMsg[i].IsEmpty() {

			G_thread_msg_pool.PostMsg(i, this.evt_threadMsg[i])
		}
	}
}

// 运行一次定时器事件(一个线程心跳可以处理多次)
func (this *Thread) runEvents() {
	all_time := (this.last_time - this.start_time) / int64(time.Millisecond)

	all_count := uint64((all_time + jiekou.Evt_gap_time - 1) >> jiekou.Evt_gap_bit)

	for i := this.evt_lastRunCount; i <= all_count; i++ {
		// 执行第一层事件
		this.runExec(this.evt_lay1[this.evt_lay1Cursor])

		// 执行第二层事件
		if _, ok := this.evt_lay2[this.evt_currRunCount]; ok {
			this.runExec(this.evt_lay2[this.evt_currRunCount])
			delete(this.evt_lay2, this.evt_currRunCount)
		}

		this.evt_currRunCount++
		this.evt_lay1Cursor++
		if this.evt_lay1Cursor >= this.evt_lay1Size {
			this.evt_lay1Cursor = 0
		}
	}

	this.evt_lastRunCount = this.evt_currRunCount
}

// 运行一条定时器事件链表, 每次都执行第一个事件, 直到链表为空
func (this *Thread) runExec(header event.IEvent) {
	for {
		// 每次得到链表第一个事件(非)
		evt := header.GetNextTimer()
		if evt.IsHeader() {
			break
		}

		// 执行事件, 返回true, 删除这个事件, 返回false表示用户自己处理
		if evt.Exec(this.self) {
			this.RemoveEvent(evt)
		} else if header.GetNextTimer() == evt {
			// 防止使用者没有删除使用过的事件, 造成死循环, 该事件, 用户要么重新投递到其他链表, 要么删除
			this.RemoveEvent(evt)
		}
	}
}

// 打印事件池现状
func (this *Thread) PrintAll() {

	fmt.Printf(
		`粒度:%d
		粒度移位:%d
		第一层池容量:%d
		第一层游标:%d
		运行次数%d
		`, jiekou.Evt_gap_time, jiekou.Evt_gap_bit, this.evt_lay1Size, this.evt_lay1Cursor, this.evt_currRunCount)

	for k, v := range this.evt_names {
		fmt.Println(k, v)
	}
}

// // Debug logs a message at debug level.
// func Debug(v ...interface{}) {
// 	BeeLogger.Debug(generateFmtStr(len(v)), v...)
// }

// // Trace logs a message at trace level.
// // compatibility alias for Warning()
// func Trace(v ...interface{}) {
// 	BeeLogger.Trace(generateFmtStr(len(v)), v...)
// }

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

// 线程日志 : 调试[D]级别日志
func (this *Thread) LogDebug(f string, v ...interface{}) {
	info := this.log_TimeString + " [D] " + fmt.Sprintf(f, v...) + "\n"
	this.log_Buffer.WriteString(info)
}

// 线程日志 : 信息[I]级别日志
func (this *Thread) LogInfo(f string, v ...interface{}) {
	info := this.log_TimeString + " [I] " + fmt.Sprintf(f, v...) + "\n"
	this.log_Buffer.WriteString(info)
}

// 线程日志 : 警告[W]级别日志
func (this *Thread) LogWarn(f string, v ...interface{}) {
	info := this.log_TimeString + " [W] " + fmt.Sprintf(f, v...) + "\n"
	this.log_Buffer.WriteString(info)
}

// 线程日志 : 错误[E]级别日志
func (this *Thread) LogError(f string, v ...interface{}) {
	info := this.log_TimeString + " [E] " + fmt.Sprintf(f, v...) + "\n"
	this.log_Buffer.WriteString(info)
}

// 线程日志 : 致命[F]级别日志
func (this *Thread) LogFatal(f string, v ...interface{}) {
	info := this.log_TimeString + " [F] " + fmt.Sprintf(f, v...) + "\n"
	this.log_Buffer.WriteString(info)
}
