package server

import (
	"fmt"
	"net"

	"login"
	"config"
)

type TcpServer struct {
	Addr      string
	TcpListener	 net.Listener
	LoginManager *login.LoginManager
}

func NewTcpServer(cfg *config.Config, loginManager *login.LoginManager) (*TcpServer, error) {
	tcpListener, err := net.Listen("tcp", cfg.ServerAddr)
	if err != nil {
		return nil, err
	}

	return &TcpServer {
		Addr: cfg.ServerAddr,
		TcpListener: tcpListener,
		LoginManager: loginManager,
	}, nil
}

func (ts *TcpServer) Start() {
	fmt.Println("[TcpServer] started.")
	for {
		if conn, err := ts.TcpListener.Accept(); err == nil{
			go ts.handleRequest(conn)
		}
	}
}

func (ts *TcpServer) Stop() {
	fmt.Println("[TcpServer] stopped.")
	ts.TcpListener.Close()
}

func (ts *TcpServer) handleRequest(conn net.Conn) {
	client := "tcp:" + conn.RemoteAddr().String()
	fmt.Printf("[TcpServer] new connected client: %v\n", client)
	ts.LoginManager.StartClient(client, conn)
}

