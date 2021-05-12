package main

import "net"

//定义用户结构
type User struct {
	Name       string
	Addr       string
	ClientChan chan string
	conn       net.Conn
}

//创建用户
func NewUser(conn net.Conn) *User {

	//获取Client的地址
	userAddr := conn.RemoteAddr().String()
	//创建用户实例
	user := &User{
		Name:       userAddr,
		Addr:       userAddr,
		ClientChan: make(chan string),
		conn:       conn,
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
