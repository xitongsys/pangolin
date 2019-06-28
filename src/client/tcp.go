package client

import (
	"fmt"
	"net"

	"comp"
	"tun"
	"header"
	"util"
	"config"
)

type TcpClient struct {
	ServerAdd string
	Cfg *config.Config
	TcpConn   *net.TCPConn
	TunConn   tun.Tun
}

func NewTcpClient(cfg *config.Config) (*TcpClient, error) {
	saddr, tname, mtu := cfg.ServerAddr, cfg.TunName, cfg.Mtu
	addr, err := net.ResolveTCPAddr("", saddr)
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
		ServerAdd: saddr,
		Cfg: cfg,
		TcpConn:   conn,
		TunConn:   tun,
	}, nil
}

func (tc *TcpClient) writeToServer() {
	data := make([]byte, tc.TunConn.GetMtu()*2)
	for {
		if n, err := tc.TunConn.Read(data); err == nil && n > 0 {
			if protocol, src, dst, err := header.GetBase(data); err == nil {
				cmpData := comp.CompressGzip(data[:n])
				util.WritePacket(tc.TcpConn, cmpData)
				fmt.Printf("[TcpClient][writeToServer] protocol:%v, len:%v, src:%v, dst:%v\n", protocol, n, src, dst)
			}
		}
	}
}

func (tc *TcpClient) readFromServer() error {
	for {
		if data, err := util.ReadPacket(tc.TcpConn); err == nil {
			if data, err := comp.UncompressGzip(data); err == nil && len(data) > 0 {
				if protocol, src, dst, err := header.GetBase(data); err == nil {
					tc.TunConn.Write(data)
					fmt.Printf("[TcpClient][readFromServer] protocol:%v, len:%v, src:%v, dst:%v\n", protocol, len(data), src, dst)
				}
			}
		}
	}
}

func (tc *TcpClient) login() error {
	if len(tc.Cfg.Tokens) <= 0 {
		return fmt.Errorf("no token provided")
	}
	if _, err := util.WritePacket(tc.TcpConn, []byte(tc.Cfg.Tokens[0])); err!=nil{
		return err
	}
	return nil
}

func (tc *TcpClient) Start() error {
	fmt.Println("[TcpClient] started.")
	if err := tc.login(); err != nil {
		return err
	}
	go tc.writeToServer()
	go tc.readFromServer()
	return nil
}

func (tc *TcpClient) Stop() error {
	fmt.Println("[TcpClient] stopped.")
	tc.TcpConn.Close()
	tc.TunConn.Close()
	return nil
}
