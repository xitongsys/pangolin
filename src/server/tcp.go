package server

import (
	"comp"
	"net"

	"config"
	"util"
	"logging"
)

type TcpServer struct {
	Addr         string
	Cfg          *config.Config
	TcpListener  net.Listener
	LoginManager *LoginManager
}

func NewTcpServer(cfg *config.Config, loginManager *LoginManager) (*TcpServer, error) {
	tcpListener, err := net.Listen("tcp", cfg.ServerAddr)
	if err != nil {
		return nil, err
	}

	return &TcpServer{
		Addr:         cfg.ServerAddr,
		Cfg:          cfg,
		TcpListener:  tcpListener,
		LoginManager: loginManager,
	}, nil
}

func (ts *TcpServer) Start() {
	logging.Log.Info("TcpServer started")
	go func() {
		for {
			if conn, err := ts.TcpListener.Accept(); err == nil {
				go ts.handleRequest(conn)
			}
		}
	}()
}

func (ts *TcpServer) Stop() {
	logging.Log.Info("TcpServer stopped")
	ts.TcpListener.Close()
}

func (ts *TcpServer) handleRequest(conn net.Conn) {
	client := "tcp:" + conn.RemoteAddr().String()
	logging.Log.Infof("New connected client: %v", client)
	if err := ts.login(client, conn); err != nil {
		logging.Log.Errorf("client %v login failed: %v", client, err)
		return
	}
	ts.LoginManager.StartClient(client, conn)
}

func (ts *TcpServer) login(client string, conn net.Conn) error {
	if data, err := util.ReadPacket(conn); err != nil {
		return err

	} else {
		if data, err = comp.UncompressGzip(data); err != nil || len(data) <= 0 {
			return err

		} else {
			return ts.LoginManager.Login(client, string(data))
		}
	}
}
