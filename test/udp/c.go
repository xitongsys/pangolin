package main

import (
	"fmt"
	"net"
)

func main() {
	raddr,_ := net.ResolveUDPAddr("udp","139.180.132.42:12345")
	laddr,_ := net.ResolveUDPAddr("udp","10.0.0.2:12345")

	conn, err := net.DialUDP("udp", laddr, raddr)
	//conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		n, err := conn.Write([]byte("from client"))
		if err != nil {
			fmt.Println(err)
			continue
		}

		data := make([]byte, 1024)
		n, err = conn.Read(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(string(data[:n]))
	}
}
