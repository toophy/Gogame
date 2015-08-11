package thread

import (
	"errors"
	"fmt"
	"github.com/toophy/Gogame/event"
	"time"
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
	Tid_last
)

const (
	Evt_gap_time  = 16     // 心跳时间(毫秒)
	Evt_gap_bit   = 4      // 心跳时间对应得移位(快速运算使用)
	Evt_lay1_time = 160000 // 第一层事件池最大支持时间(毫秒)
)

var (
	PoolSizePerTick = 10
	ErrTaskNotFound = errors.New("The task was not found.")
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
}

// 线程基本功能
// procs      msg_list               // 待处理消息链表
// sends      [Tid_last]bind_fn_list // 发送消息链表组
// recvs      msg_list               // 接收消息链表
type Thread struct {
	id               int32                   // Id号
	name             string                  // 线程名称
	heart_time       int64                   // 心跳时间(毫秒)
	start_time       int64                   // 线程开启时间戳
	last_time        int64                   // 最近一次线程运行时间戳
	heart_rate       float64                 // 本次心跳比率
	pre_stop         bool                    // 预备停止
	self             IThread                 // 自己, 初始化之后, 不要操作
	first_run        bool                    // 线程首次运行
	evt_lay1         []event.IEvent          // 第一层事件池
	evt_lay2         map[uint64]event.IEvent // 第二层事件池
	evt_names        map[string]event.IEvent // 别名
	evt_lay1Size     uint64                  // 第一层池容量
	evt_lay1Cursor   uint64                  // 第一层游标
	evt_lastRunCount uint64                  // 最近一次运行次数
	evt_currRunCount uint64                  // 当前运行次数
	evt_threadMsg    [Tid_last]event.IEvent  // 保存将要发给其他线程的事件(消息)
	evt_recvMsg      event.EventHeader       // 接收线程间消息
}

// 初始化线程(必须调用)
// usage : Init_thread(Tid_master, "主线程", 100)
func (this *Thread) Init_thread(self IThread, id int32, name string, heart_time int64, lay1_time uint64) error {
	if id < Tid_master || id >= Tid_last {
		return errors.New("[E] 线程ID超出范围 [Tid_master,Tid_last]")
	}
	if self == nil {
		return errors.New("[E] 线程自身指针不能为nil")
	}

	if lay1_time < Evt_gap_time || lay1_time > Evt_lay1_time {
		return errors.New("[E] 第一层支持16毫秒到160000毫秒")
	}

	if len(this.evt_names) > 0 {
		return errors.New("[E] EventHome 已经初始化过")
	}

	this.id = id
	this.name = name
	this.heart_time = heart_time
	this.start_time = time.Now().UnixNano() / int64(time.Millisecond)
	this.last_time = this.start_time
	this.heart_rate = 1.0
	this.self = self
	this.first_run = true

	// 初始化事件池
	this.evt_lay1Size = lay1_time >> Evt_gap_bit
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

	for i := 0; i < Tid_last; i++ {
		this.evt_threadMsg[i] = new(event.EventHeader)
		this.evt_threadMsg[i].Init("", 100)
	}

	this.evt_recvMsg.Init("", 100)

	return nil
}

// 运行线程
func (this *Thread) Run_thread() {
	// 计算心跳误差值, 决定心跳滴答(小数), heart_time, last_time, heart_rate
	// 处理线程间接收消息, 分配到水表定时器
	// 执行水表定时器
	go func() {
		GetMaster().Add_run_thread(this.self)

		this.start_time = time.Now().UnixNano() / int64(time.Millisecond)
		this.last_time = this.start_time
		next_time := time.Duration(this.heart_time * int64(time.Millisecond))
		run_time := int64(0)

		this.self.on_first_run()

		for {
			time.Sleep(next_time)
			this.last_time = time.Now().UnixNano() / int64(time.Millisecond)
			this.runThreadMsg()
			this.runEvents()
			this.runOnce()
			this.self.on_run()

			this.sendThreadMsg()

			// 计算下一次运行的时间
			run_time = int64((time.Now().UnixNano() / int64(time.Millisecond)) - this.last_time)
			if run_time >= this.heart_time {
				run_time = this.heart_time - 10
			} else if run_time < 0 {
				run_time = 0
			}

			next_time = time.Duration((this.heart_time - run_time) * int64(time.Millisecond))

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
	pos := (a.GetTouchTime() + Evt_gap_time - 1) >> Evt_gap_bit
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
		return false
	}
	if tid >= Tid_master && tid < Tid_last {
		header := this.evt_threadMsg[tid]
		old_pre := header.GetPreTimer()
		header.SetPreTimer(a)
		a.SetNextTimer(header)
		a.SetPreTimer(old_pre)
		old_pre.SetNextTimer(a)
		return true
	}
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

	G_thread_msg_pool.GetMsg(this.Get_thread_id(), &this.evt_recvMsg)

	for {
		// 每次得到链表第一个事件(非)
		evt := this.evt_recvMsg.GetNextTimer()
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
	for i := int32(Tid_master); i < Tid_last; i++ {
		if !this.evt_threadMsg[i].IsEmpty() {
			println("sendThreadMsg")
			G_thread_msg_pool.PostMsg(i, this.evt_threadMsg[i])
		}
	}
}

// 运行一次定时器事件(一个线程心跳可以处理多次)
func (this *Thread) runEvents() {
	all_time := this.last_time - this.start_time

	all_count := uint64((all_time + Evt_gap_time - 1) >> Evt_gap_bit)

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
		`, Evt_gap_time, Evt_gap_bit, this.evt_lay1Size, this.evt_lay1Cursor, this.evt_currRunCount)

	for k, v := range this.evt_names {
		fmt.Println(k, v)
	}
}
