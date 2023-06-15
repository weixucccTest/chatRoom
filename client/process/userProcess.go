package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"common/message"
	"study/shangguigu/chatroom/client/utils"
)

type UserProcess struct {
	// 暂时不需要字段
}

func (this *UserProcess)Register(UserName, Password, name string) (err error) {
	// 1.连接到服务器
	coon, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return err
	}
	// 延时关闭
	defer coon.Close()
	// 2.准备通过coon发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterResMesType
	// 3.创建一个registerMes
	var registerMes message.RegisterMes
	registerMes.User.UserName = UserName
	registerMes.User.UserPwd = Password
	registerMes.User.Name = name

	// 4.将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
	}

	// 5.把data赋给mes.Data字段
	mes.Data = string(data)

	// 6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Maeshal err = ", err)
		return err
	}

	// 创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: coon,
	}

	// 发送data到服务器
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("tf.WritePkg err :", err)
		return
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("tf.ReadPkg err :", err)
		return
	}

	// 将mes的data部分反序列化成registerMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if  registerResMes.Code == 200 {
		fmt.Println("注册成功，请登录")
	}else{
		fmt.Println(registerResMes.Error)
	}

	return
}

// 写一个函数，完成登录操作
func (this *UserProcess)Login(UserName, Password string) (err error) {

	// 1.连接到服务器
	coon, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return err
	}
	// 延时关闭
	defer coon.Close()

	// 2.准备通过coon发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	// 3.创建一个loginMes
	var loginMes message.LoginMessage
	loginMes.UserName = UserName
	loginMes.UserPwd = Password

	// 4.将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err =", err)
	}
	// 5.把data赋给mes.Data字段
	mes.Data = string(data)

	// 6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Maeshal err = ", err)
		return err
	}

	// 7.到这里，data就是要发送的数据
	// 7.1 先把data的长度发送给服务器
	// 先获取到data的长度，转换成一个表示长度的byte切片
	var pkgLen uint32

	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	// 发送长度
	n, err := coon.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("coon.Write err = ", err)
		return err
	}

	// fmt.Println("客户端发送消息的长度成功,长度为：", len(data))
	fmt.Printf("内容为：%s", string(data))

	// 发送消息本身
	_, err = coon.Write(data)
	if err != nil {
		fmt.Println("coon.Write(data) err = ", err)
		return err
	}

	// 休眠20s
	// time.Sleep(20 * time.Second)
	// fmt.Println("休眠20s")

	// 这里还需要处理服务器端返回的消息
	// 创建一个Transfer实例
	tf := &utils.Transfer{
		Conn : coon,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) fail ,err = ", err)
		return
	}
	fmt.Println(mes)

	// 将mes的data部分反序列化成loginResMes
	var loginRes message.LoginRes
	err = json.Unmarshal([]byte(mes.Data), &loginRes)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err = ", err)
		return
	}
	if loginRes.Code == 200 {

		//fmt.Println("登录成功")

		// 这里需要再客户端启动一个协成
		// 该协成保持和服务器的通讯，如果服务器有数据推送给客户端
		// 则接受并显示在客户端的终端
		go ProcessServerMes(coon)

		fmt.Println("当前在线成员：")
		for _, user := range loginRes.Users {
			fmt.Println("Name:",user.Name,"userName:",user.UserName ,"passWord:",user.UserPwd)
		}

		// 1.显示登录成功后的菜单
		for {
			ShowMenu(loginRes.LogMes.Name)
		}
	} else {
		fmt.Println(loginRes.Error)
	}
	return nil
}

// 测试

