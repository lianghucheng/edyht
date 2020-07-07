package param

type FlowDataHistoryReq struct {
	Start     int64    `json:"start"`     //开始时间
	End       int64    `json:"end"`       //结束时间
	Per       int      `json:"per"`       //页数
	Page      int      `json:"page"`      //页码
	Condition []string `json:"condition"` //查询条件
}

type FlowData struct {
	ID           int     `json:"id"`            //唯一id
	Accountid    int     `json:"accountid"`     //所属用户id
	ChangeAmount float64 `json:"change_amount"` //变动金额
	FlowType     int     `json:"flow_type"`     //流水类型
	MatchID      string  `json:"match_id"`      //比赛ID
	Status       int     `json:"status"`        //状态
	CreatedAt    int64   `json:"created_at"`    //日期
	Realname     string  `json:"realname"`      //实名昵称
	TakenFee     float64 `json:"taken_fee"`     //已提现金额
	AtferTaxFee  float64 `json:"atfer_tax_fee"` //税后奖金
	Desc         string  `json:"desc"`          //备注说明
}
type FlowDataHistoryResp struct {
	Per       int         `json:"per"`        //当前页数
	Page      int         `json:"page"`       //当前页码
	Total     int         `json:"total"`      //当前总数
	FlowDatas *[]FlowData `json:"flow_datas"` //数据
}
