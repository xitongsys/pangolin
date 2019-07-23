package header

func BuildUdpPacket(src string, dst string, data []byte) []byte {
	srcIp, srcPort := ParseAddr(src)
	dstIp, dstPort := ParseAddr(dst)

	ipv4Header := IPv4{
		VerIHL:   0x45,
		Tos:      0,
		Len:      uint16(20 + 8 + len(data)),
		Id:       0,
		Offset:   0,
		TTL:      255,
		Protocol: 0x11,
		Checksum: 0,
		Src:      Str2IP(srcIp),
		Dst:      Str2IP(dstIp),
	}
	ipv4Header.ResetChecksum()

	udpHeader := UDP{
		SrcPort:  uint16(srcPort),
		DstPort:  uint16(dstPort),
		Len:      uint16(8 + len(data)),
		Checksum: 0,
	}

	bs := []byte{}
	bs = append(bs, ipv4Header.Marshal()...)
	bs = append(bs, udpHeader.Marshal()...)
	bs = append(bs, data...)
	ReCalUdpCheckSum(bs)

	return bs
}

func BuildTcpPacket(src string, dst string, data []byte) []byte {
	srcIp, srcPort := ParseAddr(src)
	dstIp, dstPort := ParseAddr(dst)

	ipv4Header := IPv4{
		VerIHL:   0x45,
		Tos:      0,
		Len:      uint16(20 + 20 + len(data)),
		Id:       0,
		Offset:   0,
		TTL:      255,
		Protocol: 0x06,
		Checksum: 0,
		Src:      Str2IP(srcIp),
		Dst:      Str2IP(dstIp),
	}
	ipv4Header.ResetChecksum()

	tcpHeader := TCP{
		SrcPort:    uint16(srcPort),
		DstPort:    uint16(dstPort),
		Seq:        1,
		Ack:        1,
		Offset:     0x50,
		Flags:      0,
		Win:        0x10,
		Checksum:   0,
		UrgPointer: 0,
	}

	bs := []byte{}
	bs = append(bs, ipv4Header.Marshal()...)
	bs = append(bs, tcpHeader.Marshal()...)
	ReCalTcpCheckSum(bs)

	bs = append(bs, data...)
	return bs
}
