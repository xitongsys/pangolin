package server

import (
	"fmt"
	"net"
	"time"

	"github.com/xitongsys/pangolin/config"
	"github.com/xitongsys/pangolin/logging"
	"github.com/xitongsys/pangolin/protocol"
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
		conn.Close()
		return
	}
	ts.LoginManager.StartClient(client, conn)
}

func (ts *PTcpServer) login(client string, conn net.Conn) error {
	buf := make([]byte, 1024)
	var n int
	var err error

	after := time.After(time.Second * 20)
	for {
		select {
		case <-after:
			return fmt.Errorf("timeout")
		default:
		}

		n, err = conn.Read(buf)
		if err != nil {
			return err
		}
		if n <= 1 {
			continue
		}

		if buf[0] == protocol.PTCP_PACKETTYPE_LOGIN {
			data := buf[1:n]
			if err = ts.LoginManager.Login(client, "ptcp", string(data)); err != nil {
				return err
			}
			break
		}
	}

	timeout := time.Second * 20
	data := []byte{protocol.PTCP_PACKETTYPE_LOGIN, protocol.PTCP_LOGINMSG_SUCCESS}
	_, err = util.WriteUntil(conn, 10240, data, timeout,
		func(ds []byte) bool {
			if len(ds) <= 1 || ds[0] != protocol.PTCP_PACKETTYPE_DATA {
				return false
			}
			return true
		})

	return err
}
