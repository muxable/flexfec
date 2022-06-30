package main

import (
	"flexfec/recover"
	"flexfec/util"
	"fmt"
)

func main() {
	srcBlock := util.GenerateRTP(5, 1)
	util.PadPackets(&srcBlock)
	repairPacket := recover.GenerateRepair(&srcBlock, 5, 1)
	fmt.Println(repairPacket)
}
