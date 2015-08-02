package thread

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
)

func (this *Screen) Tolua_OnInitScreen() /*int*/ {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r.(error).Error())
		}
	}()

	ud := this.LuaState.NewUserData()
	ud.Value = this
	this.LuaState.SetMetatable(ud, this.LuaState.GetTypeMetatable(struct_name))

	if err := this.LuaState.CallByParam(lua.P{
		Fn:      this.LuaState.GetGlobal("OnInitScreen"), // 调用的Lua函数
		NRet:    0,                                       // 返回值的数量
		Protect: true,                                    // 保护?
	}, ud); err != nil {
		panic(err)
	}

	// ret := this.LuaState.Get(-1)
	// this.LuaState.Pop(1)
	// return int(ret.(lua.LNumber))
}
