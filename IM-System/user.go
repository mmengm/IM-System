package main

import "net"

//定义用户结构
type User struct {
	Name       string
	Addr       string
	ClientChan chan string
	conn       net.Conn
	server     *Server
}

//创建用户
func NewUser(conn net.Conn, server *Server) *User {

	//获取Client的地址
	userAddr := conn.RemoteAddr().String()
	//创建用户实例
	user := &User{
		Name:       userAddr,
		Addr:       userAddr,
		ClientChan: make(chan string),
		conn:       conn,
		server:     server,
	}
	//开启一个goroutine，监听用户chan中是否有数据
	go user.ListenMessage()

	return user
}

//监听chan管道中数据
func (this *User) ListenMessage() {
	for {
		msg := <-this.ClientChan
		this.conn.Write([]byte(msg + "\n"))
	}

}

//用户上线
func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BoradCast(this, "login in")
}

//用户下线

func (this *User) Offline() {

	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.BoradCast(this, "login out")
}

//发送消息

func (this *User) DoMessage(msg string) {
	//如果用户输入 who 就将在线人数的集合发送给他
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			OnlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "loginng ing \n "
			this.SendMsg(OnlineMsg)

		}
		this.server.mapLock.Unlock()
	} else {

		this.server.BoradCast(this, msg)
	}

}

//往当前这个连接中发送数据
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))

}
