package server

import (
	"comp"
	"fmt"
	"net"

	"config"
	"util"
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
	fmt.Println("[TcpServer] started.")
	go func() {
		for {
			if conn, err := ts.TcpListener.Accept(); err == nil {
				go ts.handleRequest(conn)
			}
		}
	}()
}

func (ts *TcpServer) Stop() {
	fmt.Println("[TcpServer] stopped.")
	ts.TcpListener.Close()
}

func (ts *TcpServer) handleRequest(conn net.Conn) {
	client := "tcp:" + conn.RemoteAddr().String()
	fmt.Printf("[TcpServer] new connected client: %v\n", client)
	if ts.login(client, conn) != nil {
		fmt.Printf("[TcpServer][Login] login failed: %v\n", client)
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
