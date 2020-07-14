package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type OfflinePaymentListReq struct {
	base.DivPage
	base.Condition
	base.TimeRange
}

func (ctx *OfflinePaymentListReq) GetPipeline() []bson.M {
	pipeline := ctx.TimeRange.GetPipeline()
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	return pipeline
}

func (ctx *OfflinePaymentListReq) GetDataPipeline() []bson.M {
	pipeline := ctx.GetPipeline()
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdat": -1}})
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type OfflinePaymentData struct {
	ID         int     `bson:"_id"`
	Nickname   string  `json:"nickname"`
	Accountid  int     `json:"accountid"`
	ActionType int     `json:"actiontype"` //0，点券 1，税后奖金
	BeforFee   float64 `json:"beforfee"`
	ChangeFee  float64 `json:"changefee"`
	AfterFee   float64 `json:"afterfee"`
	Createdat  int64   `json:"createdat"`
	Operator   string  `json:"operator"`
	Desc       string  `json:"desc"`
}

type OfflinePaymentListResp struct {
	Page                int                   `json:"page"`
	Per                 int                   `json:"per"`
	Total               int                   `json:"total"`
	OfflinePaymentDatas *[]OfflinePaymentData `json:"offline_payment_datas"`
}

type OfflinePaymentAddReq struct {
	Accountid  int     `json:"accountid"`
	ActionType int     `json:"action_type"`
	ChangeFee  float64 `json:"change_fee"`
	Desc       string  `json:"desc"`
}
