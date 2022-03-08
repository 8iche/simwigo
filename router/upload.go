package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"simwigo/internal/logger"
	"strconv"
)

func (server *Server) uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Internal error"})
		return
	}

	filename := filepath.Base(file.Filename)

	dst := server.TempDir + filename

	err = c.SaveUploadedFile(file, dst)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 200, "message": "OK"})
}

func (server *Server) uploadAndShare(c *gin.Context) {
	var scheme string

	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Internal error"})
		return
	}

	count, err := strconv.Atoi(c.DefaultQuery("count", "-1"))
	if err != nil {
		count = -1
	}

	filename := filepath.Base(file.Filename)

	dst := server.TempDir + filename

	err = server.Share.SaveFile(file, filename, dst, count)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Internal error"})
		return
	}

	if server.TLS.AutoTLS || server.TLS.SelfTLS {
		scheme = "https"
	} else {
		scheme = "http"
	}

	link := fmt.Sprintf("%s://%s/share/%s", scheme, c.Request.Host, server.Share[filename].UUID)

	c.JSON(http.StatusOK, gin.H{"status": 200, "link": link})

	logger.Info.Logfln("share link created: %s", link)

}
