package main

import (
	"flexfec/util"
	"fmt"
	"math/rand"
	"time"

	"github.com/pion/rtp"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	pkts := []byte{
		0x90, 0xE0, 0x69, 0x8F, 0xD9, 0xC2, 0x93, 0xDA,
		0x1C, 0x64, 0x27, 0x82, 0x45, 0xB1, 0x34, 0xF1,
		0xFF, 0xFF, 0xAF, 0xFF, 0x98, 0x36, 0xbE, 0x88,
	}

	csrc := []uint32{
		0x64, 0x27, 0x82, 0x45, 0x90, 0xE0,
	}

	size := len(pkts)
	endIndex := rand.Intn(size)
	fmt.Println("endIndex:", endIndex)
	csrcIndex := rand.Intn(5)
	fmt.Println("csrcIndex:", csrcIndex)

	isExtension := true

	packet := rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			Padding:        false,
			Extension:      false,
			Marker:         false,
			PayloadType:    15,
			SequenceNumber: 10000,
			Timestamp:      54243243,
			SSRC:           112313,
			CSRC:           csrc[:csrcIndex],
		},
		Payload: pkts[:endIndex],
	}

	if isExtension == true {
		packet.Header.Extension = true
		packet.Header.ExtensionProfile = 0x1000
		packet.Header.SetExtension(uint8(1), []byte{0xA5, 0x45})
		packet.Header.SetExtension(uint8(2), []byte{0xBB, 0xCC})
		packet.Header.SetExtension(uint8(3), []byte{0xCC, 0xCC, 0xCC})
	}

	fmt.Println(util.PrintPkt(packet))

	buf, _ := packet.Marshal()

	util.PrintBytes(buf)

	recoveredPacket := rtp.Packet{}

	recoveredPacket.Unmarshal(buf)
	fmt.Println(util.PrintPkt(recoveredPacket))

	fmt.Println(packet.SequenceNumber, packet.Header.SequenceNumber)
}
