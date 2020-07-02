package rpc

import (
	"fmt"
	"github.com/szxby/tools/log"
	"net/http"
)

func AddFee(id int, amount float64, feeType string) {
	log.Debug("【rpc更新用户奖金】")
	_, err := http.Get(fmt.Sprintf(`http://localhost:9084/edyht-add-fee?data={"userid":%v,"amount":%v,"fee_type":"%v"}`, id, amount, feeType))
	if err != nil {
		log.Error(err.Error())
	}
}
