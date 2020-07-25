package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/szxby/tools/log"
	"net/http"
)

const updateHeadImgUri = "/update-headimg"

type UpdateHeadImgReq struct {
	Secret    string `json:"secret"`
	Accountid int    `json:"accountid"`
	HeadImg   string `json:"headImg"`
}

func RpcUpdateHeadImg(aid int, headimg string) {
	log.Debug("【rpc更新用户头像】")
	j := new(UpdateHeadImgReq)
	j.Secret = secret
	j.Accountid = aid
	j.HeadImg = headimg
	log.Debug("*************%+v", *j)
	b, err := json.Marshal(j)
	if err != nil {
		log.Error(err.Error())
		return
	}
	req, err := http.NewRequest("GET", host+port+updateHeadImgUri, bytes.NewBuffer(b))
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
