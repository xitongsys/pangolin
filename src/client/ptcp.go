package client

import (
	"fmt"
	"net"

	"comp"
	"config"
	"encrypt"
	"header"
	"logging"
	"tun"
	"util"

	"github.com/xitongsys/ptcp/ptcp"
)

type PTcpClient struct {
	ServerAdd string
	Cfg       *config.Config
	PTcpConn  net.Conn
	TunConn   tun.Tun
}

func getPTcpAddr(addr string) string {
	ip, port := util.ParseAddr(addr)
	return fmt.Sprintf("%v:%v", ip, port+1)
}

func NewPTcpClient(cfg *config.Config) (*PTcpClient, error) {
	ptcp.Init("eth0")
	addr, tname, mtu := getPTcpAddr(cfg.ServerAddr), cfg.TunName, cfg.Mtu
	conn, err := ptcp.Dial("ptcp", addr)
	if err != nil {
		return nil, err
	}

	tun, err := tun.NewLinuxTun(tname, mtu)
	if err != nil {
		return nil, err
	}

	return &PTcpClient{
		ServerAdd: addr,
		Cfg:       cfg,
		PTcpConn:  conn,
		TunConn:   tun,
	}, nil
}

func (tc *PTcpClient) writeToServer() {
	encryptKey := encrypt.GetAESKey([]byte(tc.Cfg.Tokens[0]))
	data := make([]byte, tc.TunConn.GetMtu()*2)
	for {
		if n, err := tc.TunConn.Read(data); err == nil && n > 0 {
			if protocol, src, dst, err := header.GetBase(data); err == nil {
				if endata, err := encrypt.EncryptAES(data[:n], encryptKey); err == nil {
					cmpData := comp.CompressGzip(endata)
					util.WritePacket(tc.PTcpConn, cmpData)
					logging.Log.Debugf("ToServer: protocol:%v, len:%v, src:%v, dst:%v", protocol, n, src, dst)
				}
			}
		}
	}
}

func (tc *PTcpClient) readFromServer() error {
	encryptKey := encrypt.GetAESKey([]byte(tc.Cfg.Tokens[0]))
	for {
		if data, err := util.ReadPacket(tc.PTcpConn); err == nil {
			if data, err := comp.UncompressGzip(data); err == nil && len(data) > 0 {
				if data, err = encrypt.DecryptAES(data, encryptKey); err == nil {
					if protocol, src, dst, err := header.GetBase(data); err == nil {
						tc.TunConn.Write(data)
						logging.Log.Debugf("FromServer: protocol:%v, len:%v, src:%v, dst:%v", protocol, len(data), src, dst)
					}
				}
			}
		}
	}
}

func (tc *PTcpClient) login() error {
	return nil
	if len(tc.Cfg.Tokens) <= 0 {
		return fmt.Errorf("no token provided")
	}
	data := comp.CompressGzip([]byte(tc.Cfg.Tokens[0]))
	if _, err := util.WritePacket(tc.PTcpConn, data); err != nil {
		return err
	}
	return nil
}

func (tc *PTcpClient) Start() error {
	logging.Log.Info("PTcpClient started")
	if err := tc.login(); err != nil {
		return err
	}
	go tc.writeToServer()
	go tc.readFromServer()
	return nil
}

func (tc *PTcpClient) Stop() error {
	logging.Log.Info("PTcpClient stopped")
	tc.PTcpConn.Close()
	tc.TunConn.Close()
	return nil
}
