package tun

import (
	"sync"
	"fmt"
	"time"

	"cache"
	"header"
)

var INPUTCHANNELBUF = 1024
var OUTPUTCHANNELBUF = 1024
var TunInput = make(chan string, INPUTCHANNELBUF)
var TunOutputs = sync.Map{}
//Key: clientProtocol:clientIP:clientPort
var TunUdpOutput = make(chan string, OUTPUTCHANNELBUF)

type TunServer struct {
	TunConn Tun
	//Key: proto:src->dst   Value: clientProtocol:clientIP:clientPort
	ClientMap *cache.Cache
}

func NewTunServer(tname string, mtu int) (*TunServer, error){
	ts := &TunServer{
		ClientMap : cache.NewCache(time.Minute * 10),
	}
	if tun, err := NewLinuxTun(tname, mtu); err!=nil {
		return nil, err
	}else{
		ts.TunConn = tun
	}
	return ts, nil
}

func (ts *TunServer) Start() {
	go ts.toTun()
	go ts.fromTun()
	fmt.Println("tun server started")
}

func (ts *TunServer) Stop() {
	close(TunInput)
	TunOutputs.Range(func (_, value interface{}) bool {
		close(value.(chan string))
		return true
	})
	fmt.Println("tun server stopped")
}

func (ts *TunServer) GetClientAddr(key string) (protocol string, addr string) {
	s :=  ts.ClientMap.Get(key)
	if len(s) <= 4 {
		return "", ""
	}
	return s[:3], s[4:]
}

func (ts *TunServer) WriteToChannel(clientProtocol string, clientAddr string, data []byte){
	if proto, src, dst, err := header.GetBase(data); err == nil {
		key := proto + ":" + src + ":" + dst
		if ts.ClientMap.Get(key) == "" {
			ts.ClientMap.Put(key, clientProtocol + ":" + clientAddr)
		}

		key = clientProtocol + ":" + clientAddr
		if clientProtocol == "tcp" {
			if _, ok := TunOutputs.Load(key); !ok {
				TunOutputs.Store(key, make(chan string, OUTPUTCHANNELBUF))
			}
		}
		TunInput <- string(data)
		fmt.Printf("[TunServer][WriteToChannel] protocol:%v, src:%v, dst:%v\n", proto, src, dst)
	}
}

func (ts *TunServer) ReadFromUdpChannel() []byte {
	s := <- TunUdpOutput
	return []byte(s)
}

func (ts *TunServer) ReadFromChannel(clientAddr string) []byte {
	data := []byte{}
	if value, ok := TunOutputs.Load(clientAddr); ok {
		s := <- value.(chan string)
		data = []byte(s)
	}
	return data
}

func (ts *TunServer) toTun() {
	for {
		s := <- TunInput
		bs := []byte(s)
		for len(bs) > 0 {
			n, _ := ts.TunConn.Write(bs)
			bs = bs[n:]
		}
		fmt.Printf("[TunServer][toTun] data len: %v\n", len(s))
	}
}

func (ts *TunServer) fromTun() {
	for {
		data := make([]byte, ts.TunConn.GetMtu()*2)
		if n, err := ts.TunConn.Read(data); err==nil && n > 0 {
			if proto, src, dst, err := header.GetBase(data); err == nil {
				key := proto + ":" + dst + ":" + src
				fmt.Printf("[TunServer][fromTun] data len: %v, protocol:%v, src:%v, dst:%v\n", n, proto, src, dst)
				if caddr := ts.ClientMap.Get(key); caddr != "" {
					clientProtocol := caddr[:3]
					if clientProtocol == "tcp" {
						go func() {
							if tunOutput, ok := TunOutputs.Load(caddr); ok {
								tunOutput.(chan string) <- string(data[:n])
							}
						}()

					}else if clientProtocol == "udp" {
						TunUdpOutput <- string(data[:n])
					}
					fmt.Printf("[TunServer][fromTun] clientProtocol:%v, client:%v src:%v dst:%v proto:%v\n", clientProtocol, caddr, src, dst, proto)
				}
			}
		}
	}
}

