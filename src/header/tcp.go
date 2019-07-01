package header

import (
	"encoding/binary"
	"fmt"
)

type TCP struct {
	SrcPort    uint16
	DstPort    uint16
	Seq        uint32
	Ack        uint32
	Offset     uint8
	Flags      uint8
	Win        uint16
	Checksum   uint16
	UrgPointer uint16
	Opt        uint32
}

func (h TCP) String() string {
	res := `
{
	SrcPort: %d,
	DstPort: %d,
	Seq: %d,
	Ack: %d,
	Offset: %d,
	Flags: %d,
	Win: %d,
	Checksum: %d,
	UrgPointer: %d,
}
`
	return fmt.Sprintf(res, h.SrcPort, h.DstPort,
		h.Seq, h.Ack,
		h.Offset, h.Flags, h.Win,
		h.Checksum, h.UrgPointer)
}

func (h *TCP) HeaderLen() uint16 {
	return (uint16(h.Offset) >> 4) * 4
}

func (h *TCP) Marshal() []byte {
	res := make([]byte, 20)
	binary.BigEndian.PutUint16(res, h.SrcPort)
	binary.BigEndian.PutUint16(res[2:], h.DstPort)
	binary.BigEndian.PutUint32(res[4:], h.Seq)
	binary.BigEndian.PutUint32(res[8:], h.Ack)
	res[12] = byte(h.Offset)
	res[13] = byte(h.Flags)
	binary.BigEndian.PutUint16(res[14:], h.Win)
	binary.BigEndian.PutUint16(res[16:], h.Checksum)
	binary.BigEndian.PutUint16(res[18:], h.UrgPointer)
	return res
}

func (h *TCP) Unmarshal(bs []byte) error {
	if len(bs) < 20 {
		return fmt.Errorf("too short")
	}
	h.SrcPort = binary.BigEndian.Uint16(bs[0:2])
	h.DstPort = binary.BigEndian.Uint16(bs[2:4])
	h.Seq = binary.BigEndian.Uint32(bs[4:8])
	h.Ack = binary.BigEndian.Uint32(bs[8:12])
	h.Offset = uint8(bs[12])
	h.Flags = uint8(bs[13])
	h.Win = binary.BigEndian.Uint16(bs[14:16])
	h.Checksum = binary.BigEndian.Uint16(bs[16:18])
	h.UrgPointer = binary.BigEndian.Uint16(bs[18:20])

	return nil
}

func ReCalTcpCheckSum(bs []byte) error {
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
	tcpbs := bs[ipvh.HeaderLen():ipvh.LenBytes()]
	tcpbs[16] = 0
	tcpbs[17] = 0

	if len(tcpbs) % 2 == 1 {
		tcpbs = append(tcpbs, byte(0))
	}

	s := uint32(0)
	for i := 0; i<len(ippsbs); i+=2 {
		s +=  uint32(binary.BigEndian.Uint16(ippsbs[i : i+2]))
	}
	for i := 0; i<len(tcpbs); i+=2 {
		s +=  uint32(binary.BigEndian.Uint16(tcpbs[i : i+2]))
	}
	
	for (s>>16) > 0 {
		s = (s>>16) + (s&0xffff)
	}
	binary.BigEndian.PutUint16(tcpbs[16:], ^uint16(s))
	return nil
}
