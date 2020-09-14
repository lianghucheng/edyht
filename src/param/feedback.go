package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type FeedbackInsertReq struct {
	AccountID       int    //用户id
	Title           string //反馈标题
	Content         string //反馈内容
	PhoneNum        string //联系方式
	ReadStatus      bool   //false是未查看，true是已查看
	Nickname        string //昵称
	ReplyStatus     bool   //false是未回复，true是已回复
	MailServiceType int    //0是系统邮件，1是赛事邮件，2是活动邮件
	ReplyTitle      string //回复标题
	AwardType       int    //0是未选择，10002是报名券，10003是报名券碎片
	AwardNum        int    //奖励数量
	MailContent     string //邮箱内容
}

type FeedbackDeleteReq struct {
	base.OID //记录唯一标识
}

type FeedbackReadReq struct {
	base.OID //记录唯一标识
}

type FeedbackReadResp struct {
	Feedback
}

type FeedbackListReq struct {
	base.DivPage
	base.TimeRange
	base.Condition
}

func (ctx *FeedbackListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": 0}}}
	pipeline = append(pipeline, ctx.TimeRange.GetPipeline()...)
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	return pipeline //[]bson.M{{"$match": bson.M{"deletedat": 0}}}
}

func (ctx *FeedbackListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": 0}})
	pipeline = append(pipeline, ctx.TimeRange.GetPipeline()...)
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdat": -1}})
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type Feedback struct {
	ID              int    `bson:"_id"` //唯一标识
	AccountID       int    //用户id
	Title           string //反馈标题
	Content         string //反馈内容
	PhoneNum        string //联系方式
	ReadStatus      bool   //false是未查看，true是已查看
	Nickname        string //昵称
	ReplyStatus     bool   //false是未回复，true是已回复
	MailServiceType int    //0是系统邮件，1是赛事邮件，2是活动邮件
	ReplyTitle      string //回复标题
	AwardType       int    //0是未选择，10002是报名券，10003是报名券碎片
	AwardNum        int    //奖励数量
	MailContent     string //邮箱内容
	Operator 		string //操作人

	CreatedAt int64 //创建时间戳，0表示未创建
	UpdatedAt int64 //更新时间戳，0表示未更新
}

type FeedbackListResp struct {
	Page      int
	Per       int
	Total     int
	Feedbacks *[]Feedback
}

type FeedbackUpdateReq struct {
	base.OID               //记录唯一标识
	MailServiceType int    //0是系统邮件，1是赛事邮件，2是活动邮件
	ReplyTitle      string //回复标题
	AwardType       int    //0是未选择，10002是报名券，10003是报名券碎片
	AwardNum        int    //奖励数量
	MailContent     string //邮箱内容
	ReadStatus      bool   //false是未查看，true是已查看
	ReplyStatus     bool   //false是未回复，true是已回复
}
