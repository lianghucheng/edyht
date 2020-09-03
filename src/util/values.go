package util

import (
	"os"
)

// 后台账号管理员分组
const (
	Admin = iota
	Business
	Operate
	Normal
)

// 文件目录
const (
	MatchIconDir  = string(os.PathSeparator) + "upload" + string(os.PathSeparator) + "matchIcon" + string(os.PathSeparator)
	PlayerIconDir = string(os.PathSeparator) + "upload" + string(os.PathSeparator) + "playerIcon" + string(os.PathSeparator)
)

// 一些常量
const (
	CouponRate            = 1 // CouponRate 点券与rmb比例
	ServerStartTime int64 = 1594983600
)

// 赛事状态
const (
	Signing = iota // 报名中
	Playing        // 比赛中
	Ending         // 结算中
	Cancel         // 下架赛事
	Delete         // 删除赛事
)

// 赛事来源
const (
	MatchSourceSportsCenter = iota + 1 // 体总
	MatchSourceBackstage               // 后台
)

// 财务报表首页图
const (
	FirstViewMapLastMoney    = iota + 1 // 剩余数额图
	FirstViewMapTotalCharge             // 总充值图
	FirstViewMapTotalAward              // 总奖金发放
	FirstViewMapTotalCashout            // 总提现图
)

// 财务报表首页请求周期
const (
	FirstViewMapDay = iota + 1
	FirstViewMapWeek
	FirstViewMapMonth
	FirstViewMapYear
)

// User 用户类
type User struct {
	Account  string `bson:"account"`
	Password string `bson:"password"`
	Role     int    `bson:"role"`
}

// MatchManager 比赛类
type MatchManager struct {
	MatchSource   int    `bson:"matchsource"`   // 比赛来源,1体总,2自己后台
	MatchLevel    int    `bson:"matchlevel"`    // 体总赛事级别
	MatchID       string `bson:"matchid"`       // 赛事id号（与赛事管理的matchid不是同一个，共用一个字段）
	SonMatchID    string `bson:"sonmatchid"`    // 子赛事id
	MatchType     string `bson:"matchtype"`     // 赛事类型
	MatchName     string `bson:"matchname"`     // 赛事名称
	MatchIcon     string `bson:"matchicon"`     // 赛事图标
	RoundNum      string `bson:"roundnum"`      // 赛制制(2局1副)
	StartTime     int64  `bson:"starttime"`     // 比赛开始时间
	StartType     int    `bson:"starttype"`     // 开赛条件(1表示满足三人即可开赛,2表示倒计时多久开赛判断,3表示比赛到点开赛) '添加赛事时的必填字段'
	LimitPlayer   int    `bson:"limitplayer"`   // 比赛开始的最少人数
	Recommend     string `bson:"recommend"`     // 赛事推荐介绍(在赛事列表界面倒计时左侧的文字信息)
	TotalMatch    int    `bson:"totalmatch"`    // 后台配置的该种比赛可创建的比赛次数
	UseMatch      int    `bson:"usematch"`      // 已使用次数
	LastMatch     int    `bson:"-"`             // 剩余次数
	Eliminate     []int  `bson:"eliminate"`     // 每轮淘汰人数
	EnterFee      int64  `bson:"enterfee"`      // 报名费
	ShelfTime     int64  `bson:"shelftime"`     // 上架时间
	DownShelfTime int64  `bson:"downshelftime"` // 下架时间
	EndTime       int64  `bson:"endtime"`       // 结束时间
	Sort          int    `bson:"sort"`          // 赛事排序
	State         int    `bson:"state"`         // 赛事状态
	ShowHall      bool   `bson:"showhall"`      // 是否首页展示
	AwardList     string `bson:"awardlist"`     // 奖励列表 '添加赛事时的必填字段'
	CreateTime    int64  `bson:"createtime"`    // 比赛创建时间
}

const (
	FlowTypeAward    = 1
	FlowTypeWithDraw = 2
	FlowTypeGift     = 3
	FlowTypeSign     = 4
)

const (
	FlowDataStatusNormal = 0
	FlowDataStatusAction = 1
	FlowDataStatusOver   = 2
	FlowDataStatusBack   = 3
	FlowDataStatusGift   = 4
	FlowDataStatusSign   = 5
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
	PassStatus   int //1是已通过，0是未通过
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
	BankCard          *BankCard
	ChargeAmount      string          // 充值金额
	LoginTime         int64           `bson:"logintime"`
	MatchCount        int             `bson:"-"` // 参赛次数
	AwardTotal        string          `bson:"-"` // 累计获得奖金
	AwardAvailable    string          `bson:"-"` // 可提现奖金
	SportCenter       SportCenterData // 体总数据
	Remark            string          // 备注信息
}

type SportCenterData struct {
	BlueScore       float64 // 蓝分
	RedScore        float64 // 红分
	SilverScore     float64 // 银分
	GoldScore       float64 // 金分
	Level           string  // 等级称号
	Ranking         int     // 排名
	SyncTime        int64   // 与体总同步时间
	LastBlueScore   float64 // 同步前蓝分
	LastRedScore    float64 // 同步前红分
	LastSilverScore float64 // 同步前银分
	LastGoldScore   float64 // 同步前金分
	LastLevel       string  // 同步前等级称号
	LastRanking     int     // 同步前排名
	WalletStatus    int     // 钱包状态 0锁定,1正常
}

type BankCard struct {
	Userid        int
	BankName      string
	BankCardNo    string
	Province      string
	City          string
	OpeningBank   string
	OpeningBankNo string
}

// UserMatchReview 用户后台的赛事列表总览
type UserMatchReview struct {
	UID            int
	AccountID      int
	MatchID        string
	MatchType      string
	MatchName      string
	MatchTotal     int
	MatchWins      int
	MatchFails     int
	AverageBatting int
	Coupon         int64
	AwardMoney     int64
	PersonalProfit int64
}

// ItemLog 物品日志
type ItemLog struct {
	UID        int    `bson:"uid"`
	Item       string `bson:"item"`       // 物品名称
	Amount     int64  `bson:"amount"`     // 物品数量
	Way        string `bson:"way"`        // 增加物品的方式
	CreateTime int64  `bson:"createtime"` // 创建时间
	Before     int64  `bson:"before"`     // 操作前余额
	After      int64  `bson:"after"`      // 操作后余额
	OptType    int    `bson:"opttype"`    // 操作类型
	MatchID    string `bson:"matchid"`    // 赛事id
	ShowAmount string // 显示数量
	ShowBefore string // 显示数量
	ShowAfter  string // 显示数量
}

type OfflinePayment struct {
	Nickname   string  `json:"nickname"`
	Accountid  int     `json:"accountid"`
	ActionType int     `json:"actiontype"` //0，点券 1，税后奖金
	BeforFee   float64 `json:"beforfee"`
	ChangeFee  float64 `json:"changefee"`
	AfterFee   float64 `json:"afterfee"`
	Createdat  int64   `json:"createdat"`
	Operator   string  `json:"operator"`
	Desc       string  `json:"desc"`
}

type OfflinePaymentCol struct {
	ID         int     `bson:"_id"`
	Nickname   string  `json:"nickname"`
	Accountid  int     `json:"accountid"`
	ActionType int     `json:"actiontype"` //0，点券 1，税后奖金
	BeforFee   float64 `json:"beforfee"`
	ChangeFee  float64 `json:"changefee"`
	AfterFee   float64 `json:"afterfee"`
	Createdat  int64   `json:"createdat"`
	Operator   string  `json:"operator"`
	Desc       string  `json:"desc"`
}

type DataCount struct {
	Count int
}

const (
	PayStatusAction = iota
	PayStatusSuccess
	PayStatusFail
)

type EdyOrder struct {
	ID             int `bson:"_id",json:"-"`
	Accountid      int
	TradeNo        string
	TradeNoReceive string
	Merchant       int //商户
	Status         bool
	Fee            int64
	Createdat      int64
	PayStatus      int //0表示支付中， 1表示支付成功， 2表示支付失败
	GoodsType      int //商品类型。1表示点券，2表示碎片
	Amount         int //商品数量
}

type RobotMatchNum struct {
	ID          int `bson:"_id"`
	MatchID     string
	MatchType   string
	MatchName   string
	PerMaxNum   int
	Total       int
	JoinNum     int
	Desc        string
	Status      int
	RobotStatus int
}

type KnapsackProp struct {
	ID        int `bson:"_id"`
	Accountid int
	PropID    int
	Name      string
	Num       int
	IsAdd     bool
	IsUse     bool
	Expiredat int64
	Desc      string
	Createdat int64
}

// MatchRecord 记录一局比赛所有玩家的手牌，输赢信息等
type MatchRecord struct {
	RoundCount int    // 第几局
	CardCount  int    // 第几副牌
	RoomCount  int    // 房间编号
	UID        int    // 用户id
	Identity   int    //0 防守方 1 进攻方
	Name       string // 玩家姓名
	HandCards  []int  //手牌
	ThreeCards []int  //底牌
	Event      int    //0:失败 1:胜利
	Score      int64  //得分
	Multiples  string //倍数
}

type MatchAwardRecord struct {
	MatchName    string
	AwardContent string
	ID           int `bson:"_id"`
	Userid       int
	Accountid    int
	MatchType    string
	MatchID      string
	CreatedAt    int64
	Realname     string
	Desc         string
}

// WhiteListConfig 白名单配置
type WhiteListConfig struct {
	Config      string `bson:"config"`
	WhiteSwitch bool   `bson:"whiteswitch"`
	WhiteList   []int  `bson:"whitelist"`
}

// RestartConfig 服务器重启配置
type RestartConfig struct {
	Config         string `bson:"config"`
	ID             string `bson:"id"`
	TipsTime       int64  `bson:"tipstime"`
	RestartTime    int64  `bson:"restarttime"`
	EndTime        int64  `bson:"endtime"`
	RestartTitle   string `bson:"restarttitle"`
	RestartType    string `bson:"restarttype"`
	Status         int    `bson:"status"`
	RestartContent string `bson:"restartcontent"`
	CreateTime     int64  `bson:"createtime"`
}

// 服务器重启更新状态
const (
	RestartStatusWait = iota + 1
	RestartStatusIng
	RestartStatusFinish
)

const (
	MerchantSportCentralAthketicAssociation = 1
)

var MerchantIDs = []int{MerchantSportCentralAthketicAssociation}

type MerchantPayBranch struct {
	ID           int `bson:"_id"`
	MerchantNo   string
	MerchantName string
	MerchantID   int
	PayBranch    []int
}

var MerchantPay = []int{MerchantSportCentralAthketicAssociation}

const (
	DownStatus = 0
	UpStatus   = 1
)

const (
	PayBranchWX  = 1
	PayBranchAli = 2
	PayBranchIOS = 3
)

type ShopMerchant struct {
	ID             int    `bson:"_id"`
	MerchantType   int    //商户类型。1是体总
	MerchantNo     string //商户编号
	PayMin         int    //支付最低值，百分制
	PayMax         int    //支付最高值，百分制
	PublicKey      string //公钥
	PrivateKey     string //私钥
	Order          int    //次序
	UpPayBranchs   []int  //上架支付类型，1是微信，2是支付宝，3是IOS
	DownPayBranchs []int  //下架支付类型，1是微信，2是支付宝，3是IOS
	UpDownStatus   int    //上下架状态。0是下架，1是上架
	UpdatedAt      int    //更新时间戳
	CreatedAt      int
	DeletedAt      int
}

type PayAccount struct {
	ID         int    `bson:"_id"`
	MerchantID int    //商户唯一标识
	PayBranch  int    //支付渠道标识
	Order      int    //次序
	Account    string //账户
	UpdatedAt  int    //更新时间戳
	CreatedAt  int    //更新时间戳
	DeletedAt  int    //删除时间戳
}

type GoodsType struct {
	ID         int    `bson:"_id"` //唯一标识
	MerchantID int    //商户唯一标识
	TypeName   string //商品名称
	ImgUrl     string //商品图标
	Order      int    //次序
	UpdatedAt  int    //更新时间戳
	CreatedAt  int    //创建时间戳
	DeletedAt  int    //删除时间戳
}

const (
	TakenTypeRMB = 1
)

type Goods struct {
	ID          int    `bson:"_id"`
	GoodsTypeID int    //商品类型唯一标识
	TakenType   int    //花费类型。1是RMB
	Price       int    //花费数量（价格，百分制）
	PropType    int    //道具类型。1是点券
	GetAmount   int    //获得数量
	GiftAmount  int    //赠送数量
	Expire      int    //过期时间，单位秒，-1为永久
	ImgUrl      string //商品图标
	Order       int    //次序
	UpdatedAt   int    //更新时间戳
	CreatedAt   int    //创建时间戳
	DeletedAt   int    //删除时间戳
}

type FeedBack struct {
	ID        int `bson:"_id"`
	AccountID int
	Title     string
	Content   string
	PhoneNum  string //联系方式
	Nickname  string //昵称

	MailType        int    //邮箱邮件类型
	MailServiceType int    //0是系统邮件，1是赛事邮件，2是活动邮件
	ReplyTitle      string //回复标题
	AwardType       int    //0是未选择，10002是报名券，10003是报名券碎片
	AwardNum        int    //奖励数量
	MailContent     string //邮箱内容
	ReadStatus      bool   //false是未查看，true是已查看
	ReplyStatus     bool   //false是未回复，true是已回复

	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64
}

//sundries const
const (
	PropTypeCoupon     = 1 //物件类型，点券
	PropTypeAward      = 2 //物件类型，点券
	PropTypeCouponFrag = 3 //物件类型，点券
	PropTypeRedScore   = 4 //物件类型，点券
)

var PropID2Type = map[int]int{
	10001: PropTypeAward,
	20001: PropTypeCoupon,
	20002: PropTypeCouponFrag,
	30001: PropTypeRedScore,
}

//prop_base_conf 道具基本配置
type PropBaseConfig struct {
	ID       int    `bson:"_id"` //唯一标识
	PropType int    //道具类型, 1是点券，2是奖金，3点券碎片
	PropID   int    //道具id
	Name     string //名称
	ImgUrl   string //图片url
	Operator string //操作人

	CreatedAt int //创建时间戳
	UpdatedAt int //更新时间戳
	DeletedAt int //删除时间戳，0表示没有删除
}

// OneDailyWelfareConfig 单条每日福利配置
type OneDailyWelfareConfig struct {
	WelfareType int             `bson:"WelfareType"` // 福利类型
	AwardList   []OneItemConfig `bson:"AwardList"`
}

// OneItemConfig 单条配置
type OneItemConfig struct {
	Item         int   `bson:"Item"`         // 物品ID
	AwardAmount  int   `bson:"AwardAmount"`  // 奖励数量
	TargetAmount int64 `bson:"TargetAmount"` // 达成条件
}

// DDZGameRecord 赛事记录
type DDZGameRecord struct {
	UserId    int    //用户ID
	MatchId   string //赛事ID
	MatchType string //赛事类型
	Desc      string //赛事
	Level     int    //名次
	Award     string //奖励
	Count     int    //完成局数
	Total     int64  //总得分
	Last      int64  //尾副得分
	Wins      int    //获胜次数
	Period    int64  //累计时长
	// Result    []Result //牌局详细
	CreateDat int64 //时间
	Status    int   // 战绩发奖状态
}

type Annex struct {
	PropType int //1是点券，2是奖金，3点券碎片，4是红分
	Num      float64
}

const (
	MailcontrolStatusNotSend     = 0
	MailcontrolStatusAlreadySend = 1
)

type Mailcontrol struct {
	ID              int     `bson:"_id"` //唯一标识
	TargetID        []int   //目标用户
	MailServiceType int     //0是系统邮件，1是赛事邮件，2是活动邮件
	Title           string  //标题
	Content         string  //内容
	Annexes         []Annex //附件
	Expire          int     //过期时间（单位：分钟）
	Status          int     //状态，0是未发送，1是已发送
	Operator        string  //操作人

	CreatedAt int //创建时间戳，对应添加时间
	UpdatedAt int //更新时间戳，对应发送时间
	DeletedAt int //删除时间戳
}
