package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)
func (server *Server) deleteShare(c *gin.Context) {
	link := c.Param("link")

	filename, err := server.Share.GetShareFromUUID(link)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Internal error"})
		return
	}

	server.Share.DeleteShare(filename)

	c.JSON(http.StatusOK, gin.H{"status": 200, "message": "OK"})
}
