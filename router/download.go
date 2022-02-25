package router

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func (server *Server) downloadFile(c *gin.Context) {
	file := c.Param("filename")

	filename := filepath.Base(file)
	dst := server.TempDir + filename

	f, err := os.Open(dst)
	defer f.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(f)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal error"})
		return
	}

	contentLength := int64(buf.Len())
	contentType := "application/octet-stream"

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", file),
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, buf, extraHeaders)
}
