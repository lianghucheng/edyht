package base

import "gopkg.in/mgo.v2/bson"

type CountPipeline interface {
	GetPipeline() []bson.M
}

func GetCountPipeline(countPipeline CountPipeline) []bson.M {
	pipeline := countPipeline.GetPipeline()
	pipeline = append(pipeline, bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}})
	return pipeline
}

type DivPage struct {
	Page int `json:"page"`
	Per  int `json:"per"`
}

func (ctx *DivPage) GetPipeline() []bson.M {
	if ctx.Page <= 0 {
		ctx.Page = 1
	}
	if ctx.Per <= 0 {
		ctx.Per = 10
	}
	return []bson.M{
		{"$skip": (ctx.Page - 1) * ctx.Per},
		{"$limit": ctx.Per},
	}
}

type TimeRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (ctx *TimeRange) GetPipeline() []bson.M {
	return []bson.M{
		{"$match": bson.M{"createdat": bson.M{"$gte": ctx.Start, "$lt": ctx.End + 86400}}},
	}
}

type Condition interface {
}

func GetPipeline(cond Condition) []bson.M {
	if cond == nil {
		return nil
	}
	return []bson.M{
		{"$match": cond.(map[string]interface{})},
	}
}
