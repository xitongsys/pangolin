package util

import (
	"strconv"
	"strings"
)

func ParseAddr(addr string) (ip string, port int) {
	as := strings.Split(addr, ":")
	ip = as[0]
	port, _ = strconv.Atoi(as[1])
	return
}
