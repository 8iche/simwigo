package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": 200, "message": "pong"})
}
