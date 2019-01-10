package main

import (
	"bytes"
	"fmt"
	"github.com/xtaci/kcp-go"
	"net"
	"time"
	"../../NetManager"
)

func main() {
	conn, err := kcp.Dial("127.0.0.1:3000")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	go handleConnC(conn)

	for {
		fmt.Println("send ------> ")
		sendData := NetManager.Enpack([]byte("hello kcp!!"))
		//sendData = []byte("hello kcp!!")
		ret, err2 := conn.Write(sendData)
		if err2 != nil {
			fmt.Println("err2:", err2)
			return
		} else {
			fmt.Println("send ret length:", ret)
		}
		time.Sleep(time.Duration(5)*time.Second)
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

		fmt.Print("datas : ")
		fmt.Println(datas.Bytes(), ":", string(datas.Bytes()))

	}
}
