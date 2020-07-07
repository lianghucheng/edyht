package param

type FlowDataExportReq struct {
	Start     int64    `json:"start"`     //开始时间
	End       int64    `json:"end"`       //结束时间
	Condition []string `json:"condition"` //查询条件
}

type FlowExports struct {
	Accountid    int     `json:"accountid"`      //用户id
	PhoneNum     string  `json:"phone_num"`      //电话号码
	Realname     string  `json:"realname"`       //真实名字
	BankCardNo   string  `json:"bank_card_no"`   //银行卡号
	OpenBankName string  `json:"open_bank_name"` //开户行名称
	ChangeAmount float64 `json:"change_amount"`  //变动金额
}
type FlowDataExportResp struct {
	FlowExports *[]FlowExports `json:"flow_exports"` //数据
}
