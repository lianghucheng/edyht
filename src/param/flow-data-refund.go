package param

type FlowDataRefundReq struct {
	ID   int    `json:"id"`   //流水id
	Desc string `json:"desc"` //备注描述
}

type FlowDataRefundsReq struct {
	Ids  []int  `json:"ids"`  //选中的流水id
	Desc string `json:"desc"` //备注描述
}
