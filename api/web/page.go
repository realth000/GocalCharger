package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunServer() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello gin")
	})
	r.StaticFile("/settings", "./public/page/main.html")
	r.GET("/api/configs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"configs": getConfigs(),
		})
	})
	r.Run(":3090")
}
