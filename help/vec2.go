package help

import ()

type Vec2 struct {
	X float32
	Y float32
}

// 在 gopher-lua 文件 state.go 中增加 Require 函数

// func (ls *LState) Require(name string) LValue {
// 	loaded := ls.GetField(ls.Get(RegistryIndex), "_LOADED")
// 	lv := ls.GetField(loaded, name)
// 	if LVAsBool(lv) {
// 		if lv == loopdetection {
// 			ls.RaiseError("loop or previous error loading module: %s", name)
// 		}
// 		return lv
// 	}
// 	loaders, ok := ls.GetField(ls.Get(RegistryIndex), "_LOADERS").(*LTable)
// 	if !ok {
// 		ls.RaiseError("package.loaders must be a table")
// 	}
// 	messages := []string{}
// 	var modasfunc LValue
// 	for i := 1; ; i++ {
// 		loader := ls.RawGetInt(loaders, i)
// 		if loader == LNil {
// 			ls.RaiseError("module %s not found:\n\t%s, ", name, strings.Join(messages, "\n\t"))
// 		}
// 		ls.Push(loader)
// 		ls.Push(LString(name))
// 		ls.Call(1, 1)
// 		ret := ls.reg.Pop()
// 		switch retv := ret.(type) {
// 		case *LFunction:
// 			modasfunc = retv
// 			goto loopbreak
// 		case LString:
// 			messages = append(messages, string(retv))
// 		}
// 	}
// loopbreak:
// 	ls.SetField(loaded, name, loopdetection)
// 	ls.Push(modasfunc)
// 	ls.Push(LString(name))
// 	ls.Call(1, 1)
// 	ret := ls.reg.Pop()
// 	modv := ls.GetField(loaded, name)
// 	if ret != LNil && modv == loopdetection {
// 		ls.SetField(loaded, name, ret)
// 		return ret
// 	} else if modv == loopdetection {
// 		ls.SetField(loaded, name, LTrue)
// 		return LTrue
// 	} else {
// 		return modv
// 	}
// 	return LNumber(1)
// }

// 定制lua脚本
// 禁止以下功能
// dofile
// require
// load
// loadfile
// loadstring
