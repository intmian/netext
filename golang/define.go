package netext

// 枚举
type (
	//NetType 用来表示网络类型，例如是面向玩家的网络，还是面向服务器的网络
	NetType int
	//ConnectType 用来表示连接类型，例如tcp|kcp|udp（udp会保留session，在用户层实现伪链接）
	ConnectType int
	//HandleType 用来表示处理类型，例如是单线程处理，还是多线程处理，还是固定线程处理，亦或者外部处理
	HandleType int
)

// 外部
type (
	// NetKey 用来表示网络
	NetKey struct {
		NetType  NetType
		ConnType ConnectType
		NetID    uint32
	}
	// NetAddr 用来表示网络地址
	NetAddr struct {
		ConnType ConnectType
		Addr     string
	}
	// Setting 主管理器的配置
	Setting struct {
		errorHandler func(error) // 用于处理错误
	}
	NetRule struct {
		//NeedAuth 如果需要认证，那么在未认证状态仅允许收到认证消息，此后需要外部调用认证接口。客户端需要再连接后
		NeedAuth       bool
		NeedSimpleAuth bool // 是否需要简单认证，简单认证仅校验下面的字段，在认证阶段，会交换随机数后hash
		SimpleAuthKey  string
		NeedHeart      bool
		NeedRc4        bool // 认证阶段会交换rc4秘钥
	}
	ListenSetting struct {
		HandleType HandleType
		WorkerNum  int
	}
	Cmd uint16
)

const (
	ConnTypeTcp ConnectType = iota
	ConnTypeKcp
	ConnTypeUdp
)

// 错误相关
const (
	ErrNetExtNotInit = Error("NetExt: not init")
)

type Error string

func (e Error) Error() string { return string(e) }
