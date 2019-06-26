package client

import (
	"fmt"
	"net"

	"comp"
	"tun"
	"header"
)

type UdpClient struct {
	ServerAdd string
	UdpConn   net.Conn
	TunConn   tun.Tun
}

func NewUdpClient(sadd string, tname string, mtu int) (*UdpClient, error) {
	conn, err := net.Dial("udp", sadd)
	if err != nil {
		return nil, err
	}
	tun, err := tun.NewLinuxTun(tname, mtu)
	if err != nil {
		return nil, err
	}

	return &UdpClient{
		ServerAdd: sadd,
		UdpConn:   conn,
		TunConn:   tun,
	}, nil
}

func (uc *UdpClient) writeToServer() {
	data := make([]byte, uc.TunConn.GetMtu()*2)
	for {
		if n, err := uc.TunConn.Read(data); err == nil && n > 0 {
			if protocol, src, dst, err := header.GetBase(data); err == nil {
				cmpData := comp.CompressGzip(data[:n])
				uc.UdpConn.Write(cmpData)
				fmt.Printf("[UdpClient][writeToServer] protocol:%v, len:%v, src:%v, dst:%v\n", protocol, n, src, dst)
			}
		}
	}
}

func (uc *UdpClient) readFromServer() error {
	data := make([]byte, uc.TunConn.GetMtu()*2)
	for {
		if n, err := uc.UdpConn.Read(data); err == nil && n > 0 {
			uncmpData, err2 := comp.UncompressGzip(data[:n])
			if err2 != nil {
				continue
			}
			if protocol, src, dst, err := header.GetBase(uncmpData); err == nil {
				uc.TunConn.Write(uncmpData)
				fmt.Printf("[UdpClient][readFromServer] protocol:%v, len:%v, src:%v, dst:%v\n", protocol, n, src, dst)
			}
		}
	}
}

func (uc *UdpClient) Start() error {
	fmt.Println("[UdpClient] startted.")
	go uc.writeToServer()
	go uc.readFromServer()
	return nil
}

func (uc *UdpClient) Stop() error {
	fmt.Println("[UdpClient] stopped.")
	uc.UdpConn.Close()
	uc.TunConn.Close()
	return nil
}
