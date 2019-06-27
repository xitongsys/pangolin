package login

import (
	"fmt"
	"strings"

	"config"
)

type Dhcp struct {
	Cfg *config.Config
	UsedPorts map[string]bool
}

func NewDhcp(cfg *config.Config) *Dhcp {
	return &Dhcp{
		Cfg: cfg,
		UsedPorts: map[string]bool{},
	}
}

func (dhcp *Dhcp) GetNewAddr() (string, error){
	ip := dhcp.Cfg.TunIp
	port := ""
	for p := 2; p < 65535; p++ {
		ps := fmt.Sprint(p)
		if _, ok := dhcp.UsedPorts[ps]; !ok {
			port = ps
			dhcp.UsedPorts[ps] = true
			break
		}
	}

	if port == "" {
		return "", fmt.Errorf("No available ip")
	}
	return ip + ":" + port, nil
}

func (dhcp *Dhcp) ReleaseAddr(clientAddr string){
	port := strings.Split(clientAddr, ":")[1]
	delete(dhcp.UsedPorts, port)
}