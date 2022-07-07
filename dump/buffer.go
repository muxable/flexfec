// Function to extract relevant packets from buffer for L, D variant
// in case of row and col

import(
	"fmt"
	"github.com/pion/rtp"
)

type Key struct{
	sequenceNumber uint16
	// SSRC uint32
}

// Assumption -> repair packet is LD variant
func extract(buffer map[key]rtp.Packet, repairPacket rtp.Packet) []rtp.Packet{
	SN_base := repairPacket.Payload[8:10]
	L := repairPacket.Payload[10]
	D := repairPacket.Payload[11]

	var receivedBlock []rtp.Packet

	for i := uint16(0) ; i < L; i++ {
		if D == 0 {
			receivedBlock = append(receivedBlock, buffer[Key{SN_base + i}])
		} else {
			receivedBlock = append(receivedBlock, buffer[Key{SN_base + i*L}])
		}
	}

	return receivedBlock
}

