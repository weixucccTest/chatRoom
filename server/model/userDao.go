package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"common/message"
)

// 我们在服务器启动后，就初始化一个userDao的实例
// 把他做成全局变量，在需要和redis交互时，就直接使用即可
var (
	MyUserDao *UserDao
)

// 定义一个userDao结构体
// 完成对user结构体的各种操作

type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个userDao的实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao){
	userDao = &UserDao{
		pool : pool,
	}

	return
}

// 应该提供哪些方法给到我们
// 1.根据一个用户ID(userName)，返回一个user实例或者error

func (this *UserDao) getUserByUserName(userName string)(user *message.User, err error) {

	conn := this.pool.Get()
	defer conn.Close()
	// 通过给定的userName，去redis查询用户
	res, err := redis.String(conn.Do("HGET", "users", userName))
	if err != nil {
		// 错误！
		if err == redis.ErrNil{  //标识在users的hash中没有找到对应的成员
			err = ERROR_USER_NOTEEXISTS
		}
		return
	}
	user = &message.User{}

	fmt.Println("redis：",res)
	// 这里我们需要把res反序列化成User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	return
}

// 完成登录校验
// 1.完成对用户的验证
// 2.如果用户名和密码都正确，则返回一个user实例
// 3.如果用户名或密码有错误，则返回对应的错误信息
func (this *UserDao)Login(userName, pwd string) (user *message.User, err error) {

	// 先从userDao的连接池中取出连接
	user, err = this.getUserByUserName(userName)
	if err != nil {
		return
	}

	// 这时证明用户是获取到了，密码是否正确还需要判断
	if user.UserPwd != pwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (this *UserDao) RegisterUser(user *message.User) (err error){
	conn := this.pool.Get()
	defer conn.Close()

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json.Marshal err :", err)
		return
	}


	_, err = conn.Do("HSET", "users", user.UserName, string(data))
	if err != nil {
		fmt.Println("注册用户失败， err = ", err)
		return
	}
	return
}

func (this *UserDao) Register (user *message.User) (err error) {

	_, err = this.getUserByUserName(user.UserName)
	if err == nil{
		err = ERROR_USER_EXISTS
		return
	}

	// 这时，说明userName还没有注册，可以完成后续操作
	err = this.RegisterUser(user)
	if err != nil {
		return
	}
	return
}