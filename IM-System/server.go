package main

import (
	"fmt"
	"net"
)

//一个服务器需要有IP和端口地址
type Server struct {
	Ip   string
	Port int
}

//创建一个Server接口
//传入IP和端口,并返回这个Server的实例
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

//Server服务开启函数
func (this *Server) Start() {
	//服务器的IP和端口
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	//如果异常不为空,就说明程序出了问题
	if err != nil {
		fmt.Println("net listen err : ", err)
	}
	//最后关闭socket
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accpet err:", err)
			continue
		}

		go this.Handler(conn)

	}

}
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("Client连接成功了")

}
