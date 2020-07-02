package route

import (
	"bs/config"
	"bs/db"
	"bs/param"
	"bs/rpc"
	"bs/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/szxby/tools/log"

	"github.com/gin-gonic/gin"
	"github.com/szxby/tools/log"
)

func login(c *gin.Context) {
	code := util.OK
	desc := "OK"
	resp := ""
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"token": resp,
		})
	}()
	data := loginData{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	user := db.GetUser(data.Account)
	if user == nil {
		code = util.Retry
		desc = "用户不存在"
		return
	}
	if user.Password != util.CalculateHash(data.Password) {
		code = util.Retry
		desc = "密码错误"
		return
	}
	token := util.RandomString(10)
	db.RedisSetToken(token, user.Role)
	resp = token
}
func matchManagerList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	total := 0
	var resp interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"list":  resp,
			"total": total,
		})
	}()
	data := matchManagerReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	resp, total = db.GetMatchManagerList(data.Page, data.Count)
}
func addMatch(c *gin.Context) {
	code := util.OK
	desc := "添加赛事成功！"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := addManagerReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	err := util.PostToGame(config.GetConfig().GameServer+"/addMatch", JSON, data)
	if err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}
func showHall(c *gin.Context) {
	code := util.OK
	desc := "修改赛事成功！"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := showHallReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/showHall", JSON, data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}
func editMatch(c *gin.Context) {
	code := util.OK
	desc := "修改赛事成功！"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := editManagerReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/editMatch", JSON, data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}
func cancelMatch(c *gin.Context) {
	code := util.OK
	desc := "下架赛事成功！"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := optMatchReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/cancelMatch", JSON, data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}
func deleteMatch(c *gin.Context) {
	code := util.OK
	desc := "删除赛事成功！"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := optMatchReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/deleteMatch", JSON, data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}
func matchReport(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var list [][]byte
	var all []byte
	total := 0
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"all":   all,
			"list":  list,
			"total": total,
		})
	}()
	data := matchReportReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if data.Page <= 0 || data.Count <= 0 {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	begin, err := time.Parse("2006-01-02", data.Start)
	over, err := time.Parse("2006-01-02", data.End)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", data.Start, data.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) >= time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}

	last := data.Page * data.Count

	// 查看redis中是否有缓存
	if ret := db.RedisGetReport(data.MatchID, data.Start, data.End); ret != nil {
		total = len(ret) - 1
		if (data.Page-1)*data.Count >= total {
			log.Error("error page:%v,count:%v", data.Page, data.Count)
			code = util.Retry
			desc = "非法请求页码！"
			return
		}
		if last > total {
			last = total
		}
		list = ret[(data.Page-1)*data.Count : last]
		all = ret[len(ret)-1]
		return
	}

	result := db.GetMatchReport(data.MatchID, begin.Unix(), over.Unix())
	if result == nil {
		code = util.Retry
		desc = "查询出错请重试！"
		return
	}

	// 最后一位是总数据
	total = len(result) - 1
	if (data.Page-1)*data.Count >= total {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}

	// ok
	// 数据存入redis
	db.RedisSetReport(result, data.MatchID, data.Start, data.End)
	if last > total {
		last = total
	}
	list = result[(data.Page-1)*data.Count : last]
	all = result[len(result)-1]
}
func matchList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	total := 0
	var resp interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"list":  resp,
			"total": total,
		})
	}()
	data := matchListReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	// 按照赛事id精准查询
	if len(data.MatchID) > 0 {
		resp = db.GetMatch(data.MatchID)
		return
	}
	// 按照matchtype和时间查询

	if data.Page <= 0 || data.Count <= 0 {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	begin, err := time.Parse("2006-01-02", data.Start)
	over, err := time.Parse("2006-01-02", data.End)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", data.Start, data.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) >= time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}

	last := data.Page * data.Count
	// 查看redis中是否有缓存
	if ret := db.RedisGetMatchList(data.MatchType, data.Start, data.End); ret != nil {
		total = len(ret)
		if (data.Page-1)*data.Count >= total {
			log.Error("error page:%v,count:%v", data.Page, data.Count)
			code = util.Retry
			desc = "非法请求页码！"
			return
		}
		if last > total {
			last = total
		}
		resp = ret[(data.Page-1)*data.Count : last]
		return
	}

	result := db.GetMatchList(data.MatchType, begin.Unix(), over.Unix())
	if result == nil {
		code = util.Retry
		desc = "查询出错请重试！"
		return
	}

	total = len(result)
	if (data.Page-1)*data.Count >= total {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}

	// ok
	// 数据存入redis
	db.RedisSetMatchList(result, data.MatchType, data.Start, data.End)
	if last > total {
		last = total
	}
	resp = result[(data.Page-1)*data.Count : last]
}
func matchDetail(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var resp []byte
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"resp": resp,
		})
	}()
	data := matchDetailReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	resp = db.GetMatchDetail(data.MatchID)
}

func flowDataHistory(c *gin.Context) {
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
	data := c.Request.FormValue("data")
	flowDataReq := new(param.FlowDataHistoryReq)
	if err := json.Unmarshal([]byte(data), flowDataReq); err != nil {
		log.Error(err.Error())
		code = util.FormatFail
		desc = util.ErrMsg[util.FormatFail]
		resp = nil
		return
	}

	flowDatas, total := db.ReadFlowDatas(flowDataReq)

	resp = pkgFlowDataHistoryResp(uflow2Pflow(flowDatas), flowDataReq, total)
}

func uflow2Pflow(c *[]util.FlowData) *[]param.FlowData {
	rt := new([]param.FlowData)
	for _, v := range *c {
		temp := param.FlowData{
			ID:           v.ID,
			Accountid:    v.Accountid,
			ChangeAmount: v.ChangeAmount,
			FlowType:     v.FlowType,
			MatchID:      v.MatchID,
			Status:       v.Status,
			CreatedAt:    v.CreatedAt,
			Realname:     v.Realname,
			TakenFee:     v.TakenFee,
			AtferTaxFee:  v.AtferTaxFee,
			Desc:         v.Desc,
		}
		*rt = append(*rt, temp)
	}

	return rt
}

func pkgFlowDataHistoryResp(flowDatas *[]param.FlowData, flowDataReq *param.FlowDataHistoryReq, total int) *param.FlowDataHistoryResp {
	resp := new(param.FlowDataHistoryResp)
	resp.Total = total
	resp.Page = flowDataReq.Page
	resp.Per = flowDataReq.Per
	resp.FlowDatas = flowDatas

	return resp
}

func flowDataPayment(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := c.Request.FormValue("data")
	pm := new(param.FlowDataPaymentReq)
	if err := json.Unmarshal([]byte(data), pm); err != nil {
		log.Error(err.Error())
		code = util.FormatFail
		desc = util.ErrMsg[util.FormatFail]
		return
	}
	code, desc = thepayment(pm.ID, pm.Desc)
	return
}

func flowDataRefund(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := c.Request.FormValue("data")
	refund := new(param.FlowDataRefundReq)
	if err := json.Unmarshal([]byte(data), refund); err != nil {
		log.Error(err.Error())
		code = util.FormatFail
		desc = util.ErrMsg[util.FormatFail]
		return
	}
	code, desc = therefund(refund.ID, refund.Desc)
	return
}

func paymentByFlowIDs(flowIDs []int) {
	for _, v := range flowIDs {
		fd := db.ReadFlowDataByID(v)
		fd.Status = util.FlowDataStatusOver
		db.SaveFlowData(fd)
	}
}

func refundByFlowIDs(flowIDs []int) {
	for _, v := range flowIDs {
		fd := db.ReadFlowDataByID(v)
		fd.Status = util.FlowDataStatusNormal
		db.SaveFlowData(fd)
	}
}

func flowDataPayments(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := c.Request.FormValue("data")
	pms := new(param.FlowDataPaymentsReq)
	if err := json.Unmarshal([]byte(data), pms); err != nil {
		log.Error(err.Error())
		code = util.FormatFail
		desc = util.ErrMsg[util.FormatFail]
		return
	}
	for _, id := range pms.Ids {
		code, desc = thepayment(id, pms.Desc)
	}
	return
}

func flowDataRefunds(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := c.Request.FormValue("data")
	refunds := new(param.FlowDataRefundsReq)
	if err := json.Unmarshal([]byte(data), refunds); err != nil {
		log.Error(err.Error())
		code = util.FormatFail
		desc = util.ErrMsg[util.FormatFail]
		return
	}
	for _, id := range refunds.Ids {
		code, desc = therefund(id, refunds.Desc)
	}
	return
}

func therefund(id int, thedesc string) (code int, desc string) {
	code = util.Success
	desc = util.ErrMsg[util.Success]
	flowData := db.ReadFlowDataByID(id)
	if flowData.Status != 1 {
		code = util.Fail
		desc = util.ErrMsg[util.Fail]
		return
	}
	flowData.Status = util.FlowDataStatusBack
	flowData.Desc = thedesc
	//db.AddUserFee(flowData)
	rpc.AddFee(flowData.Userid, flowData.ChangeAmount, "fee")
	ud := db.ReadUserDataByUID(flowData.Userid)
	flowData.AtferTaxFee = ud.Fee + flowData.ChangeAmount
	//todo:出现错误中断的情况
	db.SaveFlowData(flowData)
	refundByFlowIDs(flowData.FlowIDs)
	return
}

func thepayment(id int, thedesc string) (code int, desc string) {
	code = util.Success
	desc = util.ErrMsg[util.Success]
	flowData := db.ReadFlowDataByID(id)
	if flowData.Status != 1 {
		code = util.Fail
		desc = util.ErrMsg[util.Fail]
		return
	}
	flowData.Status = util.FlowDataStatusOver
	//db.AddUserTakenFee(flowData)
	rpc.AddFee(flowData.Userid, flowData.ChangeAmount, "takenfee") //todo:不稳定
	ud := db.ReadUserDataByUID(flowData.Userid)
	flowData.TakenFee = ud.TakenFee + flowData.ChangeAmount
	flowData.Desc = thedesc
	//todo:错误出现的情况
	db.SaveFlowData(flowData)
	paymentByFlowIDs(flowData.FlowIDs)
	return
}

func flowDataExport(c *gin.Context) {
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
	data := c.Request.FormValue("data")
	feReq := new(param.FlowDataExportReq)
	if err := json.Unmarshal([]byte(data), feReq); err != nil {
		log.Error(err.Error())
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	fds := db.ReadExports(feReq)
	feResp := new(param.FlowDataExportResp)
	feResp.FlowExports = new([]param.FlowExports)
	fes := feResp.FlowExports

	for _, v := range *fds {
		ud := db.ReadUserDataByUID(v.Userid)
		bc := db.ReadBankCardByID(ud.UserID)
		temp := param.FlowExports{
			Accountid:    ud.AccountID,
			PhoneNum:     ud.Username,
			Realname:     ud.RealName,
			BankCardNo:   ud.BankCardNo,
			OpenBankName: bc.OpeningBank,
			ChangeAmount: v.ChangeAmount,
		}
		*fes = append(*fes, temp)
	}

	resp = fes
	return
}
