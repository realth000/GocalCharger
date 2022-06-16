package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiRoot(c *gin.Context) {
	c.String(http.StatusOK, "Hello gin")
}

func ApiConfigs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"configs": getConfigs(),
	})
}

func ApiConfigsLoad(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"configs": reloadConfigs(),
	})
}
