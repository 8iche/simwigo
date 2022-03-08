package router

import (
	"github.com/gin-gonic/gin"
	"simwigo/internal/api"
	"simwigo/internal/logger"
	"simwigo/internal/share"
	"simwigo/internal/tls"
)

type Server struct {
	Domain    string
	IP        string
	Port      int
	Router    *gin.Engine
	API       *api.API
	TempDir   string
	CacheDir  string
	DirList   string
	TLS       *TLS
	WhiteList []string
	Share     share.FS
}

type TLS struct {
	SelfTLS bool
	AutoTLS bool
	IsRSA   bool
	Key     string
	Cert    string
}

func New(ip string, port int, domain string, tempDir string, cacheDir string, enableAPI bool, enableTLS bool, dirList string, autoTLS bool, isRSA bool, cert string, key string, whiteList []string) *Server {

	server := new(Server)

	if enableAPI {
		server.API = api.New()
	} else {
		server.API = &api.API{}
	}

	server.IP = ip
	server.Port = port
	server.Domain = domain
	server.TLS = &TLS{
		SelfTLS: enableTLS,
		AutoTLS: autoTLS,
		IsRSA:   isRSA,
	}

	if !autoTLS && enableTLS && (cert == "" && key == "") {
		var err error
		server.TLS.Cert, server.TLS.Key, err = tls.GenerateTLSCertificate(ip, cacheDir, isRSA)
		if err != nil {
			logger.Fatalln(err)
		}
	}

	server.TempDir = tempDir
	server.CacheDir = cacheDir
	server.DirList = dirList
	server.TLS.IsRSA = isRSA
	server.WhiteList = whiteList

	server.Share = make(share.FS)

	return server
}
