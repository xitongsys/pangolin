package main

import (
	"client"
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
	var err error
	fmt.Println("Welcome to use Pangolin!")
	defer func(){
		fmt.Println("[main] error: ", err)
	}()

	flag.Parse()
	cfg, err := config.NewConfigFromFile(*configFile)
	if err != nil {
		os.Exit(-1)
	}
	fmt.Println(cfg.String())

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
		if cfg.Protocol == "tcp" {
			tcpClient, err := client.NewTcpClient(cfg)
			if err != nil {
				os.Exit(-1)
			}

			tcpClient.Start()

		} else{
			udpClient, err := client.NewUdpClient(cfg)
			if err != nil {
				os.Exit(-1)
			}

			udpClient.Start()
		}

	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
