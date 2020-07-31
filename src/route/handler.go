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
	"gopkg.in/mgo.v2/bson"

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
	log.Debug("%v", flowDataReq.Condition)

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
		data := db.ReadKnapsackPropByAidPid(ud.AccountID, 10003)
		log.Debug("************%v", *data)
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
	if last.Status != util.RestartStatusFinish {
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
	data = db.GetFirstViewData()
}
