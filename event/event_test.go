package event

import (
	"fmt"
	"testing"
)

type Evt_eat struct {
	EventNormal
	Food string
}

func (this *Evt_eat) Exec() bool {
	fmt.Println("吃 " + this.Food)
	return true
}

type Evt_say struct {
	EventNormal
	Chat   string
	Master *Npc
}

func (this *Evt_say) Exec() bool {
	fmt.Println(this.Master.Name + " 说 : " + this.Chat)

	return true
}

type Npc struct {
	EventObj
	Name string
}

func TestEvent(t *testing.T) {
	var g_Home EventHome
	g_Home.Init(30 * 1000)

	evt := &Evt_eat{Food: "西瓜"}
	evt.Init("使劲吃")
	evt.SetTouchTime(63)
	if !g_Home.PushEvent(evt) {
		fmt.Println("增加事件失败 : " + evt.GetName())
	}

	npc := new(Npc)
	npc.Name = "黄蓉"
	npc.InitEventHeader()

	evt2 := &Evt_say{Chat: "真好吃"}
	evt2.Init("使劲吃2")
	evt2.SetTouchTime(163)
	npc.PushEvent(evt2)
	evt2.Master = npc

	if !g_Home.PushEvent(evt2) {
		fmt.Println("增加事件失败 : " + evt2.GetName())
	}

	for i := 0; i < 2000; i++ {
		g_Home.RunEvents()
		//g_Home.PrintAll()
	}
}
