package server

import (
	"fmt"
	"net"

	"tun"
	"comp"
	"util"
)

type TcpServer struct {
	Addr      string
	TcpListener	 net.Listener
	TunServer *tun.TunServer
}

func NewTcpServer(addr string, tunServer *tun.TunServer) (*TcpServer, error) {
	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &TcpServer {
		Addr: addr,
		TcpListener: tcpListener,
		TunServer: tunServer,
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

func (ts *TcpServer) handleRequest(conn net.Conn) {
	//read from client, write to channel
	go func() {
		var err error
		data := util.ReadPacket(conn)
		if data, err = comp.UncompressGzip(data); err == nil && len(data)>0{
			ts.TunServer.WriteToChannel("tcp", ts.Addr, data)
		}
	}()

	//read from channel, write to client
	go func() {
		for {
			if data := comp.CompressGzip(ts.TunServer.ReadFromChannel(ts.Addr)); len(data)>0 {
				util.WritePacket(conn, data)
			}
		}
	}()
}

