package route

import (
	"bs/config"
	"bs/db"
	"bs/util"
	"strconv"
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

// 物品购买列表排序
func sortItemUseList(list []map[string]interface{}, sort int) {
	switch sort {
	case 1:
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				if util.GetInt(list[i]["all"]) < util.GetInt(list[j]["all"]) {
					list[i], list[j] = list[j], list[i]
				}
			}
		}
	case 2:
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				if util.GetInt(list[i]["all"]) > util.GetInt(list[j]["all"]) {
					list[i], list[j] = list[j], list[i]
				}
			}
		}
	case 3:
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				s := list[i]["weekRaise"].(string)
				s2 := list[j]["weekRaise"].(string)
				si, _ := strconv.Atoi(s[:len(s)-1])
				sj, _ := strconv.Atoi(s2[:len(s)-1])
				if si < sj {
					list[i], list[j] = list[j], list[i]
				}
			}
		}
	case 4:
		for i := 0; i < len(list); i++ {
			for j := i + 1; j < len(list); j++ {
				s := list[i]["weekRaise"].(string)
				s2 := list[j]["weekRaise"].(string)
				si, _ := strconv.Atoi(s[:len(s)-1])
				sj, _ := strconv.Atoi(s2[:len(s)-1])
				if si > sj {
					list[i], list[j] = list[j], list[i]
				}
			}
		}
	}
}
