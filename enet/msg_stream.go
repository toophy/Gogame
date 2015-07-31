package enet

import (
	"fmt"
)

// 网络消息体读写
type Msg_stream struct {
	msg *Msg_packet
	pos int
}

func (t *Msg_stream) InitMsgStream(msg *Msg_packet) {
	t.pos = 0
	t.msg = msg
}

func (t *Msg_stream) Print_data() {
	fmt.Println(t.msg.Data[:t.msg.Len+2])
}

func (t *Msg_stream) Seek(pos int) {
	if t.msg == nil {
		t.pos = 0
	} else if pos >= 0 && pos < t.msg.Len {
		t.pos = pos
	}
}

func (t *Msg_stream) Read_u1() int {
	if (t.pos + 1) < (t.msg.Len + 1) {
		old_pos := t.pos
		t.pos = t.pos + 1
		return int(t.msg.Data[old_pos])
	}
	return 0
}

func (t *Msg_stream) Read_u2() int {
	if (t.pos + 2) < (t.msg.Len + 1) {
		old_pos := t.pos
		t.pos = t.pos + 2
		return int(t.msg.Data[old_pos])<<8 + int(t.msg.Data[old_pos+1])
	}
	return 0
}

func (t *Msg_stream) Read_u4() int {
	if (t.pos + 4) < (t.msg.Len + 1) {
		old_pos := t.pos
		t.pos = t.pos + 4
		return (int(t.msg.Data[old_pos]) << 24) +
			(int(t.msg.Data[old_pos+1]) << 16) +
			(int(t.msg.Data[old_pos+2]) << 8) +
			(int(t.msg.Data[old_pos+3]))
	}
	return 0
}

func (t *Msg_stream) Read_u8() int {
	if (t.pos + 8) < (t.msg.Len + 1) {
		old_pos := t.pos
		t.pos = t.pos + 8
		return (int(t.msg.Data[old_pos]) << 24) +
			(int(t.msg.Data[old_pos+1]) << 16) +
			(int(t.msg.Data[old_pos+2]) << 8) +
			(int(t.msg.Data[old_pos+3]))
	}
	return 0
}

func (t *Msg_stream) Read_string() string {
	data_len := t.Read_u2()
	if data_len > 0 && (t.pos+data_len) < (t.msg.Len+1) {
		old_pos := t.pos
		t.pos = t.pos + data_len
		return string(t.msg.Data[old_pos : old_pos+data_len])
	}
	return ""
}

func (t *Msg_stream) Write_u1(d int) bool {
	if t.pos+1 < MaxDataLen {
		t.msg.Data[t.pos] = byte(d & 0xFF)
		t.pos = t.pos + 1
		t.msg.Len = t.msg.Len + 1
		return true
	}

	return false
}

func (t *Msg_stream) Write_u2(d int) bool {
	if t.pos+2 < MaxDataLen {
		// 65280
		t.msg.Data[t.pos] = byte((d & 0xFF00) >> 8)
		//
		t.msg.Data[t.pos+1] = byte(d & 0xFF)
		t.pos = t.pos + 2
		t.msg.Len = t.msg.Len + 2
		return true
	}

	return false
}

func (t *Msg_stream) Write_u4(d int) bool {
	nd := uint(d)
	if t.pos+4 < MaxDataLen {
		// 4278190080
		t.msg.Data[t.pos] = byte((nd & 0xFF000000) >> 24)
		// 16711680
		t.msg.Data[t.pos+1] = byte((nd & 0xFF0000) >> 16)
		// 65280
		t.msg.Data[t.pos+2] = byte((nd & 0xFF00) >> 8)
		//
		t.msg.Data[t.pos+3] = byte(nd & 0xFF)
		t.pos = t.pos + 4
		t.msg.Len = t.msg.Len + 4
		return true
	}

	return false
}

func (t *Msg_stream) Write_string(d *string) bool {
	d_len := len(*d)

	if t.pos+2+d_len < MaxDataLen {
		if t.Write_u2(d_len) {
			ds := (*d)[:]
			dx := t.msg.Data[t.pos : t.pos+d_len]
			copy(dx, ds)
			t.pos = t.pos + d_len
			t.msg.Len = t.msg.Len + d_len
			return true
		}
	}

	return false
}
