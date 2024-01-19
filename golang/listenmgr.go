package netext

import (
	"context"
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/intmian/netext/golang/mod"
	"github.com/xtaci/kcp-go"
	"net"
	"sync"
)

type ListenSetting struct {
}

type IListenMgr interface {
	Init(s ListenSetting, c ListenContext) error
	Add(addr NetAddr, rule NetRule) error
	Close(addr NetAddr) error
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
	mu       sync.Mutex
	misc.InitTag
}

// Init 初始化
func (l *ListenMgr) Init(s ListenSetting, c ListenContext) error {
	if l.IsInitialized() {
		return ErrListenMgrAlreadyInit
	}
	l.setting = s
	l.Context = c
	l.SetInitialized()
	return nil
}

// Add 新增监听 addr 地址，rule 规则。
// 不能重复监听同一个地址
func (l *ListenMgr) Add(addr NetAddr, rule NetRule) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	// 校验
	if !l.IsInitialized() {
		return ErrListenMgrNotInit
	}
	if !addr.IsValid() {
		return ErrAddrInvalid
	}

	if _, ok := l.listener[addr.GetAddr()]; ok {
		return ErrAddrAlreadyListen
	}

	var listen net.Listener
	var err error
	switch addr.ConnType {
	case ConnTypeTcp:
		listen, err = net.Listen("tcp", addr.GetAddr())
		if err != nil {
			return errors.Join(err, ErrListenTcpFailed)
		}
	case ConnTypeKcp:
		listen, err = kcp.Listen(addr.GetAddr())
		if err != nil {
			return errors.Join(err, ErrListenKcpFailed)
		}
	case ConnTypeUdp:
		listen, err = mod.ListenUdp(mod.UdpListenerSetting{
			IP:   addr.IP,
			Port: addr.port,
		})
	default:
		return ErrConnTypeNotSupport
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
		return errors.Join(err, ErrGoListenFailed)
	}
	return nil
}

// Close 关闭监听
func (l *ListenMgr) Close(addr NetAddr) error {
	if !l.IsInitialized() {
		return ErrListenMgrNotInit
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.listener[addr.GetAddr()]; !ok {
		return ErrAddrNotListen
	}
	l.listener[addr.GetAddr()].cancel()
	err := l.listener[addr.GetAddr()].Close()
	delete(l.listener, addr.GetAddr())
	if err != nil {
		return errors.Join(err, ErrCloseFailed)
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
				l.Context.OnErr(errors.Join(err, ErrAcceptFailed))
				continue
			}
			l.Context.OnAccept(conn, rule)
		}
	}()
	return nil
}
