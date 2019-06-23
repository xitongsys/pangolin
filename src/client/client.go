package client

import (
	"fmt"
	"net"

	"comp"
	"tun"
	"header"
)

type PClient struct {
	ServerAdd string
	UdpConn   net.Conn
	TunConn   tun.Tun
}

func NewPClient(sadd string, tname string, mtu int) (*PClient, error) {
	conn, err := net.Dial("udp", sadd)
	if err != nil {
		return nil, err
	}
	tun, err := tun.NewLinuxTun(tname, mtu)
	if err != nil {
		return nil, err
	}

	return &PClient{
		ServerAdd: sadd,
		UdpConn:   conn,
		TunConn:   tun,
	}, nil
}

func (c *PClient) sendToServer() {
	data := make([]byte, c.TunConn.GetMtu()*2)
	for {
		if n, err := c.TunConn.Read(data); err == nil && n > 0 {
			if proto, src, dst, err := header.GetBase(data); err == nil {
				ipv4Header := header.IPv4{}
				ipv4Header.Unmarshal(data)
				ipv4Header.Src = header.Str2IP("10.0.0.12")
				newData := ipv4Header.Marshal()
				for i := 0; i<len(newData); i++ {
					data[i] = newData[i]
				}

				if proto == "tcp" {
					header.ReplaceTcpCheckSum(data)
				}else if proto == "udp" {
					header.ReplaceUdpCheckSum(data)
				}

				cmpData := comp.CompressGzip(data[:n])
				c.UdpConn.Write(cmpData)
				fmt.Printf("[send] Len:%d src:%s dst:%s proto:%s\n", n, src, dst, proto)
			}
		}
	}
}

func (c *PClient) recvFromServer() error {
	data := make([]byte, c.TunConn.GetMtu()*2)
	for {
		if n, err := c.UdpConn.Read(data); err == nil && n > 0 {
			uncmpData, err2 := comp.UncompressGzip(data[:n])
			if err2 != nil {
				continue
			}
			if proto, src, dst, err := header.GetBase(uncmpData); err == nil {
				ipv4Header := header.IPv4{}
				ipv4Header.Unmarshal(uncmpData)
				ipv4Header.Dst = header.Str2IP("10.0.0.12")
				newData := ipv4Header.Marshal()
				for i := 0; i<len(newData); i++ {
					uncmpData[i] = newData[i]
				}

				if proto == "tcp" {
					header.ReplaceTcpCheckSum(uncmpData)
				}else if proto == "udp" {
					header.ReplaceUdpCheckSum(uncmpData)
				}

				c.TunConn.Write(uncmpData)
				fmt.Printf("[recv] Len:%d src:%s dst:%s proto:%s\n", len(uncmpData), src, dst, proto)
			}
		}
	}
}

func (c *PClient) Start() error {
	go c.sendToServer()
	go c.recvFromServer()
	return nil
}

func (c *PClient) Stop() error {
	c.UdpConn.Close()
	c.TunConn.Close()
	return nil
}
