package route

import (
	"bs/db"
	"bs/param"
	"bs/util"
	"github.com/gin-gonic/gin"
	"github.com/name5566/leaf/log"
	"net/http"
	"time"
)

func shopMerchantInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopMerchantInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopMerchant := new(util.ShopMerchant)
	if err := transfer(req, shopMerchant); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	if shopMerchant.MerchantType <= 0 {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error("The merchant type can not is nil")
		return
	}

	id, err := db.MongoDBNextSeq("shopmerchant")
	if err != nil {
		code = util.MongoDBCreFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	shopMerchant.ID = id
	now := int(time.Now().Unix())
	shopMerchant.UpdatedAt = now
	shopMerchant.CreatedAt = now
	shopMerchant.DeletedAt = -1
	db.SaveShopMerchant(shopMerchant)
	return
}
func shopMerchantList(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	var resp interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"resp": resp,
		})
	}()
	req := new(param.ShopMerchantListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	ret := db.ReadShopMerchantList(req)
	total := db.ReadShopMerchantCount(req)

	rt := new([]param.ShopMerchant)
	if err := transfer(ret, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	resp = &param.ShopMerchantListResp{
		Page:          req.Page,
		Per:           req.Per,
		Total:         total,
		ShopMerchants: rt,
	}
}
func shopMerchantUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopMerchantUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopMerchant := db.ReadShopMerchant(req)

	shopMerchant.MerchantNo = req.MerchantNo
	shopMerchant.MerchantType = req.MerchantType
	shopMerchant.DownPayBranchs = req.DownPayBranchs
	shopMerchant.UpPayBranchs = req.UpPayBranchs
	shopMerchant.PayMin = req.PayMin
	shopMerchant.PayMax = req.PayMax
	shopMerchant.PublicKey = req.PublicKey
	shopMerchant.PrivateKey = req.PrivateKey
	shopMerchant.Order = req.Order

	now := int(time.Now().Unix())
	shopMerchant.UpdatedAt = now
	db.SaveShopMerchant(shopMerchant)
	return
}
func shopMerchantDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopMerchantDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopMerchant := db.ReadShopMerchant(req)

	now := int(time.Now().Unix())
	shopMerchant.DeletedAt = now
	db.SaveShopMerchant(shopMerchant)
	return
}
func shopPayAccountInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopPayAccountInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopPayAccount := new(util.PayAccount)
	if err := transfer(req, shopPayAccount); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	if shopPayAccount.MerchantID <= 0 {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error("The merchant type can not is nil")
		return
	}

	id, err := db.MongoDBNextSeq("shoppayaccount")
	if err != nil {
		code = util.MongoDBCreFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	shopPayAccount.ID = id
	now := int(time.Now().Unix())
	shopPayAccount.UpdatedAt = now
	shopPayAccount.CreatedAt = now
	shopPayAccount.DeletedAt = -1
	db.SavePayAccount(shopPayAccount)
	return
}
func shopPayAccountDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopPayAccountDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopPayAccount := db.ReadPayAccount(req)

	now := int(time.Now().Unix())
	shopPayAccount.DeletedAt = now
	db.SavePayAccount(shopPayAccount)
	return
}
func shopPayAccountList(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	var resp interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"resp": resp,
		})
	}()
	req := new(param.ShopPayAccountListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	ret := db.ReadPayAccountList(req)
	total := db.ReadPayAccountCount(req)

	rt := new([]param.PayAccount)
	if err := transfer(ret, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	resp = &param.ShopPayAccountListResp{
		Page:        req.Page,
		Per:         req.Per,
		Total:       total,
		PayAccounts: rt,
	}
}
func shopPayAccountUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopPayAccountUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	payAccount := db.ReadPayAccount(req)

	payAccount.Order = req.Order
	payAccount.Account = req.Account

	now := int(time.Now().Unix())
	payAccount.UpdatedAt = now
	db.SavePayAccount(payAccount)
	return
}
func shopGoodsTypeInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopGoodsTypeInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopGoodsType := new(util.GoodsType)
	if err := transfer(req, shopGoodsType); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}
	if shopGoodsType.MerchantID <= 0 {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error("The merchant type can not is nil")
		return
	}
	id, err := db.MongoDBNextSeq("shopgoodstype")
	if err != nil {
		code = util.MongoDBCreFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	shopGoodsType.ID = id
	now := int(time.Now().Unix())
	shopGoodsType.UpdatedAt = now
	shopGoodsType.CreatedAt = now
	shopGoodsType.DeletedAt = -1
	db.SaveGoodsType(shopGoodsType)
	return
}
func shopGoodsTypeUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopGoodsTypeUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	goodsType := db.ReadGoodsType(req)

	goodsType.Order = req.Order
	goodsType.TypeName = req.TypeName
	goodsType.ImgUrl = req.ImgUrl
	now := int(time.Now().Unix())
	goodsType.UpdatedAt = now
	db.SaveGoodsType(goodsType)
	return
}
func shopGoodsTypeDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopGoodsTypeDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopGoodsType := db.ReadGoodsType(req)

	now := int(time.Now().Unix())
	shopGoodsType.DeletedAt = now
	db.SaveGoodsType(shopGoodsType)
	return
}
func shopGoodsTypeList(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	var resp interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"resp": resp,
		})
	}()
	req := new(param.ShopGoodsTypeListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	ret := db.ReadGoodsTypeList(req)
	total := db.ReadGoodsTypeCount(req)

	rt := new([]param.GoodsType)
	if err := transfer(ret, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	goods := new(param.ShopGoodsListReq)
	m := make(map[string]interface{})
	for k, v := range *rt {
		m["goodstypeid"] = v.ID
		(*rt)[k].Num = db.ReadGoodsCount(goods)
	}

	resp = &param.ShopGoodsTypeListResp{
		Page:       req.Page,
		Per:        req.Per,
		Total:      total,
		GoodsTypes: rt,
	}
}
func shopGoodsInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopGoodsInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopGoods := new(util.Goods)
	if err := transfer(req, shopGoods); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}
	if shopGoods.GoodsTypeID <= 0 {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error("The merchant type can not is nil")
		return
	}
	id, err := db.MongoDBNextSeq("shopgoods")
	if err != nil {
		code = util.MongoDBCreFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	shopGoods.ID = id
	now := int(time.Now().Unix())
	shopGoods.UpdatedAt = now
	shopGoods.CreatedAt = now
	shopGoods.DeletedAt = -1
	db.SaveGoods(shopGoods)
	return
}
func shopGoodsDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopGoodsDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	shopGoods := db.ReadGoods(req)

	now := int(time.Now().Unix())
	shopGoods.DeletedAt = now
	db.SaveGoods(shopGoods)
	return
}
func shopGoodsList(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	var resp interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"resp": resp,
		})
	}()
	req := new(param.ShopGoodsListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	ret := db.ReadGoodsList(req)
	total := db.ReadGoodsCount(req)

	rt := new([]param.Goods)
	if err := transfer(ret, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	for k := range *rt {
		cfg, err := db.ReadPropBaseConfigByType((*rt)[k].PropType)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		(*rt)[k].ImgUrl = cfg.ImgUrl
	}
	resp = &param.ShopGoodsListResp{
		Page:    req.Page,
		Per:     req.Per,
		Total:   total,
		Goodses: rt,
	}
}
func shopGoodsUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ShopGoodsUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	goods := db.ReadGoods(req)

	goods.Order = req.Order
	goods.TakenType = req.TakenType
	goods.Price = req.Price
	goods.PropType = req.PropType
	goods.GetAmount = req.GetAmount
	goods.GiftAmount = req.GiftAmount
	goods.Expire = req.Expire
	goods.ImgUrl = req.ImgUrl
	now := int(time.Now().Unix())
	goods.UpdatedAt = now
	db.SaveGoods(goods)
	return
}
