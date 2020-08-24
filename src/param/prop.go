package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type PropBaseConfigInsertReq struct {
	PropID int    //道具id
	Name   string //名称
	ImgUrl string //图片url
}

type PropBaseConfigDeleteReq struct {
	base.OID //记录唯一标识
}

type PropBaseConfigReadReq struct {
	base.OID //记录唯一标识
}

type PropBaseConfigReadResp struct {
	PropBaseConfig
}

type PropBaseConfigListReq struct {
	base.DivPage
}

func (ctx *PropBaseConfigListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": 0}}}
	return pipeline
}

func (ctx *PropBaseConfigListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": 0}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"updatedat": -1}})
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type PropBaseConfig struct {
	ID       int    `bson:"_id"` //唯一标识
	PropID   int    //道具id
	PropType int    //道具类型, 1是点券，2是奖金，3点券碎片 todo:加进wiki文档
	Name     string //名称
	ImgUrl   string //图片url
	Operator string //操作人

	UpdatedAt int //更新时间戳
}

type PropBaseConfigListResp struct {
	Page            int
	Per             int
	Total           int
	PropBaseConfigs *[]PropBaseConfig
}

type PropBaseConfigUpdateReq struct {
	base.OID        //记录唯一标识
	Name     string //名称
	ImgUrl   string //图片url
}
