package raw

import (
	"os"
	"fmt"
	"net"
	"syscall"
	"sync"

	"header"
	
)

var RAWSERVERBUFSIZE = 65535

type RawServer struct {
	ClientMap *sync.Map
}

func NewRawServer() *RawServer {
	rs := &RawServer{
		ClientMap: &sync.Map{},
	}
}

func (rs *RawServer) Start() error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return err
	}
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	buf := make([]byte, RAWSERVERBUFSIZE)
	go func() {
		for {
			if n, err := f.Read(buf); err == nil && n > 0 {
				if protocol, src, dst, err := header.GetBase(buf[:n]); err == nil {
					key := protocol + ":" + src
					if rawClienti, ok := rs.ClientMap.Load(key); ok {
						rawClient := rawClienti.(*RawClient)
						if _, _, _, _, data, err := header.Get(buf[:n]); err == nil {
							rawClient.RawConn.InputChan <- string(data)
						}
					}
				}
			}
		}
	}()
	return nil
}

func (rs *RawServer) CreateClient(client string, netConn net.Conn) net.Conn {
	rawClient := NewRawClient(client, NewConn(), netConn)
	rs.ClientMap.Store(client, rawClient)
	go func(){
		for {
			s := <- conn.OutputChan
			
		}
	}()
	return conn
}