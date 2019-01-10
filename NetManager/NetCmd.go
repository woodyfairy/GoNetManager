package NetManager

//CMD：uint16，最大不能超过65535
type CmdType uint16
const (
	cmd_None		CmdType = 0 //初始空，私有
	Cmd_Heartbeat	CmdType	= 1	//心跳，直接返回，不处理
	//...

	Cmd_Checkin		CmdType = 2	//接入
)