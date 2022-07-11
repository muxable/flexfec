package recover

import (
	"encoding/binary"
	fech "flexfec/fec_header"
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

func MissingPacket(receivedBlock *[]rtp.Packet, repairPacket rtp.Packet, SN_missing int) rtp.Packet {
	var ssrc uint32

	// Header recovery
	fecBitString := repairPacket.Payload
	fecHeaderBitString := fecBitString[:8]

	recoveredHeader := make([]byte, 8)
	var recoveredPadding byte

	for _, pkt := range *(receivedBlock) {
		buf, _ := pkt.Header.Marshal()
		buf = buf[:8]

		ssrc = pkt.Header.SSRC

		for index, BYTE := range buf {
			recoveredHeader[index] ^= BYTE
		}
		recoveredPadding ^= pkt.PaddingSize // xor of all recieved pkts
	}

	recoveredPadding ^= fecBitString[len(fecBitString)-1]

	for index, BYTE := range fecHeaderBitString {
		recoveredHeader[index] ^= BYTE
	}

	var recoveredPacket rtp.Packet

	recoveredPacket.Header.Version = 2
	recoveredPacket.Header.Padding = (recoveredHeader[0] >> 5 & 0x1) > 0
	recoveredPacket.Header.Extension = (recoveredHeader[0] >> 4 & 0x1) > 0
	recoveredPacket.Header.Marker = (recoveredHeader[1] >> 7 & 0x1) > 0
	recoveredPacket.Header.PayloadType = (recoveredHeader[1] & 0x7F)
	recoveredPacket.Header.SequenceNumber = uint16(SN_missing)
	recoveredPacket.Header.Timestamp = binary.BigEndian.Uint32(recoveredHeader[4:8])
	recoveredPacket.Header.SSRC = ssrc
	recoveredPacket.PaddingSize = 0

	// Payload recovery
	pkt := (*receivedBlock)[0]
	length := len(pkt.Payload) + len(pkt.CSRC) + len(pkt.Extensions)

	recoveredPayload := make([]byte, length)
	fecPaylodBitString := fecBitString[12 : 12+length]

	for _, pkt := range *(receivedBlock) {
		for i := 0; i < length; i++ {
			recoveredPayload[i] ^= pkt.Payload[i]
		}
	}

	for index, BYTE := range fecPaylodBitString {
		recoveredPayload[index] ^= BYTE
	}

	recoveredPacket.Payload = make([]byte, length-int(recoveredPadding))
	copy(recoveredPacket.Payload, recoveredPayload)

	return recoveredPacket
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

	if D == 0 { // row fec
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
		return MissingPacket(receivedBlock, repairPacket, missingSN), 0
	}

	// successful,  No error
	fmt.Println("All packets transmitted correctly")
	return rtp.Packet{}, 1
}
