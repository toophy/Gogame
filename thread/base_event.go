package thread

import (
	"errors"
	"github.com/toophy/Gogame/help"
)

// 事件 : 线程关闭
type Event_close_thread struct {
	help.Task
	Master IThread
}

// 事件执行
func (t *Event_close_thread) Exec() error {
	if t.Master != nil {
		t.Master.pre_close_thread()
		return nil
	}

	return errors.New("没找到线程")
}
