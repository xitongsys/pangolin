package raw

import (
	"fmt"
	"time"
	"net"
)
var CONNCHANBUFSIZE = 1024

type Conn struct {
	InputChan chan string
	OutputChan chan string
}

func NewConn() *Conn {
	return &Conn{
		InputChan: make(chan string, CONNCHANBUFSIZE),
		OutputChan: make(chan string, CONNCHANBUFSIZE),
	}
}

func (conn *Conn) Read(b []byte) (n int, err error) {
	defer func(){
		recover()
		n, err = -1, fmt.Errorf("closed")
	}()

	s := <- conn.InputChan
	ls, ln := len(s), len(b)
	l := ls
	if ln < ls {
		l = ln
	}
	sb := []byte(s)
	for i := 0; i < l; i++ {
		b[i] = sb[i]
	}
	return ls, nil	
}

func (conn *Conn) Write(b []byte) (n int, err error) {
	defer func(){
		recover()
		n, err = -1, fmt.Errorf("closed")
	}()

	conn.OutputChan <- string(b)
	return len(b), nil
}

func (conn *Conn) Close() error { 
	go func(){
		defer func(){
			recover()
		}()
		close(conn.InputChan)
	}()
	go func(){
		defer func(){
			recover()
		}()
		close(conn.OutputChan)
	}()
	return nil
}

func (conn *Conn) LocalAddr() net.Addr {
	return nil
}

func (conn *Conn) RemoteAddr() net.Addr {
	return nil
}

func (conn *Conn) SetDeadline(t time.Time) error {
	return nil
}

func (conn *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (conn *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}