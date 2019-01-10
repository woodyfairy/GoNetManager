package NetManager

import (
	"log"
	"net"
	"sync"
	"time"
)

const showSessionLog = true

const (
	HeartBeatTime	= 5
	CheckTime		= HeartBeatTime * 2
)

//.......................................SESSION.......................................
type Session struct {
	Sid uint32
	Uid uint32
	connection net.Conn
	LastTime int64 //秒
	lock  sync.Mutex
}

func NewSession(sid uint32, con net.Conn) *Session {
	if showSessionLog {
		log.Println("New Session:", sid)
	}

	return &Session{
		Sid : sid,
		connection:   con,
		LastTime: time.Now().Unix(),
	}
}
//send
func (this *Session) Send(cmd CmdType, msg string) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	data := Enpack(cmd, []byte(msg))//包
	_ ,errs := this.connection.Write(data)
	if errs != nil {
		if showSessionLog {
			log.Println("Send ERROR:", errs)
		}
	}
	return errs
}
//关闭conn
func (this *Session)Close(){

	err := this.connection.Close()
	if showSessionLog {
		if err != nil{
			log.Println("Session Close Error:", err)
		}else {
			log.Println("Session Closed:", this.Sid)
		}
	}

}
//每当收到包就刷新时间，超时的关闭
func (this *Session) updateTime() {
	this.LastTime = time.Now().Unix()
}


//.......................................SESSION管理类.......................................
//shared
var SharedSessionManager = NewSessionManager()

type SessionManager struct {
	sessions map[uint32]*Session
	num      uint32
	lock     sync.RWMutex
}
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[uint32]*Session),
		num:      0,
	}
	go sm.checkTimeout()
	return sm
}

func (this *SessionManager) checkTimeout() {
	for {
		time.Sleep(CheckTime * time.Second)
		for i,v:= range this.sessions{
			if time.Now().Unix() - CheckTime > v.LastTime {
				this.RemonveSessionById(i)
			}
		}
	}
}

func (this *SessionManager) GetSessionById(id uint32) *Session {
	if v, exist := this.sessions[id]; exist {
		return v
	}
	return nil
}

func (this *SessionManager) SetSession(id uint32, sess *Session) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.sessions[id] = sess
}

//关闭连接并删除
func (this *SessionManager) RemonveSessionById(id uint32) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v,exit := this.sessions[id];exit{
		v.Close()
	}
	delete(this.sessions, id)
}
//func (this *SessionManager) WriteByid(id uint32, msg string) bool {
//	if v, exist := this.sessions[id]; exist {
//		if err := v.Send(msg); err != nil {
//			this.RemonveSessionById(id)
//			return false
//		} else {
//			return true
//		}
//	}
//	return false
//}