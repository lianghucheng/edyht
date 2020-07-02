package util

// 后台账号管理员分组
const (
	Admin = iota
	Business
	Operate
	Normal
)

// 赛事状态
const (
	Signing = iota // 报名中
	Playing        // 比赛中
	Ending         // 结算中
	Cancel         // 下架赛事
	Delete         // 删除赛事
)

// User 用户类
type User struct {
	Account  string `bson:"account"`
	Password string `bson:"password"`
	Role     int    `bson:"role"`
}

// MatchManager 比赛类
type MatchManager struct {
	MatchID     string `bson:"matchid"`     // 赛事id号（与赛事管理的matchid不是同一个，共用一个字段）
	MatchType   string `bson:"matchtype"`   // 赛事类型
	MatchName   string `bson:"matchname"`   // 赛事名称
	RoundNum    string `bson:"roundnum"`    // 赛制制(2局1副)
	StartTime   int64  `bson:"starttime"`   // 比赛开始时间
	LimitPlayer int    `bson:"limitplayer"` // 比赛开始的最少人数
	Recommend   string `bson:"recommend"`   // 赛事推荐介绍(在赛事列表界面倒计时左侧的文字信息)
	TotalMatch  int    `bson:"totalmatch"`  // 后台配置的该种比赛可创建的比赛次数
	UseMatch    int    `bson:"usematch"`    // 已使用次数
	Eliminate   []int  `bson:"eliminate"`   // 每轮淘汰人数
	EnterFee    int64  `bson:"enterfee"`    // 报名费
	ShelfTime   int64  `bson:"shelftime"`   // 上架时间
	Sort        int    `bson:"sort"`        // 赛事排序
	State       int    `bson:"state"`       // 赛事状态
}
