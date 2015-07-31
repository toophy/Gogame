package thread

import (
	"fmt"
	"game/screen"
)

type ScreenMap map[int32]*screen.Screen

// 场景线程
type Screen struct {
	Thread

	Last_screen_id int32     // 最后一个场景id
	Screens        ScreenMap // screen 列表
}

func New_screen_thread(id int32, name string, heart_time int64) *Screen {
	a := &Screen{}
	if a.Init_screen_thread(id, name, heart_time) {
		return a
	}
	return nil
}

func (s *Screen) Init_screen_thread(id int32, name string, heart_time int64) bool {
	if id < Tid_screen_1 || id > Tid_screen_9 {
		return false
	}
	if s.Init_thread(id, name, heart_time) {
		s.Screens = make(ScreenMap, 0)
		s.Last_screen_id = (id - 1) * 10000
		return true
	}
	return false
}

func (s *Screen) Add_screen(name string, oid int32) bool {
	a := &screen.Screen{}
	a.Load(name, s.Last_screen_id, 1)
	s.Screens[1] = a

	s.Last_screen_id++

	return true
}

func (s *Screen) Del_screen(id int32) bool {
	if _, ok := s.Screens[id]; ok {
		s.Screens[id].Unload()
		delete(s.Screens, id)
		return true
	}
	return false
}

// 响应线程退出
func (this *Screen) On_thread_end() {
	fmt.Printf("场景线程(%s)正常关闭\n", this.name)
}
