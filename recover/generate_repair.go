package recover

import (
	"flexfec/bitstring"
	fech "flexfec/fec_header"
	"math/rand"

	"github.com/pion/rtp"
)

//------------------------------------------------------------------------------------
// 1 d - 1 row
func GenerateRepair(srcBlock *[]rtp.Packet, L, D int) rtp.Packet {
	var bitStrings [][]byte

	for _, pkt := range *srcBlock {
		bitStrings = append(bitStrings, bitstring.ToBitString(&pkt))
	}

	fecBitString := bitstring.ToFecBitString(bitStrings)
	fecHeader, repairPayload := fech.ToFecHeader(fecBitString)

	fecHeader.SN_base = (*srcBlock)[0].Header.SequenceNumber
	fecHeader.L = uint8(L)
	fecHeader.D = uint8(D)

	SN_base := uint16(rand.Intn(65535 - 5))
	ssrc := uint32(rand.Intn(4294967296))

	repairPacket := rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			Padding:        true,
			Extension:      false,
			Marker:         false,
			PayloadType:    15,
			SequenceNumber: SN_base,
			Timestamp:      54243243,
			SSRC:           ssrc,
			CSRC:           []uint32{},
		},
		Payload: append(fecHeader.Marshal(), repairPayload...),
	}
	return repairPacket
}

//------------------------------------------------------------------------------------