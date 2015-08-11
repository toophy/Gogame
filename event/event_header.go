package event

// 事件头
type EventHeader struct {
	pre  IEvent // 前一个
	next IEvent // 后一个
}

func (this *EventHeader) Init(name string, t uint64) {
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

func (this *EventHeader) GetName() string {
	return ""
}

func (this *EventHeader) GetPreTimer() IEvent {
	return this.pre
}

func (this *EventHeader) GetNextTimer() IEvent {
	return this.next
}

func (this *EventHeader) SetPreTimer(e IEvent) {
	this.pre = e
}

func (this *EventHeader) SetNextTimer(e IEvent) {
	this.next = e
}

func (this *EventHeader) GetPreObj() IEvent {
	return this.pre
}

func (this *EventHeader) GetNextObj() IEvent {
	return this.next
}

func (this *EventHeader) SetPreObj(e IEvent) {
	this.pre = e
}

func (this *EventHeader) SetNextObj(e IEvent) {
	this.next = e
}

func (this *EventHeader) GetTouchTime() uint64 {
	return 0
}

func (this *EventHeader) SetTouchTime(t uint64) {
}

func (this *EventHeader) SetDelayTime(d uint64, c uint64) {
}

func (this *EventHeader) IsEmpty() bool {
	return this.pre == this.next
}
