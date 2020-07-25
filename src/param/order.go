package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type OrderHistoryListReq struct {
	base.DivPage
	base.Condition
	base.TimeRange
}

func (ctx *OrderHistoryListReq) GetPipeline() []bson.M {
	pipeline := ctx.TimeRange.GetPipeline()
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	return pipeline
}

func (ctx *OrderHistoryListReq) GetDataPipeline() []bson.M {
	pipeline := ctx.GetPipeline()
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdat": -1}})
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type OrderHistory struct {
	Accountid      int    //用户id
	TradeNo        string //订单号
	TradeNoReceive string //商户订单号
	GoodsType      string //商品类型。1表示点券，2表示碎片
	Amount         int    //商品数量
	Fee            int64  //支付金额,百分制
	Createdat      int64  //支付时间
	PayStatus      string //0表示支付中， 1表示支付成功， 2表示支付失败
	Merchant       string //商户
}

type OrderHistoryListResp struct {
	Page          int             `json:"page"`
	Per           int             `json:"per"`
	Total         int             `json:"total"`
	OrderHistorys *[]OrderHistory `json:"offline_payment_datas"`
}
