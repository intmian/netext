package netext

type (
	ValidSetting struct {
	}
	ValidContext struct {
	}
)

type IValidMgr interface {
	Init(s ValidSetting, c ValidContext) error
	AddNeedValid(netType NetType, cmd CmdEnum) error
	DelNeedValid(netType NetType, cmd CmdEnum) error
	IsNeedValid(netType NetType, cmd CmdEnum) bool
	Valid(netType NetType, cmd CmdEnum, key NetKey) error
	UnValid(netType NetType, cmd CmdEnum, key NetKey) error
	IsValid(netType NetType, cmd CmdEnum, key NetKey) bool
}

type cmdKey struct {
	netType NetType
	cmd     CmdEnum
}

type ValidMgr struct {
	setting ValidSetting
	ctx     ValidContext
}
