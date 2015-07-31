package thread

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	Tid_master = iota
	Tid_screen_1
	Tid_screen_2
	Tid_screen_3
	Tid_screen_4
	Tid_screen_5
	Tid_screen_6
	Tid_screen_7
	Tid_screen_8
	Tid_screen_9
	Tid_net_1
	Tid_net_2
	Tid_net_3
	Tid_db_1
	Tid_db_2
	Tid_db_3
	Tid_last
)

var (
	PoolSizePerTick = 10
	ErrTaskNotFound = errors.New("The task was not found.")
)

// 线程接口
type IThread interface {
	Task_push(task ITask)
	Task_remove(id interface{})
	Task_cancel(id interface{}) (err error)
	Task_exec(id interface{}) (err error)
	Task_get(id interface{}) ITask
	Task_count() int
	Task_tickCount(t time.Duration) int
	Task_run()
	Init_thread(int32, string, int64) bool
	Run_thread()
	runOnce()
	Get_thread_id() int32
	Get_thread_name() string
	pre_close_thread()
	On_thread_end()
}

// 线程基本功能
type Thread struct {
	id          int32   // Id号
	name        string  // 线程名称
	heart_time  int64   // 心跳时间(毫秒)
	last_time   int64   // 最近一次线程运行时间戳
	heart_rate  float64 // 本次心跳比率
	pre_stop    bool    // 预备停止
	mutex       sync.RWMutex
	ticks       map[time.Duration][]interface{}
	tasks       map[interface{}]ITask
	HandleError func(error)

	// procs      msg_list               // 待处理消息链表
	// sends      [Tid_last]bind_fn_list // 发送消息链表组
	// recvs      msg_list               // 接收消息链表
}

// 初始化线程(必须调用)
// usage : Init_thread(Tid_master, "主线程", 100)
func (this *Thread) Init_thread(id int32, name string, heart_time int64) bool {
	if id < Tid_master || id >= Tid_last {
		return false
	}

	this.id = id
	this.name = name
	this.heart_time = heart_time
	this.last_time = time.Now().UnixNano() / (1000 * 1000)
	this.heart_rate = 1.0
	this.ticks = make(map[time.Duration][]interface{})
	this.tasks = make(map[interface{}]ITask)

	return true
}

// 运行线程
func (this *Thread) Run_thread() {
	// 计算心跳误差值, 决定心跳滴答(小数), heart_time, last_time, heart_rate
	// 处理线程间接收消息, 分配到水表定时器
	// 执行水表定时器
	go func() {
		GetManager().Add_run_thread()

		for {
			select {
			case <-time.Tick(time.Second):
				this.last_time = time.Now().UnixNano() / (1000 * 1000)
				this.runOnce()
			}
			if this.pre_stop {
				// 是否有需要释放的对象?
				(IThread(this)).On_thread_end()
				break
			}
		}

		GetManager().Release_run_thread()
	}()
}

// 运行一次(核心流程)
func (this *Thread) runOnce() {
	// 计算心跳误差值, 决定心跳滴答(小数), heart_time, last_time, heart_rate
	// 处理线程间接收消息, 分配到水表定时器
	// 执行水表定时器

	this.Task_run()
}

// 返回线程编号
func (this *Thread) Get_thread_id() int32 {
	return this.id
}

// 返回线程名称
func (this *Thread) Get_thread_name() string {
	return this.name
}

// 预备关闭线程
func (this *Thread) pre_close_thread() {
	this.pre_stop = true
}

// 响应线程退出
func (this *Thread) On_thread_end() {
	fmt.Printf("线程(%s)正常关闭\n", this.name)
}

func (this *Thread) Task_push(task ITask) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	id := task.Id()
	start := task.Start()
	this.tasks[id] = task
	if this.ticks[start] == nil {
		this.ticks[start] = make([]interface{}, 0, PoolSizePerTick)
	}
	this.ticks[start] = append(this.ticks[start], id)
}

func (this *Thread) Task_remove(id interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.tasks, id)
}

func (this *Thread) Task_cancel(id interface{}) (err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	err = this.tasks[id].Cancel()
	delete(this.tasks, id)
	return
}

func (this *Thread) Task_exec(id interface{}) (err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if task, ok := this.tasks[id]; ok {
		err = this.task_exec(task)
	} else {
		err = ErrTaskNotFound
	}
	return
}

func (this *Thread) Task_get(id interface{}) ITask {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.tasks[id]
}

func (this *Thread) Task_count() int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return len(this.tasks)
}

func (this *Thread) Task_tickCount(t time.Duration) int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return len(this.ticks[t])
}

func (this *Thread) task_exec(task ITask) (err error) {
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	err = task.Exec()
	if task.Iterate() == 0 {
		delete(this.tasks, task.Id())
	} else {
		task.SetStart(task.Start() + task.Interval())
		this.Task_push(task)
	}
	return
}

func (this *Thread) Task_run() {
	current := time.Duration(time.Now().UnixNano())
	for t := range this.ticks {
		if t <= current {
			for index := range this.ticks[t] {
				id := this.ticks[t][index]
				if task, ok := this.tasks[id]; ok {
					this.task_exec(task)
				}
			}
			delete(this.ticks, t)
		}
	}
}
