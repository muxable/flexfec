package recover

import (
	fech "flexfec/fec_header"
	"fmt"
	"flexfec/bitstring"
	"github.com/pion/rtp"
)

const (
	ssrc = uint32(2868272638)
)

var seqnum uint16 = 20000

func NewRepairPacketFlex(seqnum uint16, fecheader fech.FecHeaderFlexibleMask, repairPayload []byte) rtp.Packet {
	repairPacket := rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			Padding:        false,
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

	return repairPacket
}

func NewRepairPacketLD(seqnum uint16, fecheader fech.FecHeaderLD, repairPayload []byte) rtp.Packet {
	repairPacket := rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			Padding:        false,
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

	return repairPacket
}


// func getMaskPacktets(srcBlock *[]rtp.Packet, mask uint64, bits int, start int) []rtp.Packet {
// 	var coveredPackets []rtp.Packet
// 	for i := bits; i >= 0; i-- {
// 		if start+bits-i >= len(*srcBlock) {
// 			fmt.Println("YEAD", start+bits-i)
// 			return coveredPackets
// 		}

// 		if (mask>>i)&1 == 1 {
// 			index := uint16(bits - i)
// 			coveredPackets = append(coveredPackets, (*srcBlock)[index+uint16(start)])
// 		}
// 	}

// 	return coveredPackets
// }

// func GenerateRepairFlex(srcBlock *[]rtp.Packet, mask uint16, optionalMask1 uint32, optionalMask2 uint64) rtp.Packet {

// 	var coveredPackets []rtp.Packet
// 	var SN_base uint16 = (*srcBlock)[0].SequenceNumber

// 	isK1 := false
// 	isK2 := false

// 	// mandatory mask : 14 to 0
// 	mandMaskPackets := getMaskPacktets(srcBlock, uint64(mask), 14, 0)
// 	coveredPackets = append(coveredPackets, mandMaskPackets...)

// 	// optional mask1
// 	if optionalMask1 != 0 {
// 		isK1 = true
// 		optionalMask1Packets := getMaskPacktets(srcBlock, uint64(optionalMask1), 30, 15)
// 		coveredPackets = append(coveredPackets, optionalMask1Packets...)
// 	}

// 	if optionalMask2 != 0 {

// 		isK2 = true
// 		optionalMask2Packets := getMaskPacktets(srcBlock, optionalMask2, 63, 46)
// 		coveredPackets = append(coveredPackets, optionalMask2Packets...)
// 	}

// 	coveredBitstrings := getBlockBitstring(coveredPackets)

// 	fecBitstring := bitstring.ToFecBitString(coveredBitstrings)

// 	fecheader, repairPayload := fech.ToFecHeaderFlexibleMask(fecBitstring)

// 	// set snbase
// 	fecheader.SN_base = SN_base
// 	fecheader.Mask = mask

// 	if isK1 {
// 		fecheader.K1 = true
// 		fecheader.OptionalMask1 = optionalMask1
// 	}

// 	if isK2 {
// 		fecheader.K2 = true
// 		fecheader.OptionalMask2 = optionalMask2
// 	}

// 	return NewRepairPacketFlex(seqnum, fecheader, repairPayload)
// }

func GenerateRepairLD(srcBlkBitstrs *[][]byte, L, D int, variant int, SN_Base uint16) []rtp.Packet {
	// variant 0 -> row, 1 -> col, 2 -> 2D
	var repairPackets []rtp.Packet

	if variant == 0 {
		repairPackets = GenerateRepairRowFec(srcBlkBitstrs, L, false, SN_Base)
	} else if variant == 1 {
		repairPackets = GenerateRepairColFec(srcBlkBitstrs, L, D, SN_Base)
	} else if variant == 2 {
		repairPackets = GenerateRepair2dFec(srcBlkBitstrs, L, D, SN_Base)
	} else {
		fmt.Println("invalid variant")
	}

	return repairPackets

}

// L>0 , D=0
func GenerateRepairRowFec(srcBlkBitstrs *[][]byte, L int, is2D bool, SN_Base uint16) []rtp.Packet {
	// src [[b0], [b1],........ [b11]]
	var repairPackets []rtp.Packet
	size := len((*srcBlkBitstrs)[0])

	// seqnum := uint16(rand.Intn(65535 - L))

	for i := 0; i < len(*srcBlkBitstrs); i += L {
		rowBitstrings := make([][]byte, L)
		for j := 0; j < L; j++ {
			rowBitstrings[j] = make([]byte, size)
			copy(rowBitstrings[j], (*srcBlkBitstrs)[i + j])
		}
		
		fecBitString := bitstring.ToFecBitString(rowBitstrings)
		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)
		
		// associate src packet row with this repair packet
		fecheader.SN_base = SN_Base + uint16(i)
		fecheader.L = uint8(L)
		fecheader.D = uint8(0)
		if is2D {
			fecheader.D = uint8(1)
		}

		repairPacket := NewRepairPacketLD(seqnum, fecheader, repairPayload)
		repairPackets = append(repairPackets, repairPacket)
		seqnum++
	}

	return repairPackets

}

//  L>0 & D>0
func GenerateRepairColFec(srcBlkBitstrs *[][]byte, L, D int, SN_Base uint16) []rtp.Packet {
	var repairPackets []rtp.Packet
	// src [[b0], [b1],........ [b11]]
	// L = 0 col bitstrings = [[b0], [b4], [b8]]

	// seqnum := uint16(rand.Intn(65535 - L))
	size := len((*srcBlkBitstrs)[0])

	for j := 0; j < L; j++ { 
		colbitstrings := make([][]byte, D)
		for i := 0; i < D; i++ { 
			colbitstrings[i] = make([]byte, size)
			copy(colbitstrings[i], (*srcBlkBitstrs)[i*L+j])
		}

		fecBitString := bitstring.ToFecBitString(colbitstrings)
		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)

		// associate src packet row with this repair packet
		fecheader.SN_base = SN_Base + uint16(j)
		fecheader.L = uint8(L)
		fecheader.D = uint8(D)

		repairPacket := NewRepairPacketLD(seqnum, fecheader, repairPayload)
		repairPackets = append(repairPackets, repairPacket)
		seqnum++
	}

	return repairPackets
}

func GenerateRepair2dFec(srcBlkBitstrs *[][]byte, L, D int, SN_Base uint16) []rtp.Packet {
	is2D := true
	rowRepairPackets := GenerateRepairRowFec(srcBlkBitstrs, L, is2D, SN_Base)
	colRepairPackets := GenerateRepairColFec(srcBlkBitstrs, L, D, SN_Base)

	return append(rowRepairPackets, colRepairPackets...)
}
