package util

import (
	"fmt"
	"github.com/pion/rtp"
)

func PrintPkt(pkt rtp.Packet) string {
	result := fmt.Sprintf("Header : %v\n", pkt.Header)
	result += fmt.Sprintf("Payload : %v\n", pkt.Payload)
	result += fmt.Sprintf("PaddingSize : %d\n\n", pkt.PaddingSize)
	return result
}

