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
	server.POST("/cancelMatch", cancelMatch)
	server.POST("/deleteMatch", deleteMatch)
	server.POST("/matchReport", matchReport)
	server.POST("/matchList", matchList)
	server.POST("/matchDetail", matchDetail)
	server.GET("/flowdata/history", flowDataHistory)
	server.POST("/flowdata/payment", flowDataPayment)
	server.POST("/flowdata/refund", flowDataRefund)
	server.POST("/flowdata/payments", flowDataPayments)
	server.POST("/flowdata/refunds", flowDataRefunds)
	server.GET("/flowdata/export", flowDataExport)
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
	MatchID     string   `json:"MatchID" binding:"required"`     // 赛事id号
	MatchType   string   `json:"MatchType" binding:"required"`   // 赛事类型
	MatchName   string   `json:"MatchName" binding:"required"`   // 赛事名称
	MatchDesc   string   `json:"MatchDesc" binding:"required"`   // 赛事描述
	Round       int      `json:"Round" binding:"required"`       // 赛制几局
	Card        int      `json:"Card" binding:"required"`        // 赛制几副
	StartType   int      `json:"StartType" binding:"required"`   // 比赛开始类型
	StartTime   int64    `json:"StartTime"`                      // 比赛开始时间
	LimitPlayer int      `json:"LimitPlayer" binding:"required"` // 比赛开始的最少人数
	Recommend   string   `json:"Recommend" binding:"required"`   // 赛事推荐介绍(在赛事列表界面倒计时左侧的文字信息)
	TotalMatch  int      `json:"TotalMatch" binding:"required"`  // 后台配置的该种比赛可创建的比赛次数
	Eliminate   []int    `json:"Eliminate"`                      // 每轮淘汰人数
	EnterFee    int64    `json:"EnterFee" binding:"required"`    // 报名费
	ShelfTime   int64    `json:"ShelfTime" binding:"required"`   // 上架时间
	Sort        int      `json:"Sort" binding:"required"`        // 赛事排序
	AwardDesc   string   `json:"AwardDesc" binding:"required"`   // 奖励描述
	AwardList   string   `json:"AwardList" binding:"required"`   // 奖励列表
	TablePlayer int      `json:"TablePlayer" binding:"required"` // 一桌的游戏人数
	OfficalIDs  []string `json:"OfficalIDs"`                     // 后台配置的可用比赛id号
}

type editManagerReq struct {
	MatchID    string `json:"MatchID" binding:"required"`    // 赛事id号
	TotalMatch int    `json:"TotalMatch" binding:"required"` // 后台配置的该种比赛可创建的比赛次数
	Eliminate  []int  `json:"Eliminate"`                     // 每轮淘汰人数
	EnterFee   int64  `json:"EnterFee" binding:"required"`   // 报名费
	AwardList  string `json:"AwardList" binding:"required"`  // 奖励列表
	MatchIcon  string `json:"MatchIcon" binding:"required"`  // 赛事图标
}

type showHallReq struct {
	MatchID  string `json:"MatchID" binding:"required"`  // 赛事id号
	ShowHall bool   `json:"ShowHall" binding:"required"` // 是否首页展示
}

type optMatchReq struct {
	MatchID string `json:"MatchID" binding:"required"` // 赛事id号
}

type matchReportReq struct {
	MatchID string `json:"MatchID" binding:"required"` // 赛事id号
	Start   string `json:"Start" binding:"required"`   // 查询开始时间,格式"2006-01-02"
	End     string `json:"End" binding:"required"`     // 查询结束时间
	Page    int    `json:"Page" binding:"required"`    // 查询开始时间
	Count   int    `json:"Count" binding:"required"`   // 查询结束时间
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
