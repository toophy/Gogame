package event

import (
	"fmt"
	"testing"
)

type Evt_eat struct {
	Food string
	EventNormal
}

func (this *Evt_eat) Exec() bool {
	fmt.Println("吃 " + this.Food)
	return true
}

func TestEvent(t *testing.T) {
	var g_Home EventHome
	g_Home.Init(30 * 1000)
	evt := &Evt_eat{Food: "西瓜"}
	evt.Init("使劲吃")
	evt.SetTouchTime(0)
	g_Home.PushEvent(evt)

	for i := 0; i < 3; i++ {
		g_Home.RunEvents()
		g_Home.PrintAll()
	}
}

type IX interface {
	Name() bool
}

type A struct {
}

func (a *A) Name() bool {
	fmt.Println("A name")
	return true
}

type B struct {
	A
}

func (b *B) Name() bool {
	fmt.Println("B name")
	return true
}

func TestFunc(t *testing.T) {
	var x IX
	x = &B{}
	x.Name()
}
