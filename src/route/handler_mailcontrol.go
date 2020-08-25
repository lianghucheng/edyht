package route

import (
	"bs/db"
	"bs/param"
	"bs/param/base"
	"bs/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/szxby/tools/log"
	"net/http"
	"time"
)

func mailcontrolInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.MailcontrolInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data := new(util.Mailcontrol)
	if err := transfer(req, data); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}
	if data.Title == "" {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error("The merchant type can not is nil")
		return
	}
	id, err := db.MongoDBNextSeq("mailcontrol")
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

	if err := db.SaveMailcontrol(data); err != nil {
		code = util.MailcontrolFail
		desc = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
func mailcontrolDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.MailcontrolDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadMailcontrol(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	now := int(time.Now().Unix())
	data.DeletedAt = now
	if err := db.SaveMailcontrol(data); err != nil {
		code = util.MailcontrolFail
		desc = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
func mailcontrolRead(c *gin.Context) {
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
	req := new(param.MailcontrolReadReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	data, err := db.ReadMailcontrol(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}

	rt := new(param.Mailcontrol)
	if err := transfer(data, rt); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	resp = param.MailcontrolReadResp{
		Mailcontrol: *rt,
	}

	return
}
func mailcontrolList(c *gin.Context) {
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
	req := new(param.MailcontrolListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	datas, err := db.ReadMailcontrolList(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	total, err := db.ReadMailcontrolCount(req)
	if err != nil {
		code = util.MailcontrolFail
		desc = util.ErrMsg[code]
		return
	}
	rt := new([]param.Mailcontrol)
	if err := transfer(datas, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	resp = &param.MailcontrolListResp{
		Page:         req.Page,
		Per:          req.Per,
		Total:        total,
		Mailcontrols: rt,
	}
}
func mailcontrolUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.MailcontrolUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadMailcontrol(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}

	data.TargetID = req.TargetID
	data.Title = req.Title
	data.Content = req.Content
	data.Annexes = req.Annexes
	data.Expire = req.Expire
	data.Status = req.Status
	data.MailServiceType = req.MailServiceType
	data.Operator = db.RedisGetTokenUsrn(c.GetHeader("token"))
	now := int(time.Now().Unix())
	if data.Status == 1 {
		data.UpdatedAt = now
	}
	if err := db.SaveMailcontrol(data); err != nil {
		code = util.MailcontrolFail
		desc = err.Error()
		return
	}
	return
}


func mailcontrolSendAll(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.MailcontrolSendAllReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	fails := []string{}
	all := []string{}
	for _, v := range req.Ids {
		oid := base.OID{ID: v}
		data, err := db.ReadMailcontrol(&oid)
		if err != nil {
			code = util.MongoReadFail
			desc = util.ErrMsg[code]
			return
		}
		all =append(all, data.Title)
		data.Status = 1
		data.Operator = db.RedisGetTokenUsrn(c.GetHeader("token"))
		now := int(time.Now().Unix())
		data.UpdatedAt = now
		if err := db.SaveMailcontrol(data); err != nil {
			code = util.MailcontrolFail
			desc = err.Error()
			fails = append(fails, data.Title)
			continue
		}
	}

	if len(fails) > 0 && len(fails) != len(req.Ids) {
		desc = fmt.Sprintf("总共需要发送标题为以下邮件：%v\n标题为以下邮件发送失败：%v", all, fails)
		return
	}

	if len(fails) == len(req.Ids) {
		code = util.SendAllMailFail
		desc = util.ErrMsg[code]
		return
	}

	return
}
