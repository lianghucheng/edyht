package route

import (
	"bs/db"
	"bs/param"
	"bs/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/name5566/leaf/log"
	"time"
)

func feedbackInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.FeedbackInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	feedback := new(util.FeedBack)
	if err := transfer(req, feedback); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	feedback.ID, _ = db.MongoDBNextSeq("feedback")
	if err := db.SaveFeedback(feedback); err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}

	return
}

func feedbackDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.FeedbackDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	data, err := db.ReadFeedback(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}
	data.DeletedAt = time.Now().Unix()
	if err := db.SaveFeedback(data); err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}
	return
}

func feedbackRead(c *gin.Context) {
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
	req := new(param.FeedbackReadReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	//resp = param.FeedbackReadResp{
	//	Feedback:param.Feedback{
	//		ID:          0,
	//		AccountID:   0,
	//		Title:       "",
	//		Content:     "",
	//		PhoneNum:    "",
	//		ReadStatus:  false,
	//		Nickname:    "",
	//		ReplyStatus: false,
	//		MailType:    0,
	//		ReplyTitle:  "",
	//		AwardType:   0,
	//		AwardNum:    0,
	//		MailContent: "",
	//		CreatedAt:   0,
	//		UpdatedAt:   0,
	//	},
	//}

	data, err := db.ReadFeedback(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}

	feedBack := new(param.Feedback)
	if err := transfer(data, feedBack); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	resp = param.FeedbackReadResp{
		Feedback:*feedBack,
	}

	return
}

func feedbackList(c *gin.Context) {
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
	req := new(param.FeedbackListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	rts := new([]param.Feedback)
	datas, err := db.ReadFeedbackList(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}
	total, err := db.ReadFeedbackCount(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}

	if err := transfer(datas, rts); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	resp = &param.FeedbackListResp{
		Page:          req.Page,
		Per:           req.Per,
		Total:         total,
		Feedbacks: rts,
	}
}

func feedbackUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.FeedbackUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	data, err := db.ReadFeedback(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}
	data.MailType = req.MailServiceType
	data.ReplyTitle = req.ReplyTitle
	data.AwardType  = req.AwardType
	data.AwardNum   = req.AwardNum
	data.MailContent= req.MailContent
	data.ReadStatus = req.ReadStatus
	data.ReplyStatus = req.ReplyStatus

	log.Debug("更新反馈的请求数据：%+v", *req)
	if err := db.SaveFeedback(data); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	return
}
