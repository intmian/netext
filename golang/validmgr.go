package netext

import (
	"github.com/intmian/mian_go_lib/tool/misc"
)

type (
	ValidSetting struct {
	}
	ValidContext struct {
	}
)

// IValidMgr 是一个简单的鉴权管理器，用于管理哪些连接需要鉴权，哪些连接已经鉴权成功，需要在新增路由时先在这里注册需要鉴权。
// 在下面有一个单机版本的简单实现ValidMgr，如果需要分布式的鉴权管理器，需要自己实现
type IValidMgr interface {
	// Init 初始化
	Init(s ValidSetting, c ValidContext) error
	// AddNeedValid 增加需要鉴权的路由
	AddNeedValid(netType NetType, cmd CmdEnum) error
	// DelNeedValid 删除需要鉴权的路由
	DelNeedValid(netType NetType, cmd CmdEnum) error
	// IsNeedValid 此路由是否需要鉴权
	IsNeedValid(netType NetType, cmd CmdEnum) bool
	// Valid 鉴权成功
	Valid(netType NetType, cmd CmdEnum, key NetKey) error
	// UnValid 需要某个网络实体的鉴权
	UnValid(netType NetType, cmd CmdEnum, key NetKey) error
	// IsValid 是否鉴权成功
	IsValid(netType NetType, cmd CmdEnum, key NetKey) bool
}

type cmdKey struct {
	netType NetType
	cmd     CmdEnum
}

// ValidMgr 是一个简单的鉴权管理器，用于管理哪些连接需要鉴权，哪些连接已经鉴权成功，需要在新增路由时先在这里注册需要鉴权
type ValidMgr struct {
	setting   ValidSetting
	ctx       ValidContext
	needValid map[cmdKey]bool
	isValid   map[NetKey]bool
	misc.InitTag
}

func (v *ValidMgr) Init(s ValidSetting, c ValidContext) error {
	if v.IsInitialized() {
		return ErrValidMgrAlreadyInit
	}
	v.setting = s
	v.ctx = c
	v.SetInitialized()
	return nil
}

func (v *ValidMgr) AddNeedValid(netType NetType, cmd CmdEnum) error {
	if !v.IsInitialized() {
		return ErrValidMgrNotInit
	}
	v.needValid[cmdKey{netType: netType, cmd: cmd}] = true
	return nil
}

func (v *ValidMgr) DelNeedValid(netType NetType, cmd CmdEnum) error {
	if !v.IsInitialized() {
		return ErrValidMgrNotInit
	}
	delete(v.needValid, cmdKey{netType: netType, cmd: cmd})
	return nil
}

func (v *ValidMgr) IsNeedValid(netType NetType, cmd CmdEnum) bool {
	if !v.IsInitialized() {
		return false
	}
	return v.needValid[cmdKey{netType: netType, cmd: cmd}]
}

func (v *ValidMgr) Valid(key NetKey) error {
	if !v.IsInitialized() {
		return ErrValidMgrNotInit
	}
	v.isValid[key] = true
	return nil
}

func (v *ValidMgr) UnValid(key NetKey) error {
	if !v.IsInitialized() {
		return ErrValidMgrNotInit
	}
	v.isValid[key] = false
	return nil
}

func (v *ValidMgr) IsValid(netType NetType, cmd CmdEnum, key NetKey) bool {
	if !v.IsInitialized() {
		return false
	}
	if !v.IsNeedValid(netType, cmd) {
		return true
	}
	return v.isValid[key]
}
