package thread

import (
	"bytes"
	"fmt"
	"github.com/toophy/Gogame/jiekou"
)

const (
	LogBuffSize = 32 * 4096
)

// 场景线程
type LogThread struct {
	Thread

	Buffs bytes.Buffer // 日志总缓冲
}

// 新建场景线程
func New_log_thread(heart_time int64, lay1_time uint64) (*LogThread, error) {
	a := new(LogThread)
	err := a.Init_log_thread(heart_time, lay1_time)
	if err == nil {
		return a, nil
	}
	return nil, err
}

// 初始化场景线程
func (this *LogThread) Init_log_thread(heart_time int64, lay1_time uint64) error {
	err := this.Init_thread(this, jiekou.Tid_log, "log_thread", heart_time, lay1_time)
	if err == nil {
		return nil
	}
	return err
}

// 响应线程首次运行
func (this *LogThread) on_first_run() {
	// 处理文件
	evt := &Event_flush_log{}
	evt.Init("", 300)
	this.PostEvent(evt)
}

// 响应线程退出
func (this *LogThread) on_end() {
}

// 响应线程运行
func (this *LogThread) on_run() {
}

func (this *LogThread) Add_log(d bytes.Buffer) {
	this.Buffs.Write(d.Bytes())
}

func (this *LogThread) Flush_log() {
	fmt.Print(this.Buffs.String())
	this.Buffs.Reset()
}
