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

func getBlockBitstring(packets []rtp.Packet) [][]byte {
	var bitStrings [][]byte

	for _, pkt := range packets {
		bitStrings = append(bitStrings, bitstring.ToBitString(&pkt))
	}

	// fmt.Println(bitStrings)
	return bitStrings
}

func GenerateRepairLD(srcBlock *[]rtp.Packet, L, D int) []rtp.Packet {

	var repairPackets []rtp.Packet
	if L == 0 && D == 0 {
		fmt.Println("ignore : future use only")
		return repairPackets

	} else if L > 0 && D == 0 {
		repairPackets = GenerateRepairRowFec(srcBlock, L)
		return repairPackets
	} else if L > 0 && D == 1 {
		rowRepairPackets, colRepairPackets := GenerateRepair2dFec(srcBlock, L)
		return append(rowRepairPackets, colRepairPackets...)
	} else if L > 0 && D > 1 {
		repairPackets = GenerateRepairColFec(srcBlock, L, D)
		return repairPackets
	} else {
		fmt.Println("NOT POSSIble")
		return repairPackets
	}

}

// L>0 , D=0
func GenerateRepairRowFec(srcBlock *[]rtp.Packet, L int) []rtp.Packet {

	var repairPackets []rtp.Packet

	seqnum := uint16(rand.Intn(65535 - L))
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

	seqnum := uint16(rand.Intn(65535 - L))

	packets := make([]rtp.Packet, D)
	for j := 0; j < L; j++ {
		for i := 0; i < D; i++ {
			packets[i] = (*srcBlock)[i*D+j]
		}

		rowBitstrings := getBlockBitstring(packets)

		fecBitString := bitstring.ToFecBitString(rowBitstrings)

		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)

		// associate src packet row with this repair packet
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

func GenerateRepair2dFec(srcBlock *[]rtp.Packet, L int) ([]rtp.Packet, []rtp.Packet) {

	var rowRepairPackets, colRepairPackets []rtp.Packet
	seqnum := uint16(rand.Intn(65535 - L))

	for i := 0; i < len(*srcBlock); i += L {

		packets := (*srcBlock)[i : i+L]
		rowBitstrings := getBlockBitstring(packets)

		fecBitString := bitstring.ToFecBitString(rowBitstrings)

		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)

		// associate src packet row with this repair packet
		fecheader.SN_base = (*srcBlock)[i].Header.SequenceNumber
		fecheader.L = uint8(L)
		fecheader.D = uint8(1)

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

		rowRepairPackets = append(rowRepairPackets, repairPacket)
		seqnum++
	}

	var orderedMap [][]rtp.Packet
	for i := 0; i < len(*srcBlock); i++ {
		orderedMap[i%L] = append(orderedMap[i%L], (*srcBlock)[i])
	}

	for i := 0; i < len(orderedMap); i++ {
		packets := orderedMap[i]
		rowBitstrings := getBlockBitstring(packets)

		fecBitString := bitstring.ToFecBitString(rowBitstrings)
		fecheader, repairPayload := fech.ToFecHeaderLD(fecBitString)

		// associate src packet row with this repair packet
		fecheader.SN_base = (*srcBlock)[i].Header.SequenceNumber
		fecheader.L = uint8(L)
		fecheader.D = uint8(1)

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

		colRepairPackets = append(colRepairPackets, repairPacket)
		seqnum++
	}

	return rowRepairPackets, colRepairPackets
}
