package main

import (
	"tun"
	"flag"
	"fmt"
	"os"
	"sync"

	"client"
	"server"
)

var role = flag.String("role", "server", "")
var protocol = flag.String("protocol", "udp", "")
var saddr = flag.String("server", "0.0.0.0:12345", "")
var tunName = flag.String("tun", "tun0", "")
var mtu = flag.Int("mtu", 1500, "")


func main() {
	flag.Parse()
	fmt.Println("Welcome to use Pangolin!")
	if *role == "client" {
		if *protocol == "udp" {
			uc, err := client.NewUdpClient(*saddr, *tunName, *mtu)
			if err != nil {
				fmt.Println("[main] start udp client failed: ", err)
				os.Exit(-1)
			}
			uc.Start()

		}else if *protocol == "tcp" {
			tc, err := client.NewTcpClient(*saddr, *tunName, *mtu)
			if err != nil {
				fmt.Println("[main] start tcp client failed: ", err)
				os.Exit(-1)
			}
			tc.Start()
		} 

	} else {
		tunServer, err := tun.NewTunServer(*tunName, *mtu)
		if err != nil {
			fmt.Println("[main] tun server can't start: ", err)
			os.Exit(-1)
		}

		udpServer, err := server.NewUdpServer(*saddr, tunServer)
		if err != nil {
			fmt.Println("[main] udp server can't start: ", err)
			os.Exit(-1)
		}

		tcpServer, err := server.NewTcpServer(*saddr, tunServer)
		if err != nil {
			fmt.Println("[main] tcp server can't start: ", err)
			os.Exit(-1)
		}

		tunServer.Start()
		udpServer.Start()
		tcpServer.Start()

	}

	fmt.Printf("Run as %s, server:%s, tun:%s, mtu:%d\n", *role, *saddr, *tunName, *mtu)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
