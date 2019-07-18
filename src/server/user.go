package server

import (
	"net"

	"comp"
	"encrypt"
	"header"
	"logging"
	"util"
)

var USERCHANBUFFERSIZE = 1024
var READBUFFERSIZE = 65535

type User struct {
	Client        string
	Protocol      string
	RemoteTunIp   string
	LocalTunIp    string
	Token         string
	Key           string
	TunToConnChan chan string
	ConnToTunChan chan string
	Conn          net.Conn
	Logout        func(client string)
}

func NewUser(client string, protocol string, tun string, token string, conn net.Conn, logout func(string)) *User {
	key := string(encrypt.GetAESKey([]byte(token)))
	return &User{
		Client:        client,
		Protocol:      protocol,
		LocalTunIp:    tun,
		RemoteTunIp:   tun,
		Token:         token,
		Key:           key,
		TunToConnChan: make(chan string, USERCHANBUFFERSIZE),
		ConnToTunChan: make(chan string, USERCHANBUFFERSIZE),
		Conn:          conn,
		Logout:        logout,
	}
}

func (user *User) Start() {
	encryptKey := encrypt.GetAESKey([]byte(user.Token))
	//read from client, write to channel
	buf := make([]byte, READBUFFERSIZE)
	go func() {
		for {
			var err error
			var data []byte
			if user.Protocol == "tcp" {
				data, err = util.ReadPacket(user.Conn)
			} else {
				_, err = user.Conn.Read(buf)
			}

			if err != nil {
				user.Close()
				return
			}

			if ln := len(data); ln > 0 {
				if data, err = comp.UncompressGzip(data); err == nil && len(data) > 0 {
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
			datas, ok := <-user.TunToConnChan
			if !ok {
				user.Close()
				return
			}
			data := []byte(datas)
			if ln := len(data); ln > 0 {
				if protocol, src, dst, err := header.GetBase(data); err == nil {
					Dnat(data, user.RemoteTunIp)
					if endata, err := encrypt.EncryptAES(data, encryptKey); err == nil {
						if user.Protocol == "tcp" {
							_, err = util.WritePacket(user.Conn, comp.CompressGzip(endata))
						} else {
							_, err = user.Conn.Write(comp.CompressGzip(endata))
						}

						if err != nil {
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
	go func() {
		defer func() {
			recover()
		}()
		close(user.TunToConnChan)
	}()

	go func() {
		defer func() {
			recover()
		}()
		close(user.ConnToTunChan)
	}()

	go func() {
		defer func() {
			recover()
		}()
		user.Conn.Close()
	}()

	user.Logout(user.Client)
}
