package main

import (
	"fmt"
	"syscall"

	"header"
)

func main(){
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	//f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		src, dst, data := "192.168.35.67:12345", "47.240.40.78:12345", []byte("hello,world")
		packet := header.BuildUdpPacket(src, dst, data)
		proto, src, dst, err := header.GetBase(packet)
		fmt.Println(proto, src, dst, err)
		addr := syscall.SockaddrInet4{
			Port: 12345,
			Addr: [4]byte{47, 240, 40, 78},
		}
		err = syscall.Sendto(fd, packet, 0, &addr)
		fmt.Println(err)
	}
}