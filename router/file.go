package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	_"strconv"
)

func (server *Server) listFiles(c *gin.Context) {
	files, err := getTree(server.TempDir)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "internal server error")
	}
	c.JSON(200, files)
}

//Nécessite une amélioration
func getTree(p string) ([]string, error) {
	var files []string
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, filepath.Base(path))
		return nil
	})

	return files, err
}
