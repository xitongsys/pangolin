package server

import (
	"net"

	"util"
	"comp"
	"header"
	"encrypt"
	"logging"
)

var USERCHANBUFFERSIZE = 1024

type User struct {
	Client string
	RemoteTunIp string
	LocalTunIp string
	Token string
	Key string
	TunToConnChan chan string
	ConnToTunChan chan string
	Conn net.Conn
	Logout func(client string)
}

func NewUser(client string, tun string, token string, conn net.Conn, logout func(string)) *User {
	key := string(encrypt.GetAESKey([]byte(token)))
	return &User {
		Client: client,
		LocalTunIp: tun,
		RemoteTunIp: tun,
		Token: token,
		Key: key,
		TunToConnChan: make(chan string, USERCHANBUFFERSIZE),
		ConnToTunChan: make(chan string, USERCHANBUFFERSIZE),
		Conn: conn,
		Logout: logout,
	}
}

func (user *User) Start() {
	encryptKey := encrypt.GetAESKey([]byte(user.Token))
	//read from client, write to channel
	go func() {
		for {
			var err error
			data, err := util.ReadPacket(user.Conn)
			if err != nil {
				user.Close()
				return
			}

			if ln := len(data); ln > 0 {
				if data, err = comp.UncompressGzip(data); err == nil && len(data)>0{
					if data, err = encrypt.DecryptAES(data, encryptKey); err == nil {
						if protocol, src, dst, err := header.GetBase(data); err == nil {
							remoteIp, _ := header.ParseAddr(src)
							user.RemoteTunIp = remoteIp
							Snat(data, user.LocalTunIp)
							user.ConnToTunChan <- string(data)
							logging.Log.Debugf("TcpFromClient: client:%v, protocol:%v, len:%v, src:%v, dst:%v", user.Client, protocol, ln, src, dst)
						}
					}
				}
			}
		}
	}()

	//read from channel, write to client
	go func() {
		for {
			datas, ok := <- user.TunToConnChan
			if !ok {
				user.Close()
				return
			}
			data := []byte(datas)
			if ln := len(data); ln > 0 {
				if protocol, src, dst, err := header.GetBase(data); err == nil {
					Dnat(data, user.RemoteTunIp)
					if endata, err := encrypt.EncryptAES(data, encryptKey); err == nil {
						if _, err := util.WritePacket(user.Conn, comp.CompressGzip(endata)); err != nil {
							user.Close()
							return
						}
						logging.Log.Debugf("TcpToClient: client:%v, protocol:%v, len:%v, src:%v, dst:%v", user.Client, protocol, ln, src, dst)
					}
				}
			}
		}
	}()
}

func (user *User) Close() {
	logging.Log.Infof("Client: %v closed", user.Client)
	go func(){
		defer func(){
			recover()
		}()
		close(user.TunToConnChan)
	}()

	go func(){
		defer func(){
			recover()
		}()
		close(user.ConnToTunChan)
	}()

	go func(){
		defer func(){
			recover()
		}()
		user.Conn.Close()
	}()

	user.Logout(user.Client)
}