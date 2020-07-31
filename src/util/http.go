package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/szxby/tools/log"
)

// PostToGame http post to game
func PostToGame(url string, contentType string, send interface{}) error {
	params, err := json.Marshal(send)
	if err != nil {
		log.Error("http post call err:%v", err)
		return err
	}
	sign := CalculateHash(string(params))
	data := map[string]interface{}{"Data": string(params), "Sign": sign}
	reqStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("http post call err:%v", err)
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debug("response Body%v:", string(body))

	// 验证返回参数
	ret := map[string]interface{}{}
	json.Unmarshal(body, &ret)
	if ret["code"] == nil {
		log.Error("call game fail :%v", ret)
		return err
	}
	code, ok := ret["code"].(float64)
	if !ok || code != 0 {
		log.Error("call game fail :%v", ret)
		retMsg := "操作失败，请重试！"
		if ret["desc"] != nil {
			if msg, ok := ret["desc"].(string); ok {
				retMsg = msg
			}
		}
		return errors.New(retMsg)
	}
	return nil
}

// PostToGameResp http post to game and get resp
func PostToGameResp(url string, contentType string, send interface{}) (map[string]interface{}, error) {
	params, err := json.Marshal(send)
	if err != nil {
		log.Error("http post call err:%v", err)
		return nil, err
	}
	sign := CalculateHash(string(params))
	data := map[string]interface{}{"Data": string(params), "Sign": sign}
	reqStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("http post call err:%v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debug("response Body%v:", string(body))

	// 验证返回参数
	ret := map[string]interface{}{}
	json.Unmarshal(body, &ret)
	if ret["code"] == nil {
		log.Error("call game fail :%v", ret)
		return nil, err
	}
	code, ok := ret["code"].(float64)
	if !ok || code != 0 {
		log.Error("call game fail :%v", ret)
		retMsg := "操作失败，请重试！"
		if ret["desc"] != nil {
			if msg, ok := ret["desc"].(string); ok {
				retMsg = msg
			}
		}
		return nil, errors.New(retMsg)
	}
	return ret, nil
}
