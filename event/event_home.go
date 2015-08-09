package event

import (
	"errors"
	"fmt"
)

const (
	Evt_gap_time  = 16     // 心跳时间(毫秒)
	Evt_gap_bit   = 4      // 心跳时间对应得移位(快速运算使用)
	Evt_lay1_time = 160000 // 第一层事件池最大支持时间(毫秒)
)

type EventHome struct {
	Lay1       []IEvent          // 第一层事件池
	Lay2       map[uint64]IEvent // 第二层事件池
	Names      map[string]IEvent // 别名
	lay1Size   uint64            // 第一层池容量
	lay1Cursor uint64            // 第一层游标
	runCount   uint64            // 运行次数
}

func (this *EventHome) Init(lay1_time uint64) error {
	if lay1_time < Evt_gap_time || lay1_time > Evt_lay1_time {
		return errors.New("[E] 第一层支持16毫秒到160000毫秒")
	}
	if this.Names == nil {
		this.lay1Size = lay1_time >> Evt_gap_bit
		this.lay1Cursor = 0
		this.runCount = 1

		this.Lay1 = make([]IEvent, this.lay1Size)
		this.Lay2 = make(map[uint64]IEvent, 0)
		this.Names = make(map[string]IEvent, 0)

		for i := uint64(0); i < this.lay1Size; i++ {
			this.Lay1[i] = new(EventHeader)
			this.Lay1[i].Init("")
		}

		return nil
	}

	return errors.New("[E] EventHome 已经初始化过")
}

func (this *EventHome) PushEvent(a IEvent) bool {
	check_name := len(a.GetName()) > 0
	if check_name {
		if _, ok := this.Names[a.GetName()]; ok {
			return false
		}
	}

	if a.GetTouchTime() < 0 {
		return false
	}

	// 计算放在那一层
	pos := (a.GetTouchTime() + Evt_gap_time - 1) >> Evt_gap_bit
	if pos < 0 {
		pos = 1
	}

	a.SetEventHome(this)

	var header IEvent

	if pos < this.lay1Size {
		new_pos := this.lay1Cursor + pos
		if new_pos >= this.lay1Size {
			new_pos = new_pos - this.lay1Size
		}
		pos = new_pos
		header = this.Lay1[pos]
	} else {
		if _, ok := this.Lay2[pos]; !ok {
			this.Lay2[pos] = new(EventHeader)
			this.Lay2[pos].Init("")
		}
		header = this.Lay2[pos]
	}

	old_pre := header.getPreTimer()
	header.setPreTimer(a)
	a.setNextTimer(header)
	a.setPreTimer(old_pre)
	old_pre.setNextTimer(a)
	this.ShowLay1()

	if check_name {
		this.Names[a.GetName()] = a
	}

	return true
}

func (this *EventHome) GetEvent(name string) IEvent {
	if _, ok := this.Names[name]; ok {
		return this.Names[name]
	}
	return nil
}

func (this *EventHome) PopEvent(name string) {
	delete(this.Names, name)
}

func (this *EventHome) RunEvents() {

	// 执行第一层事件
	this.runExec(this.Lay1[this.lay1Cursor])

	// 执行第二层事件
	if _, ok := this.Lay2[this.runCount]; ok {
		this.runExec(this.Lay2[this.runCount])
		delete(this.Lay2, this.runCount)
	}

	this.runCount++
	this.lay1Cursor++
	if this.lay1Cursor >= this.lay1Size {
		this.lay1Cursor = 0
	}
}

func (this *EventHome) runExec(header IEvent) {
	for {
		// 每次得到链表第一个事件(非)
		evt := header.getNextTimer()
		if evt.IsHeader() {
			break
		}

		// 执行事件, 返回true, 删除这个事件, 返回false表示用户自己处理
		if evt.Exec() {
			evt.Remove(evt)
		} else if header.getNextTimer() == evt {
			// 防止使用者没有删除使用过的事件, 造成死循环, 该事件, 用户要么重新投递到其他链表, 要么删除
			evt.Remove(evt)
		}
	}
}

func (this *EventHome) PrintAll() {

	fmt.Printf(
		`粒度:%d
		粒度移位:%d
		第一层池容量:%d
		第一层游标:%d
		运行次数%d
		`, Evt_gap_time, Evt_gap_bit, this.lay1Size, this.lay1Cursor, this.runCount)

	for k, v := range this.Names {
		fmt.Println(k, v)
	}
}

func (this *EventHome) ShowLay1() {
}
