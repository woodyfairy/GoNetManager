package NetManager

import (
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go"
	"log"
	"net"
)

var showlog = true

func Listen(laddr string){
	//config if need
	listener, err := kcp.Listen( laddr)
	if err != nil {
		log.Panicln("Listen ERROR:", err)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept err:", err)
			//return
			continue
		}
		go handleConn(conn)
	}
}


func handleConn(conn net.Conn) {
	defer closeCoon(conn)

	tempBuff := make([]byte, 0)
	readBuff := make([]byte, 256)
	data := make([]byte, 0)

	for {
		n, err := conn.Read(readBuff)
		if err != nil {
			log.Println("Read ERROR:", err.Error())
		}

		if showlog {
			log.Println("recv -----> ")
			log.Println("length:", n)
		}
		tempBuff = append(tempBuff, readBuff[:n]...)

		errDepack := Depack(&tempBuff, &data)
		if errDepack != nil {
			if showlog {
				log.Println("Depack ERROR:", errDepack)
			}
			return
		}

		if len(data) == 0 {
			continue
		}

		_ = doData(&data)
		data = data[0:0]//清空
	}
}

func doData(data *[]byte) error {
	if showlog {
		log.Print("datas : ")
		log.Println(string((*data)[:]))
	}
	//conn.Write(datas.Bytes())

	return nil
}

func closeCoon(conn net.Conn)  {
	log.Println("CLOSE")
	conn.Close()
}

//粘包：包头格式 NETHEADER+2字节uint表示内容长度
const (
	NETHEADER       	= "~HEADER~"
	NETHEADER_LEN		= 8 //NETHEADER长度
	NETHEADER_LENSIZE	= 2 //uint16
)
func Enpack(message []byte) []byte {
	return append(append([]byte(NETHEADER), Uint16ToBytes(len(message))...), message...)
}
func Depack(buff *[]byte, data *[]byte) error {
	length := uint16(len(*buff))
	if length <  NETHEADER_LEN + NETHEADER_LENSIZE{
		return nil
	}

	//如果header不是 指定的header 说明此数据已经被污染 直接返回错误
	if string((*buff)[:NETHEADER_LEN]) != NETHEADER {
		return errors.New("NET ERROR: WRONG HEADER")
	}

	msgLength := BytesToUint16((*buff)[NETHEADER_LEN : NETHEADER_LEN + NETHEADER_LENSIZE])
	if length < NETHEADER_LEN+NETHEADER_LENSIZE+msgLength {
		return nil
	}

	*data = (*buff)[NETHEADER_LEN + NETHEADER_LENSIZE : NETHEADER_LEN + NETHEADER_LENSIZE+msgLength]
	*buff = (*buff)[NETHEADER_LEN + NETHEADER_LENSIZE + msgLength:]

	return nil
}

//将int转成四个字节
func Uint16ToBytes(n int) []byte {
	x := uint16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, x)
	if err != nil {
		log.Println("i2b ERROR:", err)
	}
	return bytesBuffer.Bytes()
}
//将四个字节转成int
func BytesToUint16(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var x uint16
	err := binary.Read(bytesBuffer, binary.BigEndian, &x)
	if err != nil {
		log.Println("b2i ERROR:", err)
	}
	return uint16(x)
}

