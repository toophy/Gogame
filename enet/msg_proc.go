package enet

// 网络消息包处理
// 暂不包含加密/解密

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	MaxDataLen     = 5080
	MaxSendDataLen = 4000
	MaxHeader      = 2
)

// 网络消息体
type Msg_packet struct {
	Data []byte
	Len  int
}

func (t *Msg_packet) Init() {
	t.Len = 0
	t.Data = make([]byte, MaxDataLen)
}

func (t *Msg_packet) Print_data() {
	fmt.Println(t.Data[:t.Len+2])
}

func (t *Msg_packet) Read_data(conn *net.TCPConn) error {

	t.Len = 0
	length, err := io.ReadFull(conn, t.Data[:2])
	if length != MaxHeader {
		fmt.Printf("Packet header : %d != %d\n", length, MaxHeader)
		return err
	}
	if err != nil {
		return err
	}

	body_len := int(t.Data[1]) + (int(t.Data[0]) << 8)

	if body_len > (MaxDataLen - 2) {
		err = errors.New("Body too much")
		return err
	}

	t.Len = body_len + 2
	return t.Read_body(conn)
}

func (t *Msg_packet) Read_body(conn *net.TCPConn) error {

	length, err := io.ReadFull(conn, t.Data[2:t.Len])
	if length != (t.Len - 2) {
		fmt.Printf("Packet length : %d != %d \n", length, t.Len-2)
		return err
	}
	if err != nil {
		return err
	}
	// 注意 : 可以解密

	return nil
}

func (t *Msg_packet) Send(conn *net.TCPConn) error {
	if t.Len > MaxHeader && t.Len < MaxSendDataLen {

		t.Data[0] = byte((t.Len & 0xFF00) >> 8)
		t.Data[1] = byte(t.Len & 0xFF)

		_, err := conn.Write(t.Data[:MaxHeader+t.Len])
		if err != nil {
			fmt.Printf(err.Error() + "\n")
			return err
		}
	}

	return nil
}
