package db

import (
	"bs/param"
	"bs/param/base"
	"bs/rpc"
	"bs/util"
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
)

func SavePropBaseConfig(data *util.PropBaseConfig) error {
	log.Debug("刷新物件配置缓存:%+v", *data)
	save(DB, data, "propbaseconfig", data.ID)
	if err := rpc.RpcSetPropBaseConfig(); err != nil {
		return err
	}
	return nil
}

func ReadPropBaseConfigList(req *param.PropBaseConfigListReq) (*[]util.PropBaseConfig, error) {
	datas := new([]util.PropBaseConfig)
	readByPipeline(DB, "propbaseconfig", req.GetDataPipeline(), datas, readTypeAll)
	return datas, nil
}

func ReadPropBaseConfigCount(req *param.PropBaseConfigListReq) (int, error) {
	cnt := new(util.DataCount)
	readByPipeline(DB, "propbaseconfig", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count, nil
}

func ReadPropBaseConfig(oid base.ObjectID) (*util.PropBaseConfig, error) {
	data := new(util.PropBaseConfig)
	readByPipeline(DB, "propbaseconfig", append(oid.GetOnePipeline(), bson.M{"$match": bson.M{"deletedat": 0}}), data, readTypeOne)
	return data, nil
}

func ReadPropBaseConfigByType(propType int) (*util.PropBaseConfig, error) {
	data := new(util.PropBaseConfig)
	readByPipeline(DB, "propbaseconfig", []bson.M{bson.M{"$match": bson.M{"proptype": propType, "deletedat": 0}}}, data, readTypeOne)
	return data, nil
}
