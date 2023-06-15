package process

import (
	"common/message"
	"encoding/json"
	"fmt"
	"net"
	"study/shangguigu/chatroom/server/model"
	"study/shangguigu/chatroom/server/utils"
)

type UserProcess struct {
	// 字段？
	Conn net.Conn
	User message.User
	// 增加一个字段，标识该conn是哪个用户的
	UserName string
}

func (this *UserProcess)ServerProcessRegister(mes *message.Message) (err error) {
	// 先从mes中取出mes.Data，并直接反序列化成RegisterMes
	var registerMes message.RegisterMes

	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err = ", err)
		return
	}

	// 1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	// 2.再声明一个RegisterResMes，并完成赋值
	var registerResMes message.RegisterResMes

	// 我们需要到redis数据库完成注册
	// 1.使用model.MyUserDao.Register去注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS{
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		}else {
			registerResMes.Code = 506
			registerResMes.Error = "未知错误"
		}
	}else{
		registerResMes.Code = 200
	}
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) fail, err = ", err)
		return
	}

	// 4.将data复制给resMes
	resMes.Data = string(data)

	// 5.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail, err = ", err)
		return
	}

	// 6.发送data，我们将其封装成一个writePkg函数中
	// 因为使用了分层模式(mvc),线创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("writePkg fail, err = ", err)
		return
	}

	return
}

// 辨析一个函数ServerProcessLogin函数，专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {

	// 核心代码
	// 先从mes中取出mes.Data, 并直接反序列化成LoginMes
	var loginMes message.LoginMessage

	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err = ", err)
		return
	}

	// 1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 2/再声明一个LoginResMes，并完成复制
	var loginResMes message.LoginRes

	// 我们需要到redis数据库去完成验证
	// 使用model.MyUserDao到redis完成验证
	user, err := model.MyUserDao.Login(loginMes.UserName, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTEEXISTS{
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		}else if err == model.ERROR_USER_PWD{
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		}else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器未知错误。。。"
		}

	}else {
		loginResMes.Code = 200
		// 用户登录成功，我们就把登录成功的用户放到userMgr中，以便以后查询在线用户
		this.User = *user
		this.User.UserPwd = ""
		this.UserName = user.UserName

		// 将当前在线的用户信息，放入到loginResMes中
		for _, process := range userMgr.GetAllOnlineUsers() {
			loginResMes.Users = append(loginResMes.Users, process.User)
		}
		userMgr.AddOnlineUser(this)
	}
	fmt.Println(user)
	loginResMes.LogMes.Name = user.Name
	// 3.将loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) fail, err = ", err)
		return
	}

	// 4.将data复制给resMes
	resMes.Data = string(data)

	// 5.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail, err = ", err)
		return
	}

	// 6.发送data，我们将其封装成一个writePkg函数中
	// 因为使用了分层模式(mvc),线创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("writePkg fail, err = ", err)
		return
	}


	return
}
