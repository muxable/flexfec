package recover

import (
	"encoding/binary"
	fech "flexfec/fec_header"
	"fmt"

	"github.com/pion/rtp"
)

//---------------------------
func SN_Missing(srcBlock *[]rtp.Packet, SN_base int, L int, D int) int {

	SN_missing := 0
	for _, pkt := range *(srcBlock) {
		SN_missing += int(pkt.Header.SequenceNumber)
	}
	SN_Sum := SN_base*L + (L*(L-1))/2
	SN_missing = (SN_Sum - SN_missing)
	return SN_missing
}

//-----------------------------
func MissingPacket(srcBlock *[]rtp.Packet, repairPacket rtp.Packet, SN_missing int) rtp.Packet {
	//SN_missing := 0
	var ssrc uint32

	// Header recovery
	fecBitString := repairPacket.Payload
	fecHeaderBitString := fecBitString[:10]
	recoveredHeader := make([]byte, 10)
	var recoveredPadding byte
	for _, pkt := range *(srcBlock) {
		buf := make([]byte, 10)
		pkt.Header.MarshalTo(buf)

		length := len(pkt.Payload)
		buf[8] = uint8(0)
		buf[7] = uint8(0)
		binary.BigEndian.PutUint16(buf[8:10], uint16(length))

		ssrc = pkt.Header.SSRC

		for index, BYTE := range buf {
			recoveredHeader[index] ^= BYTE
		}
		recoveredPadding ^= pkt.PaddingSize // xor of all recieved pkts
	}

	// recovery the actual padding size
	recoveredPadding ^= fecBitString[len(fecBitString)-1]
	// fmt.Println(recoveredPadding)

	//SN_missing = SN_Missing(srcBlock, SN_Sum)

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
	recoveredPacket.PaddingSize = recoveredPadding

	// Payload recovery
	Y := int(binary.BigEndian.Uint16(recoveredHeader[2:4])) // Y -> 16 bit representation of (length - 12)
	recoveredPayload := make([]byte, Y)

	fecPaylodBitString := fecBitString[12 : 12+Y]

	for _, pkt := range *(srcBlock) {
		for i := 0; i < Y; i++ {
			recoveredPayload[i] ^= pkt.Payload[i]
		}
	}

	for index, BYTE := range fecPaylodBitString {
		recoveredPayload[index] ^= BYTE
	}

	recoveredPacket.Payload = recoveredPayload

	// recoveredPaddingSize := len((*srcBlock)[0].Payload) - Y
	// fmt.Println(recoveredPaddingSize)
	return recoveredPacket
}

// 1d 1 row
func RecoverMissingPacket(srcBlock *[]rtp.Packet, repairPacket rtp.Packet) (rtp.Packet, int) {

	fmt.Println("I GOT CALLED")

	var fecheader fech.FecHeaderLD = fech.FecHeaderLD{}
	fecheader.Unmarshal(repairPacket.Payload[:12])

	L := int(fecheader.L)
	SN_base := int(fecheader.SN_base)
	//Here D=0
	SN_missing := SN_Missing(srcBlock, SN_base, L, 0)
	lengthofsrcBlock := len(*srcBlock)
	if lengthofsrcBlock != L {
		if (L - lengthofsrcBlock) > 1 {
			// retransmission
			fmt.Println("retransmission")
			return rtp.Packet{}, -1
		}
		// recovery
		return MissingPacket(srcBlock, repairPacket, SN_missing), 0
	}

	// successful,  No error
	fmt.Println("All packets transmitted correctly")
	return rtp.Packet{}, 1
}
