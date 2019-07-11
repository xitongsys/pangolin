package header


func BuildTcpPacket(src string, dst string, data []byte) []byte {
	ipv4 := IPv4{
		VerIHL: 0x45,
		Tos: 
	}
} 