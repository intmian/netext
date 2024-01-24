package netext

import (
	"context"
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
	"net"
)

type (
	ConnSetting struct {
		MaxNetPackSize int
	}
	/*
		ConnContext 调用方相对于ConnMgr转入的上下文。
		context.Context 用于控制附属协程的生命周期。
		OnErr 用于处理错误。并发。
		OnEntityErr 用于处理某个网络实体的错误。上层根据策略选择应对。并发。
		OnRec 用于处理接收到的数据。并发。
	*/
	ConnContext struct {
		ctx         context.Context
		OnErr       func(err error)
		OnEntityErr func(key NetKey, err error)
		OnRec       func(key NetKey, data []byte, size int)
	}
)

// IConnMgr 管理连接，负责连接的读写。收到数据后会调用OnRec返回上层。
// 支持新增网络实体的某一个类型的链接，或者断开全部链接
type IConnMgr interface {
	// Init 初始化
	Init(s ConnSetting, c ConnContext) error
	// AddConn 增加连接
	AddConn(key NetKey, connType ConnectType, conn ConnMgr) error
	// DelConn 删除连接
	DelConn(key NetKey) error
	// Send 发送数据
	Send(key NetKey, connType ConnectType, data []byte) error
	//TODO: 考虑下有没有必要 DelConnByType(connType ConnectType) error
}

/*
ConnMgr 管理连接，负责连接的读写。收到数据后会调用OnRec返回上层。
支持新增网络实体的某一个类型的链接，或者断开全部链接
*/
type ConnMgr struct {
	setting ConnSetting
	ctx     ConnContext
	misc.InitTag
	id2entity map[NetKey]*netEntity
	id2cancel map[NetKey]func()
}

func (c *ConnMgr) Init(s ConnSetting, ctx ConnContext) error {
	if c.IsInitialized() {
		return ErrConnMgrAlreadyInit
	}
	c.setting = s
	c.ctx = ctx
	c.id2entity = make(map[NetKey]*netEntity)
	c.SetInitialized()
	return nil
}

func (c *ConnMgr) AddConn(key NetKey, connType ConnectType, conn net.Conn) error {
	if !c.IsInitialized() {
		return ErrConnMgrNotInit
	}
	if _, ok := c.id2entity[key]; !ok {
		c.addEntity(key)
	}
	entity := c.id2entity[key]
	err := entity.AddConn(connType, conn)
	if err != nil {
		return errors.Join(err, ErrAddConnFailed)
	}
	return nil
}

func (c *ConnMgr) addEntity(key NetKey) {
	onRec := func(data []byte, size int) {
		dataCopy := make([]byte, size)
		copy(dataCopy, data[:size])
		c.ctx.OnRec(key, dataCopy, size)
	}
	onErr := func(err error) {
		c.ctx.OnEntityErr(key, err)
	}
	ctx, cancel := context.WithCancel(c.ctx.ctx)
	netCtx := netEntityCtx{
		ctx:   ctx,
		onRec: onRec,
		onErr: onErr,
	}
	setting := netEntitySetting{
		netPackSize: c.setting.MaxNetPackSize,
	}
	c.id2entity[key] = newNetEntity(netCtx, setting)
	c.id2cancel[key] = cancel
}

func (c *ConnMgr) DelConn(key NetKey) error {
	if !c.IsInitialized() {
		return ErrConnMgrNotInit
	}
	if _, ok := c.id2entity[key]; !ok {
		return ErrKeyNotExist
	}
	if _, ok := c.id2cancel[key]; !ok {
		return ErrCancelNotExist
	}
	c.id2cancel[key]()
	delete(c.id2entity, key)
	return nil
}

func (c *ConnMgr) Send(key NetKey, connType ConnectType, data []byte) error {
	if !c.IsInitialized() {
		return ErrConnMgrNotInit
	}
	if _, ok := c.id2entity[key]; !ok {
		return ErrKeyNotExist
	}
	entity := c.id2entity[key]
	err := entity.Send(connType, data)
	if err != nil {
		return errors.Join(err, ErrSendFailed)
	}
	return nil
}
