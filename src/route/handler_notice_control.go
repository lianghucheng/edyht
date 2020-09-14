package route

import (
	"bs/db"
	"bs/param"
	"bs/util"
	"github.com/gin-gonic/gin"
	"github.com/szxby/tools/log"
	"net/http"
	"time"
)

func noticeControlInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.NoticeControlInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data := new(util.NoticeControl)
	if err := transfer(req, data); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}
	if data.ColTitle == "" {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error("The merchant type can not is nil")
		return
	}
	id, err := db.MongoDBNextSeq("noticecontrol")
	if err != nil {
		code = util.MongoDBCreFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	data.ID = id
	now := int(time.Now().Unix())
	data.CreatedAt = now
	data.Operator = db.RedisGetTokenUsrn(c.GetHeader("token"))
	data.Status = 1

	if err := db.SaveNoticeControl(data); err != nil {
		code = util.Fail
		desc = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
func noticeControlDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.NoticeControlDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadNoticeControl(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}

	now := int(time.Now().Unix())
	if data.Status < 2 && data.PrevUpedAt < now && data.PrevDownedAt > now {
		code = util.AlreadyUp
		desc = util.ErrMsg[code]
		return
	}

	data.DeletedAt = now
	if err := db.SaveNoticeControl(data); err != nil {
		code = util.Fail
		desc = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
func noticeControlRead(c *gin.Context) {
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
	req := new(param.NoticeControlReadReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	data, err := db.ReadNoticeControl(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}

	rt := new(param.NoticeControl)
	if err := transfer(data, rt); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	resp = param.NoticeControlReadResp{
		NoticeControl: *rt,
	}

	return
}
func noticeControlList(c *gin.Context) {
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
	req := new(param.NoticeControlListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	datas, err := db.ReadNoticeControlList(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	total, err := db.ReadNoticeControlCount(req)
	if err != nil {
		code = util.MailcontrolFail
		desc = util.ErrMsg[code]
		return
	}
	rt := new([]param.NoticeControl)
	if err := transfer(datas, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	for k := range *rt {
		now := int(time.Now().Unix())
		if (*rt)[k].PrevUpedAt > now || (*rt)[k].PrevDownedAt < now {
			(*rt)[k].Status = 2
		}
	}

	resp = &param.NoticeControlListResp{
		Page:           req.Page,
		Per:            req.Per,
		Total:          total,
		NoticeControls: rt,
	}
}
func noticeControlUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.NoticeControlUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadNoticeControl(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}

	now := int(time.Now().Unix())
	if data.PrevUpedAt > now || data.PrevDownedAt < now {
		code = util.NotInTimeRange
		desc = util.ErrMsg[code]
		return
	}

	data.Order = req.Order
	data.ColTitle = req.ColTitle
	data.NoticeTitle = req.NoticeTitle
	data.PrevUpedAt = req.PrevUpedAt
	data.PrevDownedAt = req.PrevDownedAt
	data.Content = req.Content
	data.Signature = req.Signature
	data.Status = req.Status
	data.Img = req.Img
	data.Operator = db.RedisGetTokenUsrn(c.GetHeader("token"))

	now = int(time.Now().Unix())
	data.UpdatedAt = now
	if err := db.SaveNoticeControl(data); err != nil {
		code = util.MailcontrolFail
		desc = err.Error()
		return
	}
	return
}
