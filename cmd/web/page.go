package main

import (
	"github.com/gin-gonic/gin"
	webApi "gocalcharger/api/web"
	cc "gocalcharger/client/config"
	sc "gocalcharger/server/config"
	"log"
)

// test
const (
	serverConfig = `./tests/data/config/server.toml`
	clientConfig = `./tests/data/config/client.toml`
)

func loadConfig() {
	err := sc.LoadConfigFile(serverConfig)
	if err != nil {
		log.Fatalf("can not load server config:%v", err)
	}
	err = cc.LoadConfigFile(clientConfig)
	if err != nil {
		log.Fatalf("can not load client config:%v", err)
	}
}

func RunServer() {
	r := gin.Default()
	r.GET("/", webApi.ApiRoot)
	r.StaticFile("/settings", "./public/page/main.html")
	r.GET("/api/configs", webApi.ApiConfigs)
	r.Run(":3090")
}

func main() {
	loadConfig()
	RunServer()
}
