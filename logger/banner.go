package logger

import (
	"github.com/pterm/pterm"
	"time"
)

func Banner() {
	if lightMode {
		pterm.Println("Starting at: " + time.Now().Format("02 Jan 2006 - 15:04:05 MST"))
		pterm.Println()
		return
	}
	pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgLightWhite).Sprint("\\ | /"))
	ptermLogo, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Simwi", pterm.NewStyle(pterm.FgLightRed)),
		pterm.NewLettersFromStringWithStyle("Go", pterm.NewStyle(pterm.FgLightBlue))).
		Srender()

	pterm.DefaultCenter.Println(ptermLogo)
	pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgLightWhite).Sprint("/ | \\"))
	pterm.DefaultCenter.Println(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Sprint("Simple Web Server"))

	pterm.Println(pterm.LightCyan("Starting at: ") + pterm.Green(time.Now().Format("02 Jan 2006 - 15:04:05 MST")))

	pterm.Println()

}
