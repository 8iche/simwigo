package logger

import (
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"net"
	"net/http"
)

type Config struct {
	API       string
	IP        string
	Port      string
	Temp      string
	Cache     string
	WhiteList []string
	Routes    [][]string
}

func (config *Config) PrintRoutesInfo() {
	if Debug {
		serverTable := config.serverInfo()
		routesTable := config.routesInfo()
		interfacesTable := config.netIfacesInfo()

		panels, _ := pterm.DefaultPanel.WithPanels(pterm.Panels{
			{{Data: serverTable}, {Data: routesTable}},
			{{Data: interfacesTable}},
		}).Srender()

		pterm.DefaultBox.WithTitle("Server Info").WithTitleTopRight().Println(panels)
	}
}

func (config *Config) serverInfo() string {
	info := [][]string{
		{"Options", "Values"},
		{"API", config.API},
		{"IP", config.IP},
		{"PORT", string(config.Port)},
		{"Temp", config.Temp},
		{"Cache", config.Cache},
	}
	table, _ := pterm.DefaultTable.WithHasHeader().WithData(info).Srender()
	return pterm.DefaultBox.WithTitle("Config").WithTitleTopLeft().Sprint(table)
}

func (config *Config) routesInfo() string {
	table, _ := pterm.DefaultTable.WithHasHeader().WithData(config.Routes).Srender()
	return pterm.DefaultBox.WithTitle("Routes").WithTitleTopLeft().Sprint(table)
}

func (config *Config) netIfacesInfo() string {

	var (
		ipv4, ipv6 [][]string
	)
	response, _ := http.Get("http://ifconfig.so")
	publicIp, _ := io.ReadAll(response.Body)
	response.Body.Close()
	ipv4 = append(ipv4, []string{"Interfaces", "IP", "Activated"}, []string{"Public", string(publicIp), "true"})
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range ifaces {

		addrs, _ := i.Addrs()
		for _, addr := range addrs {

			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() == nil {
					if v.IP.String() == config.IP || config.IP == "0.0.0.0" {
						ipv4 = append(ipv4, []string{i.Name, v.IP.String(), "true"})
					} else {
						ipv4 = append(ipv4, []string{i.Name, v.IP.String(), "false"})
					}

				} else {
					if v.IP.String() == config.IP || config.IP == "0.0.0.0" {
						ipv6 = append(ipv6, []string{i.Name, v.IP.String(), "true"})
					} else {
						ipv6 = append(ipv6, []string{i.Name, v.IP.String(), "false"})
					}
				}
			}
		}
	}

	ipv4 = append(ipv4, ipv6...)

	ifacesTable, _ := pterm.DefaultTable.WithHasHeader().WithData(ipv4).Srender()
	return pterm.DefaultBox.WithTitle("Network").WithTitleTopLeft().Sprint(ifacesTable)
}
