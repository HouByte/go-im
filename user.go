package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// NewUser 创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{Name: userAddr, Addr: userAddr, C: make(chan string), conn: conn, server: server}
	//启动监听当前用户
	go user.ListenMessage()
	return user
}

// ListenMessage 监听用户消息
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

// Online 上线
func (this *User) Online() {
	//用户上线,加入在线列表
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//广播当前用户上线消息
	this.server.BroadCast(this, "go online")
}

// Offline 下线
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "Offline")
}

var messageStrategyFactory = NewStrategyFactory()

// DoMessage 处理消息
func (this *User) DoMessage(context string) {
	i := strings.Index(context, " ")
	var cmd string = "bc"
	var body string = ""
	//如果只有命令
	if i == -1 && len(context) > 0 {
		cmd = context
	} else if i > 0 {
		cmd = context[:i]
		if len(context) > i+1 {
			body = context[i+1:]

		}
	}
	messageStrategyFactory.HandlerMessageStrategy(cmd, this, body)
}

func (this User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}
