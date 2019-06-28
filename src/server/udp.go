package server

import (
	"fmt"
	"net"
	"time"

	"comp"
	"header"
	"cache"
	"login"
	"config"
)
var UDPCHANBUFFERSIZE = 1024

type UdpServer struct {
	Addr      string
	UdpConn   *net.UDPConn
	LoginManager *login.LoginManager
	TunToConnChan chan string
	ConnToTunChan chan string
	RouteMap *cache.Cache
}

func NewUdpServer(cfg *config.Config, loginManager *login.LoginManager) (*UdpServer, error) {
	addr := cfg.ServerAddr
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
		TunToConnChan: make(chan string, UDPCHANBUFFERSIZE),
		ConnToTunChan: make(chan string, UDPCHANBUFFERSIZE),
		RouteMap: cache.NewCache(time.Minute * 10),
	}, nil
}

func (us *UdpServer) Start() error {
	fmt.Println("[UdpServer] started.")
	us.LoginManager.TunServer.StartClient("udp", us.ConnToTunChan, us.TunToConnChan)

	//from conn to tun
	go func(){
		defer func(){
			recover()
		}()

		data := make([]byte, us.LoginManager.TunServer.TunConn.GetMtu()*2)
		for {
			if n, caddr, err := us.UdpConn.ReadFromUDP(data); err == nil && n > 0 {
				uncmpData, errc := comp.UncompressGzip(data[:n])
				if errc != nil {
					continue
				}
				if protocol, src, dst, err := header.GetBase(uncmpData); err == nil {
					key := protocol + ":" + src + ":" + dst
					us.RouteMap.Put(key, caddr)
					us.ConnToTunChan <- string(uncmpData)
					fmt.Printf("[UdpServer][readFromClient] client:%v, protocol:%v, src:%v, dst:%v\n", caddr, protocol, src, dst)
				}
			}
		}

	}()

	//from tun to conn
	go func(){
		defer func(){
			recover()
		}()

		for {
			data, ok := <- us.TunToConnChan
			if ok {
				if protocol, src, dst, err := header.GetBase([]byte(data)); err == nil {
					key := protocol + ":" + dst + ":" + src
					clientAddrI := us.RouteMap.Get(key)
					if clientAddrI != nil {
						clientAddr := clientAddrI.(string)
						if add, err := net.ResolveUDPAddr("udp", clientAddr); err == nil {
							cmpData := comp.CompressGzip([]byte(data))
							us.UdpConn.WriteToUDP(cmpData, add)
							fmt.Printf("[UdpServer][writeToClient] client:%v, protocol:%v, src:%v, dst:%v\n", clientAddr, protocol, src, dst)
						}
					}
				}
			}
		}

	}()
	
	return nil
}

func (us *UdpServer) Stop() error {
	fmt.Println("[UdpServer] stopped.")

	go func(){
		defer func(){
			recover()
		}()

		close(us.TunToConnChan)
	}()

	go func(){
		defer func(){
			recover()
		}()

		close(us.ConnToTunChan)
	}()

	go func(){
		defer func(){
			recover()
		}()

		us.UdpConn.Close()
	}()
	return nil
}
