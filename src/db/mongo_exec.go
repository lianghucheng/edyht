package db

import (
	"bs/param"
	"bs/util"
	"encoding/json"
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// GetUser 获取用户信息
func GetUser(account string) (user *util.User) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	user = new(util.User)
	err := s.DB(DB).C("users").Find(bson.M{"account": account}).One(user)
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

func readOneByQuery(rt interface{}, query bson.M, coll string) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	if err := se.DB(GDB).C(coll).Find(query).One(rt); err != nil {
		log.Error(err.Error())
	}
}

func readAllByQueryPage(rt interface{}, query bson.M, coll string, page, per int) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	if err := se.DB(GDB).C(coll).Find(query).Skip((page - 1) * per).Limit(per).Sort("-_id").All(rt); err != nil {
		log.Error(err.Error())
	}
}
func readAllByQuery(rt interface{}, query bson.M, coll string) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	if err := se.DB(GDB).C(coll).Find(query).Sort("-_id").All(rt); err != nil {
		log.Error(err.Error())
	}
}

func countByQuery(query bson.M, coll string) int {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	count, err := se.DB(GDB).C(coll).Find(query).Count()
	if err != nil {
		log.Error(err.Error())
	}
	return count
}

func save(data interface{}, coll string, id int) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	if _, err := se.DB(GDB).C(coll).Upsert(bson.M{"_id": id}, data); err != nil {
		log.Error(err.Error())
	}
}

func update(selector, update bson.M, coll string) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	if _, err := se.DB(GDB).C(coll).Upsert(selector, update); err != nil {
		log.Error(err.Error())
	}
}

func ReadFlowDatas(r *param.FlowDataHistoryReq) (*[]util.FlowData, int) {
	page, per := 1, 10
	if r.Page > 0 {
		page = r.Page
	}
	if r.Per > 0 {
		per = r.Per
	}

	query := getQueryByHistoryReq(r)

	log.Debug("【query】%v  %v", query, (page-1)*per)

	flowDatas := new([]util.FlowData)
	readAllByQueryPage(flowDatas, query, "flowdata", page, per)

	count := countByQuery(query, "flowdata")
	return flowDatas, count
}

func ReadExports(r *param.FlowDataExportReq) *[]util.FlowData {
	query := getQueryByExortReq(r)

	flowDatas := new([]util.FlowData)
	readAllByQuery(flowDatas, query, "flowdata")
	return flowDatas
}

func getQueryByHistoryReq(r *param.FlowDataHistoryReq) bson.M {
	query := bson.M{}
	if r.Start != 0 || r.End != 0 {
		query = bson.M{"createdat": bson.M{"$gte": r.Start, "$lt": r.End}}
	}

	if len(r.Condition) > 0 {
		accountid, _ := strconv.Atoi(r.Condition)
		status, _ := strconv.Atoi(r.Condition)
		query["$or"] = []bson.M{
			{"accountid": accountid},
			{"status": status},
		}
	}
	return query
}

func getQueryByExortReq(r *param.FlowDataExportReq) bson.M {
	query := bson.M{}
	if r.Start != 0 || r.End != 0 {
		query = bson.M{"createdat": bson.M{"$gte": r.Start, "$lt": r.End}}
	}

	if len(r.Condition) > 0 {
		accountid, _ := strconv.Atoi(r.Condition)
		status, _ := strconv.Atoi(r.Condition)
		query["$or"] = []bson.M{
			{"accountid": accountid},
			{"status": status},
		}
	}
	return query
}

func ReadFlowDataByID(id int) *util.FlowData {
	query := bson.M{"_id": id}
	flowData := new(util.FlowData)
	readOneByQuery(flowData, query, "flowdata")
	return flowData
}

func SaveFlowData(data *util.FlowData) {
	save(data, "flowdata", data.ID)
}

func AddUserFee(flowData *util.FlowData) {
	update(bson.M{"_id": flowData.Userid}, bson.M{"$inc": bson.M{"fee": flowData.ChangeAmount}}, "users")
}

func AddUserTakenFee(flowData *util.FlowData) {
	update(bson.M{"_id": flowData.Userid}, bson.M{"$inc": bson.M{"takenfee": flowData.ChangeAmount}}, "users")
}

func ReadUserDataByUID(id int) *util.UserData {
	ud := new(util.UserData)
	readOneByQuery(ud, bson.M{"_id": id}, "users")
	return ud
}

func ReadBankCardByID(id int) *util.BankCard {
	bc := new(util.BankCard)
	readOneByQuery(bc, bson.M{"userid": id}, "bankcard")
	return bc
}
