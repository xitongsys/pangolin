package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"

	"header"
)

const (
	IFF_NO_PI = 0x10
	IFF_TUN   = 0x01
	IFF_TAP   = 0x02
	TUNSETIFF = 0x400454CA
)

func main() {
	fd, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	ifr := make([]byte, 18)
	copy(ifr, []byte("tun0"))
	ifr[17] = IFF_NO_PI
	ifr[16] = IFF_TUN

	_, _, errn := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd.Fd()), uintptr(TUNSETIFF),
		uintptr(unsafe.Pointer(&ifr[0])))
	if errn != 0 {
		fmt.Println("ioctl err")
		return
	}

	for {
		data := make([]byte, 1500)
		n, err := fd.Read(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ipH := header.IPv4{}
		udpH := header.UDP{}
		ipH.Unmarshal(data[:20])
		udpH.Unmarshal(data[ipH.HeaderLen():])
		fmt.Println(ipH, udpH)

		conn, err2 := net.Dial("udp", "192.168.43.198:12345")
		defer conn.Close()

		if err2 != nil {
			fmt.Println(err)
			continue
		}

		n, err = conn.Write(data[:n])
		if err != nil {
			fmt.Println(err)
			continue
		}

	}

}
