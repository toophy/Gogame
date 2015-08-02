package thread

import (
	"errors"
	//"fmt"
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
	Init_thread(IThread, int32, string, int64) error // 初始化线程
	Run_thread()                                     // 运行线程
	Get_thread_id() int32                            // 获取线程ID
	Get_thread_name() string                         // 获取线程名称
	pre_close_thread()                               // 预备关闭线程
	on_first_run()                                   // 首次运行(在 on_run 前面)
	on_run()                                         // 线程运行部分
	on_end()                                         // 线程结束回调

	Task_push(task ITask)                   // 任务推送
	Task_remove(id interface{})             // 任务删除
	Task_cancel(id interface{}) (err error) // 任务取消
}

// 线程基本功能
// procs      msg_list               // 待处理消息链表
// sends      [Tid_last]bind_fn_list // 发送消息链表组
// recvs      msg_list               // 接收消息链表
type Thread struct {
	id          int32                           // Id号
	name        string                          // 线程名称
	heart_time  int64                           // 心跳时间(毫秒)
	last_time   int64                           // 最近一次线程运行时间戳
	heart_rate  float64                         // 本次心跳比率
	pre_stop    bool                            // 预备停止
	mutex       sync.RWMutex                    // 任务读写锁
	ticks       map[time.Duration][]interface{} // 任务定时器
	tasks       map[interface{}]ITask           // 任务列表
	HandleError func(error)                     // 任务异常错误
	self        IThread                         // 自己, 初始化之后, 不要操作
	first_run   bool                            // 线程首次运行
}

// 初始化线程(必须调用)
// usage : Init_thread(Tid_master, "主线程", 100)
func (this *Thread) Init_thread(self IThread, id int32, name string, heart_time int64) error {
	if id < Tid_master || id >= Tid_last {
		return errors.New("[E] 线程ID超出范围 [Tid_master,Tid_last]")
	}
	if self == nil {
		return errors.New("[E] 线程自身指针不能为nil")
	}

	this.id = id
	this.name = name
	this.heart_time = heart_time
	this.last_time = time.Now().UnixNano() / (1000 * 1000)
	this.heart_rate = 1.0
	this.ticks = make(map[time.Duration][]interface{})
	this.tasks = make(map[interface{}]ITask)
	this.self = self
	this.first_run = true

	return nil
}

// 运行线程
func (this *Thread) Run_thread() {
	// 计算心跳误差值, 决定心跳滴答(小数), heart_time, last_time, heart_rate
	// 处理线程间接收消息, 分配到水表定时器
	// 执行水表定时器
	go func() {
		GetMaster().Add_run_thread(this.self)

		this.self.on_first_run()

		for {
			select {
			case <-time.Tick(time.Second):
				this.last_time = time.Now().UnixNano() / (1000 * 1000)
				this.runOnce()
				this.self.on_run()
			}
			if this.pre_stop {
				// 是否有需要释放的对象?
				this.self.on_end()
				break
			}
		}

		GetMaster().Release_run_thread(this.self)
	}()
}

// 运行一次(核心流程)
func (this *Thread) runOnce() {
	// 计算心跳误差值, 决定心跳滴答(小数), heart_time, last_time, heart_rate
	// 处理线程间接收消息, 分配到水表定时器
	// 执行水表定时器
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

// 投递任务
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

// 删除任务
func (this *Thread) Task_remove(id interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.tasks, id)
}

// 取消任务
func (this *Thread) Task_cancel(id interface{}) (err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	err = this.tasks[id].Cancel()
	delete(this.tasks, id)
	return
}

// 执行一个任务
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
