package thread

import (
	lua "github.com/yuin/gopher-lua"
)

const struct_name = "screen_thread"

func RegLua_sct(L *lua.LState) {

	mt := L.NewTypeMetatable(struct_name)
	L.SetGlobal("screen_thread", mt)

	// 检查Lua首个参数是不是对象指针
	check := func(L *lua.LState) *Screen {
		ud := L.CheckUserData(1)
		if v, ok := ud.Value.(*Screen); ok {
			return v
		}
		L.ArgError(1, "Screen expected")

		return nil
	}

	// 成员函数
	// L.SetField(mt, "new", L.NewFunction(newScreenThread))

	// 成员变量
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(),

		map[string]lua.LGFunction{

			// 增加场景
			"Add_screen": func(L *lua.LState) int {
				p := check(L)
				name := L.CheckString(2)
				oid := int32(L.CheckInt(3))

				ret := p.Add_screen(name, oid)

				L.Push(lua.LBool(ret))
				return 1
			},

			// 获取随机数
			"Get_randNum": func(L *lua.LState) int {
				p := check(L)

				ret := p.GetRandNum()

				L.Push(lua.LNumber(ret))
				return 1

			},

			// 设置随机数
			"Set_randNum": func(L *lua.LState) int {
				p := check(L)
				num := int64(L.CheckInt64(2))

				p.SetRandNum(num)

				return 0
			},
		}))
}
