package util

const (
	Success             = 1000
	Fail                = 1001
	FormatFail          = 1002
	TaxFeeLack          = 1003
	UserNotExist        = 1004
	MatchNotExist       = 10005
	MatchRobotConfExist = 10006
	RobotNotBan         = 10007
)

var ErrMsg = map[int]string{
	Success:             "成功",
	Fail:                "失败",
	FormatFail:          "格式错误",
	TaxFeeLack:          "所剩税后奖金不足",
	UserNotExist:        "该用户不存在",
	MatchNotExist:       "该赛事不存在",
	MatchRobotConfExist: "该赛事机器人配置已存在",
	RobotNotBan:         "该赛事机器人没有金禁用",
}
