package netext

import (
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
)

// NetExt 网络管理器
type NetExt struct {
	misc.InitTag
}

func (n *NetExt) Init() error {
	n.setting = setting
	n.SetInitialized()
	return nil
}

// AddDial 试图与addr建立连接
// 需要保证双方的rule一致
func (n *NetExt) AddDial(addr NetAddr, rule NetRule) error {
	if !n.IsInitialized() {
		return errors.New("netext not init")
	}
	return nil
}

// AddListen 监听addr，并不停地
// 需要保证双方的rule一致
func (n *NetExt) AddListen(addr NetAddr, rule NetRule) error {
	if !n.IsInitialized() {
		return errors.New("netext not init")
	}
	return nil
}

// Close 关闭附属协程或逻辑
func (n *NetExt) Close() error {
	return nil
}

// Send 发送消息
// 请注意如果没有填入ID，那么会随机选择一个连接发送
func (n *NetExt) Send(key NetKey, msg interface{}) error {
	return nil
}

// SendAndRec 发送消息，并等待回复
// 请注意这个是rpc调用
func (n *NetExt) SendAndRec(key NetKey, msg interface{}, timeout int) (interface{}, error) {
	return nil, nil
}

// AddRouter 添加消息处理函数，增加某一种路由
func (n *NetExt) AddRouter(netType NetType, cmd Cmd, handler func(NetKey, interface{}, HandleContext)) error {
	return nil
}

// ManualRecMessage 手动接收消息。
// 如果模式为
func (n *NetExt) ManualRecMessage(key NetKey) error {

}
