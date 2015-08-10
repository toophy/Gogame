package event

type EventObj struct {
	event_header IEvent // 所有关联事件,便于卸载
}

// 必须调用
func (this *EventObj) InitEventHeader() {
	this.event_header = &EventHeader{}
	this.event_header.Init("", 100)
}

// 压入定时器事件
func (this *EventObj) PostEvent(e IEvent) bool {
	if e != nil && !e.IsHeader() {
		old_pre := this.event_header.GetPreObj()
		this.event_header.SetPreObj(e)
		e.SetNextObj(this.event_header)
		e.SetPreObj(old_pre)
		old_pre.SetNextObj(e)
		return true
	}
	return false
}

// 获取事件列表头
func (this *EventObj) GetEventHeader() IEvent {
	return this.event_header
}
