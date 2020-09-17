package db

import (
	"bs/param"
	"bs/param/base"
	"bs/util"
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
)

func SaveUser(data *util.User) error {
	log.Debug("SaveUser", *data)

	save(DB, data, "users", data.ID)
	return nil
}

func ReadUserList(req *param.UserListReq) (*[]util.User, error) {
	datas := new([]util.User)
	readByPipeline(DB, "users", req.GetDataPipeline(), datas, readTypeAll)
	return datas, nil
}

func ReadUserCount(req *param.UserListReq) (int, error) {
	cnt := new(util.DataCount)
	readByPipeline(DB, "users", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count, nil
}

func ReadUser(oid base.ObjectID) (*util.User, error) {
	data := new(util.User)
	readByPipeline(DB, "users", append(oid.GetOnePipeline(), bson.M{"$match": bson.M{"deletedat": 0}}), data, readTypeOne)
	return data, nil
}
