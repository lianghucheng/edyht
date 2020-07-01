package db

import "github.com/szxby/tools/log"

// 设定好redis的key
const (
	TokenKey = "token" // token
)

// 数据过期时间
const expireTime = 5 * 60

// RedisSetToken 设置会话token
func RedisSetToken(token string, role int) {
	_, err := Do("Set", token, role, "EX", expireTime)
	if err != nil {
		log.Error("set token fail:%v", err)
	}
}

// RedisGetToken 获取会话token
func RedisGetToken(token string) int {
	data, err := Do("Get", token)
	if err != nil {
		log.Error("get token fail:%v", err)
		return -1
	}
	role, ok := data.(int)
	if !ok {
		log.Error("unknown token %v, role:%v", token, data)
		return -1
	}
	return role
}
