package event

// 事件头
type EventHeader struct {
	pre  IEvent // 前一个
	next IEvent // 后一个
}

func (this *EventHeader) Init(name string) {
	this.pre = this
	this.next = this
}

func (this *EventHeader) IsHeader() bool {
	return true
}

func (this *EventHeader) Exec() bool {
	println("Header Exec")
	return true
}

func (this *EventHeader) Remove(self IEvent) bool {
	return false
}

func (this *EventHeader) GetName() string {
	return ""
}

func (this *EventHeader) PopTimer(self IEvent) {
}

func (this *EventHeader) getPreTimer() IEvent {
	return this.pre
}

func (this *EventHeader) getNextTimer() IEvent {
	return this.next
}

func (this *EventHeader) setPreTimer(e IEvent) {
	this.pre = e
}

func (this *EventHeader) setNextTimer(e IEvent) {
	this.next = e
}

func (this *EventHeader) PopObj(self IEvent) {
}

func (this *EventHeader) getPreObj() IEvent {
	return this.pre
}

func (this *EventHeader) getNextObj() IEvent {
	return this.next
}

func (this *EventHeader) setPreObj(e IEvent) {
	this.pre = e
}

func (this *EventHeader) setNextObj(e IEvent) {
	this.next = e
}

func (this *EventHeader) SetEventHome(h IEventHome) {
}

func (this *EventHeader) GetEventHome() IEventHome {
	return nil
}

func (this *EventHeader) GetTouchTime() uint64 {
	return 0
}

func (this *EventHeader) SetTouchTime(t uint64) {
}

func (this *EventHeader) SetDelayTime(d uint64, c uint64) {
}
