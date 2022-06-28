package main

import(
	"fmt"
	"math/rand"
	"time"
	"github.com/pion/rtp"
)


func printPkt(pkt rtp.Packet) {
	fmt.Println("Header :",pkt.Header)
	fmt.Println("Payload :", pkt.Payload)
	fmt.Println("PaddingSize :",pkt.PaddingSize, "\n")
}

func PadPackets(srcBlock *[]rtp.Packet) {

    maxSize := -1
    n := len(*srcBlock)

    for i := 0; i < n; i++ {
        currSize := (*srcBlock)[i].MarshalSize()
        if maxSize < currSize {
            maxSize = currSize
        }
    }

    for i := 0; i < n; i++ {
    	size := (*srcBlock)[i].MarshalSize()

        if size != maxSize {
            // leaving the length field
        	leftOverPadBytes := maxSize - size - 1
        	padding := make([]byte, leftOverPadBytes)

            (*srcBlock)[i].PaddingSize = uint8(leftOverPadBytes + 1)
            (*srcBlock)[i].Payload = append((*srcBlock)[i].Payload, padding...)
        }
    }

}

func GenerateRTP(L int, D int) []rtp.Packet {
	rand.Seed(time.Now().UnixNano())

	pkts := []byte{
		0x90, 0xE0, 0x69, 0x8F, 0xD9, 0xC2, 0x93, 0xDA, 
		0x1C, 0x64, 0x27, 0x82, 0x45, 0xB1, 0x34, 0xF1, 
		0xFF, 0xFF, 0xAF, 0xFF, 0x98, 0x36, 0xbE, 0x88,
	}

	size := len(pkts)

	n := uint16(L * D)
	SN_base := uint16(rand.Intn(65535 - int(n)))
	ssrc := uint32(rand.Intn(4294967296))

	packets := []rtp.Packet{}

	for i:=uint16(0); i<n; i++ {
		endIndex := rand.Intn(size)

		packet := rtp.Packet{
			Header: rtp.Header{
				Version 		: 2,
				Padding 		: true,
				Extension		: false,
				Marker 			: false,
				PayloadType 	: 15,
				SequenceNumber	: SN_base + i,
				Timestamp		: 54243243,
				SSRC			: ssrc,
				CSRC			: []uint32{},
			},
			Payload: pkts[:endIndex],
		}
		packets = append(packets, packet)
	}

	return packets
}

func main(){
	packets := GenerateRTP(5,1)
	PadPackets(&packets)

	for _, packet := range packets {
		printPkt(packet)
	}

}