package main

import (
	"flexfec/recover"
	"flexfec/util"
	"fmt"
	"math/rand"
	"time"

	"github.com/pion/rtp"
)

func main() {
	srcBlock := util.GenerateRTP(5, 1)

	util.PadPackets(&srcBlock)

	repairPacket := recover.GenerateRepair(&srcBlock, 5, 1)
	fmt.Println(repairPacket)

	// all except one
	rcvBlock := make([]rtp.Packet, len(srcBlock)-1)

	rand.Seed(time.Now().Unix())
	idx := rand.Intn(len(srcBlock))

	for id, val := range srcBlock {
		if id == idx {
			continue
		}

		rcvBlock[id] = val
	}

	miss_header, miss_payload := recover.RecoverMissingPacket(&rcvBlock, repairPacket)
}
