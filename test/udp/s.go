package main

import (
	"fmt"
	"net"
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
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(rAddr, string(data))

		_, err = conn.WriteToUDP([]byte("from server"), rAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}
