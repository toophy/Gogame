package screen

import (
	"fmt"
	"github.com/toophy/Gogame/actor"
	"github.com/toophy/Gogame/event"
	"github.com/toophy/Gogame/help"
	"github.com/toophy/Gogame/jiekou"
	lua "github.com/toophy/gopher-lua"
)

type Screen struct {
	event.EventObj
	Name    string
	Id      int32
	Oid     int32
	Actors  map[int64]*actor.Actor
	luaData lua.LValue
	thread  jiekou.IScreenThread
}

func (this *Screen) Load(name string, id int32, oid int32, t jiekou.IScreenThread) bool {
	this.InitEventHeader()
	config := screen_config.GetScreenConfig(oid)
	if config == nil {
		fmt.Printf("场景%s加载失败: 没有找到场景模板(%d)\n", name, oid)
		return false
	}
	if t == nil {
		fmt.Println("场景线程不存在")
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
	this.thread = t
	this.luaData = this.thread.GetLuaState().NewTable()
	fmt.Printf("场景%s加载成功\n", this.Name)
	this.Tolua_screen_init()

	evt := &Event_heart_beat{}
	evt.Init("", 3000)
	evt.Screen_ = this
	this.thread.PostEvent(evt)

	return true
}

// 场景卸载
// 场景关联的定时器, 事件, 统统要卸载掉
// 场景内的精灵呢? 有些定时器, 事件, 也是场景关联的
// 没有场景的精灵怎么进行操作呢?
func (this *Screen) Unload() {
	fmt.Printf("场景%s卸载成功\n", this.Name)
	this.thread.RemoveEventList(this.GetEventHeader())
	this.thread = nil
	this.luaData = nil
}

// 获取场景管理luaTable
func (this *Screen) Get_data() lua.LValue {
	return this.luaData
}

// !!! 获取的指针不能保存, 获取场景配置
func (this *Screen) Get_config() *Config {
	return screen_config.GetScreenConfig(this.Oid)
}

// 获取所在线程
func (this *Screen) Get_thread() jiekou.IScreenThread {
	return this.thread
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
