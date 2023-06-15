package main

import (
	"github.com/gomodule/redigo/redis"
	"time"
)


// 定义一个全局的pool
var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration){
	pool = &redis.Pool{
		MaxIdle: maxIdle, 	// 最大空闲连接数
		MaxActive: maxActive,		// 标识和数据库的最大连接数，0表示没有限制
		IdleTimeout: idleTimeout,		// 最大空闲时间
		Dial: func() (redis.Conn, error) {		// 初始化连接的代码，连接哪个ip的redis
			return redis.Dial("tcp", address,redis.DialDatabase(1))  // redis.DialDatabase(1)标识连接第一个数据库
		},
	}
}
