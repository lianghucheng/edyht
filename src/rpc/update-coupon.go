package rpc

import (
	"bs/config"
	"bytes"
	"encoding/json"
	"github.com/szxby/tools/log"
	"net/http"
)

const updateCouponUri = "/update-coupon"

type UpdateCouponReq struct {
	Secret    string `json:"secret"`
	Accountid int    `json:"accountid"`
	Amount    int    `json:"amount"`
}

func RpcUpdateCoupon(aid int, amount int) {
	log.Debug("【rpc更新用户点券】")
	j := new(UpdateCouponReq)
	j.Secret = secret
	j.Accountid = aid
	j.Amount = amount
	b, err := json.Marshal(j)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req, err := http.NewRequest("GET", config.GetConfig().GameServer+updateCouponUri, bytes.NewBuffer(b))
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = (&http.Client{}).Do(req)
	if err != nil {
		log.Error(err.Error())
	}
}
