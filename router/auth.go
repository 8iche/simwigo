package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simwigo/internal/api"
)

func validateAPIKey(api *api.API) gin.HandlerFunc {
	return func(c *gin.Context) {
		if api.IsEnabled() {
			APIKey := c.Request.Header.Get("X-API-Key")
			if APIKey != api.Get() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Permission denied"})
				return
			}
			c.Next()
		}
		c.Next()

	}
}
