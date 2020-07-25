package rpc

import (
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
