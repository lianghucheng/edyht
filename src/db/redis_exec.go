package db

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"

	"github.com/szxby/tools/log"
)

// 设定好redis的key
const (
	TokenKey          = "token"       // token
	MatchReportKey    = "matchReport" // match report
	MatchListKey      = "matchList"   // match list
	TokenExportKey    = "tokenExport" // token export
	TokenUsrn         = "usrn"
	WhiteList         = "whiteList"         // whitelist
	FirstView         = "firstView"         // 财务总览
	MapAll            = "mapAll"            // 所有图
	MapLastMoney      = "mapLastMoney"      // 剩余额图
	MapTotalCharge    = "mapTotalCharge"    // 总充值图
	MapTotalAward     = "mapTotalAward"     // 总奖金发放图
	MapTotalCashout   = "mapTotalCashout"   // 总提现图
	WeekItemBuy       = "weekItemBuy"       // 周购买
	WeekItemUse       = "weekItemUse"       // 周消耗
	ItemUseList       = "itemUseList"       // 道具购买列表
	CashoutPercent    = "CashoutPercent"    // 提现数额次数占比
	ChargeDetail      = "ChargeDetail"      // 充值明细
	CashoutDetail     = "CashoutDetail"     // 提现明细
	MatchAwardPreview = "MatchAwardPreview" // 赛事奖金总览
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

// RedisCommonSetData 通用的存储数据方法
func RedisCommonSetData(key string, data interface{}) {
	store, err := json.Marshal(data)
	if err != nil {
		log.Error("err:%v", err)
		return
	}
	_, err = Do("Set", key, store, "EX", expireTime)
	if err != nil {
		log.Error("set data fail:%v", err)
		return
	}
	return
}

// RedisCommonGetData 通用的读取数据方法
func RedisCommonGetData(key string) []byte {
	data, err := Do("Get", key)
	if err != nil {
		log.Error("get data fail:%v", err)
		return nil
	}
	ret, ok := data.([]byte)
	if !ok {
		log.Error("get data fail %v", ret)
		return nil
	}
	return ret
}

// RedisCommonDelData 通用的读取数据方法
func RedisCommonDelData(key string) {
	_, err := Do("Del", key)
	if err != nil {
		log.Error("del data fail:%v", err)
		return
	}
}

func GetCaptchaCache(account string) (captcha string, err error) {
	captcha, err = redis.String(Do("GET", "captcha:"+account))
	return
}