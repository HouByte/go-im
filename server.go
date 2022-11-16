package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播
	Message chan string
}

// NewServer 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{Ip: ip, Port: port, OnlineMap: make(map[string]*User), Message: make(chan string)}
	return server
}

// BroadCast 广播
func (this Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)
	this.Message <- sendMsg

}

// ListenMessger 监听Message广播消息channel，有消息就广播
func (this Server) ListenMessger() {
	fmt.Println("Listen for messages")
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// Handler 处理器
func (this Server) Handler(conn net.Conn) {
	fmt.Println("Establish a connection")
	user := NewUser(conn)
	//用户上线,加入在线列表
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	//广播当前用户上线消息

	this.BroadCast(user, "go online")

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.BroadCast(user, "Offline")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
			}
			msg := string(buf[:n-1])

			//广播发送的消息
			this.BroadCast(user, msg)
		}
	}()

	//当前handler阻塞
	select {}
}

// Start 启动服务
func (this *Server) Start() {
	//socket
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net Listen err:", err)
		return
	}
	//close listen socket
	defer listen.Close()
	defer fmt.Println("Server Shutdown")

	fmt.Println("Server startup")
	//启动监听Message
	go this.ListenMessger()

	for {

		//accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("net accept err:", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}

}
