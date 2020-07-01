package main

import (
	"bs/config"
	"bs/route"
)

func main() {
	// Engin
	server := route.GetServer()
	// 指定地址和端口号
	server.Run(config.GetConfig().Port)
}
