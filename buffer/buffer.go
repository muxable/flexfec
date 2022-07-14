// Function to extract relevant packets from buffer for L, D variant
// in case of row and col
package buffer

import (
	"encoding/binary"
	fech "flexfec/fec_header"
	"fmt"

	"github.com/pion/rtp"
)

type Key struct {
	sequenceNumber uint16
}

func Update(BUFFER map[Key]rtp.Packet, sourcePkt rtp.Packet) {
	src_seq := sourcePkt.SequenceNumber
	key := Key{
		sequenceNumber: src_seq,
	}
	BUFFER[key] = sourcePkt
}

// func ExtractLD(BUFFER map[Key]rtp.Packet, repairPacket rtp.Packet) []rtp.Packet {

// }

func readMask(BUFFER map[Key]rtp.Packet, receivedBlock *[]rtp.Packet, SN_base uint16, mask uint64, bits int, start uint16) {
	for i := bits; i >= 0; i-- {
		if (mask>>i)&1 == 1 {
			index := uint16(bits - i)
			_, isPresent := BUFFER[Key{SN_base + index + start}]
			if isPresent {
				*receivedBlock = append(*receivedBlock, BUFFER[Key{SN_base + index + start}])
			}
		}
	}
}

func ExtractMask(BUFFER map[Key]rtp.Packet, repairPacket rtp.Packet) []rtp.Packet {
	payload := repairPacket.Payload
	SN_base := binary.BigEndian.Uint16(payload[8:10])

	var maskheader fech.FecHeaderFlexibleMask = fech.FecHeaderFlexibleMask{}
	maskheader.Unmarshal(payload)

	var receivedBlock []rtp.Packet

	readMask(BUFFER, &receivedBlock, SN_base, uint64(maskheader.Mask), 14, 0)

	if maskheader.K1 {
		readMask(BUFFER, &receivedBlock, SN_base, uint64(maskheader.OptionalMask1), 30, 15)
	}

	if maskheader.K2 {
		readMask(BUFFER, &receivedBlock, SN_base, maskheader.OptionalMask2, 63, 15+31)
	}

	return receivedBlock
}

func Extract(BUFFER map[Key]rtp.Packet, repairPacket rtp.Packet) []rtp.Packet {
	SN_base := binary.BigEndian.Uint16(repairPacket.Payload[8:10])
	L := repairPacket.Payload[10]
	D := repairPacket.Payload[11]
	fmt.Println("SNbase,L,D :", SN_base, L, D)

	var receivedBlock []rtp.Packet

	if D == 0 || D == 1 {
		// Row fec
		for i := uint16(0); i < uint16(L); i++ {
			_, isPresent := BUFFER[Key{SN_base + i}]
			if isPresent {
				// fmt.Println(BUFFER[Key{SN_base + i}].Payload)
				receivedBlock = append(receivedBlock, BUFFER[Key{SN_base + i}])
			}
		}
	} else if D > 1 {
		// Col fec
		for i := uint16(0); i < uint16(D); i++ {
			_, isPresent := BUFFER[Key{SN_base + i*uint16(L)}]
			if isPresent {
				// fmt.Println(BUFFER[Key{SN_base + i}].Payload)
				receivedBlock = append(receivedBlock, BUFFER[Key{SN_base + i*uint16(L)}])
			}
		}
	} else {
		// NEED TO EXTEND
	}

	return receivedBlock
}
