package login

import (
	"encrypt"
	"net"

	"util"
)

var CHANBUFFERSIZE = 1024

type User struct {
	Client string
	TunAddr string
	Token string
	Key string
	Chan chan string
	Conn net.Conn
}

func NewUser(client string, tun string, token string) *User {
	key := string(encrypt.GetAESKey([]byte(token)))
	return &User {
		Client: client,
		TunAddr: tun,
		Token: token,
		Key: key,
		Chan: make(chan string, CHANBUFFERSIZE),
	}
}

func (user *User) Start() {
	//read from client, write to channel
	go func() {
		for {
			var err error
			data, err := util.ReadPacket(user.Conn)
			if err != nil {
				ts.TunServer.CloseClient(clientAddr)
				return
			}

			if ln := len(data); ln > 0 {
				if data, err = comp.UncompressGzip(data); err == nil && len(data)>0{
					if protocol, src, dst, err := header.GetBase(data); err == nil {
						ts.TunServer.WriteToChannel("tcp", clientAddr, data)
						fmt.Printf("[TcpServer][readFromClient] client:%v, protocol:%v, len:%v, src:%v, dst:%v\n", clientAddr, protocol, ln, src, dst)
					}
				}
			}
		}
	}()

	//read from channel, write to client
	go func() {
		for {
			data, err := ts.TunServer.ReadFromChannel(clientAddr)
			if err != nil {
				return
			}
			if ln := len(data); ln > 0 {
				if protocol, src, dst, err := header.GetBase(data); err == nil {
					if _, err := util.WritePacket(conn, comp.CompressGzip(data)); err != nil {
						ts.TunServer.CloseClient(clientAddr)
						return
					}
					fmt.Printf("[TcpServer][writeToClient] client:%v, protocol:%v, len:%v, src:%v, dst:%v\n", clientAddr, protocol, ln, src, dst)
				}
			}
		}
	}()
}

func (user *User) Close() {
	close(user.Chan)
	user.Conn.Close()
}