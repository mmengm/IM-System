package main

func main() {
	//创建服务器实例
	server := NewServer("127.0.0.1", 6666)
	//开启服务
	server.Start()
}
