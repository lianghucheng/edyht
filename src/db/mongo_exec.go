package db

import (
	"bs/param"
	"bs/util"
	"encoding/json"
<<<<<<< HEAD
	"time"

=======
>>>>>>> 5860ec9285e6561468ebe697ae2e7244761219cd
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2"
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
	total, _ := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"lt": util.Delete}}).Count()
	// iter := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"gte": 0}}).Sort("-shelftime").Skip((page - 1) * count).Limit(count).Iter()
	iter := s.DB(GDB).C("matchmanager").Pipe([]bson.M{
		{"$match": bson.M{"state": bson.M{"$gte": 0}}},
		{"$skip": (page - 1) * count},
		{"$limit": count},
		{"$project": bson.M{
			"MatchID":   "$matchid",
			"MatchType": "$matchtype",
			"MatchIcon": "$matchicon",
			"RoundNum":  "$roundnum",
			"StartTime": "$starttime",
			"StartType": "$matchdesc",
			"MatchInfo": "$recommend",
			"Eliminate": "$eliminate",
			"EnterFee":  "$enterfee",
			"UseCount":  "$usematch",
			"LastMatch": bson.M{"$subtract": []interface{}{"$totalmatch", "$usermatch"}},
			"ShelfTime": "$shelftime",
			"ShowHall":  "$showhall",
			"Sort":      "$sort",
		}},
		{"$sort": "-shelftime"},
	}).Iter()
	for iter.Next(&one) {
		tmp, _ := json.Marshal(one)
		list = append(list, tmp)
	}
	return list, total
}

// GetMatchReport 获取比赛报表
func GetMatchReport(matchID string, start, end int64) [][]byte {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	var oneDay int64 = 24 * 60 * 60
	len := (end - start) / oneDay
	if len <= 0 {
		log.Error("invalid time")
		return nil
	}

	// 查询时间范围内的数据总合
	allReport := struct {
		AllSignPlayer int
		AllSignFee    float64
		AllAward      float64
		AllLast       float64
	}{}

	result := make([][]byte, len)
	for i := start; i+oneDay <= end; i += oneDay {
		one := map[string]interface{}{}
		err := s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"matchid": matchID}},
			{"$match": bson.M{"createtime": bson.M{"$gt": i, "$lte": i + oneDay}}},
			{"$project": bson.M{
				"RecordTime":  time.Unix(i, 0).Format("2006-01-02"),
				"SignInCount": bson.M{"$size": "$signinplayers"}, "_id": 0, "matchid": "1001",
				"SignFee":  bson.M{"$multiply": []interface{}{bson.M{"$size": "$signinplayers"}, bson.M{"$divide": []interface{}{"$enterfee", 10}}}},
				"AwardNum": bson.M{"$size": "$award"},
				"Money":    "$moneyaward",
				"Coupon":   "$couponaward",
				"LastMoney": bson.M{"$subtract": []interface{}{bson.M{
					"$multiply": []interface{}{
						bson.M{"$size": "$signinplayers"}, bson.M{"$divide": []interface{}{"$enterfee", 10}}}},
					bson.M{"$add": []interface{}{"$moneyaward", bson.M{"$multiply": []interface{}{"$couponaward", 10}}}}}}}},
			{"$group": bson.M{
				"_id": "$matchid", "RecordTime": bson.M{"$first": "$RecordTime"}, "allMoney": bson.M{"$sum": "$Money"},
				"allCoupon": bson.M{"$sum": "$Coupon"}, "allSign": bson.M{"$sum": "$SignInCount"},
				"allSignFee": bson.M{"$sum": "$SignFee"}, "awardNum": bson.M{"$sum": "$AwardNum"},
				"lastMoney": bson.M{"$sum": "$LastMoney"}}},
			// {"$sort": bson.M{"count": -1}},
		}).One(&one)
		if err != nil {
			log.Error("get report fail:%v", err)
			return nil
		}
		data, err := json.Marshal(one)
		if err != nil {
			log.Error("get report fail:%v", err)
			return nil
		}
		result = append(result, data)
		// 数据汇总
		allReport.AllSignPlayer += one["allSign"].(int)
		allReport.AllSignFee += one["allSignFee"].(float64)
		allReport.AllAward += one["allMoney"].(float64)
		allReport.AllLast += one["lastMoney"].(float64)
	}
	// 最后一位保存汇总数据
	all, err := json.Marshal(allReport)
	if err != nil {
		log.Error("get report fail:%v", err)
		return nil
	}
	result = append(result, all)
	return result
}

// GetMatch 获取单场赛事
func GetMatch(matchID string) []byte {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	one := map[string]interface{}{}
	err := s.DB(GDB).C("match").Pipe([]bson.M{
		{"$match": bson.M{"matchid": matchID}},
		{"$project": bson.M{
			"MatchType":  "$matchtype",
			"MatchName":  "$matchname",
			"MatchID":    "$matchid",
			"CreateTime": "$createtime",
			"RoundNum":   "$roundnum",
			"StartType":  "$matchdesc",
			"MatchInfo":  "$recommend",
			"Eliminate":  "$eliminate",
			"EnterFee":   "$enterfee",
		}},
		// {"$group": bson.M{
		// 	"_id": "$matchid", "RecordTime": bson.M{"$first": "$RecordTime"}, "allMoney": bson.M{"$sum": "$Money"},
		// 	"allCoupon": bson.M{"$sum": "$Coupon"}, "allSign": bson.M{"$sum": "$SignInCount"},
		// 	"allSignFee": bson.M{"$sum": "$SignFee"}, "awardNum": bson.M{"$sum": "$AwardNum"},
		// "lastMoney": bson.M{"$sum": "$LastMoney"}}},
		// {"$sort": bson.M{"CreateTime": -1}},
	}).One(&one)
	if err != nil {
		log.Error("get match fail %v", err)
		return nil
	}
	ret, err := json.Marshal(one)
	if err != nil {
		log.Error("get match fail %v", err)
		return nil
	}
	return ret
}

// GetMatchList 获取某个时间段的赛事
func GetMatchList(matchType string, start, end int64) [][]byte {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	var oneDay int64 = 24 * 60 * 60
	length := (end - start) / oneDay
	if length <= 0 {
		log.Error("invalid time:%v,%v", start, end)
		return nil
	}

	var result [][]byte
	one := map[string]interface{}{}
	var iter *mgo.Iter
	if len(matchType) == 0 {
		iter = s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"createtime": bson.M{"$gt": start, "$lte": end}}},
			{"$project": bson.M{
				"MatchType":  "$matchtype",
				"MatchName":  "$matchname",
				"MatchID":    "$matchid",
				"CreateTime": "$createtime",
				"RoundNum":   "$roundnum",
				"StartType":  "$matchdesc",
				"MatchInfo":  "$recommend",
				"Eliminate":  "$eliminate",
				"EnterFee":   "$enterfee",
			}},
			{"$sort": bson.M{"CreateTime": 1}},
		}).Iter()
	} else {
		iter = s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"matchtype": matchType}},
			{"$match": bson.M{"createtime": bson.M{"$gt": start, "$lte": end}}},
			{"$project": bson.M{
				"MatchType":  "$matchtype",
				"MatchName":  "$matchname",
				"MatchID":    "$matchid",
				"CreateTime": "$createtime",
				"RoundNum":   "$roundnum",
				"StartType":  "$matchdesc",
				"MatchInfo":  "$recommend",
				"Eliminate":  "$eliminate",
				"EnterFee":   "$enterfee",
			}},
			{"$sort": bson.M{"CreateTime": 1}},
		}).Iter()
	}
	for iter.Next(&one) {
		data, err := json.Marshal(one)
		if err != nil {
			log.Error("get report fail:%v", err)
			return nil
		}
		result = append(result, data)
	}
	return result
}

// GetMatchDetail 获取一局战绩详情
func GetMatchDetail(matchID string) []byte {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	one := map[string]interface{}{}
	err := s.DB(GDB).C("match").Pipe([]bson.M{
		{"$match": bson.M{"matchid": matchID}},
		{"$project": bson.M{
			"Rank":        "$rank",
			"MatchRecord": "$matchrecord",
		}},
		// {"$group": bson.M{
		// 	"_id": "$matchid", "RecordTime": bson.M{"$first": "$RecordTime"}, "allMoney": bson.M{"$sum": "$Money"},
		// 	"allCoupon": bson.M{"$sum": "$Coupon"}, "allSign": bson.M{"$sum": "$SignInCount"},
		// 	"allSignFee": bson.M{"$sum": "$SignFee"}, "awardNum": bson.M{"$sum": "$AwardNum"},
		// "lastMoney": bson.M{"$sum": "$LastMoney"}}},
		// {"$sort": bson.M{"CreateTime": -1}},
	}).One(&one)
	if err != nil {
		log.Error("get detail fail %v", err)
		return nil
	}
	data, err := json.Marshal(one)
	if err != nil {
		log.Error("get detail fail %v", err)
		return nil
	}
	return data
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
