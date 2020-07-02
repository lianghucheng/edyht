package route

import (
	"bs/db"
)

// content tpye
const (
	JSON = "application/json"
)

// checkRole 检查是否越权操作 todo
func checkRole(role int, path string) bool {
	return true
}

// refreshToken 刷新token
func refreshToken(token string, role int) {
	db.RedisSetToken(token, role)
}
