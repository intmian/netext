package netext

import (
	"net"
	"strconv"
)

func (n *NetAddr) IsValid() bool {
	if n.ConnType == ConnTypeNull {
		return false
	}
	if n.Addr == "" {
		return false
	}
	ipStr, port, err := net.SplitHostPort(n.Addr)
	if err != nil {
		return false
	}
	if ipStr == "" || port == "" {
		return false
	}
	return true
}

func (n *NetAddr) GetIpAndPort() (string, int) {
	ipStr, port, err := net.SplitHostPort(n.Addr)
	if err != nil {
		return "", 0
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return "", 0
	}
	return ipStr, portInt
}
