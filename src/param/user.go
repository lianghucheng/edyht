package param

import (
	"bs/param/base"
	"gopkg.in/mgo.v2/bson"
)

type UserInsertReq struct {
	Account  string //联系方式
	Password string //密码
	Role     int    //角色，1为boss，2为运营，3为客服
	Owner 	 string //角色拥有者
}

type UserDeleteReq struct {
	base.OID //记录唯一标识
}

type UserReadReq struct {
	base.OID //记录唯一标识
}

type UserReadResp struct {
	User
}

type UserListReq struct {
	base.DivPage
	base.Condition
}

func (ctx *UserListReq) GetPipeline() []bson.M {
	pipeline := []bson.M{{"$match": bson.M{"deletedat": 0, "role" : 0}}}
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	return pipeline
}

func (ctx *UserListReq) GetDataPipeline() []bson.M {
	pipeline := []bson.M{}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"deletedat": 0, "role" : 0}})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdat": -1}})
	pipeline = append(pipeline, base.GetPipeline(ctx.Condition)...)
	pipeline = append(pipeline, ctx.DivPage.GetPipeline()...)
	return pipeline
}

type User struct {
	ID       int    `bson:"_id"`//唯一标识
	Account  string //联系方式
	Password string //密码
	Role     int    //角色，1为boss，2为运营，3为客服
	Power    int	//权限
	Owner 	 string //角色拥有者
	LastIp   string //登陆ip
	Operator string //操作人
	Status   int 	//状态，0为开启，1为关闭

	UpdatedAt int
}

type UserListResp struct {
	Page      int
	Per       int
	Total     int
	Users *[]User
}

type UserUpdateReq struct {
	base.OID        //记录唯一标识
	Account  string //联系方式
	Password string //密码
	Role     int    //角色，1为boss，2为运营，3为客服
	Power    int	//权限
	Owner 	 string //角色拥有者
	Status   int 	//状态，0为开启，1为关闭
}
