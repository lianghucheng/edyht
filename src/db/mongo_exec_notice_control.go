package db

import (
	"bs/param"
	"bs/param/base"
	"bs/util"
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
)

func SaveNoticeControl(data *util.NoticeControl) error {
	log.Debug("存跑马灯控制台配置:%+v", *data)

	save(DB, data, "noticecontrol", data.ID)
	return nil
}

func ReadNoticeControlList(req *param.NoticeControlListReq) (*[]util.NoticeControl, error) {
	datas := new([]util.NoticeControl)
	readByPipeline(DB, "noticecontrol", req.GetDataPipeline(), datas, readTypeAll)
	return datas, nil
}

func ReadNoticeControlCount(req *param.NoticeControlListReq) (int, error) {
	cnt := new(util.DataCount)
	readByPipeline(DB, "noticecontrol", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count, nil
}

func ReadNoticeControl(oid base.ObjectID) (*util.NoticeControl, error) {
	data := new(util.NoticeControl)
	readByPipeline(DB, "noticecontrol", append(oid.GetOnePipeline(), bson.M{"$match": bson.M{"deletedat": 0}}), data, readTypeOne)
	return data, nil
}
