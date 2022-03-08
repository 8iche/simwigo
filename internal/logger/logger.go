package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/pterm/pterm"
	"os"
)

var Debug = true
var lightMode = false

type Logger struct {
	pterm.PrefixPrinter
}

var (
	Error   = Logger{pterm.Error}
	Info    = Logger{pterm.Info}
	Warning = Logger{pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.WarningMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.Style{pterm.FgBlack, pterm.BgRed},
			Text:  "WARNING",
		}}}

	Success = Logger{pterm.Success}

	Example = Logger{pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.InfoPrefixStyle,
			Text:  "EXAMPLE",
		},
	}}

	Dbg = Logger{pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.DebugMessageStyle,
		Prefix: pterm.Prefix{
			Text:  " DEBUG ",
			Style: &pterm.ThemeDefault.DebugPrefixStyle,
		},
	}}
)

func DisableColor() {
	pterm.DisableStyling()
	lightMode = true
}

func (logger *Logger) Log(v ...interface{}) {
	if Debug {
		logger.Print(v...)
	}
}

func (logger *Logger) Logf(s string, v ...interface{}) {
	if Debug {
		logger.Printf(s, v...)
	}
}

func (logger *Logger) Logfln(s string, v ...interface{}) {
	if Debug {
		logger.Printfln(s, v...)
	}
}

func (logger *Logger) Logln(v ...interface{}) {
	if Debug {
		logger.Println(v...)
	}
}

func (logger *Logger) Slogln(v ...interface{}) string {
	if Debug {
		return logger.Sprintln(v...)
	}
	return ""
}

func (logger *Logger) Slogfln(s string, v ...interface{}) string {
	if Debug {
		return logger.Sprintfln(s, v...)
	}
	return ""
}

func (logger *Logger) LogFatalln(v ...interface{}) {
	if Debug {
		pterm.Error.Println(v...)
		os.Exit(1)
	}
}

func Logfln(s string, v ...interface{}) {
	if Debug {
		pterm.Printfln(s, v...)
	}

}

func Logln(v ...interface{}) {
	if Debug {
		pterm.Println(v...)
	}

}

func Fatalln(v ...interface{}) {
	pterm.Error.Println(v...)
	os.Exit(1)
}

func LogFormatter(param gin.LogFormatterParams) string {
	if Debug {
		if lightMode {
			logFormat := Info.Sprintf("%v | %3d | %-15s | %-7s | %s \n",
				param.TimeStamp.Format("02/01/2006 - 15:04:05"),
				param.StatusCode,
				param.ClientIP,
				param.Method,
				param.Path,
			)
			return logFormat
		}
		logFormat := Info.Sprintf("%v |%s %3d %s| %-15s |%s %-7s %s| %s \n",
			param.TimeStamp.Format("02/01/2006 - 15:04:05"),
			param.StatusCodeColor(), param.StatusCode, param.ResetColor(),
			param.ClientIP,
			param.MethodColor(), param.Method, param.ResetColor(),
			param.Path,
		)
		return logFormat
	}
	return ""
}

func RequestPrinter(s string) {
	pterm.DefaultBox.WithTitle("Request").WithTitleTopLeft().Println(s)
}
