package thread

import (
	//"fmt"
	"sync"
	"time"
)

// 主线程
type Master struct {
	Thread
	threadLock  sync.RWMutex
	threadCount int32
	threadIds   map[int32]IThread
}

var myMaster *Master = nil

func GetMaster() *Master {
	if myMaster == nil {
		myMaster = &Master{}
		if !myMaster.Init_master_thread(myMaster, "主线程", 100) {
			return nil
		}
		myMaster.Run_thread()
	}
	return myMaster
}

func (this *Master) Init_master_thread(self IThread, name string, heart_time int64) bool {
	if this.Init_thread(self, Tid_master, name, heart_time) {
		this.threadCount = 0
		this.threadIds = make(map[int32]IThread, 0)
		return true
	}
	return false
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

func (this *Master) Release_run_thread(a IThread) {
	this.threadLock.Lock()
	defer this.threadLock.Unlock()

	if _, ok := this.threadIds[a.Get_thread_id()]; ok == true {
		this.threadCount--
		delete(this.threadIds, a.Get_thread_id())
	}
}

// 首次运行
func (this *Master) on_first_run() {
}

// 响应线程退出
func (this *Master) on_end() {
}

// 响应线程运行
func (this *Master) on_run() {
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
