package screen

import (
	lua "github.com/toophy/gopher-lua"
)

// 注册本包所有Lua接口结构
func RegLua_all(L *lua.LState) error {

	type regLuaFunc func(string, *lua.LState) error

	regLuaStructs := map[string]regLuaFunc{
		"Screen": regLua_screen,
	}

	for k, _ := range regLuaStructs {
		regLuaStructs[k](k, L)
	}

	return nil
}

// 向Lua注册结构 : Screen
func regLua_screen(struct_name string, L *lua.LState) error {

	mt := L.NewTypeMetatable(struct_name)
	L.SetGlobal(struct_name, mt)

	// 检查Lua首个参数是不是对象指针
	check := func(L *lua.LState) *Screen {
		ud := L.CheckUserData(1)
		if v, ok := ud.Value.(*Screen); ok {
			return v
		}
		L.ArgError(1, struct_name+" expected")

		return nil
	}

	// 成员函数
	// L.SetField(mt, "new", L.NewFunction(newScreenThread))

	// 成员变量
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(),

		map[string]lua.LGFunction{

			"Get_data": func(L *lua.LState) int {
				p := check(L)

				ret := p.Get_data()

				L.Push(ret)
				return 1
			},

			"Get_thread": func(L *lua.LState) int {
				p := check(L)

				ret := L.GetUserData("ScreenThread", p.Get_thread())

				L.Push(ret)
				return 1
			},
		}))

	return nil
}
