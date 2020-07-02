package route

import (
	"github.com/gin-gonic/gin"
)

func bind(server *gin.Engine) {
	server.POST("/login", login)
	server.POST("/matchManagerList", matchManagerList)
	server.POST("/addMatch", addMatch)
	server.POST("/editMatch", editMatch)
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
	Account string `json:"account" binding:"required"`
	Pass    string `json:"password" binding:"required"`
}

type matchManagerReq struct {
	Page  int `json:"page" binding:"required"`
	Count int `json:"count" binding:"required"`
}

type addManagerReq struct {
	MatchID     string   `json:"matchid" binding:"required"`     // 赛事id号（与赛事管理的matchid不是同一个，共用一个字段）
	MatchType   string   `json:"matchtype" binding:"required"`   // 赛事类型
	MatchName   string   `json:"matchname" binding:"required"`   // 赛事名称
	MatchDesc   string   `json:"matchdesc" binding:"required"`   // 赛事描述
	Round       int      `json:"round" binding:"required"`       // 赛制几局
	Card        int      `json:"card" binding:"required"`        // 赛制几副
	StartType   int      `json:"starttype" binding:"required"`   // 比赛开始类型
	StartTime   int64    `json:"starttime"`                      // 比赛开始时间
	LimitPlayer int      `json:"limitplayer" binding:"required"` // 比赛开始的最少人数
	Recommend   string   `json:"recommend" binding:"required"`   // 赛事推荐介绍(在赛事列表界面倒计时左侧的文字信息)
	TotalMatch  int      `json:"totalmatch" binding:"required"`  // 后台配置的该种比赛可创建的比赛次数
	Eliminate   []int    `json:"eliminate"`                      // 每轮淘汰人数
	EnterFee    int64    `json:"enterfee" binding:"required"`    // 报名费
	ShelfTime   int64    `json:"shelftime" binding:"required"`   // 上架时间
	Sort        int      `json:"sort" binding:"required"`        // 赛事排序
	AwardDesc   string   `json:"awarddesc" binding:"required"`   // 奖励描述
	AwardList   string   `json:"awardlist" binding:"required"`   // 奖励列表
	TablePlayer int      `json:"tableplayer" binding:"required"` // 一桌的游戏人数
	OfficalIDs  []string `json:"officalIDs"`                     // 后台配置的可用比赛id号
}

type editManagerReq struct {
	MatchID    string `json:"matchid" binding:"required"`    // 赛事id号（与赛事管理的matchid不是同一个，共用一个字段）
	TotalMatch int    `json:"totalmatch" binding:"required"` // 后台配置的该种比赛可创建的比赛次数
	Eliminate  []int  `json:"eliminate"`                     // 每轮淘汰人数
	EnterFee   int64  `json:"enterfee" binding:"required"`   // 报名费
	AwardList  string `json:"awardlist" binding:"required"`  // 奖励列表
}

type optMatchReq struct {
	MatchID string `json:"matchid" binding:"required"` // 赛事id号（与赛事管理的matchid不是同一个，共用一个字段）
}

type matchReportReq struct {
	Start int64 `json:"start" binding:"required"` // 查询开始时间
	End   int64 `json:"end" binding:"required"`   // 查询结束时间
	Page  int64 `json:"page" binding:"required"`  // 查询开始时间
	Count int64 `json:"count" binding:"required"` // 查询结束时间
}
