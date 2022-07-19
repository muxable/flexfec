package recover

import (
	fech "flexfec/fec_header"
	"flexfec/bitstring"
	"encoding/binary"
	"flexfec/util"
	"fmt"

	"github.com/pion/rtp"
)

func SN_Missing(receivedBlock *[]rtp.Packet, SN_Sum int) int {
	SN_missing := 0

	for _, pkt := range *(receivedBlock) {
		SN_missing += int(pkt.Header.SequenceNumber)
	}

	return SN_Sum - SN_missing
}

func MissingPacket(receivedBlock *[]rtp.Packet, repairPacket rtp.Packet, SN_missing int, fecvariant string) rtp.Packet {
	ssrc := (*receivedBlock)[0].Header.SSRC
	ssrcBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(ssrcBuf, ssrc)

	bitsrings := bitstring.GetBlockBitstring(receivedBlock)
	fecBitString := []byte{}
	// remove 8-12 which contains snbase, L, D
	fecBitString = append(fecBitString, repairPacket.Payload[:8]...)
	fecBitString = append(fecBitString, repairPacket.Payload[12:]...)

	length := len(fecBitString)
	util.PadBitStrings(&bitsrings, length)

	bitstringsXOR := bitstring.ToFecBitString(&bitsrings)

	buf := make([]byte, length)

	for index, BYTE := range bitstringsXOR {
		buf[index] = BYTE ^ fecBitString[index]
	}

	recoveredBuf := []byte{}
	recoveredBuf = append(recoveredBuf, buf[:8]...)
	recoveredBuf = append(recoveredBuf, ssrcBuf...)
	recoveredBuf = append(recoveredBuf, buf[8:]...)

	lengthRecovery := binary.BigEndian.Uint16(buf[2:4])

	recoveredPkt := rtp.Packet{}
	recoveredPkt.Unmarshal(recoveredBuf[:lengthRecovery])

	recoveredPkt.Header.Version = 2
	recoveredPkt.Header.SequenceNumber = uint16(SN_missing)

	return recoveredPkt
}

// 1d 1 row
func RecoverMissingPacket(receivedBlock *[]rtp.Packet, repairPacket rtp.Packet) (rtp.Packet, int) {

	var fecheader fech.FecHeaderLD = fech.FecHeaderLD{}
	fecheader.Unmarshal(repairPacket.Payload[:12])

	SN_base := int(fecheader.SN_base)
	L := int(fecheader.L)
	D := int(fecheader.D)

	var SN_Sum int // sum of sequence numbers of row or col
	var length int // expected length of row or col

	if D == 0 || D == 1 { // row fec
		SN_Sum = SN_base*L + (L*(L-1))/2
		length = L
	} else { // col fec
		SN_Sum = (2*SN_base*D + (D-1)*L*D) / 2 // Arithematic progression
		length = D
	}

	missingSN := SN_Missing(receivedBlock, SN_Sum)
	lenReceivedBlock := len(*receivedBlock)

	if lenReceivedBlock < length {
		if (length - lenReceivedBlock) > 1 {
			fmt.Println("retransmission required")
			return rtp.Packet{}, -1
		}
		// recovery
		return MissingPacket(receivedBlock, repairPacket, missingSN, "LD"), 0
	}

	// successful,  No error
	fmt.Println("All packets transmitted correctly")
	return rtp.Packet{}, 1
}

func RecoverMissingPacketFlex(receivedBlock *[]rtp.Packet, repairPacket rtp.Packet) (rtp.Packet, int) {
	// recieved block consists of all marked packets only
	payload := repairPacket.Payload
	var fecheader fech.FecHeaderFlexibleMask = fech.FecHeaderFlexibleMask{}

	// will use only first 24 bits
	fecheader.Unmarshal(payload)

	SN_base := fecheader.SN_base
	SN_Sum := 0
	covered_count := 0

	// mandatory mask
	for i := 14; i >= 0; i-- {
		if (fecheader.Mask>>i)&1 == 1 {
			covered_count++
			SN_Sum += (int(SN_base) + 14 - i + 0) //start
		}
	}

	if fecheader.K1 {
		for i := 30; i >= 0; i-- {
			if (fecheader.OptionalMask1>>i)&1 == 1 {
				covered_count++
				SN_Sum += (int(SN_base) + 30 - i + 15)
			}
		}
	}

	if fecheader.K2 {
		for i := 63; i >= 0; i-- {
			if (fecheader.OptionalMask2>>i)&1 == 1 {
				covered_count++
				SN_Sum += (int(SN_base) + 63 - i + 46) //start
			}
		}
	}

	missingSN := SN_Missing(receivedBlock, int(SN_Sum))
	lenReceivedBlock := len(*receivedBlock)

	if lenReceivedBlock != covered_count {
		if (covered_count - lenReceivedBlock) > 1 {
			fmt.Println("retransmission required")
			return rtp.Packet{}, -1
		}

		// recovery
		return MissingPacket(receivedBlock, repairPacket, missingSN, "flexibleMask"), 0
	}

	// successful,  No error
	fmt.Println("All packets transmitted correctly")
	return rtp.Packet{}, 1
}
