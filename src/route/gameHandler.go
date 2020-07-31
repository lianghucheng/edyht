package route

// 处理一下游戏端的请求

import (
	"bs/db"
	"bs/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func getGameVersion(c *gin.Context) {
	code := util.OK
	desc := "OK"
	var version, url string
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"desc":    desc,
			"version": version,
			"url":     url,
		})
	}()
	version, url = db.GetGameVersion()
}

// 获取公告
func getNotice(c *gin.Context) {
	code := util.OK
	desc := "OK"
	notice := &util.RestartConfig{}
	defer func() {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"desc": desc,
			"info": notice,
		})
	}()
	one, err := db.GetLastestRestart()
	if err != nil || one.Status > util.RestartStatusFinish {
		code = util.Retry
		desc = "请求出错,请重试!"
		return
	}
	if one.Status <= 0 || one.Status == util.RestartStatusFinish || one.TipsTime > time.Now().Unix() {
		code = util.OK
		desc = "暂无公告!"
		notice = nil
		return
	}
	notice = &one
}
