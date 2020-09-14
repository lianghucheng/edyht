package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type NoticeControlInsertReq struct {
	Order        int    //排序
	ColTitle     string //栏目标题
	NoticeTitle  string //公告标题
	PrevUpedAt   int    //上架时间戳
	PrevDownedAt int    //下架时间戳
	Content      string //公告内容
	Signature    string //公告落款
	Img          string //公告内容
}

type NoticeControlDeleteReq struct {
	base.OID //记录唯一标识
}

type NoticeControlReadReq struct {
	base.OID //记录唯一标识
}

type NoticeControlReadResp struct {
	NoticeControl
}

type NoticeControlListReq struct {
	base.DivPage
	base.TimeRange
	base.Condition
}

func (ctx *NoticeControlListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": 0}}}
	pipeline = append(pipeline, ctx.TimeRange.GetPipeline()...)
	if ctx.Condition != nil {
		if s, ok := ctx.Condition.(map[string]interface{})["status"].(float64); ok && s == 1 {
			pipeline = append(pipeline, GetUpPipeline()...)
		}
		if s, ok := ctx.Condition.(map[string]interface{})["status"].(float64); ok && s == 2 {
			pipeline = append(pipeline, GetDownPipeline()...)
		}
	}
	return pipeline
}

func (ctx *NoticeControlListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": 0}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"order": -1}})
	pipeline = append(pipeline, ctx.TimeRange.GetPipeline()...)
	if ctx.Condition != nil {
		if s, ok := ctx.Condition.(map[string]interface{})["status"].(float64); ok && s == 1 {
			pipeline = append(pipeline, GetUpPipeline()...)
		}
		if s, ok := ctx.Condition.(map[string]interface{})["status"].(float64); ok && s == 2 {
			pipeline = append(pipeline, GetDownPipeline()...)
		}
	}
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type NoticeControl struct {
	ID           int    `bson:"_id"` //唯一标识
	Order        int    //排序
	ColTitle     string //栏目标题
	NoticeTitle  string //公告标题
	Status       int    //状态
	PrevUpedAt   int    //上架时间戳
	PrevDownedAt int    //下架时间戳
	Operator     string //操作人
	Content      string //公告内容
	Img          string //公告内容
	Signature    string //公告落款
}

type NoticeControlListResp struct {
	Page           int
	Per            int
	Total          int
	NoticeControls *[]NoticeControl
}

type NoticeControlUpdateReq struct {
	base.OID            //记录唯一标识
	Order        int    //排序
	ColTitle     string //栏目标题
	NoticeTitle  string //公告标题
	PrevUpedAt   int    //上架时间戳
	PrevDownedAt int    //下架时间戳
	Content      string //公告内容
	Signature    string //公告落款
	Status       int    //状态
	Img          string //公告内容
}

func GetUpPipeline() []bson.M {
	now := int(time.Now().Unix())
	return []bson.M{
		{
			"$match": bson.M{"status": 1, "prevdownedat": bson.M{"$gt": now}, "prevupedat": bson.M{"$lt": now}},
		},
	}
}

func GetDownPipeline() []bson.M {
	now := int(time.Now().Unix())
	return []bson.M{
		{
			"$match": bson.M{"$or": []bson.M{
				{"prevdownedat": bson.M{"$lt": now}},
				{"prevupedat": bson.M{"$gt": now}},
				{"status": 2, "prevdownedat": bson.M{"$gt": now}, "prevupedat": bson.M{"$lt": now}},
			}},
		},
	}
}
