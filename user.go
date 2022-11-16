package main

import "net"

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

// DoMessage 处理消息
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}
