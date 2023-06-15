package message

// 定义一个用户的结构体
type User struct {
	// 确定字段信息
	// 为了序列化和反序列化成功，必须保证用户信息的JSON字符串的key和结构体字段对应的tag名字一致
	UserName 	string `json:"userName"`
	UserPwd 	string `json:"userPwd"`
	Name	 	string `json:"name"`
}
