package main

import (
	"common/message"
	"fmt"
	"io"
	"net"
	process2 "study/shangguigu/chatroom/server/process"
	"study/shangguigu/chatroom/server/utils"
)

// 先创建一个process的结构体
type Processor struct {
	Conn net.Conn
}

// 编写一个ServerProcessmes 函数
// 功能：根据客户端发送的消息种类不同，决定调用哪个函数来处理
func (this *Processor) ServerProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		userPro := &process2.UserProcess{
			Conn: this.Conn,
		}
		// 处理登录
		err = userPro.ServerProcessLogin(mes)
	case message.RegisterResMesType:
		userPro := &process2.UserProcess{
			Conn: this.Conn,
		}
		// 处理注册
		err = userPro.ServerProcessRegister(mes)
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}

	return
}

func (this *Processor) Process2() (err error) {
	// 循环读客户端发送的信息
	for {
		// 这里我们将读取数据包直接封装成一个函数，readPkg()，返回message和err
		// 创建一个Transfer，完成读报的任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务端也退出。。。")
				return err
			}
			fmt.Println("readPkg err =", err)
			return err
		}
		fmt.Println("mes = ", mes)

		err = this.ServerProcessMes(&mes)
		if err != nil {
			return err
		}
		if mes.Type == message.RegisterResMesType{
			return nil
		}
	}
	return
}
