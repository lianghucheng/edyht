package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type ActivityControlInsertReq struct {
	Order        int    //排序
	Title        string //活动标题
	Img          string //图片
	Matchid      string //关联赛事id
	Link         string //活动连接
	PrevUpedAt   int    //上架时间
	PrevDownedAt int    //下架时间
}

type ActivityControlDeleteReq struct {
	base.OID //记录唯一标识
}

type ActivityControlReadReq struct {
	base.OID //记录唯一标识
}

type ActivityControlReadResp struct {
	ActivityControl
}

type ActivityControlListReq struct {
	base.DivPage
	base.TimeRange
	base.Condition
}

func (ctx *ActivityControlListReq) GetPipeline() []bson.M {
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

func (ctx *ActivityControlListReq) GetDataPipeline() []bson.M {
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

type ActivityControl struct {
	ID           int    `bson:"_id"` //唯一标识
	Order        int    //排序
	Title        string //活动标题
	Img          string //图片
	Matchid      string //关联赛事id
	Link         string //活动连接
	Status       int    //状态
	PrevUpedAt   int    //上架时间
	PrevDownedAt int    //下架时间
	ClickCnt     int    //点击量
}

type ActivityControlListResp struct {
	Page             int
	Per              int
	Total            int
	ActivityControls *[]ActivityControl
}

type ActivityControlUpdateReq struct {
	base.OID            //记录唯一标识
	Order        int    //排序
	Title        string //活动标题
	Img          string //图片
	Matchid      string //关联赛事id
	Link         string //活动连接
	Status       int    //状态
	PrevUpedAt   int    //上架时间
	PrevDownedAt int    //下架时间
}
