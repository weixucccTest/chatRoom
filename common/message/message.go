package message

const (
	LoginMesType    		= "LoginMes"
	LoginResMesType 		= "LoginResMes"
	RegisterResMesType 		= "RegisterResMes"
)

type Message struct {
	Type string `json:"type"` //消息内容
	Data string `json:"data"` //消息切片
}

// 定义两个消息，后续需要别的了在增加
type LoginMessage struct {
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
	Name	 string `json:"name"`
}

type LoginRes struct {
	Code  int    `json:"code"`  //返回状态码，500表示该用户未注册，200表示登录成功
	Error string `json:"error"` //返回的错误信息
	Users []User `json:"users"` //返回当前在线的成员信息
	LogMes LoginMessage `json:"login_message"`
}

type RegisterMes struct {
	User User `json:"user"`// 类型就是User结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码，400表示该用户已经占用了，200表示注册成功
	Error string `json:"error"` //返回的错误信息
}
