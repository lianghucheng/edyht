package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type MatchAwardRecordReq struct {
	base.DivPage
	base.TimeRange
	base.Condition
}

func (ctx *MatchAwardRecordReq) GetPipeline() []bson.M {
	pipeline := ctx.TimeRange.GetPipeline()
	pipeline = append(pipeline, base.GetUnionPipeline(ctx.Condition)...)
	return pipeline
}

func (ctx *MatchAwardRecordReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdat": -1}})
	pipeline = append(pipeline, ctx.TimeRange.GetPipeline()...)
	pipeline = append(pipeline, base.GetUnionPipeline(ctx.Condition)...)
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type MatchAward struct {
	MatchName    string //赛事昵称
	AwardContent string //奖励类型和数量
	Accountid    int    //用户id
	MatchType    string //赛事类型
	CreatedAt    int64  //日期
	Realname     string //实名昵称
	Desc         string //备注说明
}

type MatchAwardRecordResp struct {
	Page              int           `json:"page"`
	Per               int           `json:"per"`
	Total             int           `json:"total"`
	MatchAwardRecords *[]MatchAward `json:"match_award_records"`
}
