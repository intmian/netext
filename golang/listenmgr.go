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

type ListenContext struct {
	context.Context
	OnErr    func(err error)
	OnAccept func(conn net.Conn, rule NetRule)
}

type listener struct {
	net.Listener
	cancel func()
	rule   NetRule
}

type ListenMgr struct {
	setting  ListenSetting
	Context  ListenContext
	listener map[NetAddr]listener

	misc.InitTag
}

func (l *ListenMgr) Init(s ListenSetting, c ListenContext) error {
	if l.IsInitialized() {
		return errors.New("ListenMgr already init")
	}
	l.setting = s
	l.Context = c
	l.SetInitialized()
	return nil
}

func (l *ListenMgr) Add(addr NetAddr, rule NetRule) error {
	// 校验
	if !l.IsInitialized() {
		return errors.New("ListenMgr not init")
	}
	if !addr.IsValid() {
		return errors.New("addr invalid")
	}

	if _, ok := l.listener[addr]; ok {
		return errors.New("addr already listen")
	}

	var listen net.Listener
	var err error
	ip, port := addr.GetIpAndPort()
	switch addr.ConnType {
	case ConnTypeTcp:
		listen, err = net.Listen("tcp", addr.Addr)
		if err != nil {
			return errors.WithMessage(err, "listen tcp failed")
		}
	case ConnTypeKcp:
		listen, err = kcp.Listen(addr.Addr)
		if err != nil {
			return errors.WithMessage(err, "listen kcp failed")
		}
	case ConnTypeUdp:
		listen, err = mod.ListenUdp(mod.UdpListenerSetting{
			IP:   ip,
			Port: port,
		})
	default:
		return errors.New("conn type not support")
	}
	ctx, cancel := context.WithCancel(l.Context)
	l.listener[addr] = listener{
		Listener: listen,
		cancel:   cancel,
		rule:     rule,
	}
	err = l.GoListen(listen, ctx, rule)
	if err != nil {
		return errors.WithMessage(err, "GoListen failed")
	}
	return nil
}

func (l *ListenMgr) Close(addr NetAddr) error {
	if !l.IsInitialized() {
		return errors.New("ListenMgr not init")
	}
	if _, ok := l.listener[addr]; !ok {
		return errors.New("addr not listen")
	}
	l.listener[addr].cancel()
	err := l.listener[addr].Close()
	delete(l.listener, addr)
	if err != nil {
		return errors.WithMessage(err, "close failed")
	}
	return nil
}

func (l *ListenMgr) GoListen(listener net.Listener, ctx context.Context, rule NetRule) error {
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
