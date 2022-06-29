package util

import (
	"fmt"
	"github.com/pion/rtp"
)

func PrintPkt(pkt rtp.Packet) {
	fmt.Println("Header :", pkt.Header)
	fmt.Println("Payload :", pkt.Payload)
	fmt.Println("PaddingSize :", pkt.PaddingSize, "\n")
}

