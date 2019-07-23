package server

import (
	"fmt"
	"net"

	"github.com/xitongsys/pangolin/comp"
	"github.com/xitongsys/pangolin/config"
	"github.com/xitongsys/pangolin/logging"
	"github.com/xitongsys/pangolin/util"
	"github.com/xitongsys/ptcp/ptcp"
)

type PTcpServer struct {
	Addr         string
	Cfg          *config.Config
	PTcpListener net.Listener
	LoginManager *LoginManager
}

func getPTcpAddr(addr string) string {
	ip, port := util.ParseAddr(addr)
	return fmt.Sprintf("%v:%v", ip, port+1)
}

func NewPTcpServer(cfg *config.Config, loginManager *LoginManager) (*PTcpServer, error) {
	addr := getPTcpAddr(cfg.ServerAddr)
	ptcp.Init(cfg.PtcpInterface)
	ptcpListener, err := ptcp.Listen("ptcp", addr)
	if err != nil {
		return nil, err
	}

	return &PTcpServer{
		Addr:         addr,
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
	buf := make([]byte, 1024)
	for {
		if n, err := conn.Read(buf); err != nil || n <= 0 {
			continue

		} else {
			data := buf[:n]
			if data, err = comp.UncompressGzip(data); err != nil || len(data) <= 0 {
				continue

			} else {
				return ts.LoginManager.Login(client, "ptcp", string(data))
			}
		}
	}
}
