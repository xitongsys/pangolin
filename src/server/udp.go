package server

import (
	"fmt"
	"net"

	"comp"
	"header"
	"tun"
)

type UdpServer struct {
	Addr      string
	UdpConn   *net.UDPConn
	TunServer *tun.TunServer
}

func NewUdpServer(addr string, tunServer *tun.TunServer) (*UdpServer, error) {
	add, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid address", addr)
	}

	conn, err := net.ListenUDP("udp", add)
	if err != nil {
		return nil, err
	}

	return &UdpServer{
		Addr:      addr,
		UdpConn:   conn,
		TunServer: tunServer,
	}, nil
}

func (us *UdpServer) sendToClient() {
	for {
		if data := us.TunServer.ReadFromUdpChannel(); len(data) > 0 {
			if proto, src, dst, err := header.GetBase(data); err == nil {
				key := proto + ":" + dst + "->" + src
				if cprotocal, caddr := us.TunServer.GetClientAddr(key); cprotocal!= "" && caddr != "" {
					if add, err := net.ResolveUDPAddr("udp", caddr); err == nil {
						cmpData := comp.CompressGzip(data)
						us.UdpConn.WriteToUDP(cmpData, add)
						fmt.Printf("[UdpServer][sendToClient] client:%s src:%s dst:%s proto:%s\n", caddr, src, dst, proto)
					}
				}
			}
		}
	}
}

func (us *UdpServer) recvFromClient() {
	data := make([]byte, us.TunServer.TunConn.GetMtu()*2)
	for {
		if n, caddr, err := us.UdpConn.ReadFromUDP(data); err == nil && n > 0 {
			uncmpData, errc := comp.UncompressGzip(data[:n])
			if errc != nil {
				continue
			}
			if proto, src, dst, err := header.GetBase(uncmpData); err == nil {
				us.TunServer.WriteToChannel("udp", caddr.String(), uncmpData)
				fmt.Printf("[UdpServer][recvFromClient] client:%s src:%s dst:%s proto:%s\n", caddr.String(), src, dst, proto)
			}
		}
	}
}

func (us *UdpServer) Start() error {
	go us.sendToClient()
	go us.recvFromClient()
	return nil
}

func (us *UdpServer) Stop() error {
	us.UdpConn.Close()
	return nil
}
