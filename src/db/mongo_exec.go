package db

import (
	"bs/param"
	"bs/param/base"
	"bs/util"
	"encoding/json"
	"fmt"
	"strings"
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
func GetMatchManagerList(page int, count int) ([]util.MatchManager, int) {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)
	// list := []map[string]interface{}{}
	list := []util.MatchManager{}
	total, _ := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"$lt": util.Delete}}).Count()
	// iter := s.DB(GDB).C("matchmanager").Find(bson.M{"state": bson.M{"gte": 0}}).Sort("-shelftime").Skip((page - 1) * count).Limit(count).Iter()
	err := s.DB(GDB).C("matchmanager").Pipe([]bson.M{
		{"$match": bson.M{"state": bson.M{"$lt": util.Delete}}},
		// {"$project": bson.M{
		// 	"MatchSource":   "$matchsource",
		// 	"MatchLevel":    "$matchlevel",
		// 	"MatchID":       "$matchid",
		// 	"MatchName":     "$matchname",
		// 	"MatchType":     "$matchtype",
		// 	"MatchIcon":     "$matchicon",
		// 	"RoundNum":      "$roundnum",
		// 	"StartTime":     "$starttime",
		// 	"StartType":     "$starttype",
		// 	"LimitPlayer":   "$limitplayer",
		// 	"Recommend":     "$recommend",
		// 	"Eliminate":     "$eliminate",
		// 	"EnterFee":      "$enterfee",
		// 	"UseCount":      "$usematch",
		// 	"LastMatch":     bson.M{"$subtract": []interface{}{"$totalmatch", "$usematch"}},
		// 	"ShelfTime":     "$shelftime",
		// 	"DownShelfTime": "$downshelftime",
		// 	"ShowHall":      "$showhall",
		// 	"Sort":          "$sort",
		// 	"State":         "$state",
		// 	"AwardList":     "$awardlist",
		// 	"TotalMatch":    "$totalmatch",
		// 	"_id":           0,
		// }},
		{"$sort": bson.M{"sort": 1}},
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

	for i := range list {
		list[i].LastMatch = list[i].TotalMatch - list[i].UseMatch
	}
	return list, total
}

// GetMatchReport 获取比赛报表
func GetMatchReport(matchID string, start, end int64) []map[string]interface{} {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	var oneDay int64 = 24 * 60 * 60
	length := (end - start) / oneDay
	if length <= 0 {
		log.Error("invalid time")
		return nil
	}

	AllSignPlayer := 0
	var AllSignFee, AllAward, AllLast float64

	result := make([]map[string]interface{}, 0)

	for i := end; i >= start; i = i - oneDay {
		all := []map[string]interface{}{}
		err := s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"matchid": matchID}},
			{"$match": bson.M{"createtime": bson.M{"$gt": i, "$lte": i + oneDay}}},
			{"$project": bson.M{
				"SignInCount": bson.M{"$size": "$signinplayers"}, "_id": 0, "matchid": fmt.Sprintf("$%v", matchID),
				"SignFee":   bson.M{"$multiply": []interface{}{bson.M{"$size": "$signinplayers"}, bson.M{"$divide": []interface{}{"$enterfee", util.CouponRate}}}},
				"AwardNum":  bson.M{"$size": "$award"},
				"Money":     "$moneyaward",
				"AwardList": "$awardlist",
				"LastMoney": bson.M{"$subtract": []interface{}{bson.M{
					"$multiply": []interface{}{
						bson.M{"$size": "$signinplayers"}, bson.M{"$divide": []interface{}{"$enterfee", util.CouponRate}}}},
					"$moneyaward"}}}},
		}).All(&all)
		if err == mgo.ErrNotFound {
			continue
		}
		if err != nil {
			log.Error("get report fail:%v", err)
			return nil
		}
		var oneSignPlayer, oneAwardNum int
		var oneSignFee, oneAward, oneLast float64
		awardCount := map[string]float64{}
		for _, v := range all {
			oneSignPlayer += v["SignInCount"].(int)
			oneSignFee += util.GetFloat(v["SignFee"])
			oneAward += util.GetFloat(v["Money"])
			oneLast += util.GetFloat(v["LastMoney"])
			oneAwardNum += v["AwardNum"].(int)
			util.ParseAwards(util.ParseAwardItem(v["AwardList"].(string)), awardCount)
		}
		var awardStr string
		for i, v := range awardCount {
			tmp := util.FormatFloat(v, 2) + i
			if len(awardStr) == 0 {
				awardStr += tmp
			} else {
				awardStr += "," + tmp
			}
		}
		one := map[string]interface{}{}
		one["RecordTime"] = time.Unix(i, 0).Format("2006-01-02")
		one["allMoney"] = oneAward
		one["allSign"] = oneSignPlayer
		one["allSignFee"] = oneSignFee
		one["awardNum"] = oneAwardNum
		one["lastMoney"] = oneLast
		one["awardList"] = awardStr
		result = append(result, one)
	}

	allReport := map[string]interface{}{}
	allReport["AllSignPlayer"] = AllSignPlayer
	allReport["AllSignFee"] = util.FormatFloat(AllSignFee, 2)
	allReport["AllAward"] = util.FormatFloat(AllAward, 2)
	allReport["AllLast"] = util.FormatFloat(AllLast, 2)
	result = append(result, allReport)
	return result
}

// GetMatch 按条件搜索赛事
func GetMatch(selector interface{}) []util.MatchManager {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	// one := map[string]interface{}{}
	// one := []util.MatchManager{}
	ret := []util.MatchManager{}
	err := s.DB(GDB).C("match").Pipe([]bson.M{
		{"$match": selector},
		// {"$project": bson.M{
		// 	"MatchType":   "$matchtype",
		// 	"MatchName":   "$matchname",
		// 	"MatchID":     "$sonmatchid",
		// 	"CreateTime":  "$createtime",
		// 	"RoundNum":    "$roundnum",
		// 	"LimitPlayer": "$limitplayer",
		// 	"Recommend":   "$recommend",
		// 	"StartType":   "$starttype",
		// 	"Eliminate":   "$eliminate",
		// 	"EnterFee":    "$enterfee",
		// }},
		// {"$group": bson.M{
		// 	"_id": "$matchid", "RecordTime": bson.M{"$first": "$RecordTime"}, "allMoney": bson.M{"$sum": "$Money"},
		// 	"allCoupon": bson.M{"$sum": "$Coupon"}, "allSign": bson.M{"$sum": "$SignInCount"},
		// 	"allSignFee": bson.M{"$sum": "$SignFee"}, "awardNum": bson.M{"$sum": "$AwardNum"},
		// "lastMoney": bson.M{"$sum": "$LastMoney"}}},
		// {"$sort": bson.M{"CreateTime": -1}},
	}).All(&ret)
	if err != nil {
		log.Error("get match fail %v", err)
		return nil
	}
	// ret = append(ret, one)
	return ret
}

// GetMatchByAccountID 按玩家id搜索赛事
func GetMatchByAccountID(accountID int, page, count int) ([]util.MatchManager, int) {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	if page == 0 {
		page = 1
	}
	if count == 0 {
		count = 10
	}
	// one := map[string]interface{}{}
	// one := []util.MatchManager{}
	user := ReadUserDataByAID(accountID)
	log.Debug("user:%v", user)
	ret := []util.MatchManager{}
	matchRecord := []util.DDZGameRecord{}
	total, _ := s.DB(GDB).C("gamerecord").Find(bson.M{"userid": user.UserID}).Count()
	err := s.DB(GDB).C("gamerecord").Find(bson.M{"userid": user.UserID}).Sort("-createdat").Skip((page - 1) * count).Limit(count).All(&matchRecord)
	log.Debug("matchrecord:%v", matchRecord)
	if err != nil {
		log.Error("get match fail %v", err)
		return nil, total
	}
	for _, v := range matchRecord {
		ones := GetMatch(bson.M{"sonmatchid": v.MatchId})
		ret = append(ret, ones...)
	}
	// ret = append(ret, one)
	return ret, total
}

// GetMatchList 获取某个时间段的赛事
func GetMatchList(matchType string, start, end int64) []util.MatchManager {
	s := gameDB.Ref()
	defer gameDB.UnRef(s)

	var oneDay int64 = 24 * 60 * 60
	length := (end - start) / oneDay
	if length <= 0 {
		log.Error("invalid time:%v,%v", start, end)
		return nil
	}

	var result []util.MatchManager
	// one := map[string]interface{}{}
	one := util.MatchManager{}
	var iter *mgo.Iter
	if len(matchType) == 0 {
		iter = s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"createtime": bson.M{"$gt": start, "$lte": end + oneDay}}},
			// {"$project": bson.M{
			// 	"MatchType":   "$matchtype",
			// 	"MatchSource": "$matchsource",
			// 	"MatchName":   "$matchname",
			// 	"MatchID":     "$sonmatchid",
			// 	"CreateTime":  "$createtime",
			// 	"RoundNum":    "$roundnum",
			// 	"LimitPlayer": "$limitplayer",
			// 	"Recommend":   "$recommend",
			// 	"StartType":   "$starttype",
			// 	"StartTime":   "$starttime",
			// 	"Eliminate":   "$eliminate",
			// 	"EnterFee":    "$enterfee",
			// 	"_id":         0,
			// }},
			{"$sort": bson.M{"createtime": -1}},
		}).Iter()
	} else {
		iter = s.DB(GDB).C("match").Pipe([]bson.M{
			{"$match": bson.M{"matchtype": matchType}},
			{"$match": bson.M{"createtime": bson.M{"$gt": start, "$lte": end + oneDay}}},
			// {"$project": bson.M{
			// 	"MatchType":   "$matchtype",
			// 	"MatchSource": "$matchsource",
			// 	"MatchName":   "$matchname",
			// 	"MatchID":     "$sonmatchid",
			// 	"CreateTime":  "$createtime",
			// 	"RoundNum":    "$roundnum",
			// 	"LimitPlayer": "$limitplayer",
			// 	"Recommend":   "$recommend",
			// 	"StartType":   "$starttype",
			// 	"StartTime":   "$starttime",
			// 	"Eliminate":   "$eliminate",
			// 	"EnterFee":    "$enterfee",
			// 	"_id":         0,
			// }},
			{"$sort": bson.M{"createtime": -1}},
		}).Iter()
	}
	for iter.Next(&one) {
		// data, err := json.Marshal(one)
		// if err != nil {
		// 	log.Error("get report fail:%v", err)
		// 	return nil
		// }
		one.MatchID = one.SonMatchID
		result = append(result, one)
		one = util.MatchManager{}
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
			"CreateTime":  "$createtime",
			"_id":         0,
		}},
	}).One(&one)
	if err != nil {
		log.Error("get detail fail %v", err)
		return nil
	}

	// 转化一下uid为accountid
	nowYear, ok := one["CreateTime"].(int64)
	if ok {
		thisYear := time.Unix(nowYear, 0)
		record := one["MatchRecord"]
		records := [][]util.MatchRecord{}
		if slice, ok := record.([]interface{}); ok {
			for _, v := range slice {
				tmp, err := json.Marshal(v)
				if err != nil {
					log.Error("err:%v", err)
				} else {
					oneRecords := []util.MatchRecord{}
					if err := json.Unmarshal(tmp, &oneRecords); err == nil {
						for i := range oneRecords {
							oneRecords[i].UID += thisYear.Year() * 100
						}
						records = append(records, oneRecords)
					}
				}
			}
		}
		// tmp, err := json.Marshal(record)
		// if err != nil {
		// 	log.Error("err:%v", err)
		// } else {
		// 	fmt.Println(string(tmp))
		// 	if err := json.Unmarshal(tmp, &records); err == nil {
		// 		for i, v := range records {
		// 			for j := range v {
		// 				records[i][j].UID += thisYear.Year() * 100
		// 			}
		// 		}
		// 	}
		// }
		one["MatchRecord"] = records
		delete(one, "CreateTime")
	}
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
	//if len(r.Condition) == 0 {
	//	log.Debug("@@@@@@@@@@@")
	//	query["flowtype"] = bson.M{
	//		"$ne": 1,
	//	}
	//}
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

		//if status != 0 {
		//	query["flowtype"] = 2
		//}
	}
	query["flowtype"] = 2
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
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	data := []util.UserData{}
	total, _ := s.DB(GDB).C("users").Find(bson.M{"accountid": bson.M{"$lt": 10000000000}}).Count()
	iter := s.DB(GDB).C("users").Find(bson.M{"accountid": bson.M{"$lt": 10000000000}}).Sort("-createdat").Skip((page - 1) * count).Limit(count).Iter()
	// if err != nil && err != mgo.ErrNotFound {
	// 	log.Error("err:%v", err)
	// 	return nil, total
	// }
	one := util.UserData{}
	for iter.Next(&one) {
		bank := ReadBankCardByID(one.UserID)
		one.BankCard = bank
		// 查询充值
		chargeAmount := GetPlayerCharge(one.AccountID)
		one.ChargeAmount = util.FormatFloat(float64(chargeAmount/100), 2)

		// 查询累计获得奖金
		one.AwardTotal = GetPlayerAwardTotal(one.UserID)

		// 可提现奖金
		one.AwardAvailable = GetPlayerAwardAvailable(one.UserID)

		// 参赛次数
		matchCount, _ := gs.DB(GDB).C("gamerecord").Find(bson.M{"userid": one.UserID}).Count()
		one.MatchCount = matchCount

		data = append(data, one)
		one = util.UserData{}
	}
	return data, total
}

// GetOneUser 获取单个用户列表
func GetOneUser(accountID int, nickname, phone string) (*util.UserData, error) {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	data := &util.UserData{}
	var err error
	selector := bson.M{}
	if accountID > 0 {
		// err = gs.DB(GDB).C("users").Find(bson.M{"accountid": accountID}).One(data)
		selector["accountid"] = accountID
	} else if len(nickname) > 0 {
		// err = gs.DB(GDB).C("users").Find(bson.M{"nickname": nickname}).One(data)
		selector["nickname"] = nickname
	} else if len(phone) > 0 {
		selector["username"] = phone
	}
	err = gs.DB(GDB).C("users").Find(selector).One(data)
	// err := s.DB(GDB).C("users").Find(bson.M{"$or": []interface{}{bson.M{"accountid": accountID}, bson.M{"nickname": nickname}}}).One(data)
	if err != nil {
		log.Error("err:%v", err)
		return nil, err
	}
	bank := ReadBankCardByID(data.UserID)
	data.BankCard = bank

	// 查询充值
	chargeAmount := GetPlayerCharge(data.AccountID)
	data.ChargeAmount = util.FormatFloat(float64(chargeAmount/100), 2)

	// 查询累计获得奖金
	data.AwardTotal = GetPlayerAwardTotal(data.UserID)
	data.AwardAvailable = GetPlayerAwardAvailable(data.UserID)

	data.Remark = GetRemark(data.AccountID)

	return data, nil
}

// GetMatchReview 获取赛事总览
func GetMatchReview(uid int) ([]map[string]interface{}, []map[string]interface{}, map[string]interface{}) {
	s := mongoDB.Ref()
	defer mongoDB.UnRef(s)
	data := []map[string]interface{}{}
	matchs := []map[string]interface{}{}
	// ids := []map[string]interface{}{}
	s.DB(GDB).C("matchreview").Pipe([]bson.M{
		{"$match": bson.M{"accountid": uid}},
		{"$group": bson.M{"_id": "$matchtype"}},
	}).All(&matchs)
	log.Debug("matchs:%+v", matchs)
	// s.DB(GDB).C("matchreview").Pipe([]bson.M{
	// 	{"$match": bson.M{"accountid": uid}},
	// 	{"$group": bson.M{"_id": "$matchid"}},
	// }).All(&ids)
	// log.Debug("ids:%+v", ids)
	// for _, matchType := range matchs {
	// 	for _, id := range ids {
	// 		// one := util.UserMatchReview{}
	// 		one := map[string]interface{}{}
	// 		// s.DB(GDB).C("matchreview").Find(bson.M{"accountid": uid, "matchtype": matchType["_id"], "matchid": id["_id"]}).One(&one)
	// 		err := s.DB(GDB).C("matchreview").Pipe([]bson.M{
	// 			{"$match": bson.M{"accountid": uid, "matchtype": matchType["_id"], "matchid": id["_id"]}},
	// 			// {"$group": bson.M{"_id": "$matchtype", "matchtype": "$matchtype", "matchid": "$matchid", "matchtotal": "$matchtotal",
	// 			// 	"matchwins": "$matchwins", "matchfails": "$matchfails", "coupon": "$coupon", "awardmoney": "$awardmoney", "personalprofit": "$personalprofit"}},
	// 			{"$group": bson.M{"_id": "$matchtype", "matchtype": bson.M{"$first": "$matchtype"}, "matchid": bson.M{"$first": "$matchid"},
	// 				"matchtotal": bson.M{"$sum": "$matchtotal"}, "matchwins": bson.M{"$sum": "$matchwins"}, "matchfails": bson.M{"$sum": "$matchfails"},
	// 				"coupon": bson.M{"$sum": "$coupon"}, "awardmoney": bson.M{"$sum": "$awardmoney"}, "personalprofit": bson.M{"$sum": "$personalprofit"}}},
	// 		}).One(&one)
	// 		if err != nil && err != mgo.ErrNotFound {
	// 			log.Error("err:%v", err)
	// 			continue
	// 		}
	// 		// 由于游戏服采集算法有误,修改读取方式
	// 		award, ok := one["awardmoney"].(int64)
	// 		if ok {
	// 			one["awardmoney"] = util.FormatFloat(float64(award)/100, 2)
	// 			if len(one) > 0 {
	// 				data = append(data, one)
	// 			}
	// 			if coupon, ok := one["coupon"].(int64); ok {
	// 				one["personalprofit"] = util.FormatFloat(float64(award-coupon*100)/100, 2)
	// 			}
	// 		}
	// 	}
	// }

	for _, matchType := range matchs {
		one := map[string]interface{}{}
		err := s.DB(GDB).C("matchreview").Pipe([]bson.M{
			{"$match": bson.M{"accountid": uid, "matchtype": matchType["_id"]}},
			{"$group": bson.M{"_id": "$matchtype", "matchtype": bson.M{"$first": "$matchtype"}, "matchid": bson.M{"$first": "$matchid"},
				"matchtotal": bson.M{"$sum": "$matchtotal"}, "matchwins": bson.M{"$sum": "$matchwins"}, "matchfails": bson.M{"$sum": "$matchfails"},
				"coupon": bson.M{"$sum": "$coupon"}, "awardmoney": bson.M{"$sum": "$awardmoney"}, "personalprofit": bson.M{"$sum": "$personalprofit"}}},
		}).One(&one)
		if err != nil && err != mgo.ErrNotFound {
			log.Error("err:%v", err)
			continue
		}
		// 由于游戏服采集算法有误,修改读取方式
		award, ok := one["awardmoney"].(int64)
		if ok {
			one["awardmoney"] = util.FormatFloat(float64(award)/100, 2)
			if len(one) > 0 {
				data = append(data, one)
			}
			if coupon, ok := one["coupon"].(int64); ok {
				one["personalprofit"] = util.FormatFloat(float64(award-coupon*100)/100, 2)
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

	// 由于游戏服采集算法有误,修改读取方式
	award, ok := all["awardmoney"].(int64)
	if ok {
		all["awardmoney"] = util.FormatFloat(float64(award)/100, 2)
		if coupon, ok := all["coupon"].(int64); ok {
			all["personalprofit"] = util.FormatFloat(float64(award-coupon*100)/100, 2)
		}
	}

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
	// 由于游戏服采集算法有误,修改读取方式
	award, ok := all["awardmoney"].(int64)
	if ok {
		all["awardmoney"] = util.FormatFloat(float64(award)/100, 2)
		if coupon, ok := all["coupon"].(int64); ok {
			all["personalprofit"] = util.FormatFloat(float64(award-coupon*100)/100, 2)
		}
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
			// 由于游戏服采集算法有误,修改读取方式
			award, ok := one["awardmoney"].(int64)
			if ok {
				one["awardmoney"] = util.FormatFloat(float64(award)/100, 2)
				if coupon, ok := one["coupon"].(int64); ok {
					one["personalprofit"] = util.FormatFloat(float64(award-coupon*100)/100, 2)
				}
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
			Sort("-createtime").Skip((page - 1) * count).Limit(count).All(&ret)
	} else {
		total, _ = s.DB(GDB).C("itemlog").Find(bson.M{"uid": accountID, "createtime": bson.M{"$gt": start, "$lt": end}}).Count()
		err = s.DB(GDB).C("itemlog").Find(bson.M{"uid": accountID, "createtime": bson.M{"$gt": start, "$lt": end}}).
			Sort("-createtime").Skip((page - 1) * count).Limit(count).All(&ret)
	}
	if err != nil && err != mgo.ErrNotFound {
		log.Error("err:%v", err)
	}

	for i := range ret {
		if ret[i].Item == "奖金" || strings.Index(ret[i].Item, "分") != -1 {
			ret[i].ShowAmount = util.FormatFloat(float64(ret[i].Amount)/100, 2)
			ret[i].ShowBefore = util.FormatFloat(float64(ret[i].Before)/100, 2)
			ret[i].ShowAfter = util.FormatFloat(float64(ret[i].After)/100, 2)
		} else {
			ret[i].ShowAmount = strconv.FormatInt(ret[i].Amount, 10)
			ret[i].ShowBefore = strconv.FormatInt(ret[i].Before, 10)
			ret[i].ShowAfter = strconv.FormatInt(ret[i].After, 10)
		}
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

func ReadKnapsackPropByAidPtype(aid, ptype int) *util.KnapsackProp {
	rt := new(util.KnapsackProp)
	readGameByPipeline(GDB, "knapsackprop", []bson.M{{"$match": bson.M{"accountid": aid, "proptype": ptype}}}, rt, readTypeOne)
	return rt
}

func ReadMatchAwardRecord(req *param.MatchAwardRecordReq) *[]util.MatchAwardRecord {
	mar := new([]util.MatchAwardRecord)
	readByPipeline(GDB, "matchawardrecord", req.GetDataPipeline(), mar, readTypeAll)
	return mar
}

func ReadMatchAwardRecordCount(req *param.MatchAwardRecordReq) int {
	cnt := new(util.DataCount)
	readByPipeline(GDB, "matchawardrecord", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count
}

// GetWhiteList 获取白名单
func GetWhiteList() (util.WhiteListConfig, error) {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	wConfig := util.WhiteListConfig{}
	if err := gs.DB(GDB).C("serverconfig").Find(bson.M{"config": "whitelist"}).One(&wConfig); err != nil {
		log.Error("err:%v", err)
		return wConfig, err
	}
	return wConfig, nil
}

// UpdateWhiteList 更新白名单
func UpdateWhiteList(selector interface{}, update interface{}) error {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	if err := gs.DB(GDB).C("serverconfig").Update(selector, update); err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// GetRestartList 获取重启信息表
func GetRestartList(page, count int, start, end int64) ([]util.RestartConfig, int, error) {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	rConfigs := []util.RestartConfig{}
	selector := bson.M{"config": "restart"}
	if start > 0 && end > 0 {
		selector["restarttime"] = bson.M{"$gte": start, "$lt": end}
	}
	total, _ := gs.DB(GDB).C("serverconfig").Find(selector).Count()
	if err := gs.DB(GDB).C("serverconfig").Find(selector).Sort("-createtime").Skip((page - 1) * count).Limit(count).All(&rConfigs); err != nil {
		log.Error("err:%v", err)
		return rConfigs, total, err
	}
	return rConfigs, total, nil
}

// InsertRestart 新建重启信息表
func InsertRestart(data interface{}) error {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	err := gs.DB(GDB).C("serverconfig").Insert(data)
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// UpdatetRestart 更新重启信息表
func UpdatetRestart(selector interface{}, update interface{}) error {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	err := gs.DB(GDB).C("serverconfig").Update(selector, update)
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// GetOneRestart 获取单条重启信息
func GetOneRestart(selector interface{}) (util.RestartConfig, error) {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	one := util.RestartConfig{}
	err := gs.DB(GDB).C("serverconfig").Find(selector).One(&one)
	if err != nil {
		log.Error("err:%v", err)
		return one, err
	}
	return one, nil
}

// GetLastestRestart 获取最新的重启信息
func GetLastestRestart() (util.RestartConfig, error) {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	one := util.RestartConfig{}
	err := gs.DB(GDB).C("serverconfig").Find(bson.M{"config": "restart"}).Sort("-createtime").Limit(1).One(&one)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("err:%v", err)
		return one, err
	}
	return one, nil
}

// GetFirstViewData 获取首页数据
func GetFirstViewData() map[string]interface{} {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)

	// start 为游戏上线时间,只查询上线后的数据
	tmp, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-07-17 19:00:00", time.Local)
	start := tmp.Unix()

	matchData := map[string]interface{}{}
	err := gs.DB(GDB).C("match").Pipe([]bson.M{
		{"$match": bson.M{"createtime": bson.M{"$gt": start}}},
		{"$project": bson.M{
			"RecordTime":  "$createtime",
			"SignInCount": bson.M{"$size": "$signinplayers"}, "_id": 0,
			"SignFee": bson.M{"$multiply": []interface{}{bson.M{"$size": "$signinplayers"},
				bson.M{"$divide": []interface{}{"$enterfee", util.CouponRate}}}},
			"AwardNum": bson.M{"$size": "$award"},
			"Money":    "$moneyaward",
			"Coupon":   "$couponaward",
			"LastMoney": bson.M{"$subtract": []interface{}{bson.M{
				"$multiply": []interface{}{
					bson.M{"$size": "$signinplayers"}, bson.M{"$divide": []interface{}{"$enterfee", util.CouponRate}}}},
				"$moneyaward"}}}},
		{"$group": bson.M{
			"_id": "AllData", "allMoney": bson.M{"$sum": "$Money"},
			"allCoupon": bson.M{"$sum": "$Coupon"}, "allSign": bson.M{"$sum": "$SignInCount"},
			"allSignFee": bson.M{"$sum": "$SignFee"}, "awardNum": bson.M{"$sum": "$AwardNum"},
			"lastMoney": bson.M{"$sum": "$LastMoney"}}},
	}).One(&matchData)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("get report fail:%v", err)
		return nil
	}

	// chargeData := map[string]interface{}{}
	// gs.DB(GDB).C("edyorder").Pipe([]bson.M{
	// 	{"$match": bson.M{"createdat": bson.M{"$gt": start}}},
	// 	{"$match": bson.M{"status": true}},
	// 	{"$project": bson.M{
	// 		"fee": "$fee",
	// 	}},
	// 	{"$group": bson.M{
	// 		"_id": "AllCharge", "charge": bson.M{"$sum": "$fee"},
	// 	},
	// 	}}).One(&chargeData)

	totalUser, _ := gs.DB(GDB).C("users").Find(bson.M{}).Count()

	// 返回数据汇总
	ret := map[string]interface{}{}
	ret["TotalUser"] = totalUser
	ret["TotalCharge"] = GetTotalCharge(start)
	ret["TotalSignFee"] = 0
	ret["TotalAward"] = 0
	ret["TotalLast"] = 0
	// if chargeData["charge"] != nil {
	// 	ret["TotalCharge"] = chargeData["charge"]
	// }
	if matchData["allSignFee"] != nil {
		ret["TotalSignFee"] = matchData["allSignFee"]
	}
	if matchData["allMoney"] != nil {
		ret["TotalAward"] = matchData["allMoney"]
	}
	if matchData["lastMoney"] != nil {
		ret["TotalLast"] = matchData["lastMoney"]
	}
	log.Debug("ret:%+v", ret)

	return ret
}

// GetPlayerAwardAvailable 获取玩家可提现奖金
func GetPlayerAwardAvailable(uid int) string {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	// 可提现奖金
	award := map[string]interface{}{}
	gs.DB(GDB).C("flowdata").Pipe([]bson.M{
		{"$match": bson.M{"flowtype": 1, "status": 0}},
		{"$match": bson.M{"userid": uid}},
		{"$project": bson.M{
			"Total": "$changeamount",
		}},
		{"$group": bson.M{
			"_id": "$accountid",
			"all": bson.M{"$sum": "$Total"},
		}},
	}).One(&award)
	var awardAvailable string
	if awardAdd, ok := award["all"].(float64); ok {
		awardAvailable = util.FormatFloat(float64(awardAdd), 2)
	}
	return awardAvailable
}

// GetPlayerAwardTotal 获取玩家累计获得奖金
func GetPlayerAwardTotal(uid int) string {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	// 查询累计获得奖金
	award := map[string]interface{}{}
	gs.DB(GDB).C("flowdata").Pipe([]bson.M{
		{"$match": bson.M{"flowtype": 1}},
		{"$match": bson.M{"userid": uid}},
		{"$project": bson.M{
			"Total": "$changeamount",
		}},
		{"$group": bson.M{
			"_id": "$accountid",
			"all": bson.M{"$sum": "$Total"},
		}},
	}).One(&award)
	var awardAmount string
	// log.Debug("fee:%v", reflect.TypeOf(fee["all"]))
	if awardAdd, ok := award["all"].(float64); ok {
		awardAmount = util.FormatFloat(float64(awardAdd), 2)
	}
	return awardAmount
}

// GetPlayerCharge 获取某一玩家的充值
func GetPlayerCharge(accountID int) int64 {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	// 查询充值
	fee := map[string]interface{}{}
	gs.DB(GDB).C("edyorder").Pipe([]bson.M{
		{"$match": bson.M{"status": true, "accountid": accountID}},
		{"$project": bson.M{
			"TotalFee": "$fee",
		}},
		{"$group": bson.M{
			"_id": "$accountid",
			"all": bson.M{"$sum": "$TotalFee"},
		}},
	}).One(&fee)
	var chargeAmount int64
	// log.Debug("fee:%v", reflect.TypeOf(fee["all"]))
	if feeAdd, ok := fee["all"].(int64); ok {
		chargeAmount += int64(feeAdd)
	}
	fee2 := map[string]interface{}{}
	gs.DB("czddz").C("alipayresult").Pipe([]bson.M{
		{"$match": bson.M{"success": true, "userid": accountID + 1e8}},
		{"$project": bson.M{
			"TotalFee": "$totalamount",
		}},
		{"$group": bson.M{
			"_id": "$userid",
			"all": bson.M{"$sum": "$TotalFee"},
		}},
	}).One(&fee2)
	if feeAdd, ok := fee2["all"].(float64); ok {
		chargeAmount += int64(feeAdd * 100)
	}
	return chargeAmount
}

// GetTotalCharge 获取总充值
func GetTotalCharge(start int64) int64 {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)
	// 查询充值
	fee := map[string]interface{}{}
	gs.DB(GDB).C("edyorder").Pipe([]bson.M{
		{"$match": bson.M{"status": true, "createdat": bson.M{"$gt": start}}},
		{"$project": bson.M{
			"TotalFee": "$fee",
		}},
		{"$group": bson.M{
			"_id": "$accountid",
			"all": bson.M{"$sum": "$TotalFee"},
		}},
	}).One(&fee)
	var chargeAmount int64
	// log.Debug("fee:%v", reflect.TypeOf(fee["all"]))
	if feeAdd, ok := fee["all"].(int64); ok {
		chargeAmount += int64(feeAdd)
	}
	fee2 := map[string]interface{}{}
	gs.DB("czddz").C("alipayresult").Pipe([]bson.M{
		{"$match": bson.M{"success": true, "createdat": bson.M{"$gt": start}, "userid": bson.M{"$gt": 1e8}}},
		{"$project": bson.M{
			"TotalFee": "$totalamount",
		}},
		{"$group": bson.M{
			"_id": "$userid",
			"all": bson.M{"$sum": "$TotalFee"},
		}},
	}).One(&fee2)
	log.Debug("czddz charge:%v", fee2)
	if feeAdd, ok := fee2["all"].(float64); ok {
		chargeAmount += int64(feeAdd * 100)
	}
	return chargeAmount
}

// UpdateRemark 更新玩家备注
func UpdateRemark(accountID int, remark string) error {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)

	_, err := gs.DB(DB).C("remark").Upsert(bson.M{"AccountID": accountID}, bson.M{"$set": bson.M{"Remark": remark}})
	if err != nil {
		log.Error("err:%v", err)
		return err
	}
	return nil
}

// GetRemark 获取玩家备注
func GetRemark(accountID int) string {
	gs := gameDB.Ref()
	defer gameDB.UnRef(gs)

	data := map[string]interface{}{}
	err := gs.DB(DB).C("remark").Find(bson.M{"AccountID": accountID}).One(&data)
	if err != nil {
		log.Error("err:%v", err)
		return ""
	}
	ret := data["Remark"]
	if ret == nil {
		log.Error("err:%v", data)
		return ""
	}
	s, ok := ret.(string)
	if !ok {
		log.Error("err:%v", data)
		return ""
	}
	return s
}
