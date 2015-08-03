package thread

import (
	"fmt"
	lua "github.com/toophy/gopher-lua"
)

// 调用Lua函数 : OnInitScreen
func (this *ScreenThread) Tolua_OnInitScreen() (ret int) {
	defer func() {
		if r := recover(); r != nil {
			ret = -1
			fmt.Println(r.(error).Error())
		}
	}()

	if err := this.luaState.CallByParam(lua.P{
		Fn:      this.luaState.GetFunction("main", "OnInitScreen"),
		NRet:    1,
		Protect: true,
	}, this.luaState.GetUserData("ScreenThread", this)); err != nil {
		panic(err)
	}

	ret_lua := this.luaState.Get(-1)
	ret = int(ret_lua.(lua.LNumber))
	this.luaState.Pop(1)

	return
}

// 调用Lua函数 : common.PrintTable
func (this *ScreenThread) Tolua_print_table(a lua.LValue) {
	// 捕捉异常
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r.(error).Error())
		}
	}()

	// 调用Lua脚本函数
	if err := this.luaState.CallByParam(lua.P{
		Fn:      this.luaState.GetFunction("common", "PrintTable"), // 调用的Lua函数
		NRet:    0,                                                 // 返回值的数量
		Protect: true,                                              // 保护?
	}, a); err != nil {
		panic(err)
	}
}
