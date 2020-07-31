package rpc

import (
	"bytes"
	"encoding/json"
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
