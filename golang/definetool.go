package netext

import "strconv"

func (n *NetAddr) IsValid() bool {
	if n.ConnType == ConnTypeNull {
		return false
	}
	if n.IP == "" || n.port <= 0 {
		return false
	}
	return true
}

func (n *NetAddr) GetAddr() string {
	return n.IP + ":" + strconv.Itoa(n.port)
}
