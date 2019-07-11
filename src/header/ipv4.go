package header

import (
	"encoding/binary"
	"fmt"
)

type IPv4 struct {
	VerIHL   uint8
	Tos      uint8
	Len      uint16
	Id       uint16
	Offset   uint16
	TTL      uint8
	Protocol uint8
	Checksum uint16
	Src      uint32
	Dst      uint32
	Opt      []byte
}

func (h IPv4) String() string {
	res := `
{
	Version: %d,
	IHL: %d,
	Tos: %d,
	Len: %d,
	Id: %d,
	Offset: %d,
	TTL: %d,
	Protocol: %d,
	Checksum: %d,
	Src: %s,
	Dst: %s,
	Opt: %s,
}
}
`
	return fmt.Sprintf(res, h.VerIHL>>4, h.VerIHL&0xf, h.Tos,
		h.Len, h.Id, h.Offset, h.TTL,
		h.Protocol, h.Checksum,
		IP2Str(h.Src), IP2Str(h.Dst))

}

func (h *IPv4) HeaderLen() uint16 {
	return (uint16(h.VerIHL) & 0xf) * 4
}

func (h *IPv4) LenBytes() uint16 {
	return h.Len
}

func (h *IPv4) CalChecksum() uint16 {
	bs := h.MarshalRaw()
	bs[10] = 0
	bs[11] = 0
	s := uint32(0)
	for i := 0; i < int(h.HeaderLen()); i += 2 {
		s += uint32(binary.BigEndian.Uint16(bs[i : i+2]))
	}
	s = (s >> 16) + (s & 0xffff)
	return uint16(s ^ 0xffffffff)
}

func (h *IPv4) ResetChecksum() {
	h.Checksum = h.CalChecksum()
}

func (h *IPv4) Unmarshal(bs []byte) error {
	if len(bs) < 20 {
		return fmt.Errorf("too short")
	}
	h.VerIHL = uint8(bs[0])
	h.Tos = uint8(bs[1])
	h.Len = binary.BigEndian.Uint16(bs[2:4])
	h.Id = binary.BigEndian.Uint16(bs[4:6])
	h.Offset = binary.BigEndian.Uint16(bs[6:8])
	h.TTL = uint8(bs[8])
	h.Protocol = uint8(bs[9])
	h.Checksum = binary.BigEndian.Uint16(bs[10:12])
	h.Src = binary.BigEndian.Uint32(bs[12:16])
	h.Dst = binary.BigEndian.Uint32(bs[16:20])
	headerLen := int(h.HeaderLen()) 
	if headerLen > 20 {
		h.Opt = bs[20:headerLen]
	}else{
		h.Opt = []byte{}
	}
	return nil
}

func (h *IPv4) MarshalRaw() []byte {
	headerLen := int(h.HeaderLen()) 
	res := make([]byte, headerLen)
	res[0] = byte(h.VerIHL)
	res[1] = byte(h.Tos)
	binary.BigEndian.PutUint16(res[2:], h.Len)
	binary.BigEndian.PutUint16(res[4:], h.Id)
	binary.BigEndian.PutUint16(res[6:], h.Offset)
	res[8] = byte(h.TTL)
	res[9] = byte(h.Protocol)
	binary.BigEndian.PutUint16(res[10:], h.Checksum)
	binary.BigEndian.PutUint32(res[12:], h.Src)
	binary.BigEndian.PutUint32(res[16:], h.Dst)
	for i:=0; i<len(h.Opt); i++ {
		res[20 + i] = h.Opt[i]
	}
	return res
}

func (h *IPv4) Marshal() []byte {
	res := h.MarshalRaw()
	checkSum := h.CalChecksum()
	binary.BigEndian.PutUint16(res[10:], checkSum)
	return res
}
