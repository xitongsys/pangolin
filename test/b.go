package main

import (
	"fmt"
	"syscall"

	"header"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer syscall.Close(fd)

	//f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	for {
		buf := make([]byte, 65565)
		ipH := header.IPv4{}
		udpH := header.UDP{}

		n, _, err := syscall.Recvfrom(fd, buf, 0)
		//n, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		if n < 20 {
			fmt.Println("short")
		} else {
			ipH.Unmarshal(buf[:20])
			udpH.Unmarshal(buf[ipH.HeaderLen():])

			fmt.Println(n)
			fmt.Println(ipH, udpH)
			data := buf[ipH.HeaderLen()+udpH.HeaderLen() : n]
			fmt.Println(string(data))
		}
	}
}
