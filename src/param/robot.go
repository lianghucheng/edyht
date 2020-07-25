package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type RobotMatchNumReq struct {
	base.DivPage
	base.Condition
}

func (ctx *RobotMatchNumReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type MatchRobotNum struct {
	MatchID     string `json:"matchid"`     //赛事id
	MatchName   string `json:"matchname"`   //赛事名
	PerMaxNum   int    `json:"permaxnum"`   //最大机器人数
	Total       int    `json:"total"`       //机器人总数
	JoinNum     int    `json:"joinnum"`     //参加比赛数量
	Desc        string `json:"desc"`        //温馨提示
	RobotStatus int    `json:"robotstatus"` //机器人状态
}

type RobotMatchNumResp struct {
	Page           int              `json:"page"`
	Per            int              `json:"per"`
	Total          int              `json:"total"`
	MatchRobotNums *[]MatchRobotNum `json:"match_robot_nums"`
}

type RobotMatchReq struct {
	base.Condition
	base.DivPage
}

type RobotMatch struct {
	MatchType    string `json:"match_type"`     //赛事类型
	MatchNum     int    `json:"match_num"`      //赛事数量
	RobotTotal   int    `json:"robot_total"`    //机器人总数
	RobotJoinNum int    `json:"robot_join_num"` //参赛机器人总数
}

type RobotMatchResp struct {
	Page        int           `json:"page"`
	Per         int           `json:"per"`
	Total       int           `json:"total"`
	RobotMatchs *[]RobotMatch `json:"match_robots"`
	MatchTypes  []string      `json:"match_types"`
}

type RobotSaveReq struct {
	MatchID     string `json:"match_id"`
	RobotNum    int    `json:"robot_num"`
	PerMatchNum int    `json:"per_match_num"`
	Desc        string `json:"desc"`
	Type        int    `json:"type"`
}

type RobotDelReq struct {
	MatchID string `json:"match_id"`
}

type RobotStopReq struct {
	MatchID string `json:"match_id"`
}

type RobotStopAllReq struct {
	MatchTypes []string `json:"match_types"`
}

type RobotStartReq struct {
	MatchID string `json:"match_id"`
}

type RobotStartAllReq struct {
	MatchTypes []string `json:"match_types"`
}
