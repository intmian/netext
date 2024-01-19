package netext

import (
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/xtaci/kcp-go"
	"net"
)

type (
	DialSetting struct {
	}
	DialContext struct {
		//context.Context
		//OnErr     func(err error)
		OnConnect func(conn net.Conn, addr NetAddr, rule NetRule)
	}
)

type DialMgr struct {
	setting DialSetting
	ctx     DialContext
	misc.InitTag
}

func (d *DialMgr) Init(s DialSetting, c DialContext) error {
	if d.IsInitialized() {
		return ErrDialMgrAlreadyInit
	}
	d.setting = s
	d.ctx = c
	d.SetInitialized()
	return nil
}

func (d *DialMgr) Add(addr NetAddr, rule NetRule) error {
	if !d.IsInitialized() {
		return ErrDialMgrNotInit
	}
	switch addr.ConnType {
	case ConnTypeTcp:
		conn, err := net.Dial("tcp", addr.GetAddr())
		if err != nil {
			return errors.Join(err, ErrDialTcpFailed)
		}
		d.ctx.OnConnect(conn, addr, rule)
	case ConnTypeKcp:
		conn, err := kcp.Dial(addr.GetAddr())
		if err != nil {
			return errors.Join(err, ErrDialKcpFailed)
		}
		d.ctx.OnConnect(conn, addr, rule)
	case ConnTypeUdp:
		conn, err := net.Dial("udp", addr.GetAddr())
		if err != nil {
			return errors.Join(err, ErrDialUdpFailed)
		}
		d.ctx.OnConnect(conn, addr, rule)
	default:
		return ErrUnknownConnType
	}
	return nil
}
