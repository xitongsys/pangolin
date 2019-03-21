package main

import (
	"fmt"
	"syscall"
	"time"

	"header"
)

func main() {
	ipH := header.IPv4{
		VerIHL:   0x45,
		Tos:      0,
		Len:      33,
		Id:       123,
		Offset:   16384,
		TTL:      64,
		Protocol: 17,
		Checksum: 0,
		Src:      header.Str2IP("192.168.43.198"),
		Dst:      header.Str2IP("139.180.132.42"),
	}
	udpH := header.UDP{
		SrcPort:  8765,
		DstPort:  33333,
		Len:      13,
		Checksum: 0,
	}

	for {
		ipH.Id += 100
		ipH.Checksum = ipH.CalChecksum()
		udpH.SrcPort = 8765

		ipbs := ipH.Marshal()
		udpbs := udpH.Marshal()

		fbuf := make([]byte, 0)
		fbuf = append(fbuf, ipbs...)
		fbuf = append(fbuf, udpbs...)
		fbuf = append(fbuf, []byte("hello")...)

		addr := syscall.SockaddrInet4{
			Port: 0,
			Addr: [4]byte{139, 180, 132, 42},
		}

		fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
		if err := syscall.Sendto(fd, fbuf, 0, &addr); err != nil {
			fmt.Println(err)
		}
		syscall.Close(fd)

		fmt.Println(ipH, udpH)

		time.Sleep(time.Second)

	}

}
