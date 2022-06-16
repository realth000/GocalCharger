package main

import (
	"github.com/gin-gonic/gin"
	webApi "gocalcharger/api/web"
)

func RunServer() {
	r := gin.Default()
	r.GET("/", webApi.ApiRoot)
	r.StaticFile("/settings", "./public/page/main.html")
	r.GET("/api/configs", webApi.ApiConfigs)
	r.Run(":3090")
}

func main() {
	RunServer()
}
