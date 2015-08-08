package event

// 事件接口
type IEvent interface {
	Init(name string)        // 初始化(name可以为空)
	IsHeader() bool          // 是链表头
	Exec() bool              // 执行
	Remove(self IEvent) bool // 删除事件
	GetName() string         // 获取别名
	// 定时器链表
	PopTimer(self IEvent)  // 离开定时器
	getPreTimer() IEvent   // 获取前一个定时器事件
	getNextTimer() IEvent  // 获取下一个定时器事件
	setPreTimer(e IEvent)  // 设置前一个定时器事件
	setNextTimer(e IEvent) // 设置下一个定时器事件
	// 对象链表
	PopObj(self IEvent)  // 离开对象
	getPreObj() IEvent   // 获取前一个对象事件
	getNextObj() IEvent  // 获取下一个对象事件
	setPreObj(e IEvent)  // 设置前一个对象事件
	setNextObj(e IEvent) // 设置下一个对象事件
	// 事件之家
	SetEventHome(h IEventHome) // 设置事件之家
	GetEventHome() IEventHome  // 获取事件之家
	// 触发时间
	GetTouchTime() uint64            // 获取定时器触发时间戳
	SetTouchTime(t uint64)           // 设置定时器时间戳
	SetDelayTime(d uint64, c uint64) // 设置定时器相对时间, c是当前时间戳
}

// 定时器事件之家接口
type IEventHome interface {
	PushEvent(e IEvent) bool     // 压入定时器事件
	GetEvent(name string) IEvent // 通过别名获取事件
	PopEvent(name string)        // 弹出别名
	ShowLay1()                   // 显示所有事件
}

// 定时器关联对象
type IEventObj interface {
	PushEvent(e IEvent) bool // 压入关联事件
	RemoveAllEvents() bool   // 删除所有关联事件
}
