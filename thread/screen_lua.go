package thread

import (
	lua "github.com/yuin/gopher-lua"
)

func init() {

}

const regScreenThreadName = "screen_thread"

// Registers to given L.
func RegLua_screenThread(L *lua.LState) {
	mt := L.NewTypeMetatable(regScreenThreadName)
	L.SetGlobal("screen_thread", mt)

	// 成员函数
	// L.SetField(mt, "new", L.NewFunction(newScreenThread))

	// 成员变量
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), screenThreadMethods))
}

// Constructor
// func newScreenThread(L *lua.LState) int {
// 	println("newScreenThread")
// 	scr := &Screen{}
// 	ud := L.NewUserData()
// 	ud.Value = scr
// 	L.SetMetatable(ud, L.GetTypeMetatable(regScreenThreadName))
// 	L.Push(ud)
// 	return 1
// }

// 检查Lua首个参数是不是对象指针
func checkScreen(L *lua.LState) *Screen {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Screen); ok {
		return v
	}
	L.ArgError(1, "Screen expected")
	return nil
}

var screenThreadMethods = map[string]lua.LGFunction{
	"Add_screen": screenThreadAddScreen,
}

// Getter and setter for the Person#Name
func screenThreadAddScreen(L *lua.LState) int {
	p := checkScreen(L)
	name := L.CheckString(2)
	oid := int32(L.CheckInt(3))

	ret := p.Add_screen(name, oid)

	L.Push(lua.LBool(ret))
	return 1
}
