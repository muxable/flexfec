package main

import (
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/util"
	"fmt"

	"github.com/pion/rtp"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	White = "\033[37m"
	Blue  = "\033[34m"
	L     = 5
	D     = 3
)

var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)

func main() {

	// Sender
	srcBlock := util.GenerateRTP(L, D)
	util.PadPackets(&srcBlock)

	repairPacketsCol := recover.GenerateRepairColFec(&srcBlock, L, D)
	var recievedPackets []rtp.Packet

	for i := 0; i < len(srcBlock); i++ {
		if i != 1 && i != 2 && i != 6 {
			fmt.Println(string(Green), "Sending a src packet")
			util.PrintPkt(srcBlock[i])
			recievedPackets = append(recievedPackets, srcBlock[i])
		} else {
			fmt.Println(string(Red), "missing packet")
			util.PrintPkt(srcBlock[i])
		}
	}

	//receiver
	for _, pkt := range recievedPackets {
		buffer.Update(BUFFER, pkt)
	}

	for _, repairPacket := range repairPacketsCol {
		associatedSrcPackets := buffer.Extract(BUFFER, repairPacket)
		fmt.Println("num packets:", len(associatedSrcPackets))
		fmt.Println(string(White), "recovery")

		recoveredPacket, _ := recover.RecoverMissingPacket(&associatedSrcPackets, repairPacket)
		util.PrintPkt(recoveredPacket)
	}
}
