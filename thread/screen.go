package thread

import (
	"fmt"
	"github.com/toophy/Gogame/screen"
	lua "github.com/yuin/gopher-lua"
)

// 场景容器
type ScreenMap map[int32]*screen.Screen

// 场景线程
type Screen struct {
	Thread

	Last_screen_id int32       // 最后一个场景id
	RandNum        int64       //测试64位整数
	Screens        ScreenMap   // screen 列表
	LuaState       *lua.LState // Lua实体
}

// 新建场景线程
func New_screen_thread(id int32, name string, heart_time int64) *Screen {
	a := &Screen{}
	if a.Init_screen_thread(id, name, heart_time) {
		return a
	}
	return nil
}

// 初始化场景线程
func (s *Screen) Init_screen_thread(id int32, name string, heart_time int64) bool {
	if id < Tid_screen_1 || id > Tid_screen_9 {
		return false
	}
	if s.Init_thread(s, id, name, heart_time) {
		s.Screens = make(ScreenMap, 0)
		s.Last_screen_id = (id - 1) * 10000
		return true
	}
	return false
}

// 增加场景
func (s *Screen) Add_screen(name string, oid int32) bool {
	a := &screen.Screen{}
	a.Load(name, s.Last_screen_id, 1)
	s.Screens[1] = a

	s.Last_screen_id++

	return true
}

// 删除场景
func (s *Screen) Del_screen(id int32) bool {
	if _, ok := s.Screens[id]; ok {
		s.Screens[id].Unload()
		delete(s.Screens, id)
		return true
	}
	return false
}

func (this *Screen) GetRandNum() int64 {
	return this.RandNum
}

func (this *Screen) SetRandNum(a int64) {
	this.RandNum = a
}

// 响应线程首次运行
func (this *Screen) on_first_run() {
	this.LuaState = lua.NewState()
	if this.LuaState == nil {
		panic("场景线程初始化Lua失败")
	}

	this.RandNum = 12345678912345678

	// this.LuaState.SetGlobal("mysum", this.LuaState.NewFunction(Sum))

	err := this.LuaState.DoFile("data/screens/main.lua")
	if err != nil {
		fmt.Println(err.Error())
	}

	RegLua_sct(this.LuaState)

	this.Tolua_OnInitScreen()
}

// 响应线程退出
func (this *Screen) on_end() {
	if this.LuaState != nil {
		this.LuaState.Close()
		this.LuaState = nil
	}
}

// 响应线程运行
func (this *Screen) on_run() {
}
