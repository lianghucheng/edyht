package rpc

import (
	"fmt"
	"github.com/szxby/tools/log"
	"io/ioutil"
	"net/http"
)

func RobotTotalConf(total int, matchid string) {
	log.Debug("!!!!!!!!!!1rpc RobotTotalConf")
	resp, err := http.Get(fmt.Sprintf("http://localhost:9086/conf/num-limit?robot_total=%v&matchid=%v", total, matchid))
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

func RobotStatusConf(status int, matchid string) {
	log.Debug("!!!!!!!!!!!1rpc RobotStatusConf")
	resp, err := http.Get(fmt.Sprintf("http://localhost:9086/conf/robot-status?status=%v&matchid=%v", status, matchid))
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
