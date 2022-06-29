package main

import (
	"flexfec/bitstring"
	fech "flexfec/fec_header"
	"flexfec/util"
	"fmt"
	"math/rand"

	"github.com/pion/rtp"
)

//------------------------------------------------------------------------------------
// 1 d - 1 row
func generateRepair(srcBlock *[]rtp.Packet, L, D int) rtp.Packet {
	var bitStrings [][]byte

	for _, pkt := range *srcBlock {
		bitStrings = append(bitStrings, bitstring.ToBitString(&pkt))
	}

	fecBitString := bitstring.ToFecBitString(bitStrings)
	fecHeader, repairPayload := fech.UnmarshalFec(fecBitString)

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

func main() {
	srcBlock := util.GenerateRTP(5, 1)
	util.PadPackets(&srcBlock)
	repairPacket := generateRepair(&srcBlock, 5, 1)
	fmt.Println(repairPacket)
}
