package server

import (
	"comp"
	"net"

	"config"
	"logging"
	"util"

	"github.com/xitongsys/ptcp/ptcp"
)

type PTcpServer struct {
	Addr         string
	Cfg          *config.Config
	PTcpListener net.Listener
	LoginManager *LoginManager
}

func NewPTcpServer(cfg *config.Config, loginManager *LoginManager) (*PTcpServer, error) {
	ptcpListener, err := ptcp.Listen("tcp", cfg.ServerAddr)
	if err != nil {
		return nil, err
	}

	return &PTcpServer{
		Addr:         cfg.ServerAddr,
		Cfg:          cfg,
		PTcpListener: ptcpListener,
		LoginManager: loginManager,
	}, nil
}

func (ts *PTcpServer) Start() {
	logging.Log.Info("PTcpServer started")
	go func() {
		for {
			if conn, err := ts.PTcpListener.Accept(); err == nil {
				go ts.handleRequest(conn)
			}
		}
	}()
}

func (ts *PTcpServer) Stop() {
	logging.Log.Info("PTcpServer stopped")
	ts.PTcpListener.Close()
}

func (ts *PTcpServer) handleRequest(conn net.Conn) {
	client := "ptcp:" + conn.RemoteAddr().String()
	logging.Log.Infof("New connected client: %v", client)
	if err := ts.login(client, conn); err != nil {
		logging.Log.Errorf("Client %v login failed: %v", client, err)
		return
	}
	ts.LoginManager.StartClient(client, conn)
}

func (ts *PTcpServer) login(client string, conn net.Conn) error {
	return nil
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
