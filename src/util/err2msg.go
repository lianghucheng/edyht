package util

const (
	Success      = 1000
	Fail         = 1001
	FormatFail   = 1002
	TaxFeeLack   = 1003
	UserNotExist = 1004
)

var ErrMsg = map[int]string{
	Success:      "成功",
	Fail:         "失败",
	FormatFail:   "格式错误",
	TaxFeeLack:   "所剩税后奖金不足",
	UserNotExist: "该用户不存在",
}
