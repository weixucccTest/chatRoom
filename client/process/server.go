package process

import (
	"fmt"
	"net"
	"os"
	"study/shangguigu/chatroom/client/utils"
)

// 显示登录成功后的界面...
func ShowMenu(name string) {
	fmt.Println("-----------恭喜" + name + "登录成功-----------")
	fmt.Println("-----------1.显示用户在线列表-----------")
	fmt.Println("-----------2.发送消息-----------")
	fmt.Println("-----------3.信息列表-----------")
	fmt.Println("-----------4.退出系统-----------")
	fmt.Println("请选择(1-4)：")
	var key int
	fmt.Scanf("%d\n", &key)

	switch key {
	case 1:
		fmt.Println("显示用户在线列表")
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入选项不正确..")
	}
}

// 和服务器端保持通讯
func ProcessServerMes(conn net.Conn)  {

	// 创建一个Transfer实例，让他不停读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for{
		fmt.Println("客户端正在读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg，err = ", err)
			return
		}
		// 如果读取到消息，又是下一步处理逻辑
		fmt.Println("mes - ", mes)
	}
}