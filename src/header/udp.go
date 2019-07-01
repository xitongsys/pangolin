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

func ReCalUdpCheckSum(bs []byte) error {
	if len(bs) < 20 {
		return fmt.Errorf("too short")
	}
	ipvh, ipps := IPv4{}, IPv4Pseudo{}
	if err := ipvh.Unmarshal(bs); err!= nil {
		return err
	}
	ipps.Src = ipvh.Src
	ipps.Dst = ipvh.Dst
	ipps.Reserved = 0
	ipps.Protocol = ipvh.Protocol
	ipps.Len = ipvh.LenBytes() - ipvh.HeaderLen()

	ippsbs := ipps.Marshal()
	udpbs := bs[ipvh.HeaderLen():ipvh.LenBytes()]
	udpbs[6] = 0
	udpbs[7] = 0

	if len(udpbs) % 2 == 1 {
		udpbs = append(udpbs, byte(0))
	}

	s := uint32(0)
	for i := 0; i<len(ippsbs); i+=2 {
		s +=  uint32(binary.BigEndian.Uint16(ippsbs[i : i+2]))
	}
	for i := 0; i<len(udpbs); i+=2 {
		s +=  uint32(binary.BigEndian.Uint16(udpbs[i : i+2]))
	}
	
	for (s>>16) > 0 {
		s = (s>>16) + (s&0xffff)
	}
	binary.BigEndian.PutUint16(udpbs[6:], ^uint16(s))
	return nil
}
