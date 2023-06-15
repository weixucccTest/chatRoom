package process

import "fmt"

// 因为UserMgr实例在服务器端有且只有一个
// 因为在很多地方，都会使用到，因此我们将其定义为全局变量
var userMgr *UserMgr

type UserMgr struct{
	OnlineUsers map[string]*UserProcess
}

func init() {
	userMgr = &UserMgr{
		OnlineUsers: make(map[string]*UserProcess),
	}
}

// 完成对onlineUser的添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.OnlineUsers[up.UserName] = up
}

// 删除
func (this *UserMgr) DeleteOnlineUser(userName string) {
	delete(this.OnlineUsers, userName)
}

// 返回当前所有在线的用户
func (this *UserMgr) GetAllOnlineUsers() map[string]*UserProcess{
	return this.OnlineUsers
}

// 返回某个用户的信息
func (this *UserMgr) GetOnlineUserByUserName (userName string) (*UserProcess, error){
	// 如何从map中取出一个值，检测是否找到数据
	process,ok := this.OnlineUsers[userName]
	if !ok {// 说明你要查找的用户不在线或不存在
		return nil, fmt.Errorf("用户不存在或不在线")
	}
	return process,nil
}