package server

import (
	"github.com/xitongsys/pangolin/header"
)

func Snat(data []byte, src string) {
	protocol, iph, _, _, _, err := header.Get(data)
	if err != nil {
		return
	}
	iph.Src = header.Str2IP(src)
	newdata := iph.Marshal()
	for i := 0; i < len(newdata); i++ {
		data[i] = newdata[i]
	}
	if protocol == "tcp" {
		header.ReCalTcpCheckSum(data)
	} else if protocol == "udp" {
		header.ReCalUdpCheckSum(data)
	}
}

func Dnat(data []byte, dst string) {
	protocol, iph, _, _, _, err := header.Get(data)
	if err != nil {
		return
	}
	iph.Dst = header.Str2IP(dst)
	newdata := iph.Marshal()
	for i := 0; i < len(newdata); i++ {
		data[i] = newdata[i]
	}
	if protocol == "tcp" {
		header.ReCalTcpCheckSum(data)
	} else if protocol == "udp" {
		header.ReCalUdpCheckSum(data)
	}
}
