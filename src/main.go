package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"server"
	"login"
	"config"
)

var configFile = flag.String("c", "cfg.json", "")

func main() {
	fmt.Println("Welcome to use Pangolin!")

	flag.Parse()
	cfg, err := config.NewConfigFromFile(*configFile)
	if err != nil {
		os.Exit(-1)
	}

	if cfg.Role == "server" {
		loginManager, err := login.NewLoginManager(cfg)
		if err != nil {
			os.Exit(-1)
		}
		
		tcpServer, err := server.NewTcpServer(cfg, loginManager)
		if err != nil {
			os.Exit(-1)
		}

		udpServer, err := server.NewUdpServer(cfg, loginManager)
		if err != nil {
			os.Exit(-1)
		}

		loginManager.Start()
		tcpServer.Start()
		udpServer.Start()

	}else{
		
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
