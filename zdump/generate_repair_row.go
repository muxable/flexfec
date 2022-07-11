package recover

import (
	"flexfec/bitstring"
	fech "flexfec/fec_header"
	"math/rand"
	"fmt"
	"github.com/pion/rtp"
)

const (
	ssrc = uint32(2868272638)
)

func getBlockBitstring(packets *[]rtp.Packet) [][]byte {
	var bitStrings [][]byte

	for _, pkt := range *packets {
		bitStrings = append(bitStrings, bitstring.ToBitString(&pkt))
	}

	return bitStrings
}

// L>0, D=0 (in fecheader),, call with D=0 for Row fec and actual L(num cols)
func GenerateRepairRowFec(srcBlock *[]rtp.Packet, L int, D int) []rtp.Packet {

	num_packets := len(*srcBlock)

	var repairPackets []rtp.Packet
	// // Construct repair packet(another rtp packet)
	seqnum := uint16(rand.Intn(65535 - L))

	for i := 0; i < num_packets; i += L {
		fmt.Println("Row:", i)
		packets := (*srcBlock)[i : i+L]
		rowBitstrings := getBlockBitstring(&packets)

		fecBitString := bitstring.ToFecBitString(rowBitstrings)

		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)

		// associate src packet row with this repair packet
		fecheader.SN_base = (*srcBlock)[i].Header.SequenceNumber
		fecheader.L = uint8(L)
		fecheader.D = uint8(D)

		repairPacket := rtp.Packet{
			Header: rtp.Header{
				Version:        2,
				Padding:        true,
				Extension:      false,
				Marker:         false,
				PayloadType:    15,
				SequenceNumber: seqnum,
				Timestamp:      54243243,
				SSRC:           ssrc,
				CSRC:           []uint32{},
			},
			Payload: append(fecheader.Marshal(), repairPayload...),
		}

		repairPackets = append(repairPackets, repairPacket)
		seqnum++
	}

	return repairPackets
}
