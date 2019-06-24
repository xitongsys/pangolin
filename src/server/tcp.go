package server

import (
	"fmt"
	"net"

	"tun"
	"comp"
	"util"
	"header"
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

func (ts *TcpServer) Stop() {
	fmt.Println("[TcpServer] stopped.")
	ts.TcpListener.Close()
}

func (ts *TcpServer) handleRequest(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	//read from client, write to channel
	go func() {
		for {
			var err error
			data := util.ReadPacket(conn)
			if data, err = comp.UncompressGzip(data); err == nil && len(data)>0{
				if protocol, src, dst, err := header.GetBase(data); err == nil {
					ts.TunServer.WriteToChannel("tcp", ts.Addr, data)
					fmt.Printf("[TcpServer][readFromClient] client:%v, protocol:%v, len:%v, src:%v, dst:%v\n", clientAddr, protocol, len(data), src, dst)
				}
			}
		}
	}()

	//read from channel, write to client
	go func() {
		for {
			data := ts.TunServer.ReadFromChannel(ts.Addr)
			if protocol, src, dst, err := header.GetBase(data); err == nil {
				util.WritePacket(conn, comp.CompressGzip(data))
				fmt.Printf("[TcpServer][writeToClient] client:%v, protocol:%v, len:%v, src:%v, dst:%v\n", clientAddr, protocol, len(data), src, dst)
			}
		}
	}()
}

