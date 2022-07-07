package main

import (
	"fmt"
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/util"
	"github.com/pion/rtp"
)

const(
	Red        = "\033[31m"
	Green      = "\033[32m"
	White      = "\033[37m"
	Blue       = "\033[34m"
	L = 4
	D = 3
)

var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)

func main() {

	// Sender
	srcBlock := util.GenerateRTP(L, D)
	util.PadPackets(&srcBlock)

	repairPacketsRow := recover.GenerateRepairRowFec(&srcBlock, L, 0)
	var recievedPackets []rtp.Packet

	for i := 0; i < len(srcBlock); i++ {
		if i != 1 && i != 2 && i != 6 {
			fmt.Println(string(Green), "Sending a src packet")
			util.PrintPkt(srcBlock[i])
			recievedPackets = append(recievedPackets, srcBlock[i]);
		} else {
			fmt.Println(string(Red), "missing packet")
			util.PrintPkt(srcBlock[i])
		}
	}

	//receiver
	for _, pkt := range recievedPackets {
		buffer.Update(BUFFER, pkt)
	}

	for _, repairPacket := range repairPacketsRow {
		associatedSrcPackets := buffer.Extract(BUFFER, repairPacket)
		fmt.Println(string(White), "recovery")

		recoveredPacket, _ := recover.RecoverMissingPacket(&associatedSrcPackets, repairPacket)
		util.PrintPkt(recoveredPacket)
	}
}