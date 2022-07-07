// Function to extract relevant packets from buffer for L, D variant
// in case of row and col
package main

import(
	"fmt"
	"flexfec/util"
	"flexfec/recover"
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

func main() {
	buffer := make(map[Key]rtp.Packet)

	srcBlock := util.GenerateRTP(3, 4)
	util.PadPackets(&srcBlock)
	repairPackets := recover.GenerateRepairRowFec(&srcBlock, 4, 0)


	for i :=0; i < len(srcBlock) ; i++ {
		util.PrintPkt(srcBlock[i])
		if (i + 1) % 4 == 0 {
			fmt.Println("------------------------------------------")
		}
	}

	// 3 X 4
		// 0 X X 3
		// 4 5 X 7
		// 8 9 10 11

	
	// Assume packets received
	for i :=0; i < len(srcBlock) ; i++ {
		if(i == 1 || i == 2 || i == 6) {
			continue
		}
		Update(buffer, srcBlock[i])
	}

	fmt.Println(buffer)

	// repair packets received
	for _, pkt := range repairPackets {
		fmt.Println("len :",len(Extract(buffer, pkt)))
		util.PrintPkt(pkt)
	}

}

