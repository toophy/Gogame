package event

// 普通事件
type EventNormal struct {
	name       string     // 名称
	pre_timer  IEvent     // 定时器前一个
	next_timer IEvent     // 定时器后一个
	pre_obj    IEvent     // 对象前一个
	next_obj   IEvent     // 对象后一个
	home       IEventHome // 事件之家
	touch_time uint64     // 定时器触发时间戳
}

func (this *EventNormal) Init(name string) {
	this.name = name
	this.pre_timer = this
	this.pre_obj = this
	this.next_timer = this
	this.next_obj = this
}

func (this *EventNormal) IsHeader() bool {
	return false
}

func (this *EventNormal) Exec() bool {
	println("Normal Exec")
	return true
}

func (this *EventNormal) Remove(self IEvent) bool {
	if this.home == nil {
		return false
	}

	this.home.PopEvent(this.name)
	this.PopTimer(self)
	this.PopObj(self)
	this.home = nil

	return true
}

func (this *EventNormal) GetName() string {
	return this.name
}

func (this *EventNormal) PopTimer(self IEvent) {
	this.getPreTimer().setNextTimer(this.getNextTimer())
	this.getNextTimer().setPreTimer(this.getPreTimer())
	this.setNextTimer(self)
	this.setPreTimer(self)
}

func (this *EventNormal) getPreTimer() IEvent {
	return this.pre_timer
}

func (this *EventNormal) getNextTimer() IEvent {
	return this.next_timer
}

func (this *EventNormal) setPreTimer(e IEvent) {
	this.pre_timer = e
}

func (this *EventNormal) setNextTimer(e IEvent) {
	this.next_timer = e
}

func (this *EventNormal) PopObj(self IEvent) {
	this.getPreObj().setNextObj(this.getNextObj())
	this.getNextObj().setPreObj(this.getPreObj())
	this.setNextObj(self)
	this.setPreObj(self)
}

func (this *EventNormal) getPreObj() IEvent {
	return this.pre_obj
}

func (this *EventNormal) getNextObj() IEvent {
	return this.next_obj
}

func (this *EventNormal) setPreObj(e IEvent) {
	this.pre_obj = e
}

func (this *EventNormal) setNextObj(e IEvent) {
	this.next_obj = e
}

func (this *EventNormal) SetEventHome(h IEventHome) {
	this.home = h
}

func (this *EventNormal) GetEventHome() IEventHome {
	return this.home
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
