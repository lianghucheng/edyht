package edy_api

import (
	"bs/edy_api/internal/base"
	"bs/util"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/szxby/tools/log"
)

func checkCode(data []byte) error {
	log.Debug("data:%v", string(data))
	tmp := map[string]interface{}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		log.Error("err:%v", err)
		return err
	}
	if tmp["resp_code"] == nil {
		log.Error("err ret:%+v", tmp)
		return errors.New("unknow err")
	}
	s, ok := tmp["resp_code"].(string)
	if !ok {
		log.Error("err ret:%+v", tmp)
		return errors.New("unknow err")
	}
	if s != "000000" {
		log.Error("err ret:%+v", tmp)
		return fmt.Errorf("err:%v", s)
	}
	return nil
}

// MatchInfoQuery 查询比赛信息
func MatchInfoQuery(data util.MatchInfoQuery) (map[string]interface{}, error) {
	data.Cp_id = base.CpID
	str, _ := json.Marshal(data)
	log.Debug("str:%v", string(str))
	c := base.NewClient("/edy/match/info", string(str), base.ReqPost)
	c.GenerateSign(base.ReqPost)
	ret, err := c.DoPost()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// SonMatchQuery 子比赛信息查询
func SonMatchQuery(data util.SonMatchQuery) (map[string]interface{}, error) {
	data.Cp_id = base.CpID
	str, _ := json.Marshal(data)
	log.Debug("str:%v", string(str))
	c := base.NewClient("/edy/match/subInfo", string(str), base.ReqPost)
	c.GenerateSign(base.ReqPost)
	ret, err := c.DoPost()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerCashout 玩家提现
func PlayerCashout(data util.PlayerCashoutReq) (map[string]interface{}, error) {
	data.Cp_id = base.CpID
	str, _ := json.Marshal(data)
	log.Debug("str:%v", string(str))
	c := base.NewClient("/player/bonous/withdraw", string(str), base.ReqPost)
	c.GenerateSign(base.ReqPost)
	ret, err := c.DoPost()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// AwardResult 发奖结果查询接口
func AwardResult(data util.AwardResultReq) (map[string]interface{}, error) {
	// data.Cp_id = base.CpID
	// str, _ := json.Marshal(data)
	// log.Debug("str:%v", string(str))
	c := base.NewClient("/edy/bonuses/status", fmt.Sprintf("cp_id=%v&match_id=%v&page=%v&page_size=%v", base.CpID,
		data.Match_id, data.Page, data.Page_size), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerWalletInfoQuery 用户钱包余额及明细查询
func PlayerWalletInfoQuery(data util.PlayerWalletInfoQuery) (map[string]interface{}, error) {
	// data.Cp_id = base.CpID
	// str, _ := json.Marshal(data)
	// log.Debug("str:%v", string(str))
	c := base.NewClient("/wallet/balance", fmt.Sprintf("cp_id=%v&player_id=%v&page=%v&page_size=%v", base.CpID,
		data.Player_id, data.Page, data.Page_size), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerWalletBalanceQuery 用户钱包余额查询
func PlayerWalletBalanceQuery(data util.PlayerWalletBalanceQuery) (map[string]interface{}, error) {
	// data.Cp_id = base.CpID
	// str, _ := json.Marshal(data)
	// log.Debug("str:%v", string(str))
	c := base.NewClient("/wallet/accountBalance", fmt.Sprintf("cp_id=%v&player_id=%v", base.CpID, data.Player_id), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerWalletListQuery 用户钱包明細查询
func PlayerWalletListQuery(data util.PlayerWalletListQuery) (map[string]interface{}, error) {
	// data.Cp_id = base.CpID
	// str, _ := json.Marshal(data)
	// log.Debug("str:%v", string(str))
	c := base.NewClient("/wallet/accountBalanceList", fmt.Sprintf("cp_id=%v&player_id=%v&page=%v&page_size=%v", base.CpID,
		data.Player_id, data.Page, data.Page_size), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerMasterScoreQuery 玩家大师分查询
func PlayerMasterScoreQuery(data util.PlayerMasterScoreQuery) (map[string]interface{}, error) {
	// data.Cp_id = base.CpID
	// str, _ := json.Marshal(data)
	// log.Debug("str:%v", string(str))
	c := base.NewClient("/rating/by_identity_number", fmt.Sprintf("cp_id=%v&player_id_number=%v", base.CpID,
		data.Player_id_number), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerMasterScoreQuery2 玩家按身份证查询大师分
func PlayerMasterScoreQuery2(data util.PlayerMasterScoreQuery) (map[string]interface{}, error) {
	// data.Cp_id = base.CpID
	// str, _ := json.Marshal(data)
	// log.Debug("str:%v", string(str))
	c := base.NewClient("/rating/by_player_id_number", fmt.Sprintf("cp_id=%v&player_id_number=%v", base.CpID,
		data.Player_id_number), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerMasterScoreMatchQuery 按赛事ID查询大师分
func PlayerMasterScoreMatchQuery(matchID string) (map[string]interface{}, error) {
	c := base.NewClient("/rating/by_match_id", fmt.Sprintf("cp_id=%v&match_id=%v", base.CpID,
		matchID), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// AllMasterScoreQuery 厂商用户大师分排名列表查询
func AllMasterScoreQuery(page, size int) (map[string]interface{}, error) {
	c := base.NewClient("/rating/rank/list", fmt.Sprintf("cp_id=%v&page=%v&page_size=%v", base.CpID,
		page, size), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// CountryMasterScoreQuery 全国大师分排名查询
func CountryMasterScoreQuery(page, size int) (map[string]interface{}, error) {
	c := base.NewClient("/rating/rank/all/list", fmt.Sprintf("cp_id=%v&page=%v&page_size=%v", base.CpID,
		page, size), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerMasterScoreMatchDetail 玩家大师分详情查询
func PlayerMasterScoreMatchDetail(accountID, matchID, rank string) (map[string]interface{}, error) {
	c := base.NewClient("/rating/detail", fmt.Sprintf("cp_id=%v&player_id=%v&match_id=%v&rank=%v", base.CpID,
		accountID, matchID, rank), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerWalletTransaction 用户钱包交易接口
func PlayerWalletTransaction(data util.PlayerWalletTransaction) (map[string]interface{}, error) {
	data.Cp_id = base.CpID
	str, _ := json.Marshal(data)
	log.Debug("str:%v", string(str))
	c := base.NewClient("/wallet/pay", string(str), base.ReqPost)
	c.GenerateSign(base.ReqPost)
	ret, err := c.DoPost()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}

// PlayerWalletTransactionStatus 用户钱包交易状态查询接口
func PlayerWalletTransactionStatus(orderID string) (map[string]interface{}, error) {
	c := base.NewClient("/wallet/pay/status", fmt.Sprintf("cp_id=%v&order_id=%v", base.CpID,
		orderID), base.ReqGet)
	c.GenerateSign(base.ReqGet)
	ret, err := c.DoGet()
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	if err := checkCode(ret); err != nil {
		return nil, err
	}
	msg := map[string]interface{}{}
	if err := json.Unmarshal(ret, &msg); err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	return msg, nil
}
