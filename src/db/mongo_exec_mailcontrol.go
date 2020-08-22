package db

import (
	"bs/param"
	"bs/param/base"
	"bs/rpc"
	"bs/util"
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
)

func SaveMailcontrol(data *util.Mailcontrol) error {
	log.Debug("存邮件控制台配置邮件:%+v", *data)
	if data.Status == util.MailcontrolStatusAlreadySend {
		ud := ReadUserDataByAID(data.TargetID)

		annexes := []rpc.Annex{}

		mailType := rpc.MailTypeText
		if len(data.Annexes) > 0 {
			if len(data.Content) > 0 {
				mailType = rpc.MailTypeMix
			} else {
				mailType = rpc.MailTypeAward
			}

			for _, v := range data.Annexes {
				annexes = append(annexes, rpc.Annex{
					PropType: v.PropType,
					Num:      v.Num,
					Desc:     "~",
				})
			}
		}

		req := &rpc.MailBoxReq{
			TargetID:        int64(ud.UserID),
			MailType:        mailType,
			MailServiceType: rpc.MailServiceTypeOfficial,
			Title:           data.Title,
			Content:         data.Content,
			Annexes:         annexes,
			ExpireValue:     float64(data.Expire) / 1440,
		}

		if err := rpc.RpcPushMail(req); err != nil {
			return err
		}
	}
	save(DB, data, "mailcontrol", data.ID)
	return nil
}

func ReadMailcontrolList(req *param.MailcontrolListReq) (*[]util.Mailcontrol, error) {
	datas := new([]util.Mailcontrol)
	readByPipeline(DB, "mailcontrol", req.GetDataPipeline(), datas, readTypeAll)
	return datas, nil
}

func ReadMailcontrolCount(req *param.MailcontrolListReq) (int, error) {
	cnt := new(util.DataCount)
	readByPipeline(DB, "mailcontrol", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count, nil
}

func ReadMailcontrol(oid base.ObjectID) (*util.Mailcontrol, error) {
	data := new(util.Mailcontrol)
	readByPipeline(DB, "mailcontrol", append(oid.GetOnePipeline(), bson.M{"$match": bson.M{"deletedat": 0}}), data, readTypeOne)
	return data, nil
}
