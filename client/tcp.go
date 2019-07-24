package client

import (
	"fmt"
	"net"

	"github.com/xitongsys/ethernet-go/header"
	"github.com/xitongsys/pangolin/config"
	"github.com/xitongsys/pangolin/encrypt"
	"github.com/xitongsys/pangolin/logging"
	"github.com/xitongsys/pangolin/tun"
	"github.com/xitongsys/pangolin/util"
)

type TcpClient struct {
	ServerAdd string
	Cfg       *config.Config
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
		Cfg:       cfg,
		TcpConn:   conn,
		TunConn:   tun,
	}, nil
}

func (tc *TcpClient) writeToServer() {
	encryptKey := encrypt.GetAESKey([]byte(tc.Cfg.Tokens[0]))
	data := make([]byte, tc.TunConn.GetMtu()*2)
	for {
		if n, err := tc.TunConn.Read(data); err == nil && n > 0 {
			if protocol, src, dst, err := header.GetBase(data); err == nil {
				if endata, err := encrypt.EncryptAES(data[:n], encryptKey); err == nil {
					util.WritePacket(tc.TcpConn, endata)
					logging.Log.Debugf("ToServer: protocol:%v, len:%v, src:%v, dst:%v", protocol, n, src, dst)
				}
			}
		}
	}
}

func (tc *TcpClient) readFromServer() error {
	encryptKey := encrypt.GetAESKey([]byte(tc.Cfg.Tokens[0]))
	for {
		if data, err := util.ReadPacket(tc.TcpConn); err == nil {
			if data, err = encrypt.DecryptAES(data, encryptKey); err == nil {
				if protocol, src, dst, err := header.GetBase(data); err == nil {
					tc.TunConn.Write(data)
					logging.Log.Debugf("FromServer: protocol:%v, len:%v, src:%v, dst:%v", protocol, len(data), src, dst)
				}
			}
		}
	}
}

func (tc *TcpClient) login() error {
	if len(tc.Cfg.Tokens) <= 0 {
		return fmt.Errorf("no token provided")
	}
	data := []byte(tc.Cfg.Tokens[0])
	if _, err := util.WritePacket(tc.TcpConn, data); err != nil {
		return err
	}
	return nil
}

func (tc *TcpClient) Start() error {
	logging.Log.Info("TcpClient started")
	if err := tc.login(); err != nil {
		return err
	}
	go tc.writeToServer()
	go tc.readFromServer()
	return nil
}

func (tc *TcpClient) Stop() error {
	logging.Log.Info("TcpClient stopped")
	tc.TcpConn.Close()
	tc.TunConn.Close()
	return nil
}
