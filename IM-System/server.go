package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//一个服务器需要有IP和端口地址
type Server struct {
	Ip   string
	Port int

	//	在线用户列表
	OnlineMap map[string]*User
	//锁
	mapLock sync.RWMutex
	//	消息广播的chan管道
	ServerChan chan string
}

//创建一个Server接口
//传入IP和端口,并返回这个Server的实例
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:         ip,
		Port:       port,
		OnlineMap:  make(map[string]*User),
		ServerChan: make(chan string),
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

	go this.ListenMessage()
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
	//fmt.Println("Client连接成功了")

	//	创建user实例
	user := NewUser(conn, this)
	//
	////	当用户连接成功后，将用户添加到在线用户map中
	//this.mapLock.Lock()
	//this.OnlineMap[user.Name] = user
	//this.mapLock.Unlock()
	////	广播当前用户上线消息
	//this.BoradCast(user, "login in")
	user.Online()
	isLive := make(chan bool)

	//用户消息广播  群聊实现
	go func() {
		buf := make([]byte, 4096)
		for {
			read, err := conn.Read(buf)
			if read == 0 {
				//this.BoradCast(user, "login out")
				user.Offline()
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err", err)
				return
			}
			//数据从conn连接中获取
			msg := string(buf[0 : read-1])
			//this.BoradCast(user, msg)
			user.DoMessage(msg)
			//如果每次发送消息，就将当前用户生成活跃
			isLive <- true
		}

	}()

	//阻塞当前连接
	for {
		select {
		case <-isLive:

		case <-time.After(time.Second * 30):
			user.SendMsg("您长时间不活跃，已被强行踢出")
			//关闭管道
			close(user.ClientChan)
			//关闭连接
			conn.Close()
			//	推出当前handler
			return
		}

	}
}

//广播信息
func (this *Server) BoradCast(user *User, msg string) {
	//生成一条发送的string数据
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.ServerChan <- sendMsg
}

//监听广播，一有消息就发送给全部在线用户
func (this *Server) ListenMessage() {

	for {
		msg := <-this.ServerChan

		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.ClientChan <- msg
		}
		this.mapLock.Unlock()
	}
}
