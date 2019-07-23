package tun

type Tun interface {
	GetMtu() int

	Read(data []byte) (int, error)
	Write(data []byte) (int, error)
	Close() error
}
