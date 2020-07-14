package route

import (
	"bs/config"
	"bs/db"
	"strings"
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
	db.RedisSetTokenUsrn(token, db.RedisGetTokenUsrn(token))
}

// PassTokenAuth 跳过验证(一些接口不需要验证身份)
func PassTokenAuth(path string) bool {
	for _, url := range config.GetConfig().PassURL {
		if strings.Index(path, url) != -1 {
			return true
		}
	}
	return false
}

func ExportFiter(path, token string) bool {
	defer db.RedisDelTokenExport(token)
	for _, url := range config.GetConfig().ExportURL {
		if url == path {
			if db.RedisGetTokenExport(token) {
				return true
			} else {
				return false
			}
		}
	}
	return true
}
