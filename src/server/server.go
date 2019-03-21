package server

import (
	"fmt"
	"net"

	"header"
	"tun"
	"comp"
)

type PServer struct {
	Addr      string
	ClientMap map[string]string
	TunConn   tun.Tun
	UdpConn   *net.UDPConn
}

func NewPServer(addr string, tname string, mtu int) (*PServer, error) {
	add, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid address", addr)
	}

	conn, err := net.ListenUDP("udp", add)
	if err != nil {
		return nil, err
	}
	tun, err := tun.NewLinuxTun(tname, mtu)
	if err != nil {
		return nil, err
	}

	return &PServer{
		Addr:      addr,
		ClientMap: map[string]string{},
		TunConn:   tun,
		UdpConn:   conn,
	}, nil
}

func (s *PServer) sendToClient() {
	data := make([]byte, s.TunConn.GetMtu()*2)
	for {
		if n, err := s.TunConn.Read(data); err == nil && n > 0 {
			if proto, src, dst, err := header.Get(data); err == nil {
				if caddr, ok := s.ClientMap[dst+"->"+src]; ok {
					if add, err := net.ResolveUDPAddr("udp", caddr); err == nil {
						cmpData := comp.CompressGzip(data[:n])
						s.UdpConn.WriteToUDP(cmpData, add)
						fmt.Printf("[send] client:%s src:%s dst:%s proto:%s\n", caddr, src, dst, proto)
					}
				}
			}
		}
	}
}

func (s *PServer) recvFromClient() {
	data := make([]byte, s.TunConn.GetMtu()*2)
	for {
		if n, caddr, err := s.UdpConn.ReadFromUDP(data); err == nil && n > 0 {
			uncmpData, errc := comp.UncompressGzip(data[:n])
			if errc != nil {
				continue
			}
			if proto, src, dst, err := header.Get(uncmpData); err == nil {
				s.ClientMap[src+"->"+dst] = caddr.String()
				s.TunConn.Write(uncmpData)
				fmt.Printf("[recv] client:%s src:%s dst:%s proto:%s\n", caddr.String(), src, dst, proto)
			}
		}
	}
}

func (s *PServer) Start() error {
	go s.sendToClient()
	go s.recvFromClient()
	return nil
}

func (s *PServer) Stop() error {
	s.UdpConn.Close()
	s.TunConn.Close()
	return nil
}
