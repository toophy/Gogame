package screen

import (
	"Gogame/actor"
	"Gogame/help"
	"fmt"
)

type Screen struct {
	Name   string
	Id     int32
	Oid    int32
	Actors map[int64]*actor.Actor
}

func (this *Screen) Load(name string, id int32, oid int32) bool {
	config := screen_config.GetScreenConfig(oid)
	if config == nil {
		fmt.Printf("场景%s加载失败: 没有找到场景模板(%d)\n", name, oid)
		return false
	}

	if len(name) > 0 {
		this.Name = name
	} else {
		this.Name = config.Name
	}

	this.Id = id
	this.Oid = oid
	this.Actors = make(map[int64]*actor.Actor, 0)
	fmt.Printf("场景%s加载成功\n", name)

	return true
}

// 场景卸载
// 场景关联的定时器, 事件, 统统要卸载掉
// 场景内的精灵呢? 有些定时器, 事件, 也是场景关联的
// 没有场景的精灵怎么进行操作呢?
func (this *Screen) Unload() {
	fmt.Printf("场景%s卸载成功\n", this.Name)
}

// 登录地图
func (this *Screen) Actor_enter(a *actor.Actor) bool {
	if _, ok := this.Actors[a.GetId()]; ok {
		this.Actors[a.GetId()] = a
		return true
	}
	return false
}

// 离开地图
func (this *Screen) Actor_leave(id int64) bool {
	if _, ok := this.Actors[id]; ok {
		delete(this.Actors, id)
		return true
	}
	return false
}

// 角色移动
func (this *Screen) Actor_move(id int64, pos help.Vec3, check bool) {
	// 如果 check 为 true
	// 主要是位置检查
	// 1. 边界
	// 2. 障碍检查
	// 3. Actor碰撞检查
	// 否则
	// 1. 边界检查
}

// 角色移动验证
func (this *Screen) Actor_move_check(id int64, pos help.Vec3) {
	// 主要是位置检查
	// 1. 边界
	// 2. 障碍检查
	// 3. Actor碰撞检查
}
