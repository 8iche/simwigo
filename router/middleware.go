package router

import (
	"github.com/gin-gonic/gin"
	"simwigo/internal/logger"
)

func mainHandler(whiteList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(whiteList) > 0 {
			if !contains(whiteList, c.ClientIP()) {
				logger.Warning.Logfln("Request from %s rejected", c.ClientIP())
				c.AbortWithStatus(400)
				return
			}
		}

		c.Next()

		if c.FullPath() == "/print" {
			req := formatHttpRequest(c.Request)
			logger.RequestPrinter(strHttpRequest(req))
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
