package thread

import (
	"fmt"
	"github.com/toophy/Gogame/event"
	"sync"
)

var G_thread_msg_pool ThreadMsgPool

func init() {
	G_thread_msg_pool.Init()
}

// 线程间消息存放处
type ThreadMsgPool struct {
	lock   [Tid_last]sync.RWMutex // 每个线程的消息池有一个独立的读写锁
	header [Tid_last]event.IEvent // 每个线程的消息池
}

// 初始化
func (this *ThreadMsgPool) Init() {
	for i := 0; i < Tid_last; i++ {
		this.header[i] = new(event.EventHeader)
		this.header[i].Init("", 100)
	}
}

// 投递线程间消息
func (this *ThreadMsgPool) PostMsg(tid int32, a event.IEvent) bool {
	if !a.IsHeader() || a.IsEmpty() {
		return false
	}
	if tid >= Tid_master && tid < Tid_last {
		this.lock[tid].Lock()
		defer this.lock[tid].Unlock()

		header := this.header[tid]

		a_pre := a.GetPreTimer()
		a_next := a.GetNextTimer()

		a.SetPreTimer(a)
		a.SetNextTimer(a)

		header.GetPreTimer().SetNextTimer(a_pre)
		a_pre.SetPreTimer(header.GetPreTimer())

		header.SetPreTimer(a_next)
		a_next.SetNextTimer(header)

		fmt.Printf("    PostMsg -> %d\n", tid)

		return true
	}
	return false
}

// 获取线程间消息
func (this *ThreadMsgPool) GetMsg(tid int32, a event.IEvent) bool {
	if !a.IsHeader() {
		return false
	}
	if tid >= Tid_master && tid < Tid_last {
		this.lock[tid].Lock()
		defer this.lock[tid].Unlock()

		header := this.header[tid]

		header_pre := header.GetPreTimer()
		header_next := header.GetNextTimer()

		header.SetPreTimer(header)
		header.SetNextTimer(header)

		a.GetPreTimer().SetNextTimer(header_pre)
		header_pre.SetPreTimer(a.GetPreTimer())

		a.SetPreTimer(header_next)
		header_next.SetNextTimer(a)

		return true
	}
	return false
}
