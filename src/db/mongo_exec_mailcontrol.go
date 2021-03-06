package db

import (
	"bs/param"
	"bs/param/base"
	"bs/rpc"
	"bs/util"
	"errors"
	"fmt"
	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2/bson"
)

func SaveMailcontrol(data *util.Mailcontrol) error {
	log.Debug("存邮件控制台配置邮件:%+v", *data)
	failTarget := []int{}
	if data.Status == util.MailcontrolStatusAlreadySend {
		for _, targetID := range data.TargetID {
			ud := ReadUserDataByAID(targetID)

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

			userID := ud.UserID
			if targetID == -1 {
				userID = -1
			}

			req := &rpc.MailBoxReq{
				TargetID:        int64(userID),
				MailType:        mailType,
				MailServiceType: data.MailServiceType,
				Title:           data.Title,
				Content:         data.Content,
				Annexes:         annexes,
				ExpireValue:     float64(data.Expire) / 1440,
			}

			if err := rpc.RpcPushMail(req); err != nil {
				log.Debug(err.Error())
				failTarget = append(failTarget, targetID)
				continue
			}
		}
	}
	var err error = nil
	if len(failTarget) > 0 && len(failTarget) != len(data.TargetID) {
		data.TargetID = util.Remove(data.TargetID, failTarget)
		err = errors.New(fmt.Sprintf("以下用户ID：%v 由于网络原因操作失败，请重试", failTarget))
	} else if len(failTarget) == len(data.TargetID) {
		err = errors.New("操作失败，请重试")
		return err
	}
	save(DB, data, "mailcontrol", data.ID)
	return err
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
