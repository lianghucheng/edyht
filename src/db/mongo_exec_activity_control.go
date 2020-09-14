package db

import (
	"bs/param"
	"bs/param/base"
	"bs/util"
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
)

func SaveActivityControl(data *util.ActivityControl) error {
	log.Debug("存跑马灯控制台配置:%+v", *data)

	save(DB, data, "activitycontrol", data.ID)
	return nil
}

func ReadActivityControlList(req *param.ActivityControlListReq) (*[]util.ActivityControl, error) {
	datas := new([]util.ActivityControl)
	readByPipeline(DB, "activitycontrol", req.GetDataPipeline(), datas, readTypeAll)
	return datas, nil
}

func ReadActivityControlCount(req *param.ActivityControlListReq) (int, error) {
	cnt := new(util.DataCount)
	readByPipeline(DB, "activitycontrol", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count, nil
}

func ReadActivityControl(oid base.ObjectID) (*util.ActivityControl, error) {
	data := new(util.ActivityControl)
	readByPipeline(DB, "activitycontrol", append(oid.GetOnePipeline(), bson.M{"$match": bson.M{"deletedat": 0}}), data, readTypeOne)
	return data, nil
}
