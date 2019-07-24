package client

import (
	"fmt"
	"net"
	"time"

	"github.com/xitongsys/pangolin/config"
	"github.com/xitongsys/pangolin/encrypt"
	"github.com/xitongsys/pangolin/header"
	"github.com/xitongsys/pangolin/logging"
	"github.com/xitongsys/pangolin/protocol"
	"github.com/xitongsys/pangolin/tun"
	"github.com/xitongsys/pangolin/util"
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
	ptcp.Init(cfg.PtcpInterface)
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
			if proto, src, dst, err := header.GetBase(data); err == nil {
				if endata, err := encrypt.EncryptAES(data[:n], encryptKey); err == nil {
					packet := append([]byte{protocol.PTCP_PACKETTYPE_DATA}, endata...)
					tc.PTcpConn.Write(packet)
					logging.Log.Debugf("ToServer: protocol:%v, len:%v, src:%v, dst:%v", proto, n, src, dst)
				}
			}
		}
	}
}

func (tc *PTcpClient) readFromServer() error {
	encryptKey := encrypt.GetAESKey([]byte(tc.Cfg.Tokens[0]))
	buf := make([]byte, tc.TunConn.GetMtu()*2)
	for {
		if n, err := tc.PTcpConn.Read(buf); err == nil && n > 1 && buf[0] == protocol.PTCP_PACKETTYPE_DATA {
			data := buf[1:n]
			if data, err = encrypt.DecryptAES(data, encryptKey); err == nil {
				if protocol, src, dst, err := header.GetBase(data); err == nil {
					tc.TunConn.Write(data)
					logging.Log.Debugf("FromServer: protocol:%v, len:%v, src:%v, dst:%v", protocol, len(data), src, dst)
				}
			}
		}
	}
}

func (tc *PTcpClient) login() error {
	if len(tc.Cfg.Tokens) <= 0 {
		return fmt.Errorf("no token provided")
	}
	data := append([]byte{0}, []byte(tc.Cfg.Tokens[0])...)
	timeout := time.Second * 20
	res, err := util.WriteUntil(tc.PTcpConn, 1024, data, timeout,
		func(ds []byte) bool {
			if len(ds) <= 1 || ds[0] != protocol.PTCP_PACKETTYPE_LOGIN {
				return false
			}
			return true
		})
	if err != nil {
		return err
	}

	if res[1] != protocol.PTCP_LOGINMSG_SUCCESS {
		return fmt.Errorf("login failed")
	}

	return nil
}

func (tc *PTcpClient) Start() error {
	logging.Log.Info("PTcpClient login...")
	if err := tc.login(); err != nil {
		return err
	}
	go tc.writeToServer()
	go tc.readFromServer()
	logging.Log.Info("PTcpClient started")
	return nil
}

func (tc *PTcpClient) Stop() error {
	tc.PTcpConn.Close()
	tc.TunConn.Close()
	logging.Log.Info("PTcpClient stopped")
	return nil
}
