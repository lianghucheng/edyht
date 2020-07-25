package db

import (
	"bs/param"
	"bs/param/base"
	"bs/util"
	"encoding/json"
	"fmt"
	"time"

	"strconv"

	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
func GetMatchManagerList(page int, count int) ([]map[string]interface{}, int) {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)
	// one := map[string]interface{}{}
	list := []map[string]interface{}{}
	total, _ := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"$lt": util.Delete}}).Count()
	// iter := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"gte": 0}}).Sort("-shelftime").Skip((page - 1) * count).Limit(count).Iter()
	err := s.DB(GDB).C("matchmanager").Pipe([]bson.M{
		{"$match": bson.M{"state": bson.M{"$lt": util.Delete}}},
		{"$project": bson.M{
			"MatchID":     "$matchid",
			"MatchName":   "$matchname",
			"MatchType":   "$matchtype",
			"MatchIcon":   "$matchicon",
			"RoundNum":    "$roundnum",
			"StartTime":   "$starttime",
			"StartType":   "$starttype",
			"LimitPlayer": "$limitplayer",
			"Recommend":   "$recommend",
			"Eliminate":   "$eliminate",
			"EnterFee":    "$enterfee",
			"UseCount":    "$usematch",
			"LastMatch":   bson.M{"$subtract": []interface{}{"$totalmatch", "$usematch"}},
			"ShelfTime":   "$shelftime",
			"ShowHall":    "$showhall",
			"Sort":        "$sort",
			"State":       "$state",
			"AwardList":   "$awardlist",
			"TotalMatch":  "$totalmatch",
			"_id":         0,
		}},
		{"$sort": bson.M{"Sort": 1}},
		{"$skip": (page - 1) * count},
		{"$limit": count},
	}).All(&list)
	// for iter.Next(&one) {
	// 	// tmp, _ := json.Marshal(one)
	// 	log.Debug("check:%v", one)
	// 	list = append(list, one)
	// }
	// ret, err := json.Marshal(list)
	if err != nil {
		log.Error("get match manager fail %v", err)
		return nil, 0
	}
	return list, total
}

// GetMatchReport 获取比赛报表
func GetMatchReport(matchID string, start, end int64) []map[string]interface{} {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	var oneDay int64 = 24 * 60 * 60
	len := (end - start) / oneDay
	if len <= 0 {
		log.Error("invalid time")
		return nil
	}

	// 查询时间范围内的数据总合
	// allReport := struct {
	// 	AllSignPlayer int
	// 	AllSignFee    float64
	// 	AllAward      float64
	// 	AllLast       float64
	// }{}
	allReport := map[string]interface{}{}
	allReport["AllSignPlayer"] = int(0)
	allReport["AllSignFee"] = float64(0)
	allReport["AllAward"] = float64(0)
	allReport["AllLast"] = float64(0)

	// log.Debug("check,%v,%v", start, end)
	result := make([]map[string]interface{}, 0)
	for i := start; i <= end; i = i + oneDay {
		one := map[string]interface{}{}
		// rt := time.Unix(i, 0).Format("2006-01-02")
		// log.Debug("check,%v,%v", i, i+oneDay)
		err := s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"matchid": matchID}},
			// {"$match": bson.M{"createtime": bson.M{"$gt": fmt.Sprintf("$%v", i), "$lte": fmt.Sprintf("$%v", i+oneDay)}}},
			{"$match": bson.M{"createtime": bson.M{"$gt": i, "$lte": i + oneDay}}},
			{"$project": bson.M{
				// "RecordTime":  fmt.Sprintf("$%v", time.Unix(i, 0).Format("2006-01-02")),
				"RecordTime":  "$createtime",
				"SignInCount": bson.M{"$size": "$signinplayers"}, "_id": 0, "matchid": fmt.Sprintf("$%v", matchID),
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
		if err == mgo.ErrNotFound {
			continue
		}
		if err != nil {
			log.Error("get report fail:%v", err)
			return nil
		}
		// 数据汇总
		allReport["AllSignPlayer"] = allReport["AllSignPlayer"].(int) + one["allSign"].(int)
		allReport["AllSignFee"] = allReport["AllSignFee"].(float64) + one["allSignFee"].(float64)
		allReport["AllAward"] = allReport["AllAward"].(float64) + one["allMoney"].(float64)
		allReport["AllLast"] = allReport["AllLast"].(float64) + one["lastMoney"].(float64)

		one["allSignFee"] = util.Decimal(one["allSignFee"].(float64))
		one["allMoney"] = util.Decimal(one["allMoney"].(float64))
		one["lastMoney"] = util.Decimal(one["lastMoney"].(float64))
		result = append(result, one)
	}
	// 最后一位保存汇总数据
	// all, err := json.Marshal(allReport)
	// if err != nil {
	// 	log.Error("get report fail:%v", err)
	// 	return nil
	// }

	allReport["AllSignFee"] = util.Decimal(allReport["AllSignFee"].(float64))
	allReport["AllAward"] = util.Decimal(allReport["AllAward"].(float64))
	allReport["AllLast"] = util.Decimal(allReport["AllLast"].(float64))
	result = append(result, allReport)
	// ret, _ := json.Marshal(result)
	return result
}

// GetMatch 获取单场赛事
func GetMatch(matchID string) []byte {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	one := map[string]interface{}{}
	err := s.DB(GDB).C("match").Pipe([]bson.M{
		{"$match": bson.M{"sonmatchid": matchID}},
		{"$project": bson.M{
			"MatchType":   "$matchtype",
			"MatchName":   "$matchname",
			"MatchID":     "$sonmatchid",
			"CreateTime":  "$createtime",
			"RoundNum":    "$roundnum",
			"LimitPlayer": "$limitplayer",
			"Recommend":   "$recommend",
			"StartType":   "$starttype",
			"Eliminate":   "$eliminate",
			"EnterFee":    "$enterfee",
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
func GetMatchList(matchType string, start, end int64) []map[string]interface{} {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	var oneDay int64 = 24 * 60 * 60
	length := (end - start) / oneDay
	if length <= 0 {
		log.Error("invalid time:%v,%v", start, end)
		return nil
	}

	var result []map[string]interface{}
	one := map[string]interface{}{}
	var iter *mgo.Iter
	if len(matchType) == 0 {
		iter = s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"createtime": bson.M{"$gt": start, "$lte": end + oneDay}}},
			{"$project": bson.M{
				"MatchType":   "$matchtype",
				"MatchName":   "$matchname",
				"MatchID":     "$sonmatchid",
				"CreateTime":  "$createtime",
				"RoundNum":    "$roundnum",
				"LimitPlayer": "$limitplayer",
				"Recommend":   "$recommend",
				"StartType":   "$starttype",
				"StartTime":   "$starttime",
				"Eliminate":   "$eliminate",
				"EnterFee":    "$enterfee",
				"_id":         0,
			}},
			{"$sort": bson.M{"CreateTime": 1}},
		}).Iter()
	} else {
		iter = s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"matchtype": matchType}},
			{"$match": bson.M{"createtime": bson.M{"$gt": start, "$lte": end + oneDay}}},
			{"$project": bson.M{
				"MatchType":   "$matchtype",
				"MatchName":   "$matchname",
				"MatchID":     "$sonmatchid",
				"CreateTime":  "$createtime",
				"RoundNum":    "$roundnum",
				"LimitPlayer": "$limitplayer",
				"Recommend":   "$recommend",
				"StartType":   "$starttype",
				"StartTime":   "$starttime",
				"Eliminate":   "$eliminate",
				"EnterFee":    "$enterfee",
				"_id":         0,
			}},
			{"$sort": bson.M{"CreateTime": 1}},
		}).Iter()
	}
	for iter.Next(&one) {
		// data, err := json.Marshal(one)
		// if err != nil {
		// 	log.Error("get report fail:%v", err)
		// 	return nil
		// }
		result = append(result, one)
		one = map[string]interface{}{}
	}
	log.Debug("result:%v", result)
	// if len(result) == 0 {
	// 	return nil
	// }
	// ret, err := json.Marshal(result)
	// if err != nil {
	// 	log.Error("get report fail:%v", err)
	// 	return nil
	// }
	return result
}

// GetMatchDetail 获取一局战绩详情
func GetMatchDetail(matchID string) map[string]interface{} {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	one := map[string]interface{}{}
	err := s.DB(GDB).C("match").Pipe([]bson.M{
		{"$match": bson.M{"sonmatchid": matchID}},
		{"$project": bson.M{
			"Rank":        "$rank",
			"MatchRecord": "$matchrecord",
			"_id":         0,
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
	// data, err := json.Marshal(one)
	// if err != nil {
	// 	log.Error("get detail fail %v", err)
	// 	return nil
	// }
	return one
}

func readOneByQuery(rt interface{}, query bson.M, coll string) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	if err := se.DB(GDB).C(coll).Find(query).One(rt); err != nil && err != mgo.ErrNotFound {
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

const (
	readTypeAll = 1
	readTypeOne = 2
)

func readByPipeline(db, coll string, pipeline []bson.M, rt interface{}, readtype int) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	var err error
	if readtype == readTypeAll {
		err = se.DB(db).C(coll).Pipe(pipeline).All(rt)
	} else if readtype == readTypeOne {
		err = se.DB(db).C(coll).Pipe(pipeline).One(rt)
	}
	if err != nil {
		log.Error(err.Error())
	}
}

func readGameByPipeline(db, coll string, pipeline []bson.M, rt interface{}, readtype int) {
	se := gameDB.Ref()
	defer gameDB.UnRef(se)
	var err error
	if readtype == readTypeAll {
		err = se.DB(db).C(coll).Pipe(pipeline).All(rt)
	} else if readtype == readTypeOne {
		err = se.DB(db).C(coll).Pipe(pipeline).One(rt)
	}
	if err != nil {
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

func save(db string, data interface{}, coll string, id int) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	if _, err := se.DB(db).C(coll).Upsert(bson.M{"_id": id}, data); err != nil {
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
	log.Debug("【开始时间】%v", time.Unix(r.Start, 0).Format("2006-03-04 05:06"))
	log.Debug("【结束时间】%v", time.Unix(r.End, 0).Format("2006-03-04 05:06"))
	query := bson.M{}
	if r.Start != 0 || r.End != 0 {
		query = bson.M{"createdat": bson.M{"$gte": r.Start, "$lt": r.End + 86400}}
	}

	if len(r.Condition) > 0 {
		accountid := 0
		status := 0
		if len(r.Condition) >= 2 {
			accountid, _ = strconv.Atoi(r.Condition[1])
			status, _ = strconv.Atoi(r.Condition[0])
			query["accountid"] = accountid
			query["status"] = status
		} else {
			c, _ := strconv.Atoi(r.Condition[0])
			if c > 10 {
				accountid = c
				query["accountid"] = accountid
			} else {
				status = c
				query["status"] = status
			}
		}

		if status != 0 {
			query["flowtype"] = 2
		}
	}
	return query
}

func getQueryByExortReq(r *param.FlowDataExportReq) bson.M {
	query := bson.M{}
	if r.Start != 0 || r.End != 0 {
		query = bson.M{"createdat": bson.M{"$gte": r.Start, "$lt": r.End + 86400}}
	}

	if len(r.Condition) > 0 {
		accountid := 0
		status := 0
		if len(r.Condition) >= 2 {
			accountid, _ = strconv.Atoi(r.Condition[1])
			status, _ = strconv.Atoi(r.Condition[0])
			query["accountid"] = accountid
			query["status"] = status
		} else {
			c, _ := strconv.Atoi(r.Condition[0])
			if c > 10 {
				accountid = c
				query["accountid"] = accountid
			} else {
				status = c
				query["status"] = status
			}
		}
		if status != 0 {
			query["flowtype"] = 2
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
	save(GDB, data, "flowdata", data.ID)
}

func ReadUserDataByUID(id int) *util.UserData {
	ud := new(util.UserData)
	readOneByQuery(ud, bson.M{"_id": id}, "users")
	return ud
}

func ReadUserDataByAID(aid int) *util.UserData {
	ud := new(util.UserData)
	readByPipeline(GDB, "users", []bson.M{{"$match": bson.M{"accountid": aid}}}, ud, readTypeOne)
	return ud
}

func ReadBankCardByID(id int) *util.BankCard {
	bc := new(util.BankCard)
	readOneByQuery(bc, bson.M{"userid": id}, "bankcard")
	return bc
}

// GetGameVersion 获取游戏版本号,下载地址
func GetGameVersion() (version string, url string) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	data := map[string]interface{}{}
	s.DB(DB).C("GameConfig").Find(bson.M{"GameName": "edy"}).One(&data)
	if data["GameVersion"] == nil || data["URL"] == nil {
		log.Error("no config:%v", data)
		return
	}
	version, ok := data["GameVersion"].(string)
	if !ok {
		log.Error("no config:%v", data)
		return
	}
	url, ok = data["URL"].(string)
	if !ok {
		log.Error("no config:%v", data)
		return
	}
	return
}

// GetUserList 获取用户列表
func GetUserList(page, count int) ([]util.UserData, int) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	data := []util.UserData{}
	total, _ := s.DB(GDB).C("users").Find(bson.M{}).Count()
	iter := s.DB(GDB).C("users").Find(bson.M{}).Sort("-createdat").Skip((page - 1) * count).Limit(count).Iter()
	// if err != nil && err != mgo.ErrNotFound {
	// 	log.Error("err:%v", err)
	// 	return nil, total
	// }
	one := util.UserData{}
	for iter.Next(&one) {
		bank := ReadBankCardByID(one.UserID)
		one.BankCard = bank
		// 查询充值
		fee := map[string]interface{}{}
		s.DB("czddz").C("wxpayresult").Pipe([]bson.M{
			{"$match": bson.M{"success": true, "userid": one.AccountID + 1e8}},
			{"$project": bson.M{
				"TotalFee": "$totalfee",
			}},
			{"$group": bson.M{
				"_id": "$userid",
				"all": bson.M{"$sum": "$TotalFee"},
			}},
		}).One(&fee)
		var chargeAmount int64
		if feeAdd, ok := fee["all"].(int); ok {
			chargeAmount = int64(feeAdd)
		}
		one.ChargeAmount = chargeAmount
		data = append(data, one)
		one = util.UserData{}
	}
	return data, total
}

// GetOneUser 获取单个用户列表
func GetOneUser(accountID int, nickname string) (*util.UserData, error) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	data := &util.UserData{}
	var err error
	if accountID > 0 {
		err = s.DB(GDB).C("users").Find(bson.M{"accountid": accountID}).One(data)
	} else if len(nickname) > 0 {
		err = s.DB(GDB).C("users").Find(bson.M{"nickname": nickname}).One(data)
	}
	// err := s.DB(GDB).C("users").Find(bson.M{"$or": []interface{}{bson.M{"accountid": accountID}, bson.M{"nickname": nickname}}}).One(data)
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	bank := ReadBankCardByID(data.UserID)
	data.BankCard = bank
	return data, nil
}

// GetMatchReview 获取赛事总览
func GetMatchReview(uid int) ([]map[string]interface{}, []map[string]interface{}, map[string]interface{}) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	data := []map[string]interface{}{}
	matchs := []map[string]interface{}{}
	ids := []map[string]interface{}{}
	s.DB(GDB).C("matchreview").Pipe([]bson.M{
		{"$match": bson.M{"accountid": uid}},
		{"$group": bson.M{"_id": "$matchtype"}},
	}).All(&matchs)
	log.Debug("matchs:%+v", matchs)
	s.DB(GDB).C("matchreview").Pipe([]bson.M{
		{"$match": bson.M{"accountid": uid}},
		{"$group": bson.M{"_id": "$matchid"}},
	}).All(&ids)
	log.Debug("ids:%+v", ids)
	for _, matchType := range matchs {
		for _, id := range ids {
			// one := util.UserMatchReview{}
			one := map[string]interface{}{}
			// s.DB(GDB).C("matchreview").Find(bson.M{"accountid": uid, "matchtype": matchType["_id"], "matchid": id["_id"]}).One(&one)
			err := s.DB(GDB).C("matchreview").Pipe([]bson.M{
				{"$match": bson.M{"accountid": uid, "matchtype": matchType["_id"], "matchid": id["_id"]}},
				// {"$group": bson.M{"_id": "$matchtype", "matchtype": "$matchtype", "matchid": "$matchid", "matchtotal": "$matchtotal",
				// 	"matchwins": "$matchwins", "matchfails": "$matchfails", "coupon": "$coupon", "awardmoney": "$awardmoney", "personalprofit": "$personalprofit"}},
				{"$group": bson.M{"_id": "$matchtype", "matchtype": bson.M{"$first": "$matchtype"}, "matchid": bson.M{"$first": "$matchid"},
					"matchtotal": bson.M{"$sum": "$matchtotal"}, "matchwins": bson.M{"$sum": "$matchwins"}, "matchfails": bson.M{"$sum": "$matchfails"},
					"coupon": bson.M{"$sum": "$coupon"}, "awardmoney": bson.M{"$sum": "$awardmoney"}, "personalprofit": bson.M{"$sum": "$personalprofit"}}},
			}).One(&one)
			if err != nil && err != mgo.ErrNotFound {
				log.Error("err:%v", err)
				continue
			}
			if len(one) > 0 {
				data = append(data, one)
			}
		}
	}
	all := map[string]interface{}{}
	if err := s.DB(GDB).C("matchreview").Pipe([]bson.M{
		{"$match": bson.M{"accountid": uid}},
		{"$group": bson.M{"_id": "$accountid",
			// "matchtype": bson.M{"$first": "$matchtype"}, "matchid": bson.M{"$first": "$matchid"},
			"matchtotal": bson.M{"$sum": "$matchtotal"}, "matchwins": bson.M{"$sum": "$matchwins"}, "matchfails": bson.M{"$sum": "$matchfails"},
			"coupon": bson.M{"$sum": "$coupon"}, "awardmoney": bson.M{"$sum": "$awardmoney"}, "personalprofit": bson.M{"$sum": "$personalprofit"},
			// "winrate": bson.M{"$avg": []interface{}{bson.M{"$sum": "$matchwins"}, bson.M{"$sum": "$matchtotal"}}},
		}},
	}).One(&all); err != nil && err != mgo.ErrNotFound {
		log.Error("err:%v", err)
	}
	log.Debug("all:%+v", all)

	return matchs, data, all
}

// GetMatchReviewByName 根据赛事名称获取赛事总览
func GetMatchReviewByName(uid int, matchType string) (map[string]interface{}, []map[string]interface{}) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	all := map[string]interface{}{}
	if err := s.DB(GDB).C("matchreview").Pipe([]bson.M{
		{"$match": bson.M{"accountid": uid, "matchtype": matchType}},
		{"$group": bson.M{"_id": "$accountid",
			// "matchtype": bson.M{"$first": "$matchtype"}, "matchid": bson.M{"$first": "$matchid"},
			"matchtotal": bson.M{"$sum": "$matchtotal"}, "matchwins": bson.M{"$sum": "$matchwins"}, "matchfails": bson.M{"$sum": "$matchfails"},
			"coupon": bson.M{"$sum": "$coupon"}, "awardmoney": bson.M{"$sum": "$awardmoney"}, "personalprofit": bson.M{"$sum": "$personalprofit"},
			// "winrate": bson.M{"$divide": []interface{}{bson.M{"$sum": "$matchwins"}, bson.M{"$sum": "$matchtotal"}}},
		}},
	}).One(&all); err != nil && err != mgo.ErrNotFound {
		log.Error("err:%v", err)
	}
	log.Debug("all:%+v", all)

	matchs := []map[string]interface{}{}
	s.DB(GDB).C("matchreview").Pipe([]bson.M{
		{"$match": bson.M{"accountid": uid, "matchtype": matchType}},
		{"$group": bson.M{"_id": "$matchname"}},
	}).All(&matchs)
	log.Debug("matchs:%+v", matchs)

	ids := []map[string]interface{}{}
	s.DB(GDB).C("matchreview").Pipe([]bson.M{
		{"$match": bson.M{"accountid": uid}},
		{"$group": bson.M{"_id": "$matchid"}},
	}).All(&ids)
	log.Debug("ids:%+v", ids)

	list := []map[string]interface{}{}

	for _, matchName := range matchs {
		for _, id := range ids {
			one := map[string]interface{}{}
			if err := s.DB(GDB).C("matchreview").Pipe([]bson.M{
				{"$match": bson.M{"accountid": uid, "matchid": id["_id"], "matchtype": matchType, "matchname": matchName["_id"]}},
				{"$group": bson.M{"_id": "$accountid",
					"matchname":  bson.M{"$first": "$matchname"},
					"matchtotal": bson.M{"$sum": "$matchtotal"}, "matchwins": bson.M{"$sum": "$matchwins"}, "matchfails": bson.M{"$sum": "$matchfails"},
					"coupon": bson.M{"$sum": "$coupon"}, "awardmoney": bson.M{"$sum": "$awardmoney"}, "personalprofit": bson.M{"$sum": "$personalprofit"},
					// "winrate": bson.M{"$divide": []interface{}{bson.M{"$sum": "$matchwins"}, bson.M{"$sum": "$matchtotal"}}},
				}},
			}).One(&one); err != nil && err != mgo.ErrNotFound {
				log.Error("err:%v", err)
			}
			if len(one) > 0 {
				list = append(list, one)
			}
		}
	}
	log.Debug("list:%+v", list)
	return all, list
}

// GetUserOptLog 获取玩家操作日志
func GetUserOptLog(accountID, page, count, optType int, start, end int64) ([]util.ItemLog, int) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	ret := []util.ItemLog{}
	total := 0
	var err error
	var oneDay int64 = 24 * 60 * 60
	end += oneDay
	if optType > 0 {
		total, _ = s.DB(GDB).C("itemlog").Find(bson.M{"uid": accountID, "opttype": optType, "createtime": bson.M{"$gt": start, "$lt": end}}).Count()
		err = s.DB(GDB).C("itemlog").Find(bson.M{"uid": accountID, "opttype": optType, "createtime": bson.M{"$gt": start, "$lt": end}}).
			Sort("-createtime").Skip((page - 1) * count).All(&ret)
	} else {
		total, _ = s.DB(GDB).C("itemlog").Find(bson.M{"uid": accountID, "createtime": bson.M{"$gt": start, "$lt": end}}).Count()
		err = s.DB(GDB).C("itemlog").Find(bson.M{"uid": accountID, "createtime": bson.M{"$gt": start, "$lt": end}}).
			Sort("-createtime").Skip((page - 1) * count).All(&ret)
	}
	if err != nil && err != mgo.ErrNotFound {
		log.Error("err:%v", err)
	}
	return ret, total
}

func ReadOfflinePaymentList(req *param.OfflinePaymentListReq) *[]util.OfflinePaymentCol {
	op := new([]util.OfflinePaymentCol)
	readByPipeline(GDB, "offlinepayment", req.GetDataPipeline(), op, readTypeAll)
	return op
}

func ReadOfflinePaymentCount(req *param.OfflinePaymentListReq) int {
	cnt := new(util.DataCount)
	readByPipeline(GDB, "offlinepayment", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count
}

func SaveOfflinePayment(data *util.OfflinePaymentCol) {
	data.ID, _ = MongoDBNextSeq("offlinepayment")
	data.Createdat = time.Now().Unix()
	save(GDB, data, "offlinepayment", data.ID)
}

func ReadOrderHistoryList(req *param.OrderHistoryListReq) *[]util.EdyOrder {
	eo := new([]util.EdyOrder)
	readByPipeline(GDB, "edyorder", req.GetDataPipeline(), eo, readTypeAll)
	return eo
}

func ReadOrderHistoryCount(req *param.OrderHistoryListReq) int {
	cnt := new(util.DataCount)
	readByPipeline(GDB, "edyorder", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count
}

func saveBC(data interface{}, coll string, id int) {
	se := mongoDB.Ref()
	defer mongoDB.UnRef(se)
	if _, err := se.DB(DB).C(coll).Upsert(bson.M{"_id": id}, data); err != nil {
		log.Error(err.Error())
	}
}

func ReadRobotMatchNumList(req *param.RobotMatchNumReq) *[]util.RobotMatchNum {
	rt := new([]util.RobotMatchNum)
	readByPipeline(DB, "robotmatchnum", req.GetDataPipeline(), rt, readTypeAll)
	log.Debug("pipeline：%+v   Data：%v", req.GetDataPipeline(), *rt)
	return rt
}

func ReadRobotMatchNumCount(req *param.RobotMatchNumReq) int {
	cnt := new(util.DataCount)
	readByPipeline(DB, "robotmatchnum", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count
}

func ReadRobotMatchNum(condition base.Condition) *util.RobotMatchNum {
	rt := new(util.RobotMatchNum)
	readByPipeline(DB, "robotmatchnum", base.GetPipeline(condition), rt, readTypeOne)
	return rt
}

func SaveRobotMatchNum(data *util.RobotMatchNum) {
	save(DB, data, "robotmatchnum", data.ID)
}

func ReadMatchConfig(condition base.Condition) *util.MatchManager {
	rt := new(util.MatchManager)
	log.Debug("%v", base.GetPipeline(condition))
	readGameByPipeline(GDB, "matchmanager", base.GetPipeline(condition), rt, readTypeOne)
	return rt
}

func ReadAllMatchConfig(condition base.Condition) *[]util.MatchManager {
	rt := new([]util.MatchManager)
	log.Debug("%v", base.GetPipeline(condition))
	readGameByPipeline(GDB, "matchmanager", base.GetPipeline(condition), rt, readTypeAll)
	return rt
}
