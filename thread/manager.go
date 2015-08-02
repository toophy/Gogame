package thread

import (
	//"fmt"
	"errors"
	lua "github.com/toophy/gopher-lua"
	"sync"
	"time"
)

// 主线程
type Master struct {
	Thread
	threadLock  sync.RWMutex
	threadCount int32
	threadIds   map[int32]IThread
	luaState    *lua.LState
}

var myMaster *Master = nil

// 获取主线程
func GetMaster() *Master {
	if myMaster == nil {
		myMaster = &Master{}
		err := myMaster.Init_master_thread(myMaster, "主线程", 100)
		if err != nil {
			panic(err.Error())
		}
		myMaster.Run_thread()
	}
	return myMaster
}

// 初始化主线程
func (this *Master) Init_master_thread(self IThread, name string, heart_time int64) error {
	err := this.Init_thread(self, Tid_master, name, heart_time)
	if err == nil {
		this.threadCount = 0
		this.threadIds = make(map[int32]IThread, 0)

		return nil
	}
	return err
}

// 增加运行的线程
func (this *Master) Add_run_thread(a IThread) {
	this.threadLock.Lock()
	defer this.threadLock.Unlock()

	if _, ok := this.threadIds[a.Get_thread_id()]; ok == false {
		this.threadCount++
		this.threadIds[a.Get_thread_id()] = a
	}
}

// 释放运行的线程
func (this *Master) Release_run_thread(a IThread) {
	this.threadLock.Lock()
	defer this.threadLock.Unlock()

	if _, ok := this.threadIds[a.Get_thread_id()]; ok == true {
		this.threadCount--
		delete(this.threadIds, a.Get_thread_id())
	}
}

// 等待所有线程结束
func (this *Master) Wait_thread_over() {
	for {
		select {
		case <-time.Tick(10 * time.Second):
			this.threadLock.Lock()

			if this.threadCount <= 0 {
				this.threadLock.Unlock()
				time.Sleep(2 * time.Second)
				return
			} else if this.threadCount == 1 {
				n := time.Duration(time.Now().UnixNano())
				this.Task_push(&Event_close_thread{
					Task: Task{
						Id_:       3,
						Start_:    n + 2*time.Second,
						Interval_: time.Second,
						Iterate_:  0,
					},
					Master: this,
				})
			}

			this.threadLock.Unlock()
		}
	}
}

// 首次运行
func (this *Master) on_first_run() {

	errInit := this.ReloadLuaState()
	if errInit != nil {
		panic(errInit.Error())
	}

	sc1, err := New_screen_thread(Tid_screen_1, "场景线程1", 100)
	if err == nil && sc1 != nil {
		sc1.Run_thread()

		n := time.Duration(time.Now().UnixNano())
		sc1.Task_push(&Event_open_screen{
			Task: Task{
				Id_:       1,
				Start_:    n + time.Second,
				Interval_: time.Second,
				Iterate_:  0,
			},
			Screen_oid_:    1,
			Screen_name_:   "",
			Screen_thread_: sc1,
			Open:           true,
		})

		sc1.Task_push(&Event_open_screen{
			Task: Task{
				Id_:       2,
				Start_:    n + 5*time.Second,
				Interval_: time.Second,
				Iterate_:  0,
			},
			Screen_oid_:    1,
			Screen_name_:   "",
			Screen_thread_: sc1,
			Open:           false,
		})

		sc1.Task_push(&Event_close_thread{
			Task: Task{
				Id_:       3,
				Start_:    n + 10*time.Second,
				Interval_: time.Second,
				Iterate_:  0,
			},
			Master: sc1,
		})
	} else {
		if err != nil {
			println("[E] 新建场景线程失败:" + err.Error())
		} else {
			println("[E] 新建场景线程失败:")
		}
	}
}

// 响应线程退出
func (this *Master) on_end() {
	if this.luaState != nil {
		this.luaState.Close()
		this.luaState = nil
	}
}

// 响应线程运行
func (this *Master) on_run() {
}

// 初始化LuaState, 可以用来 Reload LuaState
func (this *Master) ReloadLuaState() error {

	if this.luaState != nil {
		this.luaState.Close()
		this.luaState = nil
	}

	this.luaState = lua.NewState()
	if this.luaState == nil {
		return errors.New("[E] 主线程初始化Lua失败")
	}

	// RegLua_all(this.luaState)

	// Require所有 master 文件夹里面的 *.lua 文件

	return nil
}
