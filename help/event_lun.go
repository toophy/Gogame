package help

const (
	Evt_uc         = 32                     //
	Evt_list       = 5                      //
	Evt_lay_1_bit  = 5 * 0                  // 移位
	Evt_lay_2_bit  = 5 * 1                  //
	Evt_lay_3_bit  = 5 * 2                  //
	Evt_lay_4_bit  = 5 * 3                  //
	Evt_lay_5_bit  = 5 * 4                  //
	Evt_lay_1_max  = Evt_uc                 // 上限
	Evt_lay_2_max  = Evt_uc * Evt_lay_1_max //
	Evt_lay_3_max  = Evt_uc * Evt_lay_2_max //
	Evt_lay_4_max  = Evt_uc * Evt_lay_3_max //
	Evt_lay_5_max  = Evt_uc * Evt_lay_4_max //
	Evt_lay_1_mask = 0                      //掩码
	Evt_lay_2_mask = Evt_lay_1_max - 1      //
	Evt_lay_3_mask = Evt_lay_2_max - 1      //
	Evt_lay_4_mask = Evt_lay_3_max - 1      //
	Evt_lay_5_mask = Evt_lay_4_max - 1      //
	Evt_lay_1_sn   = 5                      //
	Evt_lay_2_sn   = 4                      //
	Evt_lay_3_sn   = 3                      //
	Evt_lay_4_sn   = 2                      //
	Evt_lay_5_sn   = 1                      //
	Evt_max_time   = Evt_lay_5_sn           // 最大时间单位
)

type Event_list struct {
	Headers    [Evt_uc]IEvent // 事件链表头
	CurrHead   int32          // 当前列表头下标
	RemainTime int32          // 剩余时间
	RestTime   int32          // 重置时间
}

type EvnetLun struct {
	Lists     [Evt_list]Event_list // 事件轮盘
	Names     map[string]IEvent    // 别名事件
	TickCount int64                // 点滴
	UnitTime  int64                // ?
}

func (this *EvnetLun) Init() {
	this.TickCount = 0
	this.UnitTime = 100

	this.initEventList(Evt_lay_5_sn, Evt_lay_5_mask)
	this.initEventList(Evt_lay_4_sn, Evt_lay_4_mask)
	this.initEventList(Evt_lay_3_sn, Evt_lay_3_mask)
	this.initEventList(Evt_lay_2_sn, Evt_lay_2_mask)
	this.initEventList(Evt_lay_1_sn, Evt_lay_1_mask)

	this.Names = make(map[string]IEvent)
}

func (this *EvnetLun) initEventList(sn, mask int) {

	if sn >= Evt_lay_5_sn || sn <= Evt_lay_1_sn {
		this.Lists[sn-1].CurrHead = -1
		this.Lists[sn-1].RestTime = mask + 1
		this.Lists[sn-1].RemainTime = mask + 1
		for i := 0; i < Evt_uc; i++ {
			this.Lists[sn-1].Headers[i].SetHeader()
		}
	}
}

func (this *EvnetLun) Setup(unit_time int64) {
	if unit_time > 0 && unit_time < Evt_max_time {
		this.UnitTime = unit_time
	}
}

// 投递事件
func (this *EvnetPool) Push(e IEvent) bool {
	// 有同名事件存在
	if len(e.Name()) > 0 {
		if _, ok := this.Names[e.Name()]; ok {
			return false
		}

	}

	time := e.GetStart()
	ret := false

	if time < Evt_lay_1_max {
		ret = this.PushToList(e, time>>Evt_lay_1_bit, 0, Evt_lay_1_sn)
	} else if time < Evt_lay_2_max {
		ret = this.PushToList(e, time>>Evt_lay_2_bit, time&Evt_lay_2_mask-this.Lists[Evt_lay_2_sn-1].RemainTime, Evt_lay_2_sn)
	} else if time < Evt_lay_3_max {
		ret = this.PushToList(3, time>>Evt_lay_3_bit, time&Evt_lay_3_mask-this.Lists[Evt_lay_3_sn-1].RemainTime, Evt_lay_3_sn)
	} else if time < Evt_lay_4_max {
		ret = this.PushToList(3, time>>Evt_lay_4_bit, time&Evt_lay_4_mask-this.Lists[Evt_lay_4_sn-1].RemainTime, Evt_lay_4_sn)
	} else if time < Evt_lay_5_max {
		ret = this.PushToList(3, time>>Evt_lay_5_bit, time&Evt_lay_5_mask-this.Lists[Evt_lay_5_sn-1].RemainTime, Evt_lay_5_sn)
	}

	if ret {
		this.Names[e.Name()] = e
	}

	return ret
}

// 投递事件给时间轮盘
func (this *EvnetLun) PushToList(e IEvent, list_head int32, remain_time int64, list_id int32) {
	list_head = this.Lists[list_id-1].CurrHead + list_head

	if list_head > (Evt_uc - 1) {
		list_head = list_head - Evt_uc
	}

	e.Time_push(this.Lists[list_id-1].Headers[list_head])
}
