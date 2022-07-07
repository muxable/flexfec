// Function to extract relevant packets from buffer for L, D variant
// in case of row and col
package buffer

import(
	"fmt"
	"github.com/pion/rtp"
	"encoding/binary"
)

type Key struct{
	sequenceNumber uint16
}


func Update(buffer map[Key]rtp.Packet, sourcePkt rtp.Packet) {
	src_seq := sourcePkt.SequenceNumber
	key := Key{
		sequenceNumber: src_seq,
	}
	buffer[key] = sourcePkt
}


func Extract(buffer map[Key]rtp.Packet, repairPacket rtp.Packet) []rtp.Packet{
	SN_base := binary.BigEndian.Uint16(repairPacket.Payload[8:10]) 
	L := repairPacket.Payload[10]
	D := repairPacket.Payload[11]
	fmt.Println("SNbase,L,D :",SN_base, L, D)

	var receivedBlock []rtp.Packet

	for i := uint16(0); i < uint16(L); i++ {
		if D == 0 {
			_, isPresent := buffer[Key{SN_base + i}]
			if isPresent {
				receivedBlock = append(receivedBlock, buffer[Key{SN_base + i}])
			}
		} else {
			_, isPresent := buffer[Key{SN_base + i*uint16(L)}]
			if isPresent {
				receivedBlock = append(receivedBlock, buffer[Key{SN_base + i*uint16(L)}])
			}
		}
	}

	return receivedBlock
}

