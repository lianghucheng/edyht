package util

const (
	Success    = 1000
	Fail       = 1001
	FormatFail = 1002
)

var ErrMsg = map[int]string{
	1000: "成功",
	1001: "失败",
	1002: "格式错误",
}
