package header

import (
	"fmt"
)

func Get(data []byte) (proto string, src string, dst string, err error) {
	iph, udph, tcph := IPv4{}, UDP{}, TCP{}
	if len(data) < 20 {
		err = fmt.Errorf("Packet too short")
		return
	}

	iph.Unmarshal(data[:20])
	if iph.Protocol == uint8(UDPID) {
		proto = "udp"
		udph.Unmarshal(data[iph.HeaderLen():])
		src = fmt.Sprintf("%s:%d", IP2Str(iph.Src), udph.SrcPort)
		dst = fmt.Sprintf("%s:%d", IP2Str(iph.Dst), udph.DstPort)

	} else if iph.Protocol == uint8(TCPID) {
		proto = "tcp"
		tcph.Unmarshal(data[iph.HeaderLen():])
		src = fmt.Sprintf("%s:%d", IP2Str(iph.Src), tcph.SrcPort)
		dst = fmt.Sprintf("%s:%d", IP2Str(iph.Dst), tcph.DstPort)

	} else {
		err = fmt.Errorf("Protocol Unsupported: id=%d", iph.Protocol)
	}
	return

}
