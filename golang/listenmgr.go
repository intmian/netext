package netext

import (
	"context"
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/intmian/netext/golang/mod"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go"
	"net"
)

type ListenSetting struct {
}

/*
ListenContext 调用方相对于listen转入的上下文。
context.Context 用于控制附属协程的生命周期。
OnErr 用于处理错误。请注意回调可能是并发的。
OnAccept 用于处理新连接。请注意回调可能是并发的。接受新连接后，再进行鉴权或者基础通信获取NetType，和NetID
*/
type ListenContext struct {
	context.Context
	OnErr    func(err error)
	OnAccept func(conn net.Conn, rule NetRule)
}

type listener struct {
	addr NetAddr
	net.Listener
	cancel func()
	rule   NetRule
}

/*
ListenMgr 管理监听，并将监听到的连接转发给调用方
Init 后才能使用，不可重复调用
Add 新增监听
Close 关闭监听
*/
type ListenMgr struct {
	setting  ListenSetting
	Context  ListenContext
	listener map[string]listener

	misc.InitTag
}

// Init 初始化
func (l *ListenMgr) Init(s ListenSetting, c ListenContext) error {
	if l.IsInitialized() {
		return errors.New("ListenMgr already init")
	}
	l.setting = s
	l.Context = c
	l.SetInitialized()
	return nil
}

// Add 新增监听 addr 地址，rule 规则。
// 不能重复监听同一个地址
func (l *ListenMgr) Add(addr NetAddr, rule NetRule) error {
	// 校验
	if !l.IsInitialized() {
		return errors.New("ListenMgr not init")
	}
	if !addr.IsValid() {
		return errors.New("addr invalid")
	}

	if _, ok := l.listener[addr.GetAddr()]; ok {
		return errors.New("addr already listen")
	}

	var listen net.Listener
	var err error
	switch addr.ConnType {
	case ConnTypeTcp:
		listen, err = net.Listen("tcp", addr.GetAddr())
		if err != nil {
			return errors.WithMessage(err, "listen tcp failed")
		}
	case ConnTypeKcp:
		listen, err = kcp.Listen(addr.GetAddr())
		if err != nil {
			return errors.WithMessage(err, "listen kcp failed")
		}
	case ConnTypeUdp:
		listen, err = mod.ListenUdp(mod.UdpListenerSetting{
			IP:   addr.IP,
			Port: addr.port,
		})
	default:
		return errors.New("conn type not support")
	}
	ctx, cancel := context.WithCancel(l.Context)
	l.listener[addr.GetAddr()] = listener{
		addr:     addr,
		Listener: listen,
		cancel:   cancel,
		rule:     rule,
	}
	err = l.goListen(listen, ctx, rule)
	if err != nil {
		return errors.WithMessage(err, "goListen failed")
	}
	return nil
}

// Close 关闭监听
func (l *ListenMgr) Close(addr NetAddr) error {
	if !l.IsInitialized() {
		return errors.New("ListenMgr not init")
	}
	if _, ok := l.listener[addr.GetAddr()]; !ok {
		return errors.New("addr not listen")
	}
	l.listener[addr.GetAddr()].cancel()
	err := l.listener[addr.GetAddr()].Close()
	delete(l.listener, addr.GetAddr())
	if err != nil {
		return errors.WithMessage(err, "close failed")
	}
	return nil
}

func (l *ListenMgr) goListen(listener net.Listener, ctx context.Context, rule NetRule) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			conn, err := listener.Accept()
			if err != nil {
				l.Context.OnErr(errors.WithMessage(err, "accept failed"))
				continue
			}
			l.Context.OnAccept(conn, rule)
		}
	}()
	return nil
}
