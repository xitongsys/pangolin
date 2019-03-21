package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"client"
	"server"
)

var role = flag.String("role", "server", "")
var saddr = flag.String("server", "0.0.0.0:12345", "")
var tun = flag.String("tun", "tun0", "")
var mtu = flag.Int("mtu", 1500, "")

func main() {
	flag.Parse()
	fmt.Println("Welcome to use Pangolin!")
	if *role == "client" {
		cp, err := client.NewPClient(*saddr, *tun, *mtu)
		if err != nil {
			fmt.Println("start client failed: ", err)
			os.Exit(-1)
		}
		cp.Start()

	} else {
		sp, err := server.NewPServer(*saddr, *tun, *mtu)
		if err != nil {
			fmt.Println("start server failed: ", err)
			os.Exit(-1)
		}
		sp.Start()
	}
	fmt.Printf("Run as %s, server:%s, tun:%s, mtu:%d\n", *role, *saddr, *tun, *mtu)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
