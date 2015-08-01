package thread

import (
	"Gogame/screen"
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

// 场景容器
type ScreenMap map[int32]*screen.Screen

// 场景线程
type Screen struct {
	Thread

	Last_screen_id int32       // 最后一个场景id
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

// 测试Lua调用Go函数
func Sum(L *lua.LState) int {
	a := L.ToInt(1)
	b := L.ToInt(2)

	L.Push(lua.LNumber(a + b))

	return 1
}

func (this *Screen) Tolua_OnInitScreen() /*int*/ {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r.(error).Error())
		}
	}()

	RegLua_screenThread(this.LuaState)

	ud := this.LuaState.NewUserData()
	ud.Value = this
	this.LuaState.SetMetatable(ud, this.LuaState.GetTypeMetatable(regScreenThreadName))

	if err := this.LuaState.CallByParam(lua.P{
		Fn:      this.LuaState.GetGlobal("OnInitScreen"), // 调用的Lua函数
		NRet:    0,                                       // 返回值的数量
		Protect: true,                                    // 保护?
	}, this.LuaState.GetTypeMetatable(regScreenThreadName)); err != nil {
		panic(err)
	}

	// ret := this.LuaState.Get(-1)
	// this.LuaState.Pop(1)

	// return int(ret.(lua.LNumber))

	// 	this.LuaState.SetGlobal("g_thread", luar.New(this.LuaState, this))

	// if err := this.LuaState.CallByParam(lua.P{
	// 	Fn:      this.LuaState.GetGlobal("OnInitScreen"), // 调用的Lua函数
	// 	NRet:    0,                                       // 返回值的数量
	// 	Protect: true,                                    // 保护?
	// }, luar.New(this.LuaState, this)); err != nil {
	// 	panic(err)
	// }
}

// 响应线程首次运行
func (this *Screen) on_first_run() {
	this.LuaState = lua.NewState()
	if this.LuaState == nil {
		panic("场景线程初始化Lua失败")
	}

	this.LuaState.SetGlobal("mysum", this.LuaState.NewFunction(Sum))

	err := this.LuaState.DoFile("data/screens/main.lua")
	if err != nil {
		fmt.Println(err.Error())
	}

	this.Tolua_OnInitScreen()

	// Sub(this.LuaState, 100, 88)
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
