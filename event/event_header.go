package event

// 事件头
type EventHeader struct {
	pre  IEvent // 前一个
	next IEvent // 后一个
}

////////////////////////////////////////
func (this *EventHeader) IsHeader() bool {
	return true
}

func (this *EventHeader) Exec() {
}

func (this *EventHeader) Remove() bool {
	return false
}

func (this *EventHeader) GetName() string {
	return ""
}

func (this *EventHeader) Init() {
	this.pre = this
	this.next = this
}

func (this *EventHeader) PushTimer(header IEvent) {
}

func (this *EventHeader) PopTimer() {
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

func (this *EventHeader) PushObj(header IEvent) {
}

func (this *EventHeader) PopObj() {
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

func (this *EventNormal) GetTouchTime() int64 {
	return -1
}
