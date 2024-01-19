package netext

import "context"

type (
	CmdRouter struct {
		NetType NetType
		cmd     CmdEnum
		handle  func(ID NetID, msg Msg) error
	}
	MsgMgrSetting struct {
		handleType HandleType // 目前所有NetType的handleType都是一样的，后续有需求再拓展，建议当前在业务层组合
		threadNum  int        // 如果是固定线程池，那么需要指定线程数
		workNum    int        // 如果是协程池，那么需要指定工作数
	}
	MsgMgrContext struct {
		ctx   context.Context
		onErr func(err error)
	}
)

type IMsgMgr interface {
	Init(s MsgMgrSetting, c MsgMgrContext) error
	AddRouter(router CmdRouter) error
	DelRouter(cmd CmdEnum) error
	OnRecMsg(key NetKey, msg Msg) error
}

// TODO:
