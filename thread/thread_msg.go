package thread

import (
	"github.com/toophy/Gogame/help"
	"github.com/toophy/Gogame/jiekou"
	"sync"
)

var G_thread_msg_pool ThreadMsgPool

func init() {
	G_thread_msg_pool.Init()
}

// 线程间消息存放处
type ThreadMsgPool struct {
	lock   [jiekou.Tid_last]sync.RWMutex  // 每个线程的消息池有一个独立的读写锁
	header [jiekou.Tid_last]help.ListNode // 每个线程的消息池
}

// 初始化
func (this *ThreadMsgPool) Init() {
	for i := 0; i < jiekou.Tid_last; i++ {
		this.header[i].Init()
	}
}

// 投递线程间消息
func (this *ThreadMsgPool) PostMsg(tid int32, a *help.ListNode) bool {
	if a != nil {
		if tid >= jiekou.Tid_master && tid < jiekou.Tid_last {
			this.lock[tid].Lock()
			defer this.lock[tid].Unlock()

			header := &this.header[tid]

			a_pre := a.Pre
			a_next := a.Next

			a.Init()

			header.Pre.Next = a_pre
			a_pre.Pre = header.Pre

			header.Pre = a_next
			a_next.Next = header

			return true
		}
	}
	return false
}

// 获取线程间消息
func (this *ThreadMsgPool) GetMsg(tid int32, a *help.ListNode) bool {
	if a != nil {
		if tid >= jiekou.Tid_master && tid < jiekou.Tid_last {
			this.lock[tid].Lock()
			defer this.lock[tid].Unlock()

			header := &this.header[tid]

			header_pre := header.Pre
			header_next := header.Next

			header.Init()

			a.Pre.Next = header_pre
			header_pre.Pre = a.Pre

			a.Pre = header_next
			header_next.Next = a

			return true
		}
	}
	return false
}
