package main

import (
	"fmt"
	"time"

	"header"
	"tun"
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
		Dst:      header.Str2IP("10.0.0.2"),
	}
	udpH := header.UDP{
		SrcPort:  8765,
		DstPort:  33333,
		Len:      13,
		Checksum: 0,
	}

	tuncon, err := tun.NewLinuxTun("tun0", 1500)
	if err!=nil{
		fmt.Println("err", err)
		return
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

		tuncon.Write(fbuf)

		time.Sleep(time.Second)

	}

}
