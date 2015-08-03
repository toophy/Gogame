package thread

import (
	"errors"
	"github.com/toophy/Gogame/screen"
	lua "github.com/toophy/gopher-lua"
)

// 场景容器
type ScreenMap map[int32]*screen.Screen

// 场景线程
type ScreenThread struct {
	Thread

	lastScreenId int32       // 最后一个场景id
	screens      ScreenMap   // screen 列表
	luaState     *lua.LState // Lua实体
}

// 新建场景线程
func New_screen_thread(id int32, name string, heart_time int64) (*ScreenThread, error) {
	a := new(ScreenThread)
	err := a.Init_screen_thread(id, name, heart_time)
	if err == nil {
		return a, nil
	}
	return nil, err
}

// 初始化场景线程
func (this *ScreenThread) Init_screen_thread(id int32, name string, heart_time int64) error {
	if id < Tid_screen_1 || id > Tid_screen_9 {
		return errors.New("[E] 线程ID超出范围 [Tid_screen_1,Tid_screen_9]")
	}
	err := this.Init_thread(this, id, name, heart_time)
	if err == nil {
		this.screens = make(ScreenMap, 0)
		this.lastScreenId = (id - 1) * 10000
		return nil
	}
	return err
}

// 增加场景
func (this *ScreenThread) Add_screen(name string, oid int32) bool {
	a := new(screen.Screen)
	a.Load(name, this.lastScreenId, 1, this)
	this.screens[this.lastScreenId] = a

	this.lastScreenId++

	return true
}

// 删除场景
func (this *ScreenThread) Del_screen(id int32) bool {
	if _, ok := this.screens[id]; ok {
		this.screens[id].Unload()
		delete(this.screens, id)
		return true
	}
	return false
}

// 响应线程首次运行
func (this *ScreenThread) on_first_run() {

	errInit := this.reloadLuaState()
	if errInit != nil {
		println(errInit.Error())
		return
	}

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

// 初始化LuaState, 可以用来 Reload LuaState
func (this *ScreenThread) reloadLuaState() error {

	if this.luaState != nil {
		this.luaState.Close()
		this.luaState = nil
	}

	this.luaState = lua.NewState()
	if this.luaState == nil {
		return errors.New("[E] 场景线程初始化Lua失败")
	}

	RegLua_all(this.luaState)
	screen.RegLua_all(this.luaState)

	// 加载所有 screens 文件夹里面的 *.lua 文件
	this.luaState.RequireDir("data/screens")

	return nil
}

// !!!只能获取, 不准许保存指针, 获取LState
func (this *ScreenThread) GetLuaState() *lua.LState {
	return this.luaState
}
