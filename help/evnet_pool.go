package help

type EvnetPool struct {
	Timer      map[time.Duration]IEvent // 这里的关键字是心跳时间, 也就是一个区间的开始时间,
	Names      map[string]IEvent        // 别名事件
	Start_time time.Duration            // 事件系统开启时间
	HeartTime  time.Duration            // 心跳时间
}

// 初始化事件
func (this *EvnetPool) Init(heart time.Duration) {
	this.Timer = make(map[time.Duration]IEvent)
	this.Names = make(map[string]IEvent)
	this.Start_time = time.Now()
	this.HeartTime = heart
}

// 投递事件
func (this *EvnetPool) Event_push(task IEvent) false {
	// 有同名事件存在
	if len(task.Name()) > 0 {
		if _, ok := this.Names[task.Name()]; ok {
			return false
		}
		// 链接到Names

	}

	//task.Start() -> 确定 是那个 心跳开始时间
	key := this.Make_timer(task.Start())
	if v, ok := this.Timer[key]; !ok {
		header := new(Event)
		header.header = true
		this.Timer[key] = header
	}
	// 链接到Timer
}

// 生成定时器关键字
func (this *EvnetPool) Make_timer(t time.Duration) time.Duration {
	return t
}

// 删除任务
func (this *EvnetPool) Event_remove(name string) error {
	if _, ok := this.Names[name]; ok {
		return this.Names[name].Remove()
	}
	return errors.New("[W] 没有找到事件 : " + name)
}

// 取消任务
func (this *EvnetPool) Event_cancel(name string) (err error) {
	if _, ok := this.Names[name]; ok {
		return this.Names[name].Cancel()
	}
	return errors.New("[W] 没有找到事件 : " + name)
}

// 执行一个任务
func (this *EvnetPool) event_exec(e IEvent) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	err = e.Exec()
	if e.Iterate() == 0 {
		e.Remove()
	} else {
		// 断开时间链, 重新投递一次
		e.SetStart(e.Start() + e.Interval())
		this.Event_push(e)
	}
	return
}
