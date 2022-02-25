package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func (server *Server) uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Internal error"})
		return
	}
	filename := filepath.Base(file.Filename)

	filename = filepath.Clean(filename)

	dst := server.TempDir + filename

	err = c.SaveUploadedFile(file, dst)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "message": "OK"})
}
