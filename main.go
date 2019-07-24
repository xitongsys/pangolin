package main

import (
	"flag"
	"os"
	"sync"

	"github.com/xitongsys/pangolin/client"
	"github.com/xitongsys/pangolin/config"
	"github.com/xitongsys/pangolin/logging"
	"github.com/xitongsys/pangolin/server"
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
		var tcpClient *client.TcpClient
		var ptcpClient *client.PTcpClient
		var udpClient *client.UdpClient

		if cfg.Protocol == "tcp" {
			if tcpClient, err = client.NewTcpClient(cfg); err == nil {
				err = tcpClient.Start()
			}

		} else if cfg.Protocol == "ptcp" {
			if ptcpClient, err = client.NewPTcpClient(cfg); err == nil {
				err = ptcpClient.Start()
			}

		} else {
			if udpClient, err = client.NewUdpClient(cfg); err == nil {
				err = udpClient.Start()
			}
		}

	}

	if err != nil {
		logging.Log.Error(err)
		os.Exit(-1)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
