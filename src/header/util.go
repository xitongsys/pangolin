package header

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	TCPID = 6
	UDPID = 17
)

func IP2Str(ip uint32) string {
	res := "%d.%d.%d.%d"
	return fmt.Sprintf(res, (ip>>24)&0xff, (ip>>16)&0xff, (ip>>8)&0xff, ip&0xff)
}

func Str2IP(s string) uint32 {
	ns := strings.Split(s, ".")
	res := uint32(0)
	for i := 0; i < 4; i++ {
		n, _ := strconv.ParseInt(ns[3-i], 10, 16)
		res += (uint32(n) << uint32(i*8))
	}
	return res
}

//src: IP:PORT
func ParseAddr(src string) (string, int) {
	res := strings.Split(src, ":")
	if len(res) == 0 {
		return "",-1
	}else if len(res) == 1 {//":port"
		port, _ := strconv.Atoi(res[0])
		return "127.0.0.1", port
	}else{
		port, _ := strconv.Atoi(res[1])
		return res[0], port
	}
}

//src: ip/mask
func ParseNet(src string) (string, int) {
	res := strings.Split(src, "/")
	if len(res) == 0 {
		return "",-1
	}else if len(res) == 1 {//"/mask"
		mask, _ := strconv.Atoi(res[0])
		return "127.0.0.1", mask
	}else{
		mask, _ := strconv.Atoi(res[1])
		return res[0], mask
	}
}

func MaskNumber2Mask(mask int) uint32 {
	res := uint32(0)
	for i := 0; i<mask; i++ {
		res |= (uint32(1)<<uint32(i))
	}
	return res<<uint32(32 - mask)
}