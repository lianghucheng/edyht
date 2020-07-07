package util

import "os"

// 后台账号管理员分组
const (
	Admin = iota
	Business
	Operate
	Normal
)

// 文件目录
const (
	MatchIconDir = string(os.PathSeparator) + "upload" + string(os.PathSeparator) + "matchIcon" + string(os.PathSeparator)
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

const (
	FlowDataStatusNormal = 0 //比赛获得
	FlowDataStatusAction = 1 //提奖中
	FlowDataStatusOver   = 2 //已提奖
	FlowDataStatusBack   = 3 //已退款
)

type FlowData struct {
	ID           int `bson:"_id"`
	Userid       int
	Accountid    int
	ChangeAmount float64
	FlowType     int
	MatchType    string
	MatchID      string
	Status       int
	CreatedAt    int64
	FlowIDs      []int
	Realname     string
	TakenFee     float64
	AtferTaxFee  float64
	Desc         string
}

type UserData struct {
	UserID            int `bson:"_id"`
	AccountID         int
	Nickname          string
	Headimgurl        string
	Sex               int // 1 男性，2 女性
	LoginIP           string
	Token             string
	ExpireAt          int64 // token 过期时间
	Role              int   // 1 玩家、2 代理、3 管理员、4 超管
	Username          string
	Password          string
	Coupon            int64 // 点券
	Wins              int   // 胜场
	CreatedAt         int64
	UpdatedAt         int64
	PlayTimes         int     //当天对局次数
	Online            bool    //玩家是否在线
	Channel           int     //渠道号。0：圈圈   1：搜狗   2:IOS
	Fee               float64 //税后余额
	SignTimes         int
	DailySign         bool
	DailySignDeadLine int64
	LastTakenMail     int64
	RealName          string
	IDCardNo          string
	BankCardNo        string
	SetNickNameCount  int
	TakenFee          float64
	FirstLogin        bool
}

type BankCard struct {
	Userid      int
	BankName    string
	BankCardNo  string
	Province    string
	City        string
	OpeningBank string
}