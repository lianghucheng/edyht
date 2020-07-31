package main

import (
	"bs/config"
	"bs/route"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	// Engin
	server := route.GetServer()
	// 指定地址和端口号
	server.Run(config.GetConfig().Port)
}
