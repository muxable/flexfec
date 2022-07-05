package recover

import (
	"flexfec/bitstring"
	fech "flexfec/fec_header"
	"fmt"
	"math/rand"

	"github.com/pion/rtp"
)

// // L>0, D=0
func GenerateRepairColFec(srcBlock *[]rtp.Packet, L, D int) []rtp.Packet {
	num_packets := len(*srcBlock)

	// to map a row of packets, can use mapping in repair packet construction
	repairMap := make(map[int][]rtp.Packet)

	for i := 0; i < num_packets; i++ {
		// row of current packet
		c := i % L

		repairMap[(c + 1)] = append(repairMap[(c+1)], (*srcBlock)[i])
	}

	var repairPackets []rtp.Packet

	// Construct repair packet(another rtp packet)
	seqnum := uint16(rand.Intn(65535 - L))

	for col, packets := range repairMap {
		fmt.Println("col:", col)

		colBitstrings := getBlockBitstring(&packets)

		fecBitString := bitstring.ToFecBitString(colBitstrings)

		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)

		// associate src packet col with this repair packet
		fecheader.SN_base = packets[0].Header.SequenceNumber
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
