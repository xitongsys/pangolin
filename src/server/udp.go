package server
/*
import (
	"fmt"
	"net"

	"comp"
	"header"
	"login"
)

type UdpServer struct {
	Addr      string
	UdpConn   *net.UDPConn
	LoginManager *login.LoginManager
}

func NewUdpServer(addr string, loginManager *login.LoginManager) (*UdpServer, error) {
	add, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid address", addr)
	}

	conn, err := net.ListenUDP("udp", add)
	if err != nil {
		return nil, err
	}

	return &UdpServer{
		Addr:      addr,
		UdpConn:   conn,
		LoginManager: loginManager,
	}, nil
}

func (us *UdpServer) writeToClient() {
	for {
		if data := us.TunServer.ReadFromUdpChannel(); len(data) > 0 {
			if protocol, src, dst, err := header.GetBase(data); err == nil {
				key := protocol + ":" + dst + ":" + src
				if cprotocal, caddr := us.TunServer.GetClientAddr(key); cprotocal!= "" && caddr != "" {
					if add, err := net.ResolveUDPAddr("udp", caddr); err == nil {
						cmpData := comp.CompressGzip(data)
						us.UdpConn.WriteToUDP(cmpData, add)
						fmt.Printf("[UdpServer][writeToClient] client:%v, protocol:%v, src:%v, dst:%v\n", caddr, protocol, src, dst)
					}
				}
			}
		}
	}
}

func (us *UdpServer) readFromClient() {
	data := make([]byte, us.LoginManager.TunServer.TunConn.GetMtu()*2)
	for {
		if n, caddr, err := us.UdpConn.ReadFromUDP(data); err == nil && n > 0 {
			if user, ok := us.LoginManager.GetUser("udp:" + caddr); ok {
				uncmpData, errc := comp.UncompressGzip(data[:n])
				if errc != nil {
					continue
				}
				if protocol, src, dst, err := header.GetBase(uncmpData); err == nil {
					
					us.TunServer.WriteToChannel("udp", caddr.String(), uncmpData)
					fmt.Printf("[UdpServer][readFromClient] client:%v, protocol:%v, src:%v, dst:%v\n", caddr, protocol, src, dst)
				}
			}
		}
	}
}

func (us *UdpServer) Start() error {
	fmt.Println("[UdpServer] started.")
	go us.writeToClient()
	go us.readFromClient()
	return nil
}

func (us *UdpServer) Stop() error {
	fmt.Println("[UdpServer] stopped.")
	us.UdpConn.Close()
	return nil
}
*/