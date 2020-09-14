package db

import (
	"bs/param"
	"bs/param/base"
	"bs/rpc"
	"bs/util"
	"errors"
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
)

func SaveHorseLamp(data *util.HorseRaceLampControl) error {
	log.Debug("存跑马灯控制台配置:%+v", *data)
	if data.Status == 0 {
		if err := rpc.RpcHorseStart(data); err != nil {
			log.Debug(err.Error())
			return err
		}
	} else if data.Status == 1 {
		if err := rpc.RpcHorseStop(data); err != nil {
			log.Debug(err.Error())
			return err
		}
	} else {
		log.Debug("暂未处理")
		return errors.New("暂未处理. ")
	}
	save(DB, data, "horselampcontrol", data.ID)
	return nil
}

func ReadHorseLampList(req *param.HorseLampListReq) (*[]util.HorseRaceLampControl, error) {
	datas := new([]util.HorseRaceLampControl)
	readByPipeline(DB, "horselampcontrol", req.GetDataPipeline(), datas, readTypeAll)
	return datas, nil
}

func ReadHorseLampCount(req *param.HorseLampListReq) (int, error) {
	cnt := new(util.DataCount)
	readByPipeline(DB, "horselampcontrol", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count, nil
}

func ReadHorseLamp(oid base.ObjectID) (*util.HorseRaceLampControl, error) {
	data := new(util.HorseRaceLampControl)
	readByPipeline(DB, "horselampcontrol", append(oid.GetOnePipeline(), bson.M{"$match": bson.M{"deletedat": 0}}), data, readTypeOne)
	return data, nil
}
