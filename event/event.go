package event

// 定时器事件之家接口
type IEventHome interface {
	PushEvent(e IEvent) bool
	GetEvent(name string) IEvent
	PopEvent(name string)
}

// 事件接口
type IEvent interface {
	IsHeader() bool      // 是链表头
	Exec()               // 执行
	Remove() bool        // 删除事件
	GetName() string     // 获取别名
	Init()               // 初始化
	GetTouchTime() int64 // 获取定时器触发时间戳
	// 定时器链表
	PushTimer(header IEvent) // 加入定时器
	PopTimer()               // 离开定时器
	getPreTimer() IEvent     // 获取前一个定时器事件
	getNextTimer() IEvent    // 获取下一个定时器事件
	setPreTimer(e IEvent)    // 设置前一个定时器事件
	setNextTimer(e IEvent)   // 设置下一个定时器事件
	// 对象链表
	PushObj(header IEvent) // 加入对象
	PopObj()               // 离开对象
	getPreObj() IEvent     // 获取前一个对象事件
	getNextObj() IEvent    // 获取下一个对象事件
	setPreObj(e IEvent)    // 设置前一个对象事件
	setNextObj(e IEvent)   // 设置下一个对象事件
	// 事件之家
	SetEventHome(h IEventHome) // 设置事件之家
	GetEventHome() IEventHome  // 获取事件之家
}
