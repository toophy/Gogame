package thread

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Manager struct {
	RunThreadCount int32
}

var myManager *Manager = nil

func GetManager() *Manager {
	if myManager == nil {
		myManager = &Manager{
			RunThreadCount: 0,
		}
	}
	return myManager
}

// 增加运行的线程
func (this *Manager) Add_run_thread() {
	atomic.AddInt32(&this.RunThreadCount, 1)
}

func (this *Manager) Release_run_thread() {
	atomic.AddInt32(&this.RunThreadCount, -1)
}

// 等待所有线程结束
func (this *Manager) Wait_thread_over() {
	for {
		select {
		case <-time.Tick(5 * time.Second):
			opsFinal := atomic.LoadInt32(&this.RunThreadCount)
			if opsFinal <= 0 {
				fmt.Println("所有线程都已经停止")
				time.Sleep(2 * time.Second)
				return
			} else {
				//fmt.Printf("有%d个线程正在运行\n", opsFinal)
			}
		}
	}
}
