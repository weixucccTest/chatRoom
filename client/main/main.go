package main

import (
	"fmt"
	"os"
	"study/shangguigu/chatroom/client/process"
)

// 定义两个全局变量，标识用户名和密码
var  Password, Username, Name string

func main() {
	// 接受用户选择
	var key int
	// 判断是否还显示菜单
	up := &process.UserProcess{}
	loop := true
	for loop {
		fmt.Println("-----------------------欢迎登录------------------")
		fmt.Println("\t\t\t1 登录聊天室")
		fmt.Println("\t\t\t2 注册用户")
		fmt.Println("\t\t\t3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3):")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Printf("username:")
			fmt.Scanf("%s\n", &Username)
			fmt.Printf("password:")
			fmt.Scanf("%s\n", &Password)
			// 完成登录
			// 1.创建一个UserProcess的实例
			err := up.Login(Username,Password)
			if err != nil {
				fmt.Println("Login err :", err)
			}

		case 2:
			fmt.Println("注册用户")
			fmt.Printf("请输入用户名:")
			fmt.Scanf("%s\n", &Username)
			fmt.Printf("请输入密码:")
			fmt.Scanf("%s\n", &Password)
			fmt.Printf("请输入名称:")
			fmt.Scanf("%s\n", &Name)
			// 调用UserProcess实例，完成注册请求
			err := up.Register(Username, Password, Name)
			if err != nil {
				fmt.Println("注册失败，err:", err)
			}
		case 3:
			fmt.Println("退出登录")
			//loop = false
			os.Exit(0)
		default:
			fmt.Println("输入有误")
		}
	}
	// 根据用户输入显示新的提示信息

	//if key == 1 {
	//	// 说明用户开始的登录
	//	fmt.Printf("username:")
	//	fmt.Scanf("%v\n", &Username)
	//	fmt.Printf("password:")
	//	fmt.Scanf("%s\n", &Password)
	//
	//	// 这里需要重新调用
	//	err := Login(Username, Password)
	//	fmt.Println(err)
	//} else if key == 2 {
	//	fmt.Println("开始注册新用户")
	//}

}
