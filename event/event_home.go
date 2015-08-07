package event

import (
	"errors"
	"fmt"
)

type EventHome struct {
	Lay1       []IEvent          // 第一层事件池
	Lay2       map[int64]IEvent  // 第二层事件池
	Names      map[string]IEvent // 别名
	gap        int64             // 粒度
	gapBit     int64             // 粒度移位
	lay1Size   int64             // 第一层池容量
	lay1Cursor int64             // 第一层游标
	runCount   int64             // 运行次数
}

func (this *EventHome) Init(lay1_time int64) error {
	if lay1_time < 32 || lay1_time > 320000 {
		return errors.New("[E] 第一层支持32毫秒到320000毫秒")
	}
	if this.Names == nil {
		this.gap = 32   // 32毫秒
		this.gapBit = 5 // 2^5 => 32
		this.lay1Size = lay1_time >> this.gapBit
		this.lay1Cursor = -1
		this.runCount = 0

		this.Lay1 = make([]IEvent, this.lay1Size)
		this.Lay2 = make(map[string]IEvent, 0)
		this.Names = make(map[string]IEvent, 0)

		for i := 0; i < this.lay1Size; i++ {
			this.Lay1[i] = new(EventHeader)
			this.Lay1[i].Init()
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
	pos := a.GetTouchTime() >> this.gapBit
	if pos < this.lay1Size {
		new_pos := this.lay1Cursor + pos
		if new_pos >= this.lay1Size {
			new_pos = new_pos - this.lay1Size
		}
		a.PushTimer(this.Lay1[new_pos])
	} else {
		if _, ok := this.Lay2[pos]; ok {
			a.PushTimer(this.Lay2[pos])
		} else {
			this.Lay2[pos] = make(map[string]IEvent, 0)
			this.Lay2[pos] = new(EventHeader)
			this.Lay2[pos].Init()
			a.PushTimer(this.Lay2[pos])
		}
	}

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
	this.runCount++
	this.lay1Cursor++
	if this.lay1Cursor >= this.lay1Size {
		this.lay1Cursor = 0
	}

	// 执行第一层事件
	// 链表中, 依次执行

	// 执行第二层事件
	if _, ok := this.Lay2[this.runCount]; ok {
		// 链表中, 依次执行
	}
}

func (this *EventHome) PrintAll() {

	fmt.Printf(
		`粒度:%d
		粒度移位:%d
		第一层池容量:%d
		第一层游标:%d
		运行次数%d
		`, this.gap, this.gapBit, this.lay1Size, this.lay1Cursor, this.runCount)

	for k, v := range this.Names {
		fmt.Println(v)
	}
}
