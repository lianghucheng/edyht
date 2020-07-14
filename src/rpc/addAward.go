package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/szxby/tools/log"
	"net/http"
)

const (
	addAwardUri = "/addaward"
)

type AddAwardReq struct {
	Secret 		string `json:"secret"`
	Uid 	int `json:"uid"`
	Amount 	float64 `json:"amount"`
}

func AddAward(aid int, amount float64) {
	log.Debug("【rpc更新用户税后奖金】")
	addAward := new(AddAwardReq)
	addAward.Secret = secret
	addAward.Uid = aid
	addAward.Amount = amount
	b, err := json.Marshal(addAward)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req, err := http.NewRequest("GET", host+port+addAwardUri, bytes.NewBuffer(b))
	_, err = (&http.Client{}).Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
}