package param

type FlowDataPaymentReq struct {
	ID   int    `json:"id"`   //流水id
	Desc string `json:"desc"` //备注描述
}

type FlowDataPaymentsReq struct {
	Ids  []int  `json:"ids"`  //选中的流水id
	Desc string `json:"desc"` //备注描述
}
