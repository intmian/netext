package netext

import (
	"context"
	"errors"
	"github.com/intmian/mian_go_lib/tool/misc"
)

type (
	CmdRouterKey struct {
		NetType NetType
		Cmd     CmdEnum
	}
	Handler       func(ID NetID, msg Msg, ctx HandleContext) error
	HandleContext struct {
		Ctx          context.Context
		OnDisconnect func()
	}
	HandleMgrSetting struct {
		handleType HandleType // 目前所有NetType的handleType都是一样的，后续有需求再拓展，建议当前在业务层组合
		workNum    int        // 如果是协程池，那么需要指定工作数
		queueSize  int        // 队列大小，在外部接管模式下特别有用
	}
	MsgMgrContext struct {
		Ctx          context.Context
		OnErr        func(err error)
		OnWarn       func(str string)
		OnPanic      func(err error)  // 协程发生错误会recover
		OnDisConnect func(key NetKey) // 网络实体断开连接
	}
)

type IHandleMgr interface {
	Init(s HandleMgrSetting, c MsgMgrContext) error

	AddRouter(router CmdRouterKey, handle Handler) error
	DelRouter(router CmdRouterKey) error

	OnRecMsg(key NetKey, msg Msg) error
	OnHandle() error
}

type HandleMgr struct {
	setting       HandleMgrSetting
	ctx           MsgMgrContext
	msgChan       chan Msg
	threadMsgChan map[int]chan Msg
	start         bool
	misc.InitTag
}

func (h *HandleMgr) Init(s HandleMgrSetting, c MsgMgrContext) error {
	if h.IsInitialized() {
		return errors.New("handle mgr already init")
	}
	h.setting = s
	h.ctx = c
	h.msgChan = make(chan Msg, s.queueSize)
	h.SetInitialized()
	return nil
}

func (h *HandleMgr) startWorker() error {
	newCtx := context.WithoutCancel(h.ctx.Ctx)
	go func() {
		for {
			select {
			case <-newCtx.Done():
				return
			case msg := <-h.msgChan:

			}
		}
	}()
}
