// main.go
package main

import (
	//"Gogame/thread"

	lua "github.com/yuin/gopher-lua"
)

// Gogame framework version.
const VERSION = "0.0.1"

type Person struct {
	Name string
}

const luaPersonTypeName = "person"

// Registers my person type to given L.
func registerPersonType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaPersonTypeName)
	L.SetGlobal(luaPersonTypeName, mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newPerson))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), personMethods))
}

// Constructor
func newPerson(L *lua.LState) int {
	person := &Person{L.CheckString(1)}
	ud := L.NewUserData()
	ud.Value = person
	L.SetMetatable(ud, L.GetTypeMetatable(luaPersonTypeName))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkPerson(L *lua.LState) *Person {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Person); ok {
		return v
	}
	L.ArgError(1, "person expected")
	return nil
}

var personMethods = map[string]lua.LGFunction{
	"name": personGetSetName,
}

// Getter and setter for the Person#Name
func personGetSetName(L *lua.LState) int {
	p := checkPerson(L)
	if L.GetTop() == 2 {
		p.Name = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(p.Name))
	return 1
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			println(r.(error).Error())
		}
	}()

	L := lua.NewState()
	defer L.Close()
	registerPersonType(L)

	err := L.DoFile("data/screens/main.lua")
	if err != nil {
		println(err.Error())
	}

	person := &Person{"ooxx"}
	ud := L.NewUserData()
	ud.Value = person
	println("(------")
	L.SetMetatable(ud, L.GetTypeMetatable(luaPersonTypeName))
	println("------)")

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("OnInitScreen"), // 调用的Lua函数
		NRet:    0,                           // 返回值的数量
		Protect: true,                        // 保护?
	}, ud.Metatable); err != nil {
		panic(err)
	}

}

// func main() {
// 	thread.GetMaster().Wait_thread_over()
// }
