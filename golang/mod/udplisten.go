package mod

import (
	"github.com/intmian/mian_go_lib/tool/misc"
	"github.com/pkg/errors"
	"net"
)

type UdpListenerSetting struct {
	IP   string
	Port int
}

type UdpListener struct {
	setting UdpListenerSetting
	misc.InitTag
	closed bool
}

func NewUdpListener(setting UdpListenerSetting) (*UdpListener, error) {
	u := &UdpListener{}
	err := u.Init(setting)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func ListenUdp(setting UdpListenerSetting) (*UdpListener, error) {
	return NewUdpListener(setting)
}

func (u *UdpListener) Init(s UdpListenerSetting) error {
	if u.IsInitialized() {
		return errors.New("UdpListener already init")
	}
	u.setting = s
	u.SetInitialized()
	return nil
}

func (u *UdpListener) Accept() (net.Conn, error) {
	if !u.IsInitialized() {
		return nil, errors.New("UdpListener not init")
	}
	if u.closed {
		return nil, errors.New("UdpListener already closed")
	}
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(u.setting.IP),
		Port: u.setting.Port,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "listen udp failed")
	}
	return conn, nil
}

func (u *UdpListener) Close() error {
	u.closed = true
	return nil
}

func (u *UdpListener) Addr() net.Addr {
	return &net.UDPAddr{
		IP:   net.ParseIP(u.setting.IP),
		Port: u.setting.Port,
	}
}
