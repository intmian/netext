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
onGetNetPackBytes 用于获取一个可用的[]byte，等处理以后会通过回调返回
*/
type netEntityCtx struct {
	ctx               context.Context
	onRec             func([]byte, int)
	onErr             func(error)
	onGetNetPackBytes func() []byte // 从上层获取一个可用的[]byte，等处理以后会通过回调返回
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
}

func (n *netEntity) Init(ctx netEntityCtx) {
	n.entityMap = make(map[ConnectType]net.Conn)
	n.cancelMap = make(map[ConnectType]func())
	n.netEntityCtx = ctx
}

func newNetEntity(ctx netEntityCtx) *netEntity {
	n := &netEntity{}
	n.Init(ctx)
	return n
}

func (n *netEntity) AddConn(connType ConnectType, conn net.Conn) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.entityMap[connType]; ok {
		return ErrConntypeAlreadyExist
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
		return ErrConntypeNotExist
	}
	n.cancelMap[connType]()
	delete(n.entityMap, connType)
	delete(n.cancelMap, connType)
	return nil
}

func (n *netEntity) Send(connType ConnectType, data []byte) error {
	if _, ok := n.entityMap[connType]; !ok {
		return ErrConntypeNotExist
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
				if connType == ConnTypeTcp {
					// TCP 需要处理粘包的问题
					size, err := conn.Read(sizePack)
					if size != 2 {
						n.onErr(ErrReadSizeFailed)
						return
					}
					if err != nil {
						n.onErr(errors.Join(err, ErrReadSizeFailed))
						return
					}
					netPack := n.onGetNetPackBytes()
					size, err = conn.Read(netPack)
					if binary.BigEndian.Uint16(sizePack) != uint16(size) {
						n.onErr(ErrReadNetpackFailed)
					}
					if err != nil {
						n.onErr(errors.Join(err, ErrReadNetpackFailed))
						return
					}
					n.onRec(netPack, size)
				} else {
					netPack := n.onGetNetPackBytes()
					size, err := conn.Read(netPack)
					if err != nil {
						n.onErr(errors.Join(err, ErrReadNetpackFailed))
						return
					}
					n.onRec(netPack, size)
				}
			}
		}
	}()
}
