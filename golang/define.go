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
	// NetID 用来表示网络ID
	NetID uint64
	// NetKey 用来表示网络逻辑层地址
	NetKey struct {
		NetType NetType
		ID      NetID
	}
	// NetAddr 用来表示网络地底层址
	NetAddr struct {
		ConnType ConnectType
		IP       string
		port     int
	}
	NetRule struct {
		//NeedAuth 如果需要认证，那么在未认证状态仅允许收到认证消息，此后需要外部调用认证接口。客户端需要再连接后
		NeedAuth       bool
		NeedSimpleAuth bool // 是否需要简单认证，简单认证仅校验下面的字段，在认证阶段，会交换随机数后hash
		SimpleAuthKey  string
		NeedHeart      bool
		NeedRc4        bool // 认证阶段会交换rc4秘钥
	}
)

const (
	ConnTypeNull ConnectType = iota
	ConnTypeTcp
	ConnTypeKcp
	ConnTypeUdp
)

const (
	HandleTypeNull      HandleType = iota
	HandleTypeModThread            // 由固定的线程处理，线程号为NetID % 线程数量
	HandleTypeWorkPool             // 由协程池处理
	HandleTypeOutHandle            // 由外部驱动
)
