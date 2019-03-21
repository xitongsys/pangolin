package main

import (
	"fmt"
	"net"
	"syscall"

	"header"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:12345")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		data := make([]byte, 1024)
		n, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ipH := header.IPv4{}
		udpH := header.UDP{}

		ipH.Unmarshal(data[:20])
		udpH.Unmarshal(data[ipH.HeaderLen():])

		//nat
		ipH.Src = header.Str2IP("192.168.43.198")
		ipH.Checksum = 0
		ipH.Checksum = ipH.CalChecksum()
		udpH.Checksum = 0

		ipbs := ipH.Marshal()
		udpbs := udpH.Marshal()

		fmt.Println(rAddr, ipH, udpH)

		fbuf := make([]byte, 0)
		fbuf = append(fbuf, ipbs...)
		fbuf = append(fbuf, udpbs...)
		fbuf = append(fbuf, data[ipH.HeaderLen()+udpH.HeaderLen():n]...)

		addr := syscall.SockaddrInet4{
			Port: 0,
			Addr: [4]byte{139, 180, 132, 42},
		}
		fd2, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
		if err := syscall.Sendto(fd2, fbuf, 0, &addr); err != nil {
			fmt.Println("sendto", err)
			continue
		}
		syscall.Close(fd2)
	}

}
