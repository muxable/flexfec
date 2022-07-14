package main

import(
	"fmt"
	"flexfec/recover"
	"flexfec/bitstring"
	"flexfec/buffer"
	"flecfec/util"
	"gtihub.com/pion/rtp"
)

BUFFER := make(map[buffer.Key]rtp.Packet)

func main() {

	//sender
	srcBlock := util.GenerateRTP(10, 10)
	util.PadPackets(&packets)

	mask := uint16(36160) // 1|000110101000000 16 bit 3,4,6,8
	optionalmask1 := uint32(3229756930) // 1|1000000100000100010111000000010 32 bit 15,22,28,32,34,35,36,44,
	optionalmask2 := uint64(13871700391609118210) // 1100000010000010001011100000001011000000100000100010111000000010 64 bit 46,47,54,60,64,66,67,68,76,78,79,86,92,96,98,99

	reapirPkt := GenerateRepairFlex(srcBlock, mask, optionalmask1, optionalmask2)

	
	var receivedPackets []rtp.packet

	for i, pkt := range srcBlock {
		if i != 3 {
			receivedPackets = append(receivedPackets, pkt)
		} 
	}

	//receiver 

	for _, pkt := range receivedPackets {
		buffer.Update(BUFFER, pkt)
	}

	associatedPkts := buffer.ExtractMask(BUFFER, repairPacket)
	



}

