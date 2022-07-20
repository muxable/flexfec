package main

import (
	"flexfec/recover"
	"flexfec/util"
	"fmt"

	"github.com/pion/rtp"
)

func main() {
	srcBlock := util.GenerateRTP(5, 1)
	util.PadPackets(&srcBlock)

	util.PrintPkt(srcBlock[2])

	// removing srcBlock[2] in new Block
	var newBlock []rtp.Packet
	newBlock = append(newBlock, srcBlock[:2]...)
	newBlock = append(newBlock, srcBlock[3:]...)

	fmt.Println()

	repairPacket := recover.GenerateRepair(&srcBlock, 5, 1)
	recoveredPacket, _ := recover.RecoverMissingPacket(&newBlock, repairPacket)
	util.PrintPkt(recoveredPacket)
}
