package db

import (
	"bs/util"
	"encoding/json"

	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
)

// GetUser 获取用户信息
func GetUser(account string) (user *util.User) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	err := s.DB(DB).C("users").Find(bson.M{"user": account}).One(user)
	if err != nil {
		log.Error("get user %v err:%v", account, err)
		return nil
	}
	return
}

// GetMatchManagerList 获取比赛类型列表
func GetMatchManagerList(page int, count int) ([][]byte, int) {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)
	one := util.MatchManager{}
	list := [][]byte{}
	total, _ := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"gte": 0}}).Count()
	iter := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"gte": 0}}).Sort("-shelftime").Skip((page - 1) * count).Limit(count).Iter()
	for iter.Next(&one) {
		tmp, _ := json.Marshal(one)
		list = append(list, tmp)
	}
	return list, total
}

// GetMatchReport 获取比赛报表
func GetMatchReport(start, end int64, page int, count int) {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)
	// one := util.MatchManager{}
	// list := [][]byte{}
	// total, _ := s.DB(GDB).C("match").Find(bson.M{"state": bson.M{"gte": 0}}).Count()
	// iter := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"gte": 0}}).Sort("-shelftime").Skip((page - 1) * count).Limit(count).Iter()
	// for iter.Next(&one) {
	// 	tmp, _ := json.Marshal(one)
	// 	list = append(list, tmp)
	// }
	// return list, total
}
