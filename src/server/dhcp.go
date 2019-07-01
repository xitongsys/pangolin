 package server

import (
	"fmt"

	"config"
	"header"
)

type Dhcp struct {
	Cfg *config.Config
	Ip uint32
	Mask uint32
	UsedIps map[uint32]bool
}

func NewDhcp(cfg *config.Config) *Dhcp {
	ip, mask := header.ParseNet(cfg.Tun)
	return &Dhcp{
		Cfg: cfg,
		Ip: header.Str2IP(ip),
		Mask: header.MaskNumber2Mask(mask),
		UsedIps: map[uint32]bool{},
	}
}

func (dhcp *Dhcp) ApplyIp() (string, error){
	for ip := dhcp.Ip + 1; ip < ((dhcp.Ip & dhcp.Mask) ^ (^dhcp.Mask)); ip++ {
		if _, ok := dhcp.UsedIps[ip]; !ok {
			dhcp.UsedIps[ip] = true
			return header.IP2Str(ip), nil
		}
	}
	return "", fmt.Errorf("no enough ip")
}

func (dhcp *Dhcp) ReleaseIp(ip string){
	delete(dhcp.UsedIps, header.Str2IP(ip))
}