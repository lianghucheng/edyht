package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type HorseLampInsertReq struct {
	Name         string //通告名称
	Level        int    //等级排序，1：A，2：B，3：C，4：D
	ExpiredAt    int    //过期时间戳
	TakeEffectAt int    //发布时间戳
	Duration     int    //间隔时长，单位s
	LinkMatchID  string //关联赛事id
	Content      string //内容
}

type HorseLampDeleteReq struct {
	base.OID //记录唯一标识
}

type HorseLampReadReq struct {
	base.OID //记录唯一标识
}

type HorseLampReadResp struct {
	HorseLamp
}

type HorseLampListReq struct {
	base.DivPage
	base.TimeRange
	base.Condition
}

func (ctx *HorseLampListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": 0}}}
	pipeline = append(pipeline, ctx.TimeRange.GetPipeline()...)
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	return pipeline
}

func (ctx *HorseLampListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": 0}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdat": -1}})
	pipeline = append(pipeline, ctx.TimeRange.GetPipeline()...)
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type HorseLamp struct {
	ID           int    `bson:"_id"` //唯一标识
	Name         string //通告名称
	Level        int    //等级排序，1：A，2：B，3：C，4：D
	ExpiredAt    int    //过期时间戳
	TakeEffectAt int    //发布时间戳
	Duration     int    //间隔时长，单位s
	LinkMatchID  string //关联赛事id
	Content      string //内容
	Operator     string //操作人
	Status       int    //0表示发布，1表示暂停，2表示过期

	CreatedAt int //创建时间戳
	UpdatedAt int //更新时间戳，0表示未更新，对应着操作时间
}

type HorseLampListResp struct {
	Page      int
	Per       int
	Total     int
	HorseLamp *[]HorseLamp
}

type HorseLampUpdateReq struct {
	base.OID            //记录唯一标识
	Name         string //通告名称
	Level        int    //等级排序，1：A，2：B，3：C，4：D
	ExpiredAt    int    //过期时间戳
	TakeEffectAt int    //发布时间戳
	Duration     int    //间隔时长，单位s
	LinkMatchID  string //关联赛事id
	Content      string //内容
	Status       int    //0表示发布，1表示暂停，2表示过期
}
