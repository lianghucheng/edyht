package db

import (
	"bs/param"
	"bs/param/base"
	"bs/rpc"
	"bs/util"
	"github.com/name5566/leaf/log"
)

func SaveFeedback(data *util.FeedBack) error {
	log.Debug("更新反馈:%+v", *data)
	if data.ReplyStatus {
		ud := ReadUserDataByAID(data.AccountID)

		annexes := []rpc.Annex{}

		mailType := rpc.MailTypeText
		if data.AwardNum > 0 && data.AwardType != 0 {
			if len(data.Content) > 0 {
				mailType = rpc.MailTypeMix
			} else {
				mailType = rpc.MailTypeAward
			}
			propType := data.AwardType
			annexes = append(annexes, rpc.Annex{
				PropType: propType,
				Num:      float64(data.AwardNum),
				Desc:     "~",
			})
		}

		req := &rpc.MailBoxReq{
			TargetID:        int64(ud.UserID),
			MailType:        mailType,
			MailServiceType: data.MailServiceType,
			Title:           data.ReplyTitle,
			Content:         data.MailContent,
			Annexes:         annexes,
			ExpireValue:     30,
		}

		if err := rpc.RpcPushMail(req); err != nil {
			return err
		}
	}

	save(DB, data, "feedback", data.ID)
	return nil
}

func ReadFeedbackList(req *param.FeedbackListReq) (*[]util.FeedBack, error) {
	datas := new([]util.FeedBack)
	log.Debug("反馈列表,管道表达式：%v", req.GetDataPipeline())
	readByPipeline(DB, "feedback", req.GetDataPipeline(), datas, readTypeAll)
	return datas, nil
}

func ReadFeedbackCount(req *param.FeedbackListReq) (int, error) {
	cnt := new(util.DataCount)
	readByPipeline(DB, "feedback", base.GetCountPipeline(req), cnt, readTypeOne)
	return cnt.Count, nil
}

func ReadFeedback(oid base.ObjectID) (*util.FeedBack, error) {
	data := new(util.FeedBack)
	readByPipeline(DB, "feedback", oid.GetOnePipeline(), data, readTypeOne)
	return data, nil
}
