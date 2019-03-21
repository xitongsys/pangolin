package main

import (
	"fmt"
	"os"
	"syscall"

	"header"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer syscall.Close(fd)

	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	for {
		buf := make([]byte, 65565)
		ipH := header.IPv4{}
		tcpH := header.TCP{}

		n, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		if n < 20 {
			fmt.Println("short")
		} else {
			ipH.Unmarshal(buf[:20])
			tcpH.Unmarshal(buf[ipH.HeaderLen():])
			fmt.Println(ipH, tcpH)
		}
	}
}
