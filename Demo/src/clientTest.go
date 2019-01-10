package main

import (
	"bytes"
	"fmt"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
	"time"
	"../../NetManager"
)
const (
	HeartBeatTime	= 5
)

//var conn net.Conn

func main() {
	//var err error
	conn, err := kcp.Dial("127.0.0.1:3000")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	go handleConnC(conn)

	go checkTimeOut(&conn)

	for {
		fmt.Println("send -----> ")
		sendString1 := "HELLO_1"
		sendData1 := NetManager.Enpack(NetManager.Cmd_Checkin, []byte(sendString1))
		sendString2 := "HELLO_2"
		sendData2 := NetManager.Enpack(NetManager.Cmd_Checkin, []byte(sendString2))
		sendFull := append(sendData1, sendData2...)//测试如果有多个包同时发送了。

		ret, err2 := (conn).Write(sendFull)
		if err2 != nil {
			fmt.Println("err2:", err2)
			return
		} else {
			fmt.Println(sendString, "\nlength:", ret)
		}
		time.Sleep(time.Duration(20)*time.Second)
	}
}

func handleConnC(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [512]byte

		n, err := conn.Read(buf[0:])
		fmt.Println("recv -----> ")
		datas := bytes.NewBuffer(nil)
		datas.Write(buf[0:n])
		if err != nil {
			fmt.Println("read err:", err.Error())
			return
		}

		//fmt.Print("datas : ")
		fmt.Println(string(datas.Bytes()))

	}
}

var close = make(chan int)
var sendOrRcv = make(chan int)

func checkTimeOut(conn *net.Conn)  {
	for {
		select {
			case <-close:
				return
			case <-sendOrRcv:

			case <-time.After(HeartBeatTime * time.Second):
				log.Println("heartbeat ---->")
				sendString := ""
				sendData := NetManager.Enpack(NetManager.Cmd_Heartbeat, []byte(sendString))
				//log.Println("data:", sendData)
				//sendData = []byte("hello kcp!!")
				_, err := (*conn).Write(sendData)
				if err != nil {
					fmt.Println("err2:", err)
					return
				}
		}
	}
}
