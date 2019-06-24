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

type TunServer struct {
	TunConn Tun
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

func (ts *TunServer) WriteToChannel(clientAddr string, data []byte){
	if _, src, dst, err := header.GetBase(data); err == nil {
		key := src + "->" + dst
		if _, ok := TunOutputs.Load(key); !ok {
			TunOutputs.Store(key, make(chan string))
		}
		TunInput <- string(data)
	}
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
	}
}

func (ts *TunServer) fromTun() {
	for {
		data := make([]byte, ts.TunConn.GetMtu()*2)
		if n, err := ts.TunConn.Read(data); err==nil && n > 0 {
			if proto, src, dst, err := header.GetBase(data); err == nil {
				if caddr := ts.ClientMap.Get(dst + "->" + src); caddr != "" {
					go func() {
						if tunOutput, ok := TunOutputs.Load(caddr); ok {
							tunOutput.(chan string) <- string(data[:n])
							fmt.Printf("[send] client:%s src:%s dst:%s proto:%s\n", caddr, src, dst, proto)
						}
					}()
				}
			}
		}
	}
}

