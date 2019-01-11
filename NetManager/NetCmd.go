package NetManager

//CMD：uint16，最大不能超过65535
type CmdType uint16
const (
	Cmd_None		CmdType = iota	//0 默认空
	Cmd_NetError	CmdType = iota	//1 网络未知错误回复，解包失败时未知命令
	Cmd_Heartbeat	CmdType	= iota	//2 心跳，直接返回，不处理
	//...

	Cmd_Checkin		CmdType = iota	//3接入
)