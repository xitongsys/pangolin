package server

import (
	"fmt"
	"net"

	"tun"
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
	fmt.Println("[TcpServer] started")
	for {
		if conn, err := ts.TcpListener.Accept(); err == nil{
			go ts.handleRequest(conn)
		}
	}
}

func (ts *TcpServer) handleRequest(conn net.Conn) {
	//write to channel
	go func() {
		data := make([]byte, 0)
		for {
			lenBs := make([]byte, 1)
			n, err := conn.Read(lenBs)
			if n > 0 && err == nil {
				len := int(lenBs[0])
				if len == 0 {
					ts.TunServer.WriteToChannel("tcp", ts.Addr, data)
					data = make([]byte, 0)

				}else{
					cur := make([]byte, len)
					left := len
					for left > 0 {
						n, err := conn.Read(cur[len-left:])
						if n > 0 && err == nil {
							left -= n
						}
					}
					data = append(data, cur...)
				}
			}
		}
	}()

	//read from channel
	go func() {
		for {
			data := ts.TunServer.ReadFromChannel(ts.Addr)
			for len(data) > 0 {
				wc := 255
				if len(data) < wc {
					wc = len(data)
				}

				lenBs := []byte{byte(wc)}
				for {
					if n, err := conn.Write(lenBs); n>0 && err == nil {
						break
					}
				}

				wd := data[:wc]
				for len(wd) > 0 {
					if n, err := conn.Write(wd); n>0 && err==nil {
						wd = wd[n:]
					}
				}

				data = data[wc:]
			}
		}
	}()
}

