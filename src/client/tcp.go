package client

import (
	"fmt"
	"net"

	"comp"
	"tun"
	"header"
	"util"
)

type TcpClient struct {
	ServerAdd string
	TcpConn   *net.TCPConn
	TunConn   tun.Tun
}

func NewTcpClient(sadd string, tname string, mtu int) (*TcpClient, error) {
	addr, err := net.ResolveTCPAddr("", sadd)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		return nil, err
	}

	tun, err := tun.NewLinuxTun(tname, mtu)
	if err != nil {
		return nil, err
	}

	return &TcpClient{
		ServerAdd: sadd,
		TcpConn:   conn,
		TunConn:   tun,
	}, nil
}

func (tc *TcpClient) sendToServer() {
	data := make([]byte, tc.TunConn.GetMtu()*2)
	for {
		if n, err := tc.TunConn.Read(data); err == nil && n > 0 {
			if proto, src, dst, err := header.GetBase(data); err == nil {
				cmpData := comp.CompressGzip(data[:n])
				util.WritePacket(tc.TcpConn, cmpData)
				fmt.Printf("[TcpClient][sendToServer] Len:%d src:%s dst:%s proto:%s\n", n, src, dst, proto)
			}
		}
	}
}

func (tc *TcpClient) recvFromServer() error {
	for {
		if data, err := comp.UncompressGzip(util.ReadPacket(tc.TcpConn)); err == nil && len(data) > 0 {
			if proto, src, dst, err := header.GetBase(data); err == nil {
				tc.TunConn.Write(data)
				fmt.Printf("[TcpClient][recvFromServer] Len:%d src:%s dst:%s proto:%s\n", len(data), src, dst, proto)
			}
		}
	}
}

func (tc *TcpClient) Start() error {
	go tc.sendToServer()
	go tc.recvFromServer()
	return nil
}

func (tc *TcpClient) Stop() error {
	tc.TcpConn.Close()
	tc.TunConn.Close()
	return nil
}
