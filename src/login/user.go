package login

import (
	"net"
	"fmt"

	"util"
	"comp"
	"header"
	"encrypt"
)

var USERCHANBUFFERSIZE = 1024

type User struct {
	Client string
	TunAddr string
	Token string
	Key string
	TunToConnChan chan string
	ConnToTunChan chan string
	Conn net.Conn
}

func NewUser(client string, tun string, token string, conn net.Conn) *User {
	key := string(encrypt.GetAESKey([]byte(token)))
	return &User {
		Client: client,
		TunAddr: tun,
		Token: token,
		Key: key,
		TunToConnChan: make(chan string, USERCHANBUFFERSIZE),
		ConnToTunChan: make(chan string, USERCHANBUFFERSIZE),
		Conn: conn,
	}
}

func (user *User) Start() {
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
					data = encrypt.DecryptAES(data, []byte(user.Token))
					if protocol, src, dst, err := header.GetBase(data); err == nil {
						user.ConnToTunChan <- string(data)
						fmt.Printf("[User][readFromClient] client:%v, protocol:%v, len:%v, src:%v, dst:%v\n", user.Client, protocol, ln, src, dst)
					}
				}
			}
		}
	}()

	//read from channel, write to client
	go func() {
		for {
			data, ok := <- user.TunToConnChan
			if !ok {
				user.Close()
				return
			}

			if ln := len(data); ln > 0 {
				if protocol, src, dst, err := header.GetBase([]byte(data)); err == nil {
					endata := encrypt.EncryptAES([]byte(data), []byte(user.Token))
					if _, err := util.WritePacket(user.Conn, comp.CompressGzip(endata)); err != nil {
						user.Close()
						return
					}
					fmt.Printf("[User][writeToClient] client:%v, protocol:%v, len:%v, src:%v, dst:%v\n", user.Client, protocol, ln, src, dst)
				}
			}
		}
	}()
}

func (user *User) Close() {
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
}