package raw

import (
	"net"
)

type RawClient struct {
	Client string
	RawConn *Conn
	NetConn net.Conn
}

func NewRawClient(client string, rawConn *Conn, netConn net.Conn) *RawClient {
	return &RawClient{
		Client: client,
		RawConn: rawConn,
		NetConn: netConn,
	}
}