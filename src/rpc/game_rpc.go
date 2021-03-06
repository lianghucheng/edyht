package rpc

import (
	"bs/config"
	"bs/util"
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

const (
	MailServiceTypeOfficial = 0
	MailServiceTypeMatch    = 1
	MailServiceTypeActivity = 2
)

type Annex struct {
	PropType int     `json:"prop_type"`
	Num      float64 `json:"num"`
	Desc     string  `json:"desc"`
}

type MailBoxReq struct {
	TargetID        int64   `json:"target_id"`         //目标用户， -1表示所有用户
	MailType        int     `json:"mail_type"`         //邮箱邮件类型
	MailServiceType int     `json:"mail_service_type"` //邮件服务类型
	Title           string  `json:"title"`             //主题
	Content         string  `json:"content"`           //内容
	Annexes         []Annex `json:"annexes"`           //附件
	ExpireValue     float64 `json:"expire_value"`      //有效时长
}

const pushMail = "/pushmail"

func RpcPushMail(req *MailBoxReq) error {
	log.Debug("RpcPushMail   %+v", req)
	req_buf, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Get(config.GetConfig().GameServer + pushMail + "?data=" + string(req_buf))
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

const setPropBaseConfig = "/set/propbaseconfig"

func RpcSetPropBaseConfig() error {
	log.Debug("RpcSetPropBaseConfig")
	resp, err := http.Get(config.GetConfig().GameServer + setPropBaseConfig)
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

const setBankcard = "/bind/lianhanghao"

func RpcSetBankcard(accountid, bankName, bankCardNo, province, city, openingBank, openingBankNo string) error {
	log.Debug("RpcSetBankcard")
	resp, err := http.Get(config.GetConfig().GameServer + setBankcard +
		`?accountid=` + accountid +
		`&bankName=` + bankName +
		`&bankCardNo=` + bankCardNo +
		`&province=` + province +
		`&city=` + city +
		`&openingBank=` + openingBank +
		`&openingBankNo=` + openingBankNo)
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

type RPC_HorseLamp struct {
	ID           int
	Template     string
	ExpiredAt    int    //过期时间戳
	TakeEffectAt int    //发布时间戳
	Duration     int    //间隔时长，单位s
	LinkMatchID  string //关联赛事id
	Level        int
}

const HorseStart = "/horse/start"

func RpcHorseStart(data *util.HorseRaceLampControl) error {
	log.Debug("RpcHorseStart")
	t := (data.Duration / 12) * 12
	if data.Duration%12 != 0 {
		t += 12
	}
	rpcMsg := &RPC_HorseLamp{
		ID:           data.ID,
		Template:     data.Content,
		ExpiredAt:    data.ExpiredAt,
		TakeEffectAt: data.TakeEffectAt,
		Duration:     t,
		LinkMatchID:  data.LinkMatchID,
		Level:        data.Level,
	}
	b, err := json.Marshal(rpcMsg)
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	req, err := http.NewRequest("GET", config.GetConfig().GameServer+HorseStart, bytes.NewBuffer(b))
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Debug(err.Error())
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
	if m.Code != util.Success {
		return errors.New(fmt.Sprintf("request fail, the code is %v. the errmsg is %v. ", m.Code, m.Errmsg))
	}
	return nil
}

const HorseStop = "/horse/stop"

func RpcHorseStop(data *util.HorseRaceLampControl) error {
	log.Debug("RpcHorseStart")
	rpcMsg := new(RPC_HorseLamp)
	if err := util.Transfer(data, rpcMsg); err != nil {
		log.Debug(err.Error())
		return err
	}
	b, err := json.Marshal(rpcMsg)
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	req, err := http.NewRequest("GET", config.GetConfig().GameServer+HorseStop, bytes.NewBuffer(b))
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Debug(err.Error())
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
	if m.Code != util.Success {
		return errors.New(fmt.Sprintf("request fail, the code is %v. the errmsg is %v. ", m.Code, m.Errmsg))
	}
	return nil
}

const ActivityNotify = "/activity/notify"
func RpcActivityNotify() {
	http.Get(config.GetConfig().GameServer+ActivityNotify)
}


const NoticeNotify = "/notice/notify"
func RpcNoticeNotify() {
	http.Get(config.GetConfig().GameServer+NoticeNotify)
}

