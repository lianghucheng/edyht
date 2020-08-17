package rpc

import (
	"bs/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/szxby/tools/log"
	"io/ioutil"
	"net/http"
)

func AddFee(id int, amount float64, feeType string) {
	log.Debug("【rpc更新用户奖金】")
	_, err := http.Get(fmt.Sprintf(`http://localhost:9084/edyht-add-fee?data={"userid":%v,"amount":%v,"fee_type":"%v"}`, id, amount, feeType))
	if err != nil {
		log.Error(err.Error())
	}
}

func MatchMaxRobotNumConf(maxRobotNum int, matchid string) {
	log.Debug("!!!!!!!!!!!!!rpc MatchMaxRobotNumConf")
	resp, err := http.Get(fmt.Sprintf("http://localhost:9084/conf/robot-maxnum?max_robot_num=%v&matchid=%v", maxRobotNum, matchid))
	if err != nil {
		log.Error(err.Error())
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Debug(err.Error())
		return
	}
	log.Debug(string(b))
}

type RPC_AddCouponFrag struct {
	Secret    string
	Accountid int
	Amount    int
}

const rpcAddCounponFragUrl = "/add/coupon-frag"

func RpcAddCouponFrag(aid, amount int) {
	log.Debug("远程调用加点券碎片")
	data := new(RPC_AddCouponFrag)
	data.Secret = secret
	data.Amount = amount
	data.Accountid = aid
	b, errJson := json.Marshal(data)
	req, errNewReq := http.NewRequest("GET", host+port+rpcAddCounponFragUrl, bytes.NewBuffer(b))
	client := &http.Client{}
	resp, errHttpDo := client.Do(req)
	b_resp, errReadIO := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if errJson != nil || errNewReq != nil || errHttpDo != nil || errReadIO != nil {
		log.Error("errJson:%v, errNewReq:%v, errHttpDo:%v, errReadIO:%v", errJson, errNewReq, errHttpDo, errReadIO)
		return
	}
	log.Debug(string(b_resp))
}

const rpcNotifyPayAccount = "/notify/payaccount"

func RpcNotifyPayAccount() error {
	log.Debug("RpcNotifyPayAccount")
	resp, err := http.Get(host + port + rpcNotifyPayAccount)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	m := struct {
		Code   int
		Errmsg string
	}{}
	if err := json.Unmarshal(buf, &m); err != nil {
		return err
	}
	if m.Code != 10000 {
		return errors.New(fmt.Sprintf("request fail, the code is %v. ", m.Code))
	}
	return nil
}

const rpcNotidyPriceMenu = "/notify/pricemenu"

func RpcNotifyPriceMenu() error {
	log.Debug("RpcNotifyPriceMenu")
	resp, err := http.Get(config.GetConfig().GameServer + rpcNotidyPriceMenu)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	m := struct {
		Code   int
		Errmsg string
	}{}
	if err := json.Unmarshal(buf, &m); err != nil {
		return err
	}
	if m.Code != 10000 {
		return errors.New(fmt.Sprintf("request fail, the code is %v. ", m.Code))
	}
	return nil
}

const rpcNotidyGoodsType = "/notify/goodstype"

func RpcNotifyGoodsType() error {
	log.Debug("RpcNotifyGoodsType")
	resp, err := http.Get(config.GetConfig().GameServer + rpcNotidyGoodsType)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	m := struct {
		Code   int
		Errmsg string
	}{}
	if err := json.Unmarshal(buf, &m); err != nil {
		return err
	}
	if m.Code != 10000 {
		return errors.New(fmt.Sprintf("request fail, the code is %v. ", m.Code))
	}
	return nil
}


const (
	MailTypeText  = 1
	MailTypeAward = 2
	MailTypeMix   = 3
)

type MailBox struct {
	ID          int64   `bson:"_id"`          //唯一id
	TargetID    int64   `json:"target_id"`    //目标用户， -1表示所有用户
	MailType    int     `json:"mail_type"`    //邮箱邮件类型
	CreatedAt   int64   `json:"created_at"`   //收件时间
	Title       string  `json:"title"`        //主题
	Content     string  `json:"content"`      //内容
	Annexes     []Annex `json:"annexes"`      //附件
	Status      int64   `json:"status"`       //邮箱邮件状态
	ExpireValue int64   `json:"expire_value"` //有效时长
	MailServiceType int //邮件服务类型
}

const (
	AnnexTypeCoupon = 1
	AnnexTypeCouponFrag = 2
)

type Annex struct {
	Type int    `json:"type"`
	Num  int    `json:"num"`
	Desc string `json:"desc"`
}

type MailBoxReq struct {
	TargetID    int64   `json:"target_id"`    //目标用户， -1表示所有用户
	MailType    int     `json:"mail_type"`    //邮箱邮件类型
	MailServiceType int `json:"mail_service_type"`//邮件服务类型
	Title       string  `json:"title"`        //主题
	Content     string  `json:"content"`      //内容
	Annexes     []Annex `json:"annexes"`      //附件
	ExpireValue int64   `json:"expire_value"` //有效时长
}

const pushMail = "/pushmail"

func RpcPushMail(req *MailBoxReq) error {
	log.Debug("RpcPushMail   %+v", req)
	req_buf,err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Get(config.GetConfig().GameServer + pushMail+"?data="+string(req_buf))
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	m := struct {
		Code   int
		Errmsg string
	}{}
	if err := json.Unmarshal(buf, &m); err != nil {
		return err
	}
	if m.Code != 0 {
		return errors.New(fmt.Sprintf("request fail, the code is %v. the errmsg is %v. ", m.Code, m.Errmsg))
	}
	return nil
}