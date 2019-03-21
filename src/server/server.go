package server

import (
	"fmt"
	"net"
	"time"

	"cache"
	"comp"
	"header"
	"tun"
)

type PServer struct {
	Addr      string
	ClientMap *cache.Cache
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
		ClientMap: cache.NewCache(time.Minute * 10),
		TunConn:   tun,
		UdpConn:   conn,
	}, nil
}

func (s *PServer) sendToClient() {
	data := make([]byte, s.TunConn.GetMtu()*2)
	for {
		if n, err := s.TunConn.Read(data); err == nil && n > 0 {
			if proto, src, dst, err := header.Get(data); err == nil {
				if caddr := s.ClientMap.Get(dst + "->" + src); caddr != "" {
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
				s.ClientMap.Put(src+"->"+dst, caddr.String())
				//s.ClientMap[src+"->"+dst] = caddr.String()
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
