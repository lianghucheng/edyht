package route

import (
	"bs/config"
	"bs/db"
	"bs/param"
	"bs/rpc"
	"bs/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/szxby/tools/log"
	"gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "content-type")
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
	db.RedisSetTokenUsrn(token, user.Account)
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

	result := db.GetMatchReport(data.MatchID, begin.Unix(), over.Unix())
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
	resp := []map[string]interface{}{}
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
		tmp := db.GetMatch(data.MatchID)
		if tmp != nil {
			resp = append(resp, tmp)
			total = 1
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
	if redisData := db.RedisGetMatchList(data.MatchType, data.Start, data.End); redisData != nil {
		ret := []map[string]interface{}{}
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

	result := db.GetMatchList(data.MatchType, begin.Unix(), over.Unix())
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
		if v.FlowType != 1 {
			stat = v.Status
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

func getGameVersion(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var version, url string
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"desc":    desc,
			"version": version,
			"url":     url,
		})
	}()
	version, url = db.GetGameVersion()
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
	if data.AccountID <= 0 && data.Nickname == "" {
		code = util.Retry
		desc = "搜索参数不能为空！"
		return
	}
	one, err := db.GetOneUser(data.AccountID, data.Nickname)
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
	list, total = db.GetUserOptLog(data.AccountID, data.Page, data.Count, data.OptType, begin.Unix(), over.Unix())
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
	}
	db.SaveOfflinePayment(offlinePaymentCol)
	if opar.ActionType == 0 {
		rpc.RpcUpdateCoupon(opar.Accountid, int(opar.ChangeFee))
	} else if opar.ActionType == 1 {
		rpc.AddAward(opar.Accountid, opar.ChangeFee)
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
