package tun

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	IFF_NO_PI = 0x10
	IFF_TUN   = 0x01
	IFF_TAP   = 0x02
	TUNSETIFF = 0x400454CA
)

type tun struct {
	Mtu  int
	Name string
	fd   *os.File
}

func NewLinuxTun(name string, mtu int) (Tun, error) {
	fd, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	ifr := make([]byte, 18)
	copy(ifr, []byte(name))
	ifr[16] = IFF_TUN
	ifr[17] = IFF_NO_PI

	_, _, errn := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd.Fd()), uintptr(TUNSETIFF),
		uintptr(unsafe.Pointer(&ifr[0])))
	if errn != 0 {
		return nil, fmt.Errorf("ioctl open tun failed")
	}
	syscall.SetNonblock(int(fd.Fd()), false)

	return &tun{
		Mtu:  mtu,
		Name: name,
		fd:   fd,
	}, nil
}

func (t *tun) Read(data []byte) (int, error) {
	return t.fd.Read(data)
}

func (t *tun) Write(data []byte) (int, error) {
	return t.fd.Write(data)
}

func (t *tun) Close() error {
	return t.fd.Close()
}

func (t *tun) GetMtu() int {
	return t.Mtu
}
