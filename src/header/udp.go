package header

import (
	"encoding/binary"
	"fmt"
)

type UDP struct {
	SrcPort  uint16
	DstPort  uint16
	Len      uint16
	Checksum uint16
}

func (h UDP) String() string {
	res := `
{
	SrcPort: %d,
	DstPort: %d,
	Len: %d,
	Checksum: %d,
}
`
	return fmt.Sprintf(res, h.SrcPort, h.DstPort, h.Len, h.Checksum)
}

func (h *UDP) HeaderLen() uint16 {
	return 8
}

func (h *UDP) LenBytes() uint16 {
	return h.Len
}

func (h *UDP) Marshal() []byte {
	res := make([]byte, 8)
	binary.BigEndian.PutUint16(res, h.SrcPort)
	binary.BigEndian.PutUint16(res[2:], h.DstPort)
	binary.BigEndian.PutUint16(res[4:], h.Len)
	binary.BigEndian.PutUint16(res[6:], h.Checksum)
	return res
}

func (h *UDP) Unmarshal(bs []byte) error {
	if len(bs) < 8 {
		return fmt.Errorf("too short")
	}
	h.SrcPort = binary.BigEndian.Uint16(bs[0:2])
	h.DstPort = binary.BigEndian.Uint16(bs[2:4])
	h.Len = binary.BigEndian.Uint16(bs[4:6])
	h.Checksum = binary.BigEndian.Uint16(bs[6:8])
	return nil
}
