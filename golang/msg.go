package netext

import (
	"encoding/binary"
	"github.com/intmian/mian_go_lib/tool/misc"
)

type MsgFlag uint16

const (
	MsgFlagNull      MsgFlag = 0
	MsgFlagNeedReply MsgFlag = 1 << iota
)

type Msg struct {
	flag     MsgFlag // 用于标记消息的类型
	recallID uint32  // 如果是有回调的消息，服务器正常注册路由，但是返还时，需要将回调ID一起返回
	cmd      uint16  // 命令，用于区分不同的消息，小于100的为系统保留命令
	data     []byte
}

func (m *Msg) FromBytes(data []byte) error {
	if len(data) < 4 {
		return ErrMsgDataTooShort
	}
	m.flag = MsgFlag(binary.BigEndian.Uint16(data[0:2]))
	if misc.HasProperty(m.flag, MsgFlagNeedReply) {
		if len(data) < 8 {
			return ErrMsgDataTooShort
		}
		m.recallID = binary.BigEndian.Uint32(data[2:6])
		m.cmd = binary.BigEndian.Uint16(data[6:8])
		m.data = data[8:]
	} else {
		m.cmd = binary.BigEndian.Uint16(data[2:4])
		m.data = data[4:]
	}
	return nil
}

func (m *Msg) ToBytes() []byte {
	var data []byte
	if misc.HasProperty(m.flag, MsgFlagNeedReply) {
		data = make([]byte, 8+len(m.data))
		binary.BigEndian.PutUint16(data[0:2], uint16(m.flag))
		binary.BigEndian.PutUint32(data[2:6], m.recallID)
		binary.BigEndian.PutUint16(data[6:8], m.cmd)
		copy(data[8:], m.data)
	} else {
		data = make([]byte, 4+len(m.data))
		binary.BigEndian.PutUint16(data[0:2], uint16(m.flag))
		binary.BigEndian.PutUint16(data[2:4], m.cmd)
		copy(data[4:], m.data)
	}
	return data
}
