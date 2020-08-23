package route

import (
	"github.com/gin-gonic/gin"
)

func bind(server *gin.Engine) {
	server.POST("/login", login)
	server.OPTIONS("/login", login)
	server.POST("/matchManagerList", matchManagerList)
	server.POST("/addMatch", addMatch)
	server.POST("/editMatch", editMatch)
	server.POST("/showHall", showHall)
	server.POST("/editSort", editSort)
	server.POST("/optMatch", optMatch)
	server.POST("/matchReport", matchReport)
	server.POST("/matchList", matchList)
	server.POST("/matchDetail", matchDetail)
	server.POST("/flowdata/history", flowDataHistory)
	server.POST("/flowdata/payment", flowDataPayment)
	server.POST("/flowdata/refund", flowDataRefund)
	server.POST("/flowdata/payments", flowDataPayments)
	server.POST("/flowdata/refunds", flowDataRefunds)
	server.POST("/flowdata/export", flowDataExport)
	server.POST("/flowdata/pass", flowDataPass)

	server.POST("/getUserList", getUserList)
	server.POST("/getOneUser", getOneUser)
	server.POST("/optUser", optUser)
	server.POST("/getMatchReview", getMatchReview)
	server.POST("/getMatchReviewByName", getMatchReviewByName)
	server.POST("/getUserOptLog", getUserOptLog)
	server.POST("/clearRealInfo", clearRealInfo)
	server.POST("/optWhitList", optWhitList)
	server.POST("/addWhitList", addWhitList)
	server.POST("/delWhitList", delWhitList)
	server.POST("/getWhitList", getWhitList)
	server.POST("/searchWhiteList", searchWhiteList)

	server.GET("/download/matchIcon/*action", downloadMatchIcon)
	server.POST("/upload/matchIcon", uploadMatchIcon)
	server.GET("/download/playerIcon/*action", downloadPlayerIcon)
	server.POST("/upload/playerIcon", uploadPlayerIcon)
	server.GET("/getGameVersion", getGameVersion)
	server.GET("/getNotice", getNotice)
	server.POST("/editRemark", editRemark)

	server.POST("/offlinepayment/list", offlinePaymentList)
	server.POST("/offlinepayment/add", offlinePaymentAdd)

	server.POST("/order/history", OrderHistory)
	server.POST("/robot/match-detail", robotMatchDetail)
	server.POST("/robot/match", robotMatch)
	server.POST("/robot/save", robotSave)
	server.POST("/robot/delete", robotDelete)
	server.POST("/robot/stop", robotStop)
	server.POST("/robot/stop-all", robotStopAll)
	server.POST("/robot/start", robotStart)
	server.POST("/robot/start-all", robotStartAll)
	server.POST("/match/award-record", matchAwardRecord)

	server.POST("/getRestartList", getRestartList)
	server.POST("/addRestart", addRestart)
	server.POST("/editRestart", editRestart)
	server.POST("/optRestart", optRestart)
	server.GET("/getFirstViewData", getFirstViewData)
	server.POST("/shop/merchant-insert", shopMerchantInsert)
	server.POST("/shop/merchant-list", shopMerchantList)
	server.POST("/shop/merchant-update", shopMerchantUpdate)
	server.POST("/shop/merchant-delete", shopMerchantDelete)
	server.POST("/shop/payaccount-insert", shopPayAccountInsert)
	server.POST("/shop/payaccount-delete", shopPayAccountDelete)
	server.POST("/shop/payaccount-list", shopPayAccountList)
	server.POST("/shop/payaccount-update", shopPayAccountUpdate)
	server.POST("/shop/goodstype-insert", shopGoodsTypeInsert)
	server.POST("/shop/goodstype-update", shopGoodsTypeUpdate)
	server.POST("/shop/goodstype-delete", shopGoodsTypeDelete)
	server.POST("/shop/goodstype-list", shopGoodsTypeList)
	server.POST("/shop/goods-insert", shopGoodsInsert)
	server.POST("/shop/goods-delete", shopGoodsDelete)
	server.POST("/shop/goods-list", shopGoodsList)
	server.POST("/shop/goods-update", shopGoodsUpdate)

	server.POST("/feedback/list", feedbackList)
	server.POST("/feedback/update", feedbackUpdate)

	server.POST("/propbase/config-insert", propBaseConfigInsert)
	server.POST("/propbase/config-delete", propBaseConfigDelete)
	server.POST("/propbase/config-read", propBaseConfigRead)
	server.POST("/propbase/config-list", propBaseConfigList)
	server.POST("/propbase/config-update", propBaseConfigUpdate)

	server.POST("/mailcontrol/insert", mailcontrolInsert)
	server.POST("/mailcontrol/delete", mailcontrolDelete)
	server.POST("/mailcontrol/read", mailcontrolRead)
	server.POST("/mailcontrol/list", mailcontrolList)
	server.POST("/mailcontrol/update", mailcontrolUpdate)
}

type loginData struct {
	Account  string `json:"Account" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

type matchManagerReq struct {
	Page  int `json:"Page" binding:"required"`
	Count int `json:"Count" binding:"required"`
}

type addManagerReq struct {
	MatchSource int      `json:"MatchSource" binding:"required"` // 赛事来源
	MatchLevel  int      `json:"MatchLevel"`                     // 赛事级别
	MatchID     string   `json:"MatchID" binding:"required"`     // 赛事id号
	MatchType   string   `json:"MatchType" binding:"required"`   // 赛事类型
	MatchName   string   `json:"MatchName" binding:"required"`   // 赛事名称
	MatchDesc   string   `json:"MatchDesc"`                      // 赛事描述
	Round       int      `json:"Round" binding:"required"`       // 赛制几局
	Card        int      `json:"Card" binding:"required"`        // 赛制几副
	StartType   int      `json:"StartType" binding:"required"`   // 比赛开始类型
	StartTime   int64    `json:"StartTime"`                      // 比赛开始时间
	LimitPlayer int      `json:"LimitPlayer" binding:"required"` // 比赛开始的最少人数
	Recommend   string   `json:"Recommend" binding:"required"`   // 赛事推荐介绍(在赛事列表界面倒计时左侧的文字信息)
	TotalMatch  int      `json:"TotalMatch" binding:"required"`  // 后台配置的该种比赛可创建的比赛次数
	Eliminate   []int    `json:"Eliminate"`                      // 每轮淘汰人数
	EnterFee    *int64   `json:"EnterFee" binding:"required"`    // 报名费
	ShelfTime   int64    `json:"ShelfTime" binding:"required"`   // 上架时间
	Sort        int      `json:"Sort" binding:"required"`        // 赛事排序
	AwardDesc   string   `json:"AwardDesc"`                      // 奖励描述
	AwardList   string   `json:"AwardList" binding:"required"`   // 奖励列表
	TablePlayer int      `json:"TablePlayer"`                    // 一桌的游戏人数
	OfficalIDs  []string `json:"OfficalIDs"`                     // 后台配置的可用比赛id号
	MatchIcon   string   `json:"MatchIcon" binding:"required"`   // 赛事图标
}

type editManagerReq struct {
	MatchID       string `json:"MatchID" binding:"required"` // 赛事id号
	MatchName     string `json:"MatchName"`                  // 赛事名称
	TotalMatch    int    `json:"TotalMatch"`                 // 后台配置的该种比赛可创建的比赛次数
	Eliminate     []int  `json:"Eliminate"`                  // 每轮淘汰人数
	EnterFee      *int64 `json:"EnterFee"`                   // 报名费
	AwardList     string `json:"AwardList"`                  // 奖励列表
	MatchIcon     string `json:"MatchIcon"`                  // 赛事图标
	Card          int    `json:"Card"`                       // 赛制几副
	StartType     int    `json:"StartType"`                  // 比赛开始类型
	StartTime     int64  `json:"StartTime"`                  // 比赛开始时间
	ShelfTime     int64  `json:"ShelfTime"`                  // 上架时间
	DownShelfTime int64  `json:"DownShelfTime"`              // 下架时间
	LimitPlayer   int    `json:"LimitPlayer"`                // 比赛开始的最少人数 '添加赛事时的必填字段'
}

type showHallReq struct {
	MatchID  string `json:"MatchID" binding:"required"` // 赛事id号
	ShowHall bool   `json:"ShowHall"`                   // 是否首页展示
}

type editSortReq struct {
	MatchID string `json:"MatchID" binding:"required"` // 赛事id号
	Sort    int    `json:"Sort"`                       // 赛事排序
}

type optMatchReq struct {
	MatchID string `json:"MatchID" binding:"required"` // 赛事id号
	Opt     int    `json:"Opt" binding:"required"`     // 操作代码符，1上架赛事，2下架赛事，3删除赛事
}

type matchReportReq struct {
	MatchID string `json:"MatchID" binding:"required"` // 赛事id号
	Start   string `json:"Start" binding:"required"`   // 查询开始时间,格式"2006-01-02"
	End     string `json:"End" binding:"required"`     // 查询结束时间
	Page    int    `json:"Page" binding:"required"`
	Count   int    `json:"Count" binding:"required"`
}

type matchListReq struct {
	MatchID   string `json:"MatchID"`   // 赛事id号
	MatchType string `json:"MatchType"` // 赛事类型
	Start     string `json:"Start"`     // 查询开始时间,格式"2006-01-02""
	End       string `json:"End"`       // 查询结束时间
	Page      int    `json:"Page"`      // 查询开始时间
	Count     int    `json:"Count"`     // 查询结束时间
}

type matchDetailReq struct {
	MatchID string `json:"MatchID" binding:"required"` // 赛事id号
}

type getUserListReq struct {
	Page  int `json:"Page" binding:"required"`
	Count int `json:"Count" binding:"required"`
}

type getOneUserReq struct {
	AccountID int    `json:"AccountID"`
	Nickname  string `json:"Nickname"`
	Phone     string `json:"Phone"`
}

type optUserReq struct {
	UID int `json:"UID" binding:"required"`
	Opt int `json:"Opt" binding:"required"`
}

type getMatchReviewReq struct {
	AccountID int `json:"AccountID" binding:"required"`
	Page      int `json:"Page" binding:"required"`
	Count     int `json:"Count" binding:"required"`
}

type getMatchReviewByNameReq struct {
	AccountID int    `json:"AccountID" binding:"required"`
	MatchType string `json:"MatchType" binding:"required"`
	Page      int    `json:"Page" binding:"required"`
	Count     int    `json:"Count" binding:"required"`
}

type getUserOptLogReq struct {
	AccountID int    `json:"AccountID" binding:"required"`
	Start     string `json:"Start"` // 查询开始时间,格式"2006-01-02""
	End       string `json:"End"`   // 查询结束时间
	Page      int    `json:"Page" binding:"required"`
	Count     int    `json:"Count" binding:"required"`
	OptType   int    `json:"OptType"`
}

type clearInfoReq struct {
	UID int `json:"UID" binding:"required"`
	Opt int `json:"Opt" binding:"required"`
}

type optWhitListReq struct {
	Open *bool `json:"Open"`
}

type addWhitListReq struct {
	AccountID int `json:"AccountID"`
}

type normalPageListReq struct {
	Page  int `json:"Page" binding:"required"`
	Count int `json:"Count" binding:"required"`
}

type addRestartReq struct {
	TipsTime       int64  `json:"TipsTime" binding:"required"`
	RestartTime    int64  `json:"RestartTime" binding:"required"`
	EndTime        int64  `json:"EndTime" binding:"required"`
	RestartTitle   string `json:"RestartTitle" binding:"required"`
	RestartType    string `json:"RestartType" binding:"required"`
	RestartContent string `json:"RestartContent" binding:"required"`
}

type editRestartReq struct {
	ID             string `json:"ID" binding:"required"`
	TipsTime       int64  `json:"TipsTime" binding:"required"`
	RestartTime    int64  `json:"RestartTime" binding:"required"`
	EndTime        int64  `json:"EndTime" binding:"required"`
	RestartTitle   string `json:"RestartTitle" binding:"required"`
	RestartType    string `json:"RestartType" binding:"required"`
	RestartContent string `json:"RestartContent" binding:"required"`
}

type optRestartReq struct {
	ID     string `json:"ID" binding:"required"`
	Status int    `json:"Status" binding:"required"`
}

type getRestartListReq struct {
	Start int64 `json:"Start" `
	End   int64 `json:"End"`
	Page  int   `json:"Page" binding:"required"`
	Count int   `json:"Count" binding:"required"`
}

type editRemarkReq struct {
	AccountID int    `json:"AccountID" binding:"required"`
	Remark    string `json:"Remark" binding:"required"`
}
