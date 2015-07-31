package thread

import (
	"errors"
	"fmt"
)

type Event_close_thread struct {
	Task
	Master IThread
}

func (t *Event_close_thread) Exec() error {
	if t.Master != nil {
		fmt.Printf("预备关闭线程%s\n", t.Master.Get_thread_name())
		t.Master.pre_close_thread()
		return nil
	}

	return errors.New("没找到线程")
}
