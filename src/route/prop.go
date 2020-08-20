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

func propBaseConfigInsert(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.PropBaseConfigInsertReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data := new(util.PropBaseConfig)
	if err := transfer(req, data); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}
	if data.PropID <= 0 {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error("The merchant type can not is nil")
		return
	}
	id, err := db.MongoDBNextSeq("propbaseconfig")
	if err != nil {
		code = util.MongoDBCreFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	data.ID = id
	now := int(time.Now().Unix())
	data.UpdatedAt = now
	data.CreatedAt = now
	data.Operator = db.RedisGetTokenUsrn(c.GetHeader("token"))
	propType, ok := util.PropID2Type[req.PropID]
	if !ok {
		code = util.PropIDNotExist
		desc = util.ErrMsg[code]
		return
	}
	data.PropType = propType
	if err := db.SavePropBaseConfig(data); err != nil {
		code = util.PropBaseConfCacheFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}
	return
}
func propBaseConfigDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.PropBaseConfigDeleteReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadPropBaseConfig(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	now := int(time.Now().Unix())
	data.DeletedAt = now
	if err := db.SavePropBaseConfig(data); err != nil {
		code = util.PropBaseConfCacheFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}
	return
}
func propBaseConfigRead(c *gin.Context) {
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
	req := new(param.PropBaseConfigReadReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	data, err := db.ReadPropBaseConfig(req)
	if err != nil {
		code = util.Fail
		desc = err.Error()
		return
	}

	rt := new(param.PropBaseConfig)
	if err := transfer(data, rt); err != nil {
		code = util.ModelTransferFail
		desc = util.ErrMsg[code]
		log.Error(err.Error())
		return
	}

	resp = param.PropBaseConfigReadResp{
		PropBaseConfig: *rt,
	}

	return
}
func propBaseConfigList(c *gin.Context) {
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
	req := new(param.PropBaseConfigListReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	datas, err := db.ReadPropBaseConfigList(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	total, err := db.ReadPropBaseConfigCount(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}
	rt := new([]param.PropBaseConfig)
	if err := transfer(datas, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	resp = &param.PropBaseConfigListResp{
		Page:            req.Page,
		Per:             req.Per,
		Total:           total,
		PropBaseConfigs: rt,
	}
}
func propBaseConfigUpdate(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.PropBaseConfigUpdateReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	data, err := db.ReadPropBaseConfig(req)
	if err != nil {
		code = util.MongoReadFail
		desc = util.ErrMsg[code]
		return
	}

	data.Name = req.Name
	data.ImgUrl = req.ImgUrl
	data.Operator = db.RedisGetTokenUsrn(c.GetHeader("token"))
	now := int(time.Now().Unix())
	data.UpdatedAt = now
	if err := db.SavePropBaseConfig(data); err != nil {
		code = util.PropBaseConfCacheFail
		desc = util.ErrMsg[code]
		return
	}
	return
}
