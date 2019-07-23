package client

import (
	"net"

	"github.com/xitongsys/pangolin/comp"
	"github.com/xitongsys/pangolin/config"
	"github.com/xitongsys/pangolin/header"
	"github.com/xitongsys/pangolin/logging"
	"github.com/xitongsys/pangolin/tun"
)

type UdpClient struct {
	ServerAdd string
	UdpConn   net.Conn
	TunConn   tun.Tun
}

func NewUdpClient(cfg *config.Config) (*UdpClient, error) {
	saddr, tname, mtu := cfg.ServerAddr, cfg.TunName, cfg.Mtu
	conn, err := net.Dial("udp", saddr)
	if err != nil {
		return nil, err
	}
	tun, err := tun.NewLinuxTun(tname, mtu)
	if err != nil {
		return nil, err
	}

	return &UdpClient{
		ServerAdd: saddr,
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
				logging.Log.Debugf("ToServer: protocol:%v, len:%v, src:%v, dst:%v", protocol, n, src, dst)
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
				logging.Log.Debugf("FromServer: protocol:%v, len:%v, src:%v, dst:%v", protocol, n, src, dst)
			}
		}
	}
}

func (uc *UdpClient) Start() error {
	logging.Log.Info("UdpClient started")
	go uc.writeToServer()
	go uc.readFromServer()
	return nil
}

func (uc *UdpClient) Stop() error {
	logging.Log.Info("UdpClient stopped")
	uc.UdpConn.Close()
	uc.TunConn.Close()
	return nil
}
