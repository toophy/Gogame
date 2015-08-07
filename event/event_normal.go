package event

// 普通事件
type EventNormal struct {
	name       string     // 名称
	pre_timer  IEvent     // 定时器前一个
	next_timer IEvent     // 定时器后一个
	pre_obj    IEvent     // 对象前一个
	next_obj   IEvent     // 对象后一个
	home       IEventHome // 事件之家
	touch_time int64      // 定时器触发时间戳
}

func (this *EventNormal) IsHeader() bool {
	return false
}

func (this *EventNormal) Exec() {
}

func (this *EventNormal) Remove() bool {
	if this.home == nil {
		return false
	}

	this.home.PopEvent(this.name)
	this.PopTimer()
	this.PopObj()

	return true
}

func (this *EventNormal) GetName() string {
	return this.name
}

func (this *EventNormal) Init() {
	this.pre_timer = this
	this.pre_obj = this
	this.next_timer = this
	this.next_obj = this
}

func (this *EventNormal) PushTimer(header IEvent) {
	old_pre := header.getPreTimer()
	header.setPreTimer(this)
	this.setNextTimer(header)
	this.setPreTimer(old_pre)
	old_pre.setNextTimer(this)
}

func (this *EventNormal) PopTimer() {
	this.getPreTimer().setNextTimer(this.getNextTimer())
	this.getNextTimer().setPreTimer(this.getPreTimer())
	this.setNextTimer(this)
	this.setPreTimer(this)
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

func (this *EventNormal) PushObj(header IEvent) {
	old_pre := header.getPreObj()
	header.setPreObj(this)
	this.setNextObj(header)
	this.setPreObj(old_pre)
	old_pre.setNextObj(this)
}

func (this *EventNormal) PopObj() {
	this.getPreObj().setNextObj(this.getNextObj())
	this.getNextObj().setPreObj(this.getPreObj())
	this.setNextObj(this)
	this.setPreObj(this)
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

func (this *EventNormal) GetTouchTime() int64 {
	return this.touch_time
}
