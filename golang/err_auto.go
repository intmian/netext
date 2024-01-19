package netext

type ErrStr string

const (
	ErrNil                  = ErrStr("nil")
	ErrConnMgrAlreadyInit   = ErrStr("ConnMgr already init")   // auto generated from .\connmgr.go
	ErrConnMgrNotInit       = ErrStr("ConnMgr not init")       // auto generated from .\connmgr.go
	ErrAddConnFailed        = ErrStr("add conn failed")        // auto generated from .\connmgr.go
	ErrDialMgrAlreadyInit   = ErrStr("DialMgr already init")   // auto generated from .\dialmgr.go
	ErrDialMgrNotInit       = ErrStr("DialMgr not init")       // auto generated from .\dialmgr.go
	ErrDialTcpFailed        = ErrStr("dial tcp failed")        // auto generated from .\dialmgr.go
	ErrDialKcpFailed        = ErrStr("dial kcp failed")        // auto generated from .\dialmgr.go
	ErrDialUdpFailed        = ErrStr("dial udp failed")        // auto generated from .\dialmgr.go
	ErrUnknownConnType      = ErrStr("unknown conn type")      // auto generated from .\dialmgr.go
	ErrListenMgrAlreadyInit = ErrStr("ListenMgr already init") // auto generated from .\listenmgr.go
	ErrListenMgrNotInit     = ErrStr("ListenMgr not init")     // auto generated from .\listenmgr.go
	ErrAddrInvalid          = ErrStr("addr invalid")           // auto generated from .\listenmgr.go
	ErrAddrAlreadyListen    = ErrStr("addr already listen")    // auto generated from .\listenmgr.go
	ErrListenTcpFailed      = ErrStr("listen tcp failed")      // auto generated from .\listenmgr.go
	ErrListenKcpFailed      = ErrStr("listen kcp failed")      // auto generated from .\listenmgr.go
	ErrConnTypeNotSupport   = ErrStr("conn type not support")  // auto generated from .\listenmgr.go
	ErrGoListenFailed       = ErrStr("goListen failed")        // auto generated from .\listenmgr.go
	ErrAddrNotListen        = ErrStr("addr not listen")        // auto generated from .\listenmgr.go
	ErrCloseFailed          = ErrStr("close failed")           // auto generated from .\listenmgr.go
	ErrAcceptFailed         = ErrStr("accept failed")          // auto generated from .\listenmgr.go
	ErrNetextNotInit        = ErrStr("netext not init")        // auto generated from .\mgr.go
	ErrMsgDataTooShort      = ErrStr("msg data too short")     // auto generated from .\msg.go
	ErrConntypeAlreadyExist = ErrStr("connType already exist") // auto generated from .\netentity.go
	ErrConntypeNotExist     = ErrStr("connType not exist")     // auto generated from .\netentity.go
	ErrWriteSizeFailed      = ErrStr("write size failed")      // auto generated from .\netentity.go
	ErrWriteDataFailed      = ErrStr("write data failed")      // auto generated from .\netentity.go
	ErrReadSizeFailed       = ErrStr("read size failed")       // auto generated from .\netentity.go
	ErrReadNetpackFailed    = ErrStr("read netPack failed")    // auto generated from .\netentity.go
	ErrKeyNotExist          = ErrStr("key not exist")          // auto generated from .\connmgr.go
	ErrCancelNotExist       = ErrStr("cancel not exist")       // auto generated from .\connmgr.go
	ErrSendFailed           = ErrStr("send failed")            // auto generated from .\connmgr.go
)

func (e ErrStr) Error() string { return string(e) }
