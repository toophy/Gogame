package help

// 事件头
type ListNode struct {
	Pre  *ListNode   // 前一个
	Next *ListNode   // 后一个
	Data interface{} // 数据
}

func (this *ListNode) Init() {
	this.Pre = this
	this.Next = this
}

func (this *ListNode) IsHeader() bool {
	return (this.Pre == this.Next)
}

func (this *ListNode) IsEmpty() bool {
	return (this.Pre == this.Next)
}
