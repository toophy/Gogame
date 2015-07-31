// main.go
package main

import (
	"Gogame/thread"
	"time"
)

func main() {
	sc1 := thread.New_screen_thread(thread.Tid_screen_1, "场景线程1", 100)
	if sc1 != nil {
		sc1.Run_thread()

		n := time.Duration(time.Now().UnixNano())
		sc1.Task_push(&thread.Event_open_screen{
			Task: thread.Task{
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

		sc1.Task_push(&thread.Event_open_screen{
			Task: thread.Task{
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

		sc1.Task_push(&thread.Event_close_thread{
			Task: thread.Task{
				Id_:       3,
				Start_:    n + 10*time.Second,
				Interval_: time.Second,
				Iterate_:  0,
			},
			Master: sc1,
		})
	}

	thread.GetMaster().Wait_thread_over()
}
