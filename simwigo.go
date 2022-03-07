package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"net/netip"
	"os"
	"os/user"
	"runtime"
	"simwigo/logger"
	"simwigo/router"
	"strings"
)

const goos = runtime.GOOS

var version string

var opts struct {
	AutoTLS   bool   `long:"autotls" description:"Enable TLS using Let's Encrypt"`
	SelfTLS   bool   `long:"self-cert" description:"Auto generates self-signed TLS certificates (dangerous)"`
	IsRSA     bool   `long:"rsa" description:"Force the use of RSA to generate TLS certificates (default=EDCSA)"`
	EnableAPI bool   `short:"a" long:"api" description:"Enable authorization api key"`
	Debug     bool   `long:"no-debug" description:"Disable Debug mode: no output"`
	Light     bool   `long:"light" description:"Light mode: disable colored mode"`
	Port      int    `short:"p" long:"port" description:"Port to listen" default:"8000"`
	IP        string `short:"i" description:"Interface to listen" default:"0.0.0.0"`
	Cert      string `long:"cert" description:"Path to TLS certificate file"`
	Key       string `long:"key" description:"Path to TLS KEY file"`
	TempDir   string `long:"temp" description:"Path to a temporary directory"`
	CacheDir  string `long:"cache" description:"Path to a cache directory"`
	DirList   string `long:"dir" description:"Enable directory listing (dangerous)"`
	Domain    string `long:"domain" short:"d" description:"Domain to specify for Let's Encrypt"`
	WhiteList string `long:"white-list" short:"l" description:"White list IP separated by comma (,)"`
	Version   bool   `long:"version" description:"Show version"`
}

func main() {


	parser := flags.NewParser(&opts, flags.Default)

	_, err := parser.Parse()

	if err != nil {
		os.Exit(0)
	}

	if opts.Version {
		logger.Logfln("simwigo v%s", version)
		os.Exit(0)
	}
	if opts.Debug {
		logger.DisableColor()
		logger.Debug = false
	}

	if opts.Light {
		logger.DisableColor()
	}

	var whiteList []string
	if opts.WhiteList != "" {
		whiteList = strings.Split(opts.WhiteList, ",")
	}

	checkConfig(whiteList)

	server := router.New(opts.IP, opts.Port, opts.Domain, opts.TempDir, opts.CacheDir, opts.EnableAPI, opts.SelfTLS, opts.DirList, opts.AutoTLS, opts.IsRSA, opts.Cert, opts.Key, whiteList)

	err = createDir(opts.TempDir, opts.CacheDir)

	if err != nil {
		logger.Fatalln(err)
	}

	if !opts.Debug {
		logger.Banner()
	}

	if !opts.EnableAPI {
		logger.Warning.Logln("API authentication disabled. It is recommended to enable it!!!")
	}

	if opts.DirList != "" {
		logger.Warning.Logln("Directory listing enabled (dangerous)")
	}

	if opts.AutoTLS {
		if opts.Domain == "" {
			logger.Fatalln("Domain should not be empty when autotls is enabled. Use -d option to specify a domain")
		}
		server.Port = 443
	} else if opts.SelfTLS {
		logger.Warning.Logln("Self-signed certificate generated (dangerous)")
		logger.Logln()
	} else {
		logger.Warning.Logln("TLS disabled (dangerous)")
		logger.Logln()
	}

	if err = server.Run(); err != nil {
		logger.Fatalln(err)
	}
}

func checkConfig(whiteList []string) {
	if opts.IP != "0.0.0.0" {
		if _, err := netip.ParseAddr(opts.IP); err != nil {
			logger.Error.LogFatalln(err)
		}
	}

	if opts.Port < 0 || opts.Port > 65535 {
		logger.Error.LogFatalln("Invalid port")
	}

	if opts.TempDir == "" {
		switch goos {
		case "windows":
			u, err := user.Current()
			if err != nil {
				logger.Fatalln(err.Error())
			}
			opts.TempDir = fmt.Sprintf("%s\\AppData\\Local\\Temp\\simwigo\\", u.HomeDir)
		case "linux":
			opts.TempDir = "/tmp/simwigo/"
		default:
			opts.TempDir = "simwigo"
		}
	} else {
		opts.TempDir = checkDir(opts.TempDir)
	}

	if opts.CacheDir == "" {
		switch goos {
		case "windows":
			u, err := user.Current()
			if err != nil {
				logger.Fatalln(err.Error())
			}
			opts.CacheDir = fmt.Sprintf("%s\\AppData\\Local\\Temp\\.simwigo\\", u.HomeDir)
		case "linux":
			opts.CacheDir = "/tmp/.simwigo/"
		default:
			opts.CacheDir = ".simwigo"
		}
	} else {
		opts.CacheDir = checkDir(opts.CacheDir)
	}

	for _, v := range whiteList {
		if _, err := netip.ParseAddr(v); err != nil {
			logger.Error.LogFatalln(err)
		}
	}
}

func createDir(dir ...string) error {
	for _, d := range dir {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			err := os.Mkdir(d, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkDir(dir string) string {
	if lastChar := dir[len(dir)-1:]; lastChar != "/" || lastChar != "\\" {
		if goos == "windows" {
			dir += "\\"
		} else {
			dir += "/"
		}
	}
	return dir
}
