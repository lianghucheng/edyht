package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type ShopMerchantInsertReq struct {
	MerchantNo string//商户编号
	MerchantType int//商户类型。1是体总
	DownPayBranchs []int//下架支付渠道类型
	PayMin int//支付最低值，百分制
	PayMax int//支付最高值，百分制
	PublicKey string//公钥
	PrivateKey string//私钥
	Order int//次序
}

type ShopMerchantListReq struct {
	base.DivPage
}

func (ctx *ShopMerchantListReq) GetPipeline() []bson.M {
	return []bson.M{{"$match": bson.M{"deletedat": -1}}}
}

func (ctx *ShopMerchantListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": -1}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"order": 1}})
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type ShopMerchant struct {
	ID int //唯一标识
	MerchantType int//商户类型。1是体总
	MerchantNo string//商户编号
	PayMin int//支付最低值，百分制
	PayMax int//支付最高值，百分制
	PublicKey string//公钥
	PrivateKey string//私钥
	Order int//次序
	UpdatedAt int//更新时间戳
	UpPayBranchs []int//上架支付渠道类型
	DownPayBranchs []int//下架支付渠道类型
}

type ShopMerchantListResp struct {
	Page              int
	Per               int
	Total             int
	ShopMerchants *[]ShopMerchant
}

type ShopMerchantUpdateReq struct {
	base.OID //记录唯一标识
	MerchantNo string//商户编号
	MerchantType int//商户类型。1是体总
	DownPayBranchs []int//下架支付类型
	UpPayBranchs []int//上架支付类型
	PayMin int//支付最低值，百分制
	PayMax int//支付最高值，百分制
	PublicKey string//公钥
	PrivateKey string//私钥
	Order int//次序
}

type ShopMerchantDeleteReq struct {
	base.OID//记录唯一标识
}

type ShopPayAccountInsertReq struct {
	MerchantID int//商户唯一标识
	PayBranch int//支付渠道标识
	Order int //次序
	Account string //账户
}

type ShopPayAccountDeleteReq struct {
	base.OID //记录唯一标识
}

type ShopPayAccountListReq struct {
	base.DivPage
	base.Condition
}

func (ctx *ShopPayAccountListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": -1}}}
	return append(pipeline,base.GetPipeline(ctx.Condition)...)
}

func (ctx *ShopPayAccountListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": -1}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"order": 1}})
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type PayAccount struct {
	ID int //唯一标识
	Order int //次序
	Account string //账户
	CreatedAt int//更新时间戳
}

type ShopPayAccountListResp struct {
	Page              int
	Per               int
	Total             int
	PayAccounts *[]PayAccount
}

type ShopPayAccountUpdateReq struct {
	base.OID //唯一标识
	Order int //次序
	Account string //账户
}

type ShopGoodsTypeInsertReq struct {
	MerchantID int //商户唯一标识
	TypeName string//商品名称
	ImgUrl string//商品图标
	Order int//次序
}

type ShopGoodsTypeUpdateReq struct {
	base.OID//记录唯一标识
	TypeName string//商品名称
	ImgUrl string//商品图标
	Order int//次序
}

type ShopGoodsTypeDeleteReq struct {
	base.OID//记录唯一标识
}

type ShopGoodsTypeListReq struct {
	base.DivPage
	base.Condition
}

func (ctx *ShopGoodsTypeListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": -1}}}
	return append(pipeline, base.GetPipeline(ctx.Condition)...)
}

func (ctx *ShopGoodsTypeListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": -1}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"order": 1}})
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type GoodsType struct {
	ID int `bson:"_id"`//唯一标识
	TypeName string//商品名称
	ImgUrl string//商品图标
	Order int//次序
	CreatedAt int//创建时间戳
}

type ShopGoodsTypeListResp struct {
	Page              int
	Per               int
	Total             int
	GoodsTypes *[]GoodsType
}

type ShopGoodsInsertReq struct {
	GoodsTypeID int //商品类型唯一标识
	TakenType int//花费类型。1是RMB
	Price int//花费数量（价格，百分制）
	PropType int//道具类型。1是点券
	GetAmount int//获得数量
	GiftAmount int//赠送数量
	Expire int//过期时间，单位秒，-1为永久
	ImgUrl string//商品图标
	Order int//次序
}

type ShopGoodsDeleteReq struct {
	base.OID //唯一标识
}

type ShopGoodsListReq struct {
	base.DivPage
	base.Condition
}

func (ctx *ShopGoodsListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": -1}}}
	return append(pipeline,base.GetPipeline(ctx.Condition)...)
}

func (ctx *ShopGoodsListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": -1}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"order": 1}})
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type Goods struct {
	ID int//唯一标识
	TakenType int//花费类型。1是RMB
	Price int//花费数量（价格，百分制）
	PropType int//道具类型。1是点券
	GetAmount int//获得数量
	GiftAmount int//赠送数量
	Expire int//过期时间，单位秒，-1为永久
	ImgUrl string//商品图标
	Order int//次序
	CreatedAt int//创建时间戳
}

type ShopGoodsListResp struct {
	Page              int
	Per               int
	Total             int
	Goodses *[]Goods
}

type ShopGoodsUpdateReq struct {
	base.OID//唯一标识
	TakenType int//花费类型。1是RMB
	Price int//花费数量（价格，百分制）
	PropType int//道具类型。1是点券
	GetAmount int//获得数量
	GiftAmount int//赠送数量
	Expire int//过期时间，单位秒，-1为永久
	ImgUrl string//商品图标
	Order int//次序
}
