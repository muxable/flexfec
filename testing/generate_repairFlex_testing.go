package main

import (
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/util"
	"fmt"

	"github.com/pion/rtp"
)

var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)

func main() {

	//sender
	srcBlock := util.GenerateRTP(10, 10)
	util.PadPackets(&srcBlock)

	mask := uint16(36160)                         // 1|000110101000000 16 bit 3,4,6,8
	optionalmask1 := uint32(3229756930)           // 1|1000000100000100010111000000010 32 bit 15,22,28,32,34,35,36,44,
	optionalmask2 := uint64(13871700391609117696) // 1100000010000010001011100000001011000000100000100010110000000000 64 bit 46,47,54,60,64,66,67,68,76,78,79,86,92,96,98,99

	repairPkt := recover.GenerateRepairFlex(&srcBlock, mask, optionalmask1, optionalmask2)
	fmt.Println("repair packet")
	util.PrintPkt(repairPkt)

	var receivedPackets []rtp.Packet

	for i, pkt := range srcBlock {
		if i != 54 {
			receivedPackets = append(receivedPackets, pkt)
		} else {
			fmt.Println("Missing packet")
			util.PrintPkt(pkt)
		}
	}

	//receiver

	for _, pkt := range receivedPackets {
		buffer.Update(BUFFER, pkt)
	}

	associatedPkts := buffer.ExtractMask(BUFFER, repairPkt)
	recoveredPacket, status := recover.RecoverMissingPacketFlex(&associatedPkts, repairPkt)

	if status == 0 {
		util.PrintPkt(recoveredPacket)
	}
}
