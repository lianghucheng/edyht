package param

import (
	"bs/param/base"
	"bs/util"
	"gopkg.in/mgo.v2/bson"
)

type MailcontrolInsertReq struct {
	TargetID []int        //目标用户
	Title    string       //标题
	Content  string       //内容
	Annexes  []util.Annex //附件
	Expire   int          //过期时间（单位：分钟）
}

type MailcontrolDeleteReq struct {
	base.OID //记录唯一标识
}

type MailcontrolReadReq struct {
	base.OID //记录唯一标识
}

type MailcontrolReadResp struct {
	Mailcontrol
}

type MailcontrolListReq struct {
	base.DivPage
}

func (ctx *MailcontrolListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": 0}}}
	return pipeline
}

func (ctx *MailcontrolListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": 0}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdat": -1}})
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type Mailcontrol struct {
	ID       int          `bson:"_id"` //唯一标识
	TargetID int          //目标用户
	Title    string       //标题
	Annexes  []util.Annex //附件
	Expire   int          //过期时间（单位：分钟）
	Status   int          //状态，0是未发送，1是已发送
	Operator string       //操作人

	CreatedAt int //创建时间戳，对应添加时间
	UpdatedAt int //更新时间戳，对应发送时间，0表示从未更新过
}

type MailcontrolListResp struct {
	Page         int
	Per          int
	Total        int
	Mailcontrols *[]Mailcontrol
}

type MailcontrolUpdateReq struct {
	base.OID              //记录唯一标识
	TargetID []int        //目标用户
	Title    string       //标题
	Content  string       //内容
	Annexes  []util.Annex //附件
	Expire   int          //过期时间（单位：分钟）
	Status   int          //状态，0是未发送，1是已发送
}
