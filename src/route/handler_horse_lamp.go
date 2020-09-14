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

func horselampInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.HorseLampInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data := new(util.HorseRaceLampControl)
	if err := transfer(req, data); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}
	if data.Name == "" {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error("The merchant type can not is nil")
		return
	}
	id, err := db.MongoDBNextSeq("horselampcontrol")
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

	if err := db.SaveHorseLamp(data); err != nil {
		code = util.MailcontrolFail
		desc = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
func horselampDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.HorseLampDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadHorseLamp(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	now := int(time.Now().Unix())
	data.DeletedAt = now
	if err := db.SaveHorseLamp(data); err != nil {
		code = util.MailcontrolFail
		desc = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
func horselampRead(c *gin.Context) {
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
	req := new(param.HorseLampReadReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	data, err := db.ReadHorseLamp(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}

	rt := new(param.HorseLamp)
	if err := transfer(data, rt); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	resp = param.HorseLampReadResp{
		HorseLamp: *rt,
	}

	return
}
func horselampList(c *gin.Context) {
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
	req := new(param.HorseLampListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	datas, err := db.ReadHorseLampList(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	total, err := db.ReadHorseLampCount(req)
	if err != nil {
		code = util.MailcontrolFail
		desc = util.ErrMsg[code]
		return
	}
	rt := new([]param.HorseLamp)
	if err := transfer(datas, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	resp = &param.HorseLampListResp{
		Page:      req.Page,
		Per:       req.Per,
		Total:     total,
		HorseLamp: rt,
	}
}
func horselampUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.HorseLampUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadHorseLamp(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}

	data.Name = req.Name
	data.Level = req.Level
	data.ExpiredAt = req.ExpiredAt
	data.TakeEffectAt = req.TakeEffectAt
	data.Duration = req.Duration
	data.LinkMatchID = req.LinkMatchID
	data.Content = req.Content
	data.Operator = db.RedisGetTokenUsrn(c.GetHeader("token"))
	data.Status = req.Status

	now := int(time.Now().Unix())
	data.UpdatedAt = now
	if err := db.SaveHorseLamp(data); err != nil {
		code = util.MailcontrolFail
		desc = err.Error()
		return
	}
	return
}
