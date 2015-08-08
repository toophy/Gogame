package event

type EventObj struct {
	event_header IEvent // 所有关联事件,便于卸载
}

// 必须调用
func (this *EventObj) InitEventHeader() {
	this.event_header = &EventHeader{}
	this.event_header.Init("")
}

func (this *EventObj) PushEvent(e IEvent) bool {
	if e != nil && !e.IsHeader() {
		old_pre := this.event_header.getPreObj()
		this.event_header.setPreObj(e)
		e.setNextObj(this.event_header)
		e.setPreObj(old_pre)
		old_pre.setNextObj(e)
		return true
	}
	return false
}

func (this *EventObj) RemoveAllEvents() {
	for {
		// 每次得到链表第一个事件(非)
		evt := this.event_header.getNextObj()
		if evt.IsHeader() {
			break
		}
		evt.Remove(evt)
	}
}
