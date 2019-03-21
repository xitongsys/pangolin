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
