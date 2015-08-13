package thread

import (
	lua "github.com/toophy/gopher-lua"
)

// 获取用Lua类型封装结构指针  *LUserData
func (this *ScreenThread) GetLUserData(n string, a interface{}) *lua.LUserData {

	ud := this.luaState.NewUserData()
	ud.Value = a
	this.luaState.SetMetatable(ud, this.luaState.GetTypeMetatable(n))

	return ud
}

// 调用Lua函数 : OnInitScreen
func (this *ScreenThread) Tolua_OnInitScreen() (ret int) {
	// 捕捉异常
	defer func() {
		if r := recover(); r != nil {
			ret = -1
			this.LogFatal("Tolua_OnInitScreen : " + r.(error).Error())
		}
	}()

	// 调用Lua脚本函数
	if err := this.luaState.CallByParam(lua.P{
		Fn:      this.luaState.GetFunction("main", "OnInitScreen"), // 调用的Lua函数
		NRet:    1,                                                 // 返回值的数量
		Protect: true,                                              // 保护?
	}, this.GetLUserData("ScreenThread", this)); err != nil {
		panic(err)
	}

	// 处理Lua脚本函数返回值
	ret_lua := this.luaState.Get(-1)
	ret = int(ret_lua.(lua.LNumber))
	this.luaState.Pop(1)

	return
}
