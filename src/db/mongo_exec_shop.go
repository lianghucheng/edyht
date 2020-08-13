package db

import (
	"bs/param"
	"bs/param/base"
	"bs/rpc"
	"bs/util"
	"github.com/name5566/leaf/log"
)

func SaveShopMerchant(data *util.ShopMerchant) {
	if len(data.UpPayBranchs) != 0 {
		data.UpDownStatus = 1
	} else if len(data.UpPayBranchs) == 0 {
		data.UpDownStatus = 0
	}
	save(DB, data, "shopmerchant", data.ID)
	if err := rpc.RpcNotifyPayAccount(); err != nil {
		log.Error(err.Error())
	}
	if err := rpc.RpcNotifyPriceMenu(); err != nil {
		log.Error(err.Error())
	}
}

func ReadShopMerchantList(req *param.ShopMerchantListReq) *[]util.ShopMerchant {
	datas := new([]util.ShopMerchant)
	readByPipeline(DB, "shopmerchant", req.GetDataPipeline(), datas, readTypeAll)
	return datas
}

func ReadShopMerchantCount(req *param.ShopMerchantListReq) int {
	cnt := new(util.DataCount)
	readByPipeline(DB, "shopmerchant", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count
}

func ReadShopMerchant(oid base.ObjectID) *util.ShopMerchant {
	data := new(util.ShopMerchant)
	readByPipeline(DB, "shopmerchant", oid.GetOnePipeline(), data, readTypeOne)
	return data
}

func SavePayAccount(data *util.PayAccount) {
	save(DB, data, "shoppayaccount", data.ID)
	if err := rpc.RpcNotifyPayAccount(); err != nil {
		log.Error(err.Error())
	}
}

func ReadPayAccountList(req *param.ShopPayAccountListReq) *[]util.PayAccount {
	datas := new([]util.PayAccount)
	readByPipeline(DB, "shoppayaccount", req.GetDataPipeline(), datas, readTypeAll)
	return datas
}

func ReadPayAccountCount(req *param.ShopPayAccountListReq) int {
	cnt := new(util.DataCount)
	readByPipeline(DB, "shoppayaccount", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count
}

func ReadPayAccount(oid base.ObjectID) *util.PayAccount {
	data := new(util.PayAccount)
	readByPipeline(DB, "shoppayaccount", oid.GetOnePipeline(), data, readTypeOne)
	return data
}

func SaveGoodsType(data *util.GoodsType) {
	save(DB, data, "shopgoodstype", data.ID)
	if err := rpc.RpcNotifyPriceMenu(); err != nil {
		log.Error(err.Error())
	}
	if err := rpc.RpcNotifyGoodsType(); err != nil {
		log.Error(err.Error())
	}
}

func ReadGoodsTypeList(req *param.ShopGoodsTypeListReq) *[]util.GoodsType {
	datas := new([]util.GoodsType)
	log.Debug("查看商品类型：%v", req.GetDataPipeline())
	readByPipeline(DB, "shopgoodstype", req.GetDataPipeline(), datas, readTypeAll)
	return datas
}

func ReadGoodsTypeCount(req *param.ShopGoodsTypeListReq) int {
	cnt := new(util.DataCount)
	readByPipeline(DB, "shopgoodstype", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count
}

func ReadGoodsType(oid base.ObjectID) *util.GoodsType {
	data := new(util.GoodsType)
	readByPipeline(DB, "shopgoodstype", oid.GetOnePipeline(), data, readTypeOne)
	return data
}

func SaveGoods(data *util.Goods) {
	save(DB, data, "shopgoods", data.ID)
	if err := rpc.RpcNotifyPriceMenu(); err != nil {
		log.Error(err.Error())
	}
}

func ReadGoodsList(req *param.ShopGoodsListReq) *[]util.Goods {
	datas := new([]util.Goods)
	readByPipeline(DB, "shopgoods", req.GetDataPipeline(), datas, readTypeAll)
	return datas
}

func ReadGoodsCount(req *param.ShopGoodsListReq) int {
	cnt := new(util.DataCount)
	readByPipeline(DB, "shopgoods", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count
}

func ReadGoods(oid base.ObjectID) *util.Goods {
	data := new(util.Goods)
	readByPipeline(DB, "shopgoods", oid.GetOnePipeline(), data, readTypeOne)
	return data
}
