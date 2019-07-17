package main

import (
	"client"
	"flag"
	"os"
	"sync"

	"config"
	"logging"
	"server"
)

var configFile = flag.String("c", "cfg.json", "")
var logLevel = flag.String("l", "info", "")

func main() {
	var err error
	logging.Log.Info("Welcome to use Pangolin!")
	defer func() {
		logging.Log.Error(err)
	}()

	flag.Parse()
	logging.SetLevel(*logLevel)

	cfg, err := config.NewConfigFromFile(*configFile)
	if err != nil {
		os.Exit(-1)
	}
	logging.Log.Info(cfg.String())

	if cfg.Role == "server" {
		loginManager, err := server.NewLoginManager(cfg)
		if err != nil {
			logging.Log.Error(err)
			os.Exit(-1)
		}

		tcpServer, err := server.NewTcpServer(cfg, loginManager)
		if err != nil {
			logging.Log.Error(err)
			os.Exit(-1)
		}

		ptcpServer, err := server.NewPTcpServer(cfg, loginManager)
		if err != nil {
			logging.Log.Error(err)
			os.Exit(-1)
		}

		udpServer, err := server.NewUdpServer(cfg, loginManager)
		if err != nil {
			logging.Log.Error(err)
			os.Exit(-1)
		}

		loginManager.Start()
		tcpServer.Start()
		udpServer.Start()
		ptcpServer.Start()

	} else {
		if cfg.Protocol == "tcp" {
			tcpClient, err := client.NewTcpClient(cfg)
			if err != nil {
				logging.Log.Error(err)
				os.Exit(-1)
			}

			tcpClient.Start()

		} else if cfg.Protocol == "ptcp" {
			ptcpClient, err := client.NewPTcpClient(cfg)
			if err != nil {
				logging.Log.Error(err)
				os.Exit(-1)
			}

			ptcpClient.Start()

		} else {
			udpClient, err := client.NewUdpClient(cfg)
			if err != nil {
				logging.Log.Error(err)
				os.Exit(-1)
			}
			udpClient.Start()
		}

	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
