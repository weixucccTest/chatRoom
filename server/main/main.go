package main

import (
	"fmt"
	"net"
	"study/shangguigu/chatroom/server/model"
	"time"
)

// 处理和客户端的通讯
func Process(conn net.Conn) {
	// 这里需要延时关闭coon
	defer conn.Close()

	// 调用总控，先创建一个总控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.Process2()
	if err != nil {
		fmt.Println("客户端和服务器端通信携程错误， err = ", err)
		return
	}
}

// 这里我们编写一个函数，完成对USerDao的初始化任务
func initUserDao(){
	// 这里的pool，本身就是一个全局变量，在redis.go
	// 这里需要注意一个初始化的顺序问题
	// 先调用initPool，再调用initUserDao
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {

	// 当服务器刚开始时，我们就去初始化redis连接池
	initPool("127.0.0.1:6379", 16, 0 ,300*time.Second)
	initUserDao()

	// 提示信息
	fmt.Println("服务器[新的结构]在8889端口监听")
	listrn, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net listen err =", err)
		return
	}
	defer listrn.Close()

	// 一旦监听成功，就等待客户端来连接服务器
	for {
		fmt.Println("等待客户端来连接服务器")
		coon, err := listrn.Accept()
		if err != nil {
			fmt.Println("listen.accept err =", err)
		}

		// 一旦连接成功，则启动一个协成和客户端保持通讯
		go Process(coon)
	}
}
