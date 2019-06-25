package header

import (
	"encoding/binary"
	"fmt"
)

type IPv4Pseudo struct {
	Src      uint32
	Dst      uint32
	Reserved uint8
	Protocol uint8
	Len      uint16
}

func (h IPv4Pseudo) String() string {
	res := `
{
	Src: %s,
	Dst: %s,
	Reserved: %d,
	Protocol: %d,
	Len: %d,
}
`
	return fmt.Sprintf(res, IP2Str(h.Src), IP2Str(h.Dst), h.Reserved, h.Protocol, h.Len)
}

func (h *IPv4Pseudo) HeaderLen() uint16 {
	return 12
}

func (h *IPv4Pseudo) LenBytes() uint16 {
	return h.Len + 12
}

func (h *IPv4Pseudo) Unmarshal(bs []byte) error {
	if len(bs) < 12 {
		return fmt.Errorf("too short")
	}
	h.Src = binary.BigEndian.Uint32(bs[0:4])
	h.Dst = binary.BigEndian.Uint32(bs[4:8])
	h.Reserved = uint8(0)
	h.Protocol = uint8(bs[9])
	h.Len = binary.BigEndian.Uint16(bs[10:12])
	return nil
}

func (h *IPv4Pseudo) Marshal() []byte {
	headerLen := int(h.HeaderLen()) 
	res := make([]byte, headerLen)
	binary.BigEndian.PutUint32(res[0:], h.Src)
	binary.BigEndian.PutUint32(res[4:], h.Dst)
	res[8] = byte(h.Reserved)
	res[9] = byte(h.Protocol)
	binary.BigEndian.PutUint16(res[10:], h.Len)
	return res
}
