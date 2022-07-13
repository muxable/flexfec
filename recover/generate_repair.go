package recover

import (
	"flexfec/bitstring"
	fech "flexfec/fec_header"
	"fmt"
	"math/rand"

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

	return repairPacket
}

func NewRepairPacketLD(seqnum uint16, fecheader fech.FecHeaderLD, repairPayload []byte) rtp.Packet {
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

	return repairPacket
}

func getBlockBitstring(packets []rtp.Packet) [][]byte {
	var bitStrings [][]byte

	for _, pkt := range packets {
		bitStrings = append(bitStrings, bitstring.ToBitString(&pkt))
	}

	// fmt.Println(bitStrings)
	return bitStrings
}

func GenerateRepairFlex(srcBlock *[]rtp.Packet, mask uint16, optionalMask1 uint32, optionalMask2 uint64) rtp.Packet {

	var coveredPackets []rtp.Packet

	var SN_base uint16 = 0

	// mandatory mask : 14 to 0
	MSB_m := 14
	max := 14
	if optionalMask1 != 0 {
		MSB_m = 45
		max = 45
	}

	if optionalMask2 != 0 {
		MSB_m = 109
		max = 45
	}

	for MSB_m >= 0 {
		is_Covered := mask & (1 << MSB_m)

		if is_Covered > 0 {
			coveredPackets = append(coveredPackets, (*srcBlock)[max-MSB_m])
			if SN_base == 0 {
				SN_base = (*srcBlock)[14-MSB_m].SequenceNumber
			}
		}
	}

	seqnum := uint16(rand.Intn(65535))
	coveredBitstrings := getBlockBitstring(coveredPackets)

	fecBitstring := bitstring.ToFecBitString(coveredBitstrings)

	fecheader, repairPayload := fech.ToFecHeaderFlexibleMask(fecBitstring)

	// set snbase
	fecheader.SN_base = SN_base
	if max == 14 {
		fecheader.Mask = mask
		return NewRepairPacketFlex(seqnum, fecheader, repairPayload)
	} else if max == 45 {
		fecheader.Mask = mask
		fecheader.K1 = true
		fecheader.OptionalMask1 = optionalMask1
		return NewRepairPacketFlex(seqnum, fecheader, repairPayload)
	} else {
		//  both masks are used
		fecheader.Mask = mask
		fecheader.K1 = true
		fecheader.OptionalMask1 = optionalMask1
		fecheader.K2 = true
		fecheader.OptionalMask2 = optionalMask2
		return NewRepairPacketFlex(seqnum, fecheader, repairPayload)
	}

}

func GenerateRepairLD(srcBlock *[]rtp.Packet, L, D int) []rtp.Packet {

	var repairPackets []rtp.Packet
	if L == 0 && D == 0 {
		fmt.Println("ignore : future use only")
		return repairPackets

	} else if L > 0 && D == 0 {
		repairPackets = GenerateRepairRowFec(srcBlock, L,false)
		return repairPackets
	} else if L > 0 && D == 1 {
		rowRepairPackets, colRepairPackets := GenerateRepair2dFec(srcBlock, L, D)
		return append(rowRepairPackets, colRepairPackets...)

		// TODO
	} else if L > 0 && D > 1 {
		repairPackets = GenerateRepairColFec(srcBlock, L, D)
		return repairPackets
	} else {
		fmt.Println("NOT POSSIble")
		return repairPackets
	}

}

// L>0 , D=0
func GenerateRepairRowFec(srcBlock *[]rtp.Packet, L int, is2D bool) []rtp.Packet {

	var repairPackets []rtp.Packet

	// seqnum := uint16(rand.Intn(65535 - L))
	

	for i := 0; i < len(*srcBlock); i += L {
		packets := (*srcBlock)[i : i+L]
		rowBitstrings := getBlockBitstring(packets)

		fecBitString := bitstring.ToFecBitString(rowBitstrings)
		// fmt.Println("fecbtstr", fecBitString)

		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)

		// associate src packet row with this repair packet
		fecheader.SN_base = (*srcBlock)[i].Header.SequenceNumber
		fecheader.L = uint8(L)
		fecheader.D = uint8(0)
		if is2D{
			fecheader.D = uint8(1)
		}

		repairPacket := NewRepairPacketLD(seqnum, fecheader, repairPayload)
		// fmt.Println("repair")
		// util.PrintPkt(repairPacket)

		repairPackets = append(repairPackets, repairPacket)
		seqnum++
	}

	return repairPackets

}

//  L>0 & D>0
func GenerateRepairColFec(srcBlock *[]rtp.Packet, L, D int) []rtp.Packet {
	var repairPackets []rtp.Packet

	// seqnum := uint16(rand.Intn(65535 - L))
	
	packets := make([]rtp.Packet, D)
	for j := 0; j < L; j++ {
		for i := 0; i < D; i++ {
			// packets[i] = (*srcBlock)[i*D+j]
			packets[i] = (*srcBlock)[i*L+j]
		}

		rowBitstrings := getBlockBitstring(packets)

		fecBitString := bitstring.ToFecBitString(rowBitstrings)

		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)

		// associate src packet row with this repair packet
		fecheader.SN_base = packets[0].Header.SequenceNumber
		fecheader.L = uint8(L)
		fecheader.D = uint8(D)

		repairPacket := NewRepairPacketLD(seqnum, fecheader, repairPayload)

		repairPackets = append(repairPackets, repairPacket)
		seqnum++
	}

	return repairPackets
}

func GenerateRepair2dFec(srcBlock *[]rtp.Packet, L, D int) ([]rtp.Packet, []rtp.Packet) {

	is2D:=true
	rowRepairPackets := GenerateRepairRowFec(srcBlock, L,is2D)
	colRepairPackets := GenerateRepairColFec(srcBlock, L, D)
	

	return rowRepairPackets, colRepairPackets
}
