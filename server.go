package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{Ip: ip, Port: port}
	return server
}

func (this Server) Handler(conn net.Conn) {
	fmt.Println("成功建立连接")
}

func (this *Server) Start() {
	//socket
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net Listen err:", err)
		return
	}
	//close listen socket
	defer listen.Close()
	defer fmt.Println("服务器关闭")

	fmt.Println("服务器启动")
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
