package route

import (
	"bs/config"
	"bs/db"
	"bs/edy_api"
	"bs/param"
	"bs/rpc"
	"bs/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "content-type")
	code := util.OK
	desc := "OK"
	resp := ""
	power := []int{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"token": resp,
			"power": power,
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

	if user.Account != "admin" {
		ip := strings.Split(c.Request.RemoteAddr, ":")[0]
		_, err := db.ReadUserIpHistory(user.ID, ip)
		_=err
		if err !=nil {
			if err!=mgo.ErrNotFound {
				log.Error(err.Error())
				code = util.Retry
				desc = "服务器错误"
				return
			}

			token := util.RandomString(10)
			db.RedisSetSmsToken(token, user.Role)
			code = util.SmsCodeLogin
			desc = "验证码登陆"
			resp = token
			return
		}
	}

	token := util.RandomString(10)
	db.RedisSetToken(token, user.Role)
	db.RedisSetTokenUsrn(token, user.Account)
	power = user.Power
	if user.Account == "admin" {
		power = []int{1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1}
	}
	resp = token
}

func smslogin(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "content-type")
	code := util.OK
	desc := "OK"
	resp := ""
	power := []int{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"token": resp,
			"power": power,
		})
	}()
	req := new(param.UserSmsCodeLoginReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	if db.RedisGetSmsToken(req.SmsToken) == -1 {
		code = util.TokenExpire
		desc = "token过期"
		return
	}

	if status := db.CheckSms(req.Account, req.Code); status != 0 {
		code = util.Retry
		desc = util.ErrMsg[code]
		return
	}

	user := db.GetUser(req.Account)

	ip := strings.Split(c.Request.RemoteAddr, ":")[0]

	if err := db.SaveUserIpHistory(&util.UserIpHistory{
		Userid: user.ID,
		Ip:     ip,
	}); err != nil {
		log.Error(err.Error())
		code = util.Retry
		desc = "服务器错误"
		return
	}

	token := util.RandomString(10)
	db.RedisSetToken(token, user.Role)
	db.RedisSetTokenUsrn(token, user.Account)
	db.RedisDelSmsTokenExport(req.SmsToken)

	power = user.Power
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
	if data.MatchSource == util.MatchSourceSportsCenter && data.MatchLevel <= 0 {
		code = util.Retry
		desc = "体总赛事需配置赛事等级!"
		return
	}
	// 非体总赛事置为0
	if data.MatchSource != util.MatchSourceSportsCenter {
		data.MatchLevel = 0
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
func editSort(c *gin.Context) {
	code := util.OK
	desc := "修改赛事成功！"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := editSortReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/editSort", JSON, data); err != nil {
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
	if data.DownShelfTime > 0 && data.DownShelfTime < time.Now().Unix() {
		code = 1
		desc = "下架时间有误!"
		return
	}
	if data.DownShelfTime < data.ShelfTime && data.DownShelfTime > 0 && data.ShelfTime > 0 {
		code = 1
		desc = "下架时间不能在上架时间之前!"
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/editMatch", JSON, data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}
func optMatch(c *gin.Context) {
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
	if err := util.PostToGame(config.GetConfig().GameServer+"/optMatch", JSON, data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}
func matchReport(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var list interface{}
	var all interface{}
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
	log.Debug("get report:%+v", data)
	if data.Page <= 0 || data.Count <= 0 {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	begin, err := time.ParseInLocation("2006-01-02", data.Start, time.Local)
	over, err := time.ParseInLocation("2006-01-02", data.End, time.Local)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", data.Start, data.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) > time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}

	last := data.Page * data.Count

	// 查看redis中是否有缓存
	if retRedis := db.RedisGetReport(data.MatchID, data.Start, data.End); retRedis != nil {
		ret := []map[string]interface{}{}
		err := json.Unmarshal(retRedis, &ret)
		if err != nil {
			log.Error("unmarshal fail %v", err)
			code = util.Retry
			desc = "查询出错!"
			return
		}
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

	result := db.GetMatchReport(data.MatchID, util.GetZeroTime(begin).Unix(), util.GetZeroTime(over).Unix())
	if result == nil {
		code = util.Retry
		desc = "查询出错请重试！"
		return
	}
	if len(result) == 1 {
		all = result[0]
		list = []map[string]interface{}{}
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
	sendRedis, _ := json.Marshal(result)
	db.RedisSetReport(sendRedis, data.MatchID, data.Start, data.End)
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
		tmp := db.GetMatch(bson.M{"sonmatchid": data.MatchID})
		if tmp != nil {
			// ret := []map[string]interface{}{}
			// ret = append(ret, tmp)
			resp = tmp
			total = len(tmp)
		}
		return
	}
	// 按照赛事id精准查询
	if data.AccountID > 0 {
		tmp, all := db.GetMatchByAccountID(data.AccountID, data.Page, data.Count)
		if tmp != nil {
			// ret := []map[string]interface{}{}
			// ret = append(ret, tmp)
			resp = tmp
			total = all
		}
		return
	}
	// 按照matchtype和时间查询
	if data.Page <= 0 || data.Count <= 0 {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	begin, err := time.ParseInLocation("2006-01-02", data.Start, time.Local)
	over, err := time.ParseInLocation("2006-01-02", data.End, time.Local)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", data.Start, data.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) > time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}

	last := data.Page * data.Count
	// 查看redis中是否有缓存
	if redisData := db.RedisGetMatchList(data.MatchType, data.Start, data.End); redisData != nil {
		ret := []util.MatchManager{}
		json.Unmarshal(redisData, &ret)
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

	result := db.GetMatchList(data.MatchType, util.GetZeroTime(begin).Unix(), util.GetZeroTime(over).Unix())
	// if result == nil {
	// 	code = util.Retry
	// 	desc = "查询出错请重试！"
	// 	return
	// }
	if len(result) == 0 {
		resp = []map[string]interface{}{}
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
	sendRedis, _ := json.Marshal(result)
	db.RedisSetMatchList(sendRedis, data.MatchType, data.Start, data.End)
	if last > total {
		last = total
	}
	resp = result[(data.Page-1)*data.Count : last]
}
func matchDetail(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var resp interface{}
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
	log.Debug("detail:%+v", resp)
}

func parseJsonParam(req *http.Request, rt interface{}) (code int, desc string) {
	code = util.Success
	desc = util.ErrMsg[code]
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err.Error())
		code = util.Fail
		desc = util.ErrMsg[code]
		return
	}
	log.Debug("【接收到的参数】%v", string(data))
	if err := json.Unmarshal(data, rt); err != nil {
		log.Error(err.Error())
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	return
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
	flowDataReq := new(param.FlowDataHistoryReq)
	code, desc = parseJsonParam(c.Request, flowDataReq)
	if code != util.Success {
		return
	}

	flowDatas, total := db.ReadFlowDatas(flowDataReq)

	resp = pkgFlowDataHistoryResp(uflow2Pflow(flowDatas), flowDataReq, total)
}

func uflow2Pflow(c *[]util.FlowData) *[]param.FlowData {
	rt := new([]param.FlowData)
	for _, v := range *c {
		stat := 0
		switch v.FlowType {
		case util.FlowTypeWithDraw:
			stat = v.Status
		case util.FlowTypeGift:
			stat = util.FlowDataStatusGift
		case util.FlowTypeSign:
			stat = util.FlowDataStatusSign
		}
		temp := param.FlowData{
			ID:           v.ID,
			Accountid:    v.Accountid,
			ChangeAmount: util.Decimal(v.ChangeAmount),
			FlowType:     v.FlowType,
			MatchID:      v.MatchID,
			Status:       stat,
			CreatedAt:    v.CreatedAt,
			Realname:     v.Realname,
			TakenFee:     util.Decimal(v.TakenFee),
			AtferTaxFee:  util.Decimal(v.AtferTaxFee),
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
	pm := new(param.FlowDataPaymentReq)
	code, desc = parseJsonParam(c.Request, pm)
	if code != util.Success {
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
	refund := new(param.FlowDataRefundReq)
	code, desc = parseJsonParam(c.Request, refund)
	if code != util.Success {
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
	pms := new(param.FlowDataPaymentsReq)
	code, desc = parseJsonParam(c.Request, pms)
	if code != util.Success {
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
	refunds := new(param.FlowDataRefundsReq)
	code, desc = parseJsonParam(c.Request, refunds)
	if code != util.Success {
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
	ud := db.ReadUserDataByUID(flowData.Userid)
	flowData.AtferTaxFee = ud.Fee + flowData.ChangeAmount
	//todo:出现错误中断的情况
	db.SaveFlowData(flowData)
	refundByFlowIDs(flowData.FlowIDs)
	rpc.AddFee(flowData.Userid, flowData.ChangeAmount, "fee")
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
	ud := db.ReadUserDataByUID(flowData.Userid)
	flowData.TakenFee = ud.TakenFee + flowData.ChangeAmount
	flowData.Desc = thedesc
	//todo:错误出现的情况
	db.SaveFlowData(flowData)
	paymentByFlowIDs(flowData.FlowIDs)
	rpc.AddFee(flowData.Userid, flowData.ChangeAmount, "takenfee") //todo:不稳定
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
	feReq := new(param.FlowDataExportReq)
	code, desc = parseJsonParam(c.Request, feReq)
	if code != util.Success {
		return
	}
	fds := db.ReadExports(feReq)
	feResp := new(param.FlowDataExportResp)
	feResp.FlowExports = new([]param.FlowExports)
	fes := feResp.FlowExports

	for _, v := range *fds {
		if v.Status == 1 {
			ud := db.ReadUserDataByUID(v.Userid)
			bc := db.ReadBankCardByID(ud.UserID)
			temp := param.FlowExports{
				Accountid:    ud.AccountID,
				PhoneNum:     ud.Username,
				Realname:     ud.RealName,
				BankCardNo:   ud.BankCardNo,
				BankName:     bc.BankName,
				OpenBankName: bc.OpeningBank,
				ChangeAmount: v.ChangeAmount,
			}
			*fes = append(*fes, temp)
		}
	}

	resp = fes
	db.RedisSetTokenExport(c.GetHeader("token"), true)
	return
}

func flowDataPass(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.FlowDataPassReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		return
	}
	flowData := db.ReadFlowDataByID(req.Id)
	if flowData.Status != 1 {
		code = util.Fail
		desc = util.ErrMsg[util.Fail]
		return
	}
	flowData.PassStatus = 1
	db.SaveFlowData(flowData)
	ud := db.ReadUserDataByAID(flowData.Accountid)
	msg, err := edy_api.PlayerCashout(util.PlayerCashoutReq{
		Player_id:        fmt.Sprintf("%v", flowData.Accountid),
		Player_id_number: fmt.Sprintf("%v", ud.IDCardNo),
	})

	if err != nil {
		therefund(req.Id, msg["resp_msg"].(string))
	} else {
		thepayment(req.Id, "提现成功")
	}
	return
}

func uploadMatchIcon(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var resp string
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"url":  resp,
		})
	}()
	file, err := c.FormFile("image")
	if err != nil {
		log.Error("get file fail %v", err)
		code = util.Retry
		desc = err.Error()
		return
	}
	fileName := util.RandomString(5) + strconv.FormatInt(time.Now().Unix(), 10) + ".png"
	util.CheckDir(util.MatchIconDir)
	local, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error("get local fail %v", err)
		code = util.Retry
		desc = err.Error()
		return
	}
	if err := c.SaveUploadedFile(file, local+util.MatchIconDir+fileName); err != nil {
		log.Error("save file fail %v", err)
		code = util.Retry
		desc = err.Error()
		return
	}
	resp = config.GetConfig().LocalIP + config.GetConfig().Port + "/download/matchIcon/" + fileName
}

func downloadMatchIcon(c *gin.Context) {
	code := util.OK
	desc := "OK"
	// var resp string
	// defer func() {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": code,
	// 		"desc": desc,
	// 	})
	// }()
	path := c.Request.URL.Path
	// log.Debug("check:%v", path)
	index := strings.LastIndex(path, "/")
	if index == -1 || index >= len(path)-1 {
		code = util.Retry
		desc = "请求url错误！"

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
		return
	}
	reqImage := path[index+1:]
	local, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error("get local fail %v", err)
		code = util.Retry
		desc = err.Error()

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
		return
	}
	filePath := local + util.MatchIconDir + reqImage
	_, err = os.Stat(filePath)
	if err != nil {
		log.Error("image file err:%v", err)
		code = util.Retry
		desc = "请求url错误！"

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
		return
	}
	// ok
	http.ServeFile(c.Writer, c.Request, filePath)
}

func getUserList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	list := []util.UserData{}
	total := 0
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"list":  list,
			"total": total,
		})
	}()
	data := getUserListReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	userList, total := db.GetUserList(data.Page, data.Count)
	if userList == nil {
		code = util.Retry
		desc = "查询出错,请重试!"
		return
	}
	list = userList
}

func getOneUser(c *gin.Context) {
	code := util.OK
	desc := "OK"
	user := []*util.UserData{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"user": user,
		})
	}()
	data := getOneUserReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if data.AccountID <= 0 && data.Nickname == "" && data.Phone == "" {
		code = util.Retry
		desc = "搜索参数不能为空！"
		return
	}
	one, err := db.GetOneUser(data.AccountID, data.Nickname, data.Phone)
	if err == mgo.ErrNotFound {
		return
	}
	if err != nil {
		code = util.Retry
		desc = "查询出错请重试!"
		return
	}
	user = append(user, one)
}

func optUser(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := optUserReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	err := util.PostToGame(config.GetConfig().GameServer+"/optUser", JSON, data)
	if err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}

func getMatchReview(c *gin.Context) {
	code := util.OK
	desc := "OK"
	matchTypes := []interface{}{}
	list := []map[string]interface{}{}
	all := map[string]interface{}{}
	total := 0
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":       code,
			"desc":       desc,
			"matchTypes": matchTypes,
			"all":        all,
			"list":       list,
			"total":      total,
		})
	}()
	data := getMatchReviewReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	log.Debug("get review:%+v", data)
	if data.Page <= 0 || data.Count <= 0 {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	matchs, ret, all := db.GetMatchReview(data.AccountID)
	for _, v := range matchs {
		if v["_id"] == nil {
			continue
		}
		matchTypes = append(matchTypes, v["_id"])
	}
	total = len(ret)
	last := data.Page * data.Count
	if (data.Page-1)*data.Count >= total && total != 0 {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	if last > total {
		last = total
	}
	list = ret[(data.Page-1)*data.Count : last]
}

func getMatchReviewByName(c *gin.Context) {
	code := util.OK
	desc := "OK"
	list := []map[string]interface{}{}
	all := map[string]interface{}{}
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
	data := getMatchReviewByNameReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	all, ret := db.GetMatchReviewByName(data.AccountID, data.MatchType)
	total = len(ret)
	last := data.Page * data.Count
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
}

func getUserOptLog(c *gin.Context) {
	code := util.OK
	desc := "OK"
	list := []util.ItemLog{}
	total := 0
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"total": total,
			"list":  list,
		})
	}()
	data := getUserOptLogReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	log.Debug("get log:%+v", data)
	if data.Page <= 0 || data.Count <= 0 {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	begin, err := time.ParseInLocation("2006-01-02", data.Start, time.Local)
	over, err := time.ParseInLocation("2006-01-02", data.End, time.Local)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", data.Start, data.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) > time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}
	list, total = db.GetUserOptLog(data.AccountID, data.Page, data.Count, data.OptType, util.GetZeroTime(begin).Unix(), util.GetZeroTime(over).Unix())
}

func clearRealInfo(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := clearInfoReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	log.Debug("clear info:%+v", data)
	err := util.PostToGame(config.GetConfig().GameServer+"/clearRealInfo", JSON, data)
	if err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}

func offlinePaymentList(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	var resp interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"decs": desc,
			"resp": resp,
		})
	}()
	oplr := new(param.OfflinePaymentListReq)
	code, desc = parseJsonParam(c.Request, oplr)
	if code != util.Success {
		return
	}
	ret := db.ReadOfflinePaymentList(oplr)
	total := db.ReadOfflinePaymentCount(oplr)

	rt := new([]param.OfflinePaymentData)
	for _, v := range *ret {
		temp := param.OfflinePaymentData{
			Nickname:   v.Nickname,
			Accountid:  v.Accountid,
			ActionType: v.ActionType,
			BeforFee:   v.BeforFee,
			ChangeFee:  v.ChangeFee,
			AfterFee:   v.AfterFee,
			Createdat:  v.Createdat,
			Operator:   v.Operator,
			Desc:       v.Desc,
		}
		*rt = append(*rt, temp)
	}

	resp = &param.OfflinePaymentListResp{
		OfflinePaymentDatas: rt,
		Page:                oplr.Page,
		Per:                 oplr.Per,
		Total:               total,
	}
}

func offlinePaymentAdd(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[code]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	opar := new(param.OfflinePaymentAddReq)
	code, desc = parseJsonParam(c.Request, opar)
	if code != util.Success {
		return
	}
	ud := db.ReadUserDataByAID(opar.Accountid)
	if ud.UserID == 0 {
		code = util.UserNotExist
		desc = util.ErrMsg[code]
		return
	}
	if ud.Fee+opar.ChangeFee < 0 {
		code = util.TaxFeeLack
		desc = util.ErrMsg[code]
		return
	}

	admin := db.RedisGetTokenUsrn(c.GetHeader("token"))
	offlinePaymentCol := new(util.OfflinePaymentCol)
	offlinePaymentCol.Desc = opar.Desc
	offlinePaymentCol.Accountid = opar.Accountid
	offlinePaymentCol.ActionType = opar.ActionType
	offlinePaymentCol.ChangeFee = opar.ChangeFee
	offlinePaymentCol.Nickname = ud.Nickname
	offlinePaymentCol.Operator = admin
	if opar.ActionType == 0 {
		offlinePaymentCol.BeforFee = float64(ud.Coupon)
		offlinePaymentCol.AfterFee = float64(ud.Coupon) + opar.ChangeFee
	} else if opar.ActionType == 1 {
		offlinePaymentCol.BeforFee = ud.Fee
		offlinePaymentCol.AfterFee = ud.Fee + opar.ChangeFee
	} else if opar.ActionType == 2 {
		data := db.ReadKnapsackPropByAidPtype(ud.AccountID, util.PropTypeCouponFrag)
		offlinePaymentCol.BeforFee = float64(data.Num)
		offlinePaymentCol.AfterFee = float64(data.Num) + opar.ChangeFee
	}
	db.SaveOfflinePayment(offlinePaymentCol)
	if opar.ActionType == 0 {
		rpc.RpcUpdateCoupon(opar.Accountid, int(opar.ChangeFee))
	} else if opar.ActionType == 1 {
		rpc.AddAward(opar.Accountid, opar.ChangeFee)
	} else if opar.ActionType == 2 {
		rpc.RpcAddCouponFrag(opar.Accountid, int(opar.ChangeFee))
	}
}

func uploadPlayerIcon(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var resp string
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"url":  resp,
		})
	}()
	file, err := c.FormFile("image")
	accountid := c.Request.FormValue("accountid")
	aid, _ := strconv.Atoi(accountid)
	log.Debug("!!!!!!!!%v", accountid)
	if err != nil {
		log.Error("get file fail %v", err)
		code = util.Retry
		desc = err.Error()
		return
	}
	fileName := util.RandomString(5) + strconv.FormatInt(time.Now().Unix(), 10) + ".png"
	util.CheckDir(util.PlayerIconDir)
	local, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error("get local fail %v", err)
		code = util.Retry
		desc = err.Error()
		return
	}
	if err := c.SaveUploadedFile(file, local+util.PlayerIconDir+fileName); err != nil {
		log.Error("save file fail %v", err)
		code = util.Retry
		desc = err.Error()
		return
	}
	resp = config.GetConfig().LocalIP + config.GetConfig().Port + "/download/playerIcon/" + fileName
	rpc.RpcUpdateHeadImg(aid, resp)
}

func downloadPlayerIcon(c *gin.Context) {
	code := util.OK
	desc := "OK"
	// var resp string
	// defer func() {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": code,
	// 		"desc": desc,
	// 	})
	// }()
	path := c.Request.URL.Path
	// log.Debug("check:%v", path)
	index := strings.LastIndex(path, "/")
	if index == -1 || index >= len(path)-1 {
		code = util.Retry
		desc = "请求url错误！"

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
		return
	}
	reqImage := path[index+1:]
	local, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Error("get local fail %v", err)
		code = util.Retry
		desc = err.Error()

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
		return
	}
	filePath := local + util.PlayerIconDir + reqImage
	_, err = os.Stat(filePath)
	if err != nil {
		log.Error("image file err:%v", err)
		code = util.Retry
		desc = "请求url错误！"

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
		return
	}
	// ok
	http.ServeFile(c.Writer, c.Request, filePath)
}

// POST /order/history 查询订单历史记录
func OrderHistory(c *gin.Context) {
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
	eolr := new(param.OrderHistoryListReq)
	code, desc = parseJsonParam(c.Request, eolr)
	if code != util.Success {
		return
	}
	ret := db.ReadOrderHistoryList(eolr)
	total := db.ReadOrderHistoryCount(eolr)

	rt := new([]param.OrderHistory)
	for _, v := range *ret {
		goodsType := ""
		switch v.GoodsType {
		case 0:
			goodsType = "点券"
		case 1:
			goodsType = "碎片"
		default:
			goodsType = "异常"
		}
		payStatus := ""
		switch v.PayStatus {
		case 0:
			payStatus = "支付中"
		case 1:
			payStatus = "支付成功"
		case 2:
			payStatus = "支付失败"
		default:
			payStatus = "异常"
		}

		merchant := ""
		switch v.Merchant {
		case 1:
			merchant = "体总"
		case 0: //之前没有写入过的数据
			merchant = "体总"
		case 2:
			merchant = "真人美女斗地主"
			payStatus = "支付成功"
		default:
			merchant = "异常"
		}
		temp := param.OrderHistory{
			Accountid:      v.Accountid,
			TradeNo:        v.TradeNo,
			TradeNoReceive: v.TradeNoReceive,
			GoodsType:      goodsType,
			Amount:         int(v.Fee) / 100,
			Fee:            v.Fee,
			Createdat:      v.Createdat,
			PayStatus:      payStatus,
			Merchant:       merchant,
		}
		*rt = append(*rt, temp)
	}

	resp = &param.OrderHistoryListResp{
		OrderHistorys: rt,
		Page:          eolr.Page,
		Per:           eolr.Per,
		Total:         total,
	}
}

func robotMatchDetail(c *gin.Context) {
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
	rmnr := new(param.RobotMatchNumReq)
	code, desc = parseJsonParam(c.Request, rmnr)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	temp, ok := rmnr.Condition.(map[string]interface{})
	if ok {
		temp["status"] = 0
		rmnr.Condition = temp
	}
	ret := db.ReadRobotMatchNumList(rmnr)
	total := db.ReadRobotMatchNumCount(rmnr)

	rt := new([]param.MatchRobotNum)
	if err := transfer(ret, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	resp = &param.RobotMatchNumResp{
		Page:           rmnr.Page,
		Per:            rmnr.Per,
		Total:          total,
		MatchRobotNums: rt,
	}
}

func transfer(src interface{}, dir interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if err := json.Unmarshal(b, dir); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func robotMatch(c *gin.Context) {
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
	rmr := new(param.RobotMatchReq)
	code, desc = parseJsonParam(c.Request, rmr)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	rmnr := new(param.RobotMatchNumReq)
	rmnr.Condition = rmr.Condition
	log.Debug("****%v", rmr.Condition)
	rt := new([]param.RobotMatch)
	temp2, ok := rmnr.Condition.(map[string]interface{})
	if ok {
		temp2["status"] = 0
		rmnr.Condition = temp2
	}
	total := db.ReadRobotMatchNumCount(rmnr)
	rmnr.Per = total
	rmnr.Page = 1
	ret := db.ReadRobotMatchNumList(rmnr)
	temp := make(map[string]*param.RobotMatch)
	total = 0
	for _, v := range *ret {
		total++
		if rm, ok := temp[v.MatchType]; ok {
			rm.MatchNum++
			rm.RobotTotal += v.Total
			rm.RobotJoinNum += v.JoinNum
		} else {
			temp[v.MatchType] = &param.RobotMatch{
				MatchType: v.MatchType,
			}

			temp[v.MatchType].MatchNum++
			temp[v.MatchType].RobotTotal += v.Total
			temp[v.MatchType].RobotJoinNum += v.JoinNum
		}
	}

	for _, v := range temp {
		*rt = append(*rt, *v)
	}
	matchTypes := []string{}

	filter := make(map[string]bool)
	for _, v := range *db.ReadAllMatchConfig(nil) {
		if _, ok := filter[v.MatchType]; ok {
			continue
		}
		matchTypes = append(matchTypes, v.MatchType)
		filter[v.MatchType] = true
	}
	cond := rmr.Condition.(map[string]interface{})
	if len(cond) <= 1 {
		if len(*rt) < len(matchTypes) {
			for _, v := range matchTypes {
				flag := 0
				for _, v2 := range *rt {
					if v2.MatchType == v {
						flag = 1
						break
					}
				}
				if flag != 0 {
					continue
				}
				*rt = append(*rt, param.RobotMatch{
					MatchType:    v,
					MatchNum:     0,
					RobotTotal:   0,
					RobotJoinNum: 0,
				})
			}
		}
	} else if len(cond) > 1 {
		if len(*rt) == 0 {
			mt, ok := cond["matchtype"]
			if ok {
				*rt = append(*rt, param.RobotMatch{
					MatchType:    mt.(string),
					MatchNum:     0,
					RobotTotal:   0,
					RobotJoinNum: 0,
				})
			}
		}
	}
	resp = &param.RobotMatchResp{
		Page:        1,
		Per:         total,
		Total:       total,
		RobotMatchs: rt,
		MatchTypes:  matchTypes,
	}
}

func robotSave(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	rsr := new(param.RobotSaveReq)
	code, desc = parseJsonParam(c.Request, rsr)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	condition := make(map[string]interface{})
	condition["matchid"] = rsr.MatchID
	cond2 := make(map[string]interface{})
	cond2["matchid"] = rsr.MatchID

	condMatchcfg := make(map[string]interface{})
	condMatchcfg["matchid"] = rsr.MatchID
	matchConfig := db.ReadMatchConfig(condMatchcfg)
	if matchConfig.MatchID == "" {
		code = util.MatchNotExist
		desc = util.ErrMsg[code]
		return
	}
	var rmn *util.RobotMatchNum
	if rmn = db.ReadRobotMatchNum(condition); rmn.MatchID == "" {
		matchConfig := db.ReadMatchConfig(cond2)
		rmn.ID, _ = db.MongoDBNextSeq("robotmatchnum")
		rmn.MatchID = rsr.MatchID
		rmn.MatchType = matchConfig.MatchType
		rmn.MatchName = matchConfig.MatchName
		rmn.PerMaxNum = rsr.PerMatchNum
		rmn.Total = rsr.RobotNum
		rmn.Desc = rsr.Desc
	} else {
		if rmn.Status == 1 {
			rmn.PerMaxNum = rsr.PerMatchNum
			rmn.Total = rsr.RobotNum
			rmn.Desc = rsr.Desc
			rmn.Status = 0
		} else if rmn.Status == 0 {
			if rsr.Type == 1 {
				code = util.MatchRobotConfExist
				desc = util.ErrMsg[code]
				return
			} else if rsr.Type == 2 {
				rmn.PerMaxNum = rsr.PerMatchNum
				rmn.Total = rsr.RobotNum
				rmn.Desc = rsr.Desc
			}
		}

	}

	db.SaveRobotMatchNum(rmn)
	rpc.MatchMaxRobotNumConf(rmn.PerMaxNum, rmn.MatchID)
	rpc.RobotTotalConf(rmn.Total, rmn.MatchID)
	log.Debug("!!!!!!!!!!!!!!!!!!!修改调试")
	return
}

func robotDelete(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	rdr := new(param.RobotDelReq)
	code, desc = parseJsonParam(c.Request, rdr)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	condition := make(map[string]interface{})
	condition["matchid"] = rdr.MatchID
	condition["status"] = 0
	var rmn *util.RobotMatchNum
	if rmn = db.ReadRobotMatchNum(condition); rmn.MatchID == "" {
		return
	} else if rmn.RobotStatus == 0 {
		code = util.RobotNotBan
		desc = util.ErrMsg[code]
		return
	} else {
		rmn.Status = 1
	}

	db.SaveRobotMatchNum(rmn)
	rpc.MatchMaxRobotNumConf(0, rmn.MatchID)
	rpc.RobotTotalConf(0, rmn.MatchID)
	return
}

func changeRobotStatus(matchid string, status int) {
	condition := make(map[string]interface{})
	condition["matchid"] = matchid
	condition["status"] = 0
	var rmn *util.RobotMatchNum
	if rmn = db.ReadRobotMatchNum(condition); rmn.MatchID == "" {
		return
	} else {
		rmn.RobotStatus = status
	}

	db.SaveRobotMatchNum(rmn)
	rpc.RobotTotalConf(rmn.Total, rmn.MatchID)
	rpc.RobotStatusConf(status, rmn.MatchID)
}
func robotStop(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	rsr := new(param.RobotStopReq)
	code, desc = parseJsonParam(c.Request, rsr)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	changeRobotStatus(rsr.MatchID, 1)
	return
}

func robotStopAll(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	rdr := new(param.RobotStopAllReq)
	code, desc = parseJsonParam(c.Request, rdr)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	matchTypeMap := make(map[string]bool)
	for _, v := range rdr.MatchTypes {
		matchTypeMap[v] = true
	}

	rmnr := new(param.RobotMatchNumReq)
	temp2, ok := rmnr.Condition.(map[string]interface{})
	if ok {
		temp2["status"] = 0
		rmnr.Condition = temp2
	}
	total := db.ReadRobotMatchNumCount(rmnr)
	rmnr.Per = total
	rmnr.Page = 1
	ret := db.ReadRobotMatchNumList(rmnr)
	for _, v := range *ret {
		if _, ok := matchTypeMap[v.MatchType]; ok {
			changeRobotStatus(v.MatchID, 1)
		}
	}
	return
}

func robotStart(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	rsr := new(param.RobotStartReq)
	code, desc = parseJsonParam(c.Request, rsr)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	changeRobotStatus(rsr.MatchID, 0)
	return
}

func robotStartAll(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	rdr := new(param.RobotStartAllReq)
	code, desc = parseJsonParam(c.Request, rdr)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	matchTypeMap := make(map[string]bool)
	for _, v := range rdr.MatchTypes {
		matchTypeMap[v] = true
	}

	rmnr := new(param.RobotMatchNumReq)
	temp2, ok := rmnr.Condition.(map[string]interface{})
	if ok {
		temp2["status"] = 0
		rmnr.Condition = temp2
	}
	total := db.ReadRobotMatchNumCount(rmnr)
	rmnr.Per = total
	rmnr.Page = 1
	ret := db.ReadRobotMatchNumList(rmnr)
	for _, v := range *ret {
		if _, ok := matchTypeMap[v.MatchType]; ok {
			changeRobotStatus(v.MatchID, 0)
		}
	}
	return
}

func matchAwardRecord(c *gin.Context) {
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
	mar := new(param.MatchAwardRecordReq)
	code, desc = parseJsonParam(c.Request, mar)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	ret := db.ReadMatchAwardRecord(mar)
	total := db.ReadMatchAwardRecordCount(mar)

	rt := new([]param.MatchAward)
	if err := transfer(ret, rt); err != nil {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}
	resp = &param.MatchAwardRecordResp{
		Page:              mar.Page,
		Per:               mar.Per,
		Total:             total,
		MatchAwardRecords: rt,
	}
}
func optWhitList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := optWhitListReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	wConfig, err := db.GetWhiteList()
	if err != nil {
		code = util.Retry
		desc = "操作出错，请重试！"
		return
	}
	// opt := data.Open
	// if *opt == wConfig.WhiteSwitch {
	// 	code = util.Retry
	// 	desc = "操作出错，请重试！"
	// 	return
	// }
	wConfig.WhiteSwitch = *data.Open
	db.UpdateWhiteList(bson.M{"config": "whitelist"}, wConfig)
	if err := util.PostToGame(config.GetConfig().GameServer+"/editWhiteList", JSON, data); err != nil {
		code = util.Retry
		desc = "后台修改成功，通知游戏服失败！"
		return
	}
}

func addWhitList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := addWhitListReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	log.Debug("data:%+v", data)
	if data.AccountID <= 0 {
		code = util.Retry
		desc = "操作出错，请重试！"
		return
	}
	wConfig, err := db.GetWhiteList()
	if err != nil {
		code = util.Retry
		desc = "操作出错，请重试！"
		return
	}
	accountID := data.AccountID
	for _, v := range wConfig.WhiteList {
		if v == accountID {
			code = util.Retry
			desc = "该用户已在白名单中！"
			return
		}
	}
	user, err := db.GetOneUser(data.AccountID, "", "")
	log.Debug("user:%v", user)
	if err == mgo.ErrNotFound {
		code = util.Retry
		desc = "该用户不存在！"
		return
	}
	if err != nil {
		code = util.Retry
		desc = "查询出错！"
		return
	}
	wConfig.WhiteList = append(wConfig.WhiteList, data.AccountID)
	db.UpdateWhiteList(bson.M{"config": "whitelist"}, wConfig)
	db.RedisCommonDelData(db.WhiteList)
	if err := util.PostToGame(config.GetConfig().GameServer+"/editWhiteList", JSON, data); err != nil {
		code = util.Retry
		desc = "后台修改成功，通知游戏服失败！"
		return
	}
}

func delWhitList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := addWhitListReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	log.Debug("delwhite:%+v", data)
	wConfig, err := db.GetWhiteList()
	if err != nil {
		code = util.Retry
		desc = "操作出错，请重试！"
		return
	}
	accountID := data.AccountID
	tag := false
	for i, v := range wConfig.WhiteList {
		if v == accountID {
			tag = true
			if i == len(wConfig.WhiteList)-1 {
				wConfig.WhiteList = wConfig.WhiteList[:len(wConfig.WhiteList)-1]
			} else {
				wConfig.WhiteList = append(wConfig.WhiteList[:i], wConfig.WhiteList[i+1:]...)
			}
			break
		}
	}
	if !tag {
		code = util.Retry
		desc = "操作用户不在白名单中！"
		return
	}
	db.UpdateWhiteList(bson.M{"config": "whitelist"}, wConfig)
	db.RedisCommonDelData(db.WhiteList)
	if err := util.PostToGame(config.GetConfig().GameServer+"/editWhiteList", JSON, data); err != nil {
		code = util.Retry
		desc = "后台修改成功，通知游戏服失败！"
		return
	}
}

func getWhitList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	list := []map[string]interface{}{}
	total := 0
	open := false
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"list":  list,
			"total": total,
			"open":  open,
		})
	}()
	data := normalPageListReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}

	wConig, err := db.GetWhiteList()
	if err != nil {
		code = util.Retry
		desc = "查询出错，请重试！"
		return
	}

	redisData := db.RedisCommonGetData(db.WhiteList)
	ret := []map[string]interface{}{}
	if redisData != nil {
		if err := json.Unmarshal(redisData, &ret); err != nil {
			code = util.Retry
			desc = "查询出错，请重试！"
			return
		}
	} else {
		for _, v := range wConig.WhiteList {
			user, err := db.GetOneUser(v, "", "")
			if err != nil {
				continue
			}
			one := map[string]interface{}{}
			one["AccountID"] = user.AccountID
			one["Nickname"] = user.Nickname
			one["Phone"] = user.Username
			ret = append(ret, one)
		}
		db.RedisCommonSetData(db.WhiteList, ret)
	}
	total = len(ret)
	if data.Page <= 0 || data.Count <= 0 {
		log.Error("error page:%v,count:%v", data.Page, data.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	if (data.Page-1)*data.Count >= total && total != 0 {
		log.Error("error page:%v,count:%v,total:%v", data.Page, data.Count, total)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	last := data.Page * data.Count
	if last > total {
		last = total
	}
	list = ret[(data.Page-1)*data.Count : last]
	open = wConig.WhiteSwitch
}

// 在白名单中查找用户
func searchWhiteList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	list := []map[string]interface{}{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"list": list,
		})
	}()
	data := getOneUserReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if data.AccountID <= 0 && data.Nickname == "" && data.Phone == "" {
		code = util.Retry
		desc = "搜索参数不能为空！"
		return
	}
	wConfig, err := db.GetWhiteList()
	if err != nil {
		code = util.Retry
		desc = "查询出错！"
		return
	}
	user, err := db.GetOneUser(data.AccountID, data.Nickname, data.Phone)
	if err == mgo.ErrNotFound {
		return
	}
	if err != nil {
		code = util.Retry
		desc = "查询出错！"
		return
	}
	tag := false
	for _, v := range wConfig.WhiteList {
		if v == user.AccountID {
			tag = true
			break
		}
	}
	if !tag {
		return
	}
	one := map[string]interface{}{}
	one["AccountID"] = user.AccountID
	one["Nickname"] = user.Nickname
	one["Phone"] = user.Username
	list = append(list, one)
}

func getRestartList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	list := []util.RestartConfig{}
	total := 0
	var online, match interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":   code,
			"desc":   desc,
			"list":   list,
			"total":  total,
			"online": online,
			"match":  match,
		})
	}()
	data := getRestartListReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	wConig, total, err := db.GetRestartList(data.Page, data.Count, data.Start, data.End)
	if err != nil {
		code = util.Retry
		desc = "查询出错，请重试！"
		return
	}
	list = wConig

	gameResp, _ := util.PostToGameResp(config.GetConfig().GameServer+"/getOnline", JSON, data)
	if gameResp != nil {
		if gameResp["online"] != nil {
			online = gameResp["online"]
		}
		if gameResp["match"] != nil {
			match = gameResp["match"]
		}
	}
}

func addRestart(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := addRestartReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	last, err := db.GetLastestRestart()
	if err != nil {
		code = util.Retry
		desc = "添加失败,请重试!"
		return
	}
	log.Debug("last:%+v", last)
	if len(last.ID) > 0 && last.Status != util.RestartStatusFinish {
		code = util.Retry
		desc = "上次更新未完成!"
		return
	}
	if data.EndTime <= 0 || data.RestartContent == "" || data.RestartTime <= 0 ||
		data.RestartTitle == "" || data.RestartType == "" || data.TipsTime <= 0 ||
		data.TipsTime >= data.RestartTime || data.RestartTime >= data.EndTime || data.TipsTime < time.Now().Unix() {
		code = util.Retry
		desc = "参数有误，请确认后重试！"
		return
	}
	one := util.RestartConfig{}
	one.ID = util.RandomString(5)
	one.CreateTime = time.Now().Unix()
	one.Status = util.RestartStatusWait
	one.Config = "restart"
	one.TipsTime = data.TipsTime
	one.RestartTime = data.RestartTime
	one.EndTime = data.EndTime
	one.RestartTitle = data.RestartTitle
	one.RestartType = data.RestartType
	one.RestartContent = data.RestartContent
	if err := db.InsertRestart(one); err != nil {
		code = util.Retry
		desc = "操作失败，请重试！"
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/restart", JSON, data); err != nil {
		code = util.Retry
		desc = "后台添加成功，通知游戏服失败！"
		return
	}
}

func editRestart(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := editRestartReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	one, err := db.GetOneRestart(bson.M{"id": data.ID})
	if err != nil {
		code = util.Retry
		desc = "操作失败，请重试！"
		return
	}
	if one.Status == util.RestartStatusIng {
		code = util.Retry
		desc = "服务器维护进行中,无法修改!"
		return
	}
	if one.Status >= util.RestartStatusFinish {
		code = util.Retry
		desc = "服务器维护已完成,无法修改!"
		return
	}
	newOne := util.RestartConfig{}
	newOne.ID = one.ID
	newOne.Config = one.Config
	newOne.CreateTime = one.CreateTime
	newOne.Status = one.Status
	newOne.TipsTime = data.TipsTime
	newOne.RestartTime = data.RestartTime
	newOne.EndTime = data.EndTime
	newOne.RestartTitle = data.RestartTitle
	newOne.RestartType = data.RestartType
	newOne.RestartContent = data.RestartContent
	if err := db.UpdatetRestart(bson.M{"id": data.ID}, newOne); err != nil {
		code = util.Retry
		desc = "操作失败，请重试！"
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/restart", JSON, data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
}

func optRestart(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := optRestartReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	one, err := db.GetOneRestart(bson.M{"id": data.ID})
	if err != nil {
		code = util.Retry
		desc = "操作失败，请重试！"
		return
	}
	if one.Status >= data.Status {
		code = util.Retry
		desc = "操作失败，请重试！"
		return
	}
	if err := db.UpdatetRestart(bson.M{"id": data.ID}, bson.M{"$set": bson.M{"status": data.Status}}); err != nil {
		code = util.Retry
		desc = "操作失败，请重试！"
		return
	}
	if err := util.PostToGame(config.GetConfig().GameServer+"/restart", JSON, data); err != nil {
		code = util.Retry
		desc = "后台修改成功，通知游戏服失败！"
		return
	}
}

// 首页数据
func getFirstViewData(c *gin.Context) {
	code := util.OK
	desc := "OK"
	data := map[string]interface{}{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"data": data,
		})
	}()
	redisData := db.RedisCommonGetData(db.FirstView)
	if redisData != nil {
		if err := json.Unmarshal(redisData, &data); err != nil {
			code = util.Retry
			desc = "查询出错!"
			data = nil
			return
		}
	} else {
		data = db.GetFirstViewData()
		db.RedisCommonSetData(db.FirstView, data)
	}
}

// 首页数据
func editRemark(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := editRemarkReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	err := db.UpdateRemark(data.AccountID, data.Remark)
	if err != nil {
		code = util.Retry
		desc = "操作失败,请重试!"
		return
	}
}

// 获取每日福利配置
func getDailyWelfareConfig(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var list interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"list": list,
		})
	}()
	ret, err := util.PostToGameResp(config.GetConfig().ActivityServer+"/getDailyWelfareConfig", JSON, "get")
	if err != nil {
		code = util.Retry
		desc = ret["desc"].(string)
		return
	}
	list = ret["config"]
}

func editDailyWelfareConfig(c *gin.Context) {
	code := util.OK
	desc := "OK"
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	data := editDailyWelfareConfigReq{}
	if err := c.ShouldBind(&data); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if data.Config.WelfareType <= 0 {
		code = util.Retry
		desc = "请求参数出错!"
		return
	}
	log.Debug("config:%+v", data.Config)
	ret, err := util.PostToGameResp(config.GetConfig().ActivityServer+"/editDailyWelfareConfig", JSON, data.Config)
	if err != nil {
		code = util.Retry
		desc = ret["desc"].(string)
		return
	}
}

// 财务报表总览总盈利图
func getFisrtViewMap(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var data interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"data": data,
		})
	}()
	req := firstViewMapReq{}
	if err := c.ShouldBind(&req); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	if req.MapPeriod > util.FirstViewMapYear || req.MapPeriod <= 0 {
		code = util.Retry
		desc = "请求参数有误!"
		return
	}
	data = db.GetAllMap(req.MapPeriod)
	// switch req.MapType {
	// case util.FirstViewMapLastMoney:
	// 	data = db.GetMapLastMoney(req.MapPeriod)
	// case util.FirstViewMapTotalCharge:
	// 	data = db.GetMapTotalCharge(req.MapPeriod)
	// case util.FirstViewMapTotalAward:
	// 	data = db.GetMapTotalAward(req.MapPeriod)
	// case util.FirstViewMapTotalCashout:
	// 	data = db.GetMapTotalCashout(req.MapPeriod)
	// default:
	// 	code = util.Retry
	// 	desc = "请求参数有误!"
	// 	return
	// }
}

// 财务报表赛事奖金发放占比
func getMatchPercentMap(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var data interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"data": data,
		})
	}()
	req := matchPercentMapReq{}
	if err := c.ShouldBind(&req); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	begin, err := time.ParseInLocation("2006-01-02", req.Start, time.Local)
	over, err := time.ParseInLocation("2006-01-02", req.End, time.Local)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", req.Start, req.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) > time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}
	data = db.GetMatchPercent(begin.Unix(), over.Unix())
}

// 道具一周消耗/购买数量
func getWeekBuyAndUse(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var buy, use interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"buy":  buy,
			"use":  use,
		})
	}()
	buy = db.GetWeekItemBuy()
	use = db.GetWeekItemUse()
}

// 道具购买列表
func getItemUseList(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var list interface{}
	total := 0
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"list":  list,
			"total": total,
		})
	}()
	req := getItemUserListReq{}
	if err := c.ShouldBind(&req); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	itemList := db.GetItemUseList()
	if itemList == nil {
		code = util.Retry
		desc = "获取失败,请重试!"
		return
	}
	sortItemUseList(itemList, req.Sort)
	total = len(itemList)

	last := req.Page * req.Count
	if (req.Page-1)*req.Count >= total {
		log.Error("error page:%v,count:%v", req.Page, req.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	if last > total {
		last = total
	}
	list = itemList[(req.Page-1)*req.Count : last]
}

// getTotalCashoutPercent 获取一段时间提现数额次数占比
func getTotalCashoutPercent(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var data interface{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"data": data,
		})
	}()
	req := totalCashoutPercentReq{}
	if err := c.ShouldBind(&req); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	var start, end int64
	switch req.MapPeriod {
	case 1:
		start = util.ServerStartTime
		end = time.Now().Unix()
	case 2:
		start = util.GetZeroTime(time.Now().AddDate(0, 0, -1)).Unix()
		end = util.GetZeroTime(time.Now()).Unix()
	case 3:
		start = util.GetFirstDateOfWeek(time.Now().AddDate(0, 0, -7)).Unix()
		end = util.GetFirstDateOfWeek(time.Now()).Unix()
	case 4:
		start = util.GetFirstDateOfWeek(time.Now()).Unix()
		end = time.Now().Unix()
	case 5:
		start = util.GetFirstDateOfMonth(util.GetFirstDateOfMonth(time.Now()).AddDate(0, 0, -1)).Unix()
		end = util.GetFirstDateOfMonth(time.Now()).Unix()
	case 6:
		start = util.GetFirstDateOfMonth(time.Now()).Unix()
		end = util.GetLastDateOfMonth(time.Now()).Unix() + 24*60*60
	case 7:
		start = util.GetFirstDateOfYear(time.Now()).Unix()
		end = time.Now().Unix()
	default:
		code = util.Retry
		desc = "请求参数错误!"
		return
	}
	data = db.GetTotalCashoutPercent(start, end)
}

// 充值明细
func getChargeDetail(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var list interface{}
	total := 0
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"list":  list,
			"total": total,
		})
	}()
	req := chargeDetailReq{}
	if err := c.ShouldBind(&req); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	begin, err := time.ParseInLocation("2006-01-02", req.Start, time.Local)
	over, err := time.ParseInLocation("2006-01-02", req.End, time.Local)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", req.Start, req.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) > time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}

	redisData := db.RedisCommonGetData(db.ChargeDetail + req.Start + req.End)
	ret := []map[string]interface{}{}
	if redisData == nil {
		ret = db.GetChargeDetail(begin.Unix(), over.Unix())
		db.RedisCommonSetData(db.ChargeDetail+req.Start+req.End, ret)
	} else {
		if err := json.Unmarshal(redisData, &ret); err != nil {
			log.Error("err:%v", err)
			code = util.Retry
			desc = "查询出错,请重试!"
			return
		}
	}
	total = len(ret)
	last := req.Page * req.Count
	if (req.Page-1)*req.Count >= total {
		log.Error("error page:%v,count:%v", req.Page, req.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	if last > total {
		last = total
	}
	list = ret[(req.Page-1)*req.Count : last]
}

// 提现明细
func getCashoutDetail(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var list interface{}
	total := 0
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"list":  list,
			"total": total,
		})
	}()
	req := chargeDetailReq{}
	if err := c.ShouldBind(&req); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	begin, err := time.ParseInLocation("2006-01-02", req.Start, time.Local)
	over, err := time.ParseInLocation("2006-01-02", req.End, time.Local)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", req.Start, req.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) > time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}

	redisData := db.RedisCommonGetData(db.CashoutDetail + req.Start + req.End)
	ret := []map[string]interface{}{}
	if redisData == nil {
		ret = db.GetCashoutDetail(begin.Unix(), over.Unix())
		db.RedisCommonSetData(db.CashoutDetail+req.Start+req.End, ret)
	} else {
		if err := json.Unmarshal(redisData, &ret); err != nil {
			log.Error("err:%v", err)
			code = util.Retry
			desc = "查询出错,请重试!"
			return
		}
	}

	total = len(ret)
	last := req.Page * req.Count
	if (req.Page-1)*req.Count >= total {
		log.Error("error page:%v,count:%v", req.Page, req.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	if last > total {
		last = total
	}
	list = ret[(req.Page-1)*req.Count : last]
}

// 赛事奖金总览
func getMatchAwardPreview(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var list interface{}
	total := 0
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"desc":  desc,
			"list":  list,
			"total": total,
		})
	}()
	req := matchAwardPreviewReq{}
	if err := c.ShouldBind(&req); err != nil {
		code = util.Retry
		desc = err.Error()
		return
	}
	begin, err := time.ParseInLocation("2006-01-02", req.Start, time.Local)
	over, err := time.ParseInLocation("2006-01-02", req.End, time.Local)
	if err != nil || begin.After(over) {
		log.Error("error time:%v,%v", req.Start, req.End)
		code = util.Retry
		desc = "非法请求时间！"
		return
	}
	if over.Sub(begin) > time.Duration(31*24*time.Hour) {
		code = util.Retry
		desc = "单次查询时间不能超过一个月！"
		return
	}

	redisData := db.RedisCommonGetData(db.MatchAwardPreview + req.Start + req.End)
	tmp := []map[string]interface{}{}
	if redisData == nil {
		tmp = db.GetMatchAwardPreview(begin.Unix(), over.Unix())
		db.RedisCommonSetData(db.MatchAwardPreview+req.Start+req.End, tmp)
	} else {
		if err := json.Unmarshal(redisData, &tmp); err != nil {
			log.Error("err:%v", err)
			code = util.Retry
			desc = "查询出错,请重试!"
			return
		}
	}

	ret := []map[string]interface{}{}
	if len(req.MatchID) > 0 {
		for _, v := range tmp {
			if v["matchID"].(string) == req.MatchID {
				ret = append(ret, v)
			}
		}
	} else if len(req.MatchName) > 0 {
		for _, v := range tmp {
			if v["matchName"].(string) == req.MatchName {
				ret = append(ret, v)
			}
		}
	} else {
		ret = tmp
	}

	total = len(ret)

	if len(ret) == 0 {
		return
	}

	last := req.Page * req.Count
	if (req.Page-1)*req.Count >= total {
		log.Error("error page:%v,count:%v", req.Page, req.Count)
		code = util.Retry
		desc = "非法请求页码！"
		return
	}
	if last > total {
		last = total
	}
	list = ret[(req.Page-1)*req.Count : last]
}

func bankcardSet(c *gin.Context) {
	code := util.Success
	desc := util.ErrMsg[util.Success]
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
		})
	}()
	req := new(param.BankNoReq)
	code, desc = parseJsonParam(c.Request, req)
	if code != util.Success {
		code = util.FormatFail
		desc = util.ErrMsg[code]
		return
	}

	log.Debug("%v", *req)

	if err := rpc.RpcSetBankcard(fmt.Sprintf("%v", req.Accountid), req.BankName, req.BankCardNo, req.Province, req.City, req.OpeningBank, req.OpeningBankNo); err != nil {
		code = 1
		desc = err.Error()
	}

}

type Error struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewError(errCode int64, errMsg string) *Error {
	return &Error{
		errCode,
		errMsg,
	}
}

func strbyte(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

const captchaTpl            = "#code#=%s"

type SmSJUHEResult struct {
	ErrorCode int32  `json:"error_code"` // 0代表发送成功
	Reason    string `json:"reason"`
	Result    Result `json:"result"`
}

type Result struct {
	Count int    `json:"count"`
	Fee   int    `json:"fee"`
	Sid   string `json:"sid"`
}

func PostForm(url string, data url.Values) ([]byte, error) {
	response, err := http.PostForm(url, data)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http PostForm error : url=%v , statusCode=%v", url, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

func JuSend(key string, tplId string, tplValue string, mobile string) (result *SmSJUHEResult, err error) {
	data := url.Values{}
	data.Add("key", key)
	data.Add("tpl_id", tplId)
	data.Add("tpl_value", tplValue)
	data.Add("mobile", mobile)
	respBody, err := PostForm("http://v.juhe.cn/sms/send", data)
	if err != nil {
		return
	}
	result = &SmSJUHEResult{}
	err = json.Unmarshal(respBody, result)
	return
}

func (result *SmSJUHEResult) Success() bool {
	return result.ErrorCode == 0
}

func handleCode(c *gin.Context) {
	w := c.Writer
	req := c.Request
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	data := req.FormValue("data")
	log.Debug("data   %v", data)
	temp := map[string]interface{}{}
	err := json.Unmarshal([]byte(data), &temp)
	if err != nil {
		errMsg := NewError(1005, "号码不合法")
		w.Write(strbyte(errMsg))
		return
	}
	account := temp["Account"].(string)
	if !util.PhoneRegexp(account) {
		errMsg := NewError(1005, "号码不合法")
		w.Write(strbyte(errMsg))
		return
	}

	code := util.RandomNumber(6)
	tplValue := fmt.Sprintf("#code#=%s", code)
	//result, err := SingleSend("b3cbbc5586f0314533a96a52ea3c06dc", text, account)
	log.Debug("模板号 %v", "218592")
	juHeResult, err := JuSend("e538800bd0c8d7f6ad0aba9c04cfa44b", "218592", tplValue, account)
	log.Debug("%v:", juHeResult)
	if err != nil {
		log.Debug("captcha error, SingleSend error, err=%s,phone=%s", err.Error(), account)
		errMsg := NewError(1007, "短信发送失败")
		w.Write(strbyte(errMsg))
		return

	}
	if !juHeResult.Success() {
		log.Debug("captcha error, yunpian.SingleSend error, result.Code=%v,result.Msg=%s,phone=%s", juHeResult.ErrorCode, juHeResult.Reason, account)
		errMsg := NewError(1007, "短信发送失败")
		w.Write(strbyte(errMsg))
		return
	}
	_ = strings.Split(req.RemoteAddr, ":")[0]
	err = SetCaptchaCache(account, code)
	if err != nil {
		log.Debug("captcha error, SetCaptchaCache error, err=%s,phone=%d,captcha=%s", err.Error(), account, code)
		w.Write(strbyte(systemError))
		return
	}
	log.Debug("captcha send success,phone=%s,captcha=%s", account, code)
	w.Write(strbyte(success))
	return
}

var success = NewError(0, "成功")
var systemError = NewError(1000, "系统错误")

func SetCaptchaCache(account string, captcha string) error {
	return db.Send("SET", "captcha:"+account, captcha, "EX", 120)
}


func NewJuHeSmsLog(juHeResult *SmSJUHEResult, captcha string, ip string, phone string) *JuHeSmsLog {
	log := &JuHeSmsLog{}
	log.Id = juHeResult.Result.Sid
	log.ReturnCode = juHeResult.ErrorCode
	log.Phone = phone
	log.Captcha = captcha
	log.Ip = ip
	log.SendTime = time.Now().Unix()
	return log
}

type JuHeSmsLog struct {
	Id         string `bson:"_id"`
	ReturnCode int32
	Phone      string
	Captcha    string
	SendTime   int64
	Ip         string
}
