package NetManager

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go"
	"log"
)

const showNetLog = false

var sid uint32 = 0

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

		//创建session
		session := NewSession(sid, conn)
		SharedSessionManager.SetSession(sid, session)
		sid++

		go handleConn(session)
	}
}

func handleConn(session *Session) {
	defer session.Close()

	tempBuff := make([]byte, 0)
	readBuff := make([]byte, 256)
	data := make([]byte, 0)
	cmd := Cmd_None

	for {
		n, err := session.connection.Read(readBuff)
		if err != nil {
			if showNetLog {
				//主动close时也会返回错误
				log.Println("Read ERROR:", err.Error())
			}
			return
		}

		if showNetLog {
			log.Println("recv -----> ")
			log.Println("length:", n)
		}
		tempBuff = append(tempBuff, readBuff[:n]...)

		errDepack := Depack(&tempBuff, &cmd, &data)
		if errDepack != nil {
			if showNetLog {
				log.Println("Depack ERROR:", errDepack)
			}
			return
		}

		if cmd == Cmd_None {
			continue
		}

		go doData(session, cmd, &data)

		cmd = Cmd_None
		data = data[0:0]//清空
	}
}

type msgHandler func (session *Session, cmd CmdType, data string) (error, string)
var NetMsgHandler msgHandler = nil
func doData(session *Session, cmd CmdType, data *[]byte) {
	if showNetLog {
		log.Print("Data : ")
		log.Println(string((*data)[:]))
	}
	//刷新时间
	session.updateTime()
	//心跳直接返回
	if cmd == Cmd_Heartbeat {
		//log.Println("Heartbeat")
		_= session.Send(cmd, "")
		return
	}
	//处理数据
	dataString := string(*data)
	if NetMsgHandler != nil {
		err, msg :=  NetMsgHandler(session, cmd, dataString)
		if err != nil {
			log.Println("NetMsgHandler ERROR:", err)
			errMsg := fmt.Sprintf("{\"err\" : \"%s\"}",err.Error())
			_= session.Send(cmd, errMsg)
		}else if len(msg) > 0{
			_= session.Send(cmd, msg)
		}
	}
}

//粘包：包头格式 NETHEADER+2字节uint16表示内容长度
const (
	netHeader       	= "~HEADER~"
	headerLength		= 8 //NETHEADER长度
	headerCmdSize		= 2 //CmdType占2字节 //CMD编号，0-65535
	headerMsgLengthSize	= 2 //uint16 //内容最大长度65535个字节

	maxUint16 int = 0xFFFF
)

func Enpack(cmd CmdType, message []byte) []byte {
	msgL := len(message)
	if msgL > maxUint16 {
		return nil
	}

	return append(append(append([]byte(netHeader), Uint16ToBytes(uint16(cmd))...), Uint16ToBytes(uint16(msgL))...), message...)
}
func Depack(buff *[]byte, cmd *CmdType, data *[]byte) error {
	l := len(*buff)
	if l > maxUint16 {
		return errors.New("NET ERROR: MSG TOO LONG")
	}

	length := uint16(l)
	if length <  headerLength + headerCmdSize + headerMsgLengthSize{
		return nil
	}

	//如果header不是 指定的header 说明此数据已经被污染 直接返回错误
	if string((*buff)[:headerLength]) != netHeader {
		return errors.New("NET ERROR: WRONG HEADER")
	}

	msgLength := BytesToUint16((*buff)[headerLength + headerCmdSize : headerLength + headerCmdSize + headerMsgLengthSize])
	if length < headerLength + headerCmdSize + headerMsgLengthSize + msgLength {
		return nil
	}

	*cmd = CmdType(BytesToUint16((*buff)[headerLength : headerLength + headerCmdSize]))
	*data = (*buff)[headerLength + headerCmdSize + headerMsgLengthSize : headerLength + headerCmdSize + headerMsgLengthSize + msgLength]
	*buff = (*buff)[headerLength + headerCmdSize + headerMsgLengthSize + msgLength:]

	return nil
}

//将int转成四个字节
func Uint16ToBytes(n uint16) []byte {
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

