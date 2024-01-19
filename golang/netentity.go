package netext

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"sync"
)

/*
netEntityCtx 调用方相对于netEntity转入的上下文
ctx 用于控制附属协程的生命周期
onRec 用于处理接收到的数据。请注意回调是并发的
onErr 用于处理错误。回调是并发的
*/
type netEntityCtx struct {
	ctx   context.Context
	onRec func([]byte, int)
	onErr func(error)
	//TODO: 如果性能或者gc出现瓶颈考虑
	//onGetNetPackBytes func() []byte // 从上层获取一个可用的[]byte，等处理以后会通过回调返回
	//onPutNetPackBytes func([]byte)  // 如果出现意外，需要将[]byte归还给上层
}

type netEntitySetting struct {
	netPackSize int
}

/*
netEntity 管理一个网络实体的链接读写。
Init 后才能使用，不可重复调用
AddConn 新增连接
DelConn 关闭连接
Send 发送数据
*/
type netEntity struct {
	entityMap map[ConnectType]net.Conn
	cancelMap map[ConnectType]func()
	lock      sync.Mutex
	netEntityCtx
	netEntitySetting
}

func (n *netEntity) Init(ctx netEntityCtx, s netEntitySetting) {
	n.entityMap = make(map[ConnectType]net.Conn)
	n.cancelMap = make(map[ConnectType]func())
	n.netEntityCtx = ctx
	n.netEntitySetting = s
}

func newNetEntity(ctx netEntityCtx, s netEntitySetting) *netEntity {
	n := &netEntity{}
	n.Init(ctx, s)
	return n
}

func (n *netEntity) AddConn(connType ConnectType, conn net.Conn) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.entityMap[connType]; ok {
		return ErrConnTypeAlreadyExist
	}
	n.entityMap[connType] = conn
	ctx, cancel := context.WithCancel(context.Background())
	n.cancelMap[connType] = cancel
	n.goRead(ctx, connType, conn)
	return nil
}

func (n *netEntity) DelConn(connType ConnectType) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.entityMap[connType]; !ok {
		return ErrConnTypeNotExist
	}
	n.cancelMap[connType]()
	delete(n.entityMap, connType)
	delete(n.cancelMap, connType)
	return nil
}

func (n *netEntity) Send(connType ConnectType, data []byte) error {
	if _, ok := n.entityMap[connType]; !ok {
		return ErrConnTypeNotExist
	}
	if connType == ConnTypeTcp {
		err := binary.Write(n.entityMap[connType], binary.BigEndian, uint16(len(data)))
		if err != nil {
			n.onErr(errors.Join(err, ErrWriteSizeFailed))
		}
	}
	_, err := n.entityMap[connType].Write(data)
	if err != nil {
		n.onErr(errors.Join(err, ErrWriteDataFailed))
	}
	return nil
}

func (n *netEntity) goRead(ctx context.Context, connType ConnectType, conn net.Conn) {
	go func() {
		sizePack := make([]byte, 2)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if n.read(connType, conn, sizePack) {
					return
				}
			}
		}
	}()
}

func (n *netEntity) read(connType ConnectType, conn net.Conn, sizePack []byte) bool {
	if connType == ConnTypeTcp {
		// TCP 需要处理粘包的问题
		sizeBytes, err := conn.Read(sizePack)
		if sizeBytes != 2 {
			n.onErr(ErrReadSizeFailed)
			return true
		}
		if err != nil {
			n.onErr(errors.Join(err, ErrReadSizeFailed))
			return true
		}
		size := int(binary.BigEndian.Uint16(sizePack))
		netPack := make([]byte, n.netEntitySetting.netPackSize)
		sizeBytes, err = conn.Read(netPack)
		if size != sizeBytes {
			n.onErr(ErrReadNetpackFailed)
		}
		if err != nil {
			n.onErr(errors.Join(err, ErrReadNetpackFailed))
			return true
		}
		n.onRec(netPack, sizeBytes)
	} else {
		netPack := make([]byte, n.netPackSize)
		size, err := conn.Read(netPack)
		if err != nil {
			n.onErr(errors.Join(err, ErrReadNetpackFailed))
			return true
		}
		n.onRec(netPack, size)
	}
	return false
}
