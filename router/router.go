package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"simwigo/internal/logger"
)

func (server *Server) Run() error {

	gin.SetMode(gin.ReleaseMode)

	server.Router = gin.New()
	server.Router.SetTrustedProxies(nil)
	server.Router.Use(mainHandler(server.WhiteList))
	server.loggerFormat()
	server.setupRoutes()
	server.Info()

	if err := server.listenAndServe(); err != nil {
		return err
	}
	return nil
}

func (server *Server) setupRoutes() {

	root := server.Router.Group("/")
	{
		if server.DirList != "" {
			root.StaticFS("index", http.Dir(server.DirList))
		}

		root.GET("ping", ping)
		root.Any("print")
	}

	authRoutes := server.Router.Group("/file").Use(validateAPIKey(server.API))
	{
		server.Router.MaxMultipartMemory = 8 << 20

		authRoutes.POST("upload", server.uploadFile)
		authRoutes.GET("download/:filename", server.downloadFile)
		authRoutes.GET("list", server.listFiles)
	}

	shareRoutes := server.Router.Group("/share")
	{
		shareRoutes.GET(":link", server.downloadLink)
		shareRoutes.POST("upload", server.uploadAndShare).Use(validateAPIKey(server.API))
		shareRoutes.GET(":link/delete", server.deleteShare).Use(validateAPIKey(server.API))
	}
}

func (server *Server) listenAndServe() (err error) {
	host := fmt.Sprintf("%s:%d", server.IP, server.Port)

	if server.TLS.AutoTLS {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(server.Domain),
			Cache:      autocert.DirCache(server.CacheDir),
		}
		if err = autotls.RunWithManager(server.Router, &m); err != nil {
			return err
		}
	} else if server.TLS.SelfTLS {
		if err = server.Router.RunTLS(host, server.TLS.Cert, server.TLS.Key); err != nil {
			return err
		}
	} else {
		if err = server.Router.Run(host); err != nil {
			return err
		}
	}
	return errors.New("could not run server")
}

func (server *Server) loggerFormat() {
	server.Router.Use(gin.LoggerWithFormatter(logger.LogFormatter))
}

func (server *Server) Info() {

	auth := server.API.IsEnabled()

	if auth {
		logger.Info.Logfln("API Header: X-API-Key: %s", server.API.Get())
	}

	logger.Logln()

	authStr := fmt.Sprintf("%t", auth)

	config := &logger.Config{
		API:       server.API.Get(),
		IP:        server.IP,
		Port:      fmt.Sprintf("%d", server.Port),
		Temp:      server.TempDir,
		Cache:     server.CacheDir,
		WhiteList: server.WhiteList,
		Routes: [][]string{
			{"Method", "Path", "Auth"},
			{"GET", "/ping", "false"},
			{"ANY", "/print", "false"},
			{"POST", "/file/upload", authStr},
			{"GET", "/file/download/:filename", authStr},
			{"GET", "/file/list", authStr},
			{"POST", "/share/upload", authStr},
			{"GET", "/share/:UUID", "false"},
			{"GET", "/share/:UUID/delete", authStr},
		},
	}

	if server.DirList != "" {
		config.Routes = append(config.Routes, []string{"GET", "/index/", "false"})
	}

	config.PrintRoutesInfo()

	logger.Logln()

	if server.TLS.AutoTLS || server.TLS.SelfTLS {
		logger.Success.Logfln("Listening and serving HTTPS on %s:%d", server.IP, server.Port)
	} else {
		logger.Success.Logfln("Listening and serving HTTP on %s:%d", server.IP, server.Port)
	}

}
