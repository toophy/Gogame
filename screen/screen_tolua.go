package screen

import (
	"fmt"
	lua "github.com/toophy/gopher-lua"
)

// 调用Lua函数 : OnInit
func (this *Screen) Tolua_screen_init() {
	// 捕捉异常
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r.(error).Error())
		}
	}()

	// 调用Lua脚本函数
	if err := this.thread.GetLuaState().CallByParam(lua.P{
		Fn:      this.thread.GetLuaState().GetFunction(this.Get_config().ModName, "OnInit"), // 调用的Lua函数
		NRet:    0,                                                                          // 返回值的数量
		Protect: true,                                                                       // 保护?
	}, this.thread.GetLuaState().GetUserData("Screen", this)); err != nil {
		panic(err)
	}
}

// 调用Lua函数 : OnHeartBeat
func (this *Screen) Tolua_heart_beat() {
	// 捕捉异常
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r.(error).Error())
		}
	}()

	// 调用Lua脚本函数
	if err := this.thread.GetLuaState().CallByParam(lua.P{
		Fn:      this.thread.GetLuaState().GetFunction(this.Get_config().ModName, "OnHeartBeat"), // 调用的Lua函数
		NRet:    0,                                                                               // 返回值的数量
		Protect: true,                                                                            // 保护?
	}, this.thread.GetLuaState().GetUserData("Screen", this)); err != nil {
		panic(err)
	}
}
