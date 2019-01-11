package main

import (
	"../../NetManager"
	"fmt"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"time"
)

const (
	showClinetLog	= false
	showDataLog		= true
	showPing		= true

	HeartBeatTime	= 5
	serverAddr = "localhost:3000"
)

func main() {
	//var err error
	conn, err := kcp.Dial(serverAddr)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	go handleConnC(conn)

	go checkTimeOut(&conn)

	for {
		if showClinetLog {
			fmt.Println("send -----> ")
		}

		//测试如果有多个包同时发送了
		//sendString1 := "HELLO_1"
		//sendData1 := NetManager.Enpack(NetManager.Cmd_Checkin, []byte(sendString1))
		//sendString2 := "HELLO_2"
		//sendData2 := NetManager.Enpack(NetManager.Cmd_Checkin, []byte(sendString2))
		//sendFull := append(sendData1, sendData2...)
		//ret, err2 := (conn).Write(sendFull)
		//if err2 != nil {
		//	fmt.Println("err2:", err2)
		//	return
		//} else {
		//	fmt.Println(sendFull, "\nlength:", ret)
		//}


		sendString := "hello"
		err2 := sendMsg(&conn, NetManager.Cmd_Checkin, sendString)
		if err2 != nil {
			fmt.Println("err2:", err2)
			return
		} else {
			if showClinetLog {
				fmt.Println(sendString)
			}
		}

		time.Sleep(time.Duration(8)*time.Second)
	}
}
func handleConnC(conn net.Conn) {
	defer conn.Close()

	tempBuff := make([]byte, 0) //总buff
	cmd := NetManager.Cmd_None //cmd
	data := make([]byte, 0) //msg，可能=0
	readBuff := make([]byte, 128) //临时读取buff
	for {
		n, err := conn.Read(readBuff)
		if err != nil {
			log.Println("Read ERROR:", err.Error())
			return
		}

		if showClinetLog {
			log.Println("recv -----> ")
			log.Println("length:", n)
		}
		tempBuff = append(tempBuff, readBuff[:n]...)

		//若buff中存在多个完整命令包，则都需要取出来执行，否则会阻塞后面的命令包
		for len(tempBuff) > 0 {
			errDepack, finish := NetManager.Depack(&tempBuff, &cmd, &data)
			if errDepack != nil {
				log.Println("Depack ERROR:", errDepack)
				return
			}

			if finish == false{
				break // 说明只有半包，break，等待read
			}

			doData(&conn, cmd, data) //是否需要go?

			cmd = NetManager.Cmd_None
			data = data[0:0]//清空
		}
	}
}

//发送
var sendTime int64 = 0
func sendMsg(conn *net.Conn, cmd NetManager.CmdType, msg string) error {
	if showDataLog {
		log.Println("Send Cmd:", cmd, " data:", msg)
	}
	sendData := NetManager.Enpack(cmd, []byte(msg))
	_, err := (*conn).Write(sendData)
	if err != nil {
		log.Println("Send Error:", err)
	}else {
		sendTime = time.Now().UnixNano()
		setSendOrRecv() //心跳重新计时
	}
	return err
}
//接受到命令
func doData(conn *net.Conn, cmd NetManager.CmdType, data []byte) {
	setSendOrRecv() //心跳重新计时
	//ping
	if sendTime != 0 {
		now := time.Now().UnixNano()
		ping := float32(now - sendTime) / 1e6
		if showPing {
			log.Println("ping:", ping , "ms")
		}
	}
	//data
	if showDataLog {
		log.Println("Recv Cmd:", cmd, " data:", string(data))
	}
}

//当发送或者接受之后心跳重新计时
var sendOrRecv = make(chan int)
func setSendOrRecv()  {
	select {
	case <- sendOrRecv:
		//如果有则清空
	default:
		//否则不做
	}
	sendOrRecv <- 1
}
var close = make(chan int)
//检查心跳
func checkTimeOut(conn *net.Conn)  {
	for {
		select {
			case <-close:
				return
			case <-sendOrRecv:

			case <-time.After(HeartBeatTime * time.Second):
				if showClinetLog {
					log.Println("heartbeat ---->")
				}
				err := sendMsg(conn, NetManager.Cmd_Heartbeat, "")
				if err != nil {
					fmt.Println("err2:", err)
					return
				}
		}
	}
}
