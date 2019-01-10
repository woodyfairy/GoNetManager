package main

import (
	"../../NetManager"
	"log"
)

const(
	port = ":3000"

)

func main() {
	NetManager.NetMsgHandler = netHandler
	log.Println("Start Server", port)
	//最后阻断
	NetManager.Listen(port)
}

func netHandler (session *NetManager.Session, cmd NetManager.CmdType, data string) (error, string) {
	//log.Println("CMD:", cmd)
	if cmd == NetManager.Cmd_Checkin {
		return nil, "checkin:"+data
	}else {
		return nil, ""
	}
}