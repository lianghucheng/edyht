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

func activityControlInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ActivityControlInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data := new(util.ActivityControl)
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
	id, err := db.MongoDBNextSeq("activitycontrol")
	if err != nil {
		code = util.MongoDBCreFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	data.ID = id
	now := int(time.Now().Unix())
	data.CreatedAt = now
	data.Status = 1

	if err := db.SaveActivityControl(data); err != nil {
		code = util.Fail
		desc = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
func activityControlDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ActivityControlDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadActivityControl(req)
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
	if err := db.SaveActivityControl(data); err != nil {
		code = util.Fail
		desc = err.Error()
		log.Error(err.Error())
		return
	}
	return
}
func activityControlRead(c *gin.Context) {
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
	req := new(param.ActivityControlReadReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	data, err := db.ReadActivityControl(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}

	rt := new(param.ActivityControl)
	if err := transfer(data, rt); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	resp = param.ActivityControlReadResp{
		ActivityControl: *rt,
	}

	return
}
func activityControlList(c *gin.Context) {
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
	req := new(param.ActivityControlListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	datas, err := db.ReadActivityControlList(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	total, err := db.ReadActivityControlCount(req)
	if err != nil {
		code = util.MailcontrolFail
		desc = util.ErrMsg[code]
		return
	}
	rt := new([]param.ActivityControl)
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

	resp = &param.ActivityControlListResp{
		Page:             req.Page,
		Per:              req.Per,
		Total:            total,
		ActivityControls: rt,
	}
}
func activityControlUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.ActivityControlUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadActivityControl(req)
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
	data.Title = req.Title
	data.Img = req.Img
	data.Matchid = req.Matchid
	data.Link = req.Link
	data.Status = req.Status
	data.PrevUpedAt = req.PrevUpedAt
	data.PrevDownedAt = req.PrevDownedAt

	now = int(time.Now().Unix())
	data.UpdatedAt = now
	if err := db.SaveActivityControl(data); err != nil {
		code = util.MailcontrolFail
		desc = err.Error()
		return
	}
	return
}
