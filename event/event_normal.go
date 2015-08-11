package event

// 事件接口
type IEvent interface {
	Init(name string, t uint64) // 初始化(name可以为空, t是触发时间)
	IsHeader() bool             // 是链表头
	Exec() bool                 // 执行
	GetName() string            // 获取别名
	IsEmpty() bool              // 是空列表
	// 定时器链表
	GetPreTimer() IEvent   // 获取前一个定时器事件
	GetNextTimer() IEvent  // 获取下一个定时器事件
	SetPreTimer(e IEvent)  // 设置前一个定时器事件
	SetNextTimer(e IEvent) // 设置下一个定时器事件
	// 对象链表
	GetPreObj() IEvent   // 获取前一个对象事件
	GetNextObj() IEvent  // 获取下一个对象事件
	SetPreObj(e IEvent)  // 设置前一个对象事件
	SetNextObj(e IEvent) // 设置下一个对象事件
	// 触发时间
	GetTouchTime() uint64            // 获取定时器触发时间戳
	SetTouchTime(t uint64)           // 设置定时器时间戳
	SetDelayTime(d uint64, c uint64) // 设置定时器相对时间, c是当前时间戳
}

// 普通事件
type EventNormal struct {
	name       string // 名称
	pre_timer  IEvent // 定时器前一个
	next_timer IEvent // 定时器后一个
	pre_obj    IEvent // 对象前一个
	next_obj   IEvent // 对象后一个
	touch_time uint64 // 定时器触发时间戳
}

func (this *EventNormal) Init(name string, t uint64) {
	this.name = name
	this.pre_timer = this
	this.pre_obj = this
	this.next_timer = this
	this.next_obj = this
	this.touch_time = t
}

func (this *EventNormal) IsHeader() bool {
	return false
}

func (this *EventNormal) Exec() bool {
	println("Normal Exec")
	return true
}

func (this *EventNormal) GetName() string {
	return this.name
}

func (this *EventNormal) GetPreTimer() IEvent {
	return this.pre_timer
}

func (this *EventNormal) GetNextTimer() IEvent {
	return this.next_timer
}

func (this *EventNormal) SetPreTimer(e IEvent) {
	this.pre_timer = e
}

func (this *EventNormal) SetNextTimer(e IEvent) {
	this.next_timer = e
}

func (this *EventNormal) GetPreObj() IEvent {
	return this.pre_obj
}

func (this *EventNormal) GetNextObj() IEvent {
	return this.next_obj
}

func (this *EventNormal) SetPreObj(e IEvent) {
	this.pre_obj = e
}

func (this *EventNormal) SetNextObj(e IEvent) {
	this.next_obj = e
}

func (this *EventNormal) GetTouchTime() uint64 {
	return this.touch_time
}

func (this *EventNormal) SetTouchTime(t uint64) {
	this.touch_time = t
}

func (this *EventNormal) SetDelayTime(d uint64, c uint64) {
	this.touch_time = c + d
}

func (this *EventNormal) IsEmpty() bool {
	return true
}
