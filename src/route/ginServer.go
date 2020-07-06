package route

import (
	"bs/config"
	"bs/db"
	"bs/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/szxby/tools/log"
)

var server *gin.Engine

func init() {
	server = gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	// server.Use(ipAuthMiddleWare())
	server.Use(tokenAuthMiddleWare())
	bind(server)
}

// GetServer return defalut server
func GetServer() *gin.Engine {
	return server
}

func ipAuthMiddleWare() gin.HandlerFunc {
	con := config.GetConfig()
	tag := false
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		for _, l := range con.IPList {
			if l == clientIP {
				tag = true
				break
			}
		}
		if !tag {
			log.Error("ivalid ip: %v", clientIP)
			c.Abort()
		}
	}
}

func tokenAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		token := c.GetHeader("token")
		path := c.Request.URL.Path
		if path == "/favicon.ico" {
			c.Abort()
			return
		}
		if path != "/login" && strings.Index(path, "/download/matchIcon/") == -1 {
			role := db.RedisGetToken(token)
			log.Debug("redis中的权限  %v", role)
			if role == -1 {
				log.Debug("ivalid token: %v", token)
				c.JSON(http.StatusOK, gin.H{
					"code": util.TokenExpire,
					"desc": "当前会话已过期，请重新登录！",
				})
				c.Abort()
				return
			}
			tag := checkRole(role, path)
			// 接收操作后，刷新token
			refreshToken(token, role)
			if !tag {
				log.Debug("ivalid role:%v go path:%v", token, path)
				c.JSON(http.StatusOK, gin.H{
					"code": util.Retry,
					"desc": "您没有该项的操作权限！",
				})
				c.Abort()
				return
			}
		}
	}
}
