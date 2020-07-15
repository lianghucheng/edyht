package db

import (
	"github.com/szxby/tools/log"
)

// 设定好redis的key
const (
	TokenKey       = "token"       // token
	MatchReportKey = "matchReport" // match report
	MatchListKey   = "matchList"   // match list
	TokenExportKey = "tokenExport" // token export
	TokenUsrn      = "usrn"
)

// 数据过期时间
const expireTime = 5 * 60

// RedisSetToken 设置会话token
func RedisSetToken(token string, role int) {
	_, err := Do("Set", TokenKey+token, role, "EX", expireTime)
	if err != nil {
		log.Error("set token fail:%v", err)
	}
}

// RedisGetToken 设置会话token
func RedisGetToken(token string) int {
	data, err := Do("Get", TokenKey+token)
	if err != nil {
		log.Error("get token fail:%v", err)
		return -1
	}
	role, ok := data.([]uint8)
	if !ok {
		log.Error("unknown token %v, role:%v", token, data)
		return -1
	}
	return int(role[0])
}

// RedisSetReport 设置report
func RedisSetReport(data []byte, matchID, start, end string) {
	_, err := Do("Set", MatchReportKey+matchID+start+end, data, "EX", expireTime)
	if err != nil {
		log.Error("set report fail:%v", err)
		return
	}
	return
}

// RedisGetReport 获取report
func RedisGetReport(matchID, start, end string) []byte {
	data, err := Do("Get", MatchReportKey+matchID+start+end)
	if err != nil {
		log.Error("get report fail:%v", err)
		return nil
	}
	if data == nil {
		return nil
	}
	ret, ok := data.([]byte)
	if !ok {
		log.Error("get report fail %v", ret)
		return nil
	}
	return ret
}

// RedisSetMatchList 设置matchlist
func RedisSetMatchList(data []byte, matchType, start, end string) {
	_, err := Do("Set", MatchListKey+matchType+start+end, data, "EX", expireTime)
	if err != nil {
		log.Error("set matchList fail:%v", err)
		return
	}
	return
}

// RedisGetMatchList 获取matchlist
func RedisGetMatchList(matchType, start, end string) []byte {
	data, err := Do("Get", MatchListKey+matchType+start+end)
	if err != nil {
		log.Error("set matchList fail:%v", err)
		return nil
	}
	ret, ok := data.([]byte)
	if !ok {
		log.Error("get report fail %v", ret)
		return nil
	}
	return ret
}

// RedisSetToken 设置会话token
func RedisSetTokenExport(token string, active bool) {
	_, err := Do("Set", TokenExportKey+token, active, "EX", expireTime*100)
	if err != nil {
		log.Error("set token fail:%v", err)
	}
}

// RedisGetToken 设置会话token
func RedisGetTokenExport(token string) bool {
	_, err := Do("Get", TokenExportKey+token)
	if err != nil {
		log.Error("get token fail:%v", err)
		return false
	}
	return true
}

func RedisDelTokenExport(token string) {
	_, err := Do("Del", TokenExportKey+token)
	if err != nil {
		log.Error("get token fail:%v", err)
	}
}

func RedisSetTokenUsrn(token string, usrn string) {
	_, err := Do("Set", TokenUsrn+token, usrn, "EX", expireTime)
	if err != nil {
		log.Error("set token fail: %v. ", err)
	}
}

func RedisGetTokenUsrn(token string) string {
	data, err := Do("Get", TokenUsrn+token)
	if err != nil {
		log.Error("get token fail: %v. ", err)
		return ""
	}

	res, ok := data.([]uint8)
	if !ok {
		log.Error("data not []uint8")
		return ""
	}

	return string(res)
}
