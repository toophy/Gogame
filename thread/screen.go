package thread

import (
	"fmt"
	"github.com/toophy/Gogame/screen"
	lua "github.com/yuin/gopher-lua"
)

// 场景容器
type ScreenMap map[int32]*screen.Screen

// 场景线程
type ScreenThread struct {
	Thread

	lastScreenId int32       // 最后一个场景id
	randNum      int64       //测试64位整数
	screens      ScreenMap   // screen 列表
	luaState     *lua.LState // Lua实体
}

// 新建场景线程
func New_screen_thread(id int32, name string, heart_time int64) *ScreenThread {
	a := &ScreenThread{}
	if a.Init_screen_thread(id, name, heart_time) {
		return a
	}
	return nil
}

// 初始化场景线程
func (s *ScreenThread) Init_screen_thread(id int32, name string, heart_time int64) bool {
	if id < Tid_screen_1 || id > Tid_screen_9 {
		return false
	}
	if s.Init_thread(s, id, name, heart_time) {
		s.screens = make(ScreenMap, 0)
		s.lastScreenId = (id - 1) * 10000
		return true
	}
	return false
}

// 增加场景
func (s *ScreenThread) Add_screen(name string, oid int32) bool {
	a := &screen.Screen{}
	a.Load(name, s.lastScreenId, 1)
	s.screens[1] = a

	s.lastScreenId++

	return true
}

// 删除场景
func (s *ScreenThread) Del_screen(id int32) bool {
	if _, ok := s.screens[id]; ok {
		s.screens[id].Unload()
		delete(s.screens, id)
		return true
	}
	return false
}

func (this *ScreenThread) GetRandNum() int64 {
	return this.randNum
}

func (this *ScreenThread) SetRandNum(a int64) {
	this.randNum = a
}

// 响应线程首次运行
func (this *ScreenThread) on_first_run() {
	this.luaState = lua.NewState()
	if this.luaState == nil {
		panic("场景线程初始化Lua失败")
	}

	this.randNum = 12345678912345678

	// this.LuaState.SetGlobal("mysum", this.LuaState.NewFunction(Sum))

	err := this.luaState.DoFile("data/screens/main.lua")
	if err != nil {
		fmt.Println(err.Error())
	}

	RegLua_all(this.luaState)

	println(this.Tolua_OnInitScreen())
}

// 响应线程退出
func (this *ScreenThread) on_end() {
	if this.luaState != nil {
		this.luaState.Close()
		this.luaState = nil
	}
}

// 响应线程运行
func (this *ScreenThread) on_run() {
}
