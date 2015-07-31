package actor

import (
//"fmt"
)

const (
	Atype_null = iota
	Atype_player
	Atype_monster
	Atype_npc
	Atype_pupet
	Atype_bus
	Atype_item
	Atype_last
)

const (
	Amdl_BaseAtr = iota
	Amdl_ExAtr
	Amdl_Last
)

// 演员 : 地图上所有对象
// usage :
// a := Actor{}
type Actor struct {
	Mdls map[int32]interface{}
	Id   int64
	Type int32
}

// 演员初始化
// usage :
// a := Actor{}
// a.Init(Atype_player)
func (a *Actor) Init(t int32, id int64) bool {
	if t <= Atype_null || t >= Atype_last {
		return false
	}
	a.Type = t
	a.Id = id
	a.Mdls = make(map[int32]interface{}, 0)
	return true
}

// 获取Id
func (a *Actor) GetId() int64 {
	return a.Id
}

// 获取类型
func (a *Actor) GetType() int32 {
	return a.Type
}

// 增加演员功能模块
// usage :
// a := Actor{}
// b := BaseAtr{}
// a.Mdl_add(&b)
//
func (a *Actor) Mdl_add(m interface{}) bool {
	id := int32(Amdl_Last)

	switch m.(type) {
	case *BaseAtr:
		id = Amdl_BaseAtr
	case *ExAtr:
		id = Amdl_ExAtr
	}

	if id != Amdl_Last {
		if _, ok := a.Mdls[id]; ok == false {
			a.Mdls[id] = m
			return true
		}
	}

	return false
}

// 删除演员功能模块
// usage :
// a := Actor{}
// emB :=  0
// a.Mdl_del(emB)
//
func (a *Actor) Mdl_del(id int32) {
	if _, ok := a.Mdls[id]; ok {
		delete(a.Mdls, id)
	}
}

// 检查演员功能模块是否存在, 返回结构指针
// usage :
// a := Actor{}
// emB :=  0
// v := a.Mdl_check(emB)
// if v != nil {
//   调用模块
//   v.(*BaseAtr).Name = name
// }
func (a *Actor) Mdl_check(id int32) interface{} {
	if v, ok := a.Mdls[id]; ok {
		return v
	}
	return nil
}
