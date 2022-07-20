package main

import (
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/util"
	"fmt"

	"github.com/pion/rtp"
)

func main() {

	queue := make(map[buffer.Key]rtp.Packet)

	srcBlock := util.GenerateRTP(L, D)
	util.PadPackets(&srcBlock)

	repairPacketsCol := recover.GenerateRepairColFec(&srcBlock, L, D)

	// assume recieving src packets first,   skip packet 2 and 11 (0 index)
	for i := 0; i < D; i++ {
		for j := 0; j < L; j++ {
			if i*L+j == 2 || i*L+j == 9 {
				continue
			}

			buffer.Update(queue, srcBlock[i*L+j])
		}
	}

	fmt.Println("Buffer len:", len(queue), "\n")

	// assume recieving repair packets
	for k := 0; k < len(repairPacketsCol); k++ {
		associatedSrcPackets := buffer.Extract(queue, repairPacketsCol[k])
		fmt.Println("repair packet", k, ":", len(associatedSrcPackets))
		fmt.Println()
		missingPkt, status := recover.RecoverMissingPacket(&associatedSrcPackets, repairPacketsCol[k])

		if status == 1 {
			fmt.Println("success")
		} else if status == -1 {
			fmt.Println("retransmission")
		} else {
			fmt.Println("recovered using Col FEC")
		}

		util.PrintPkt(missingPkt)
	}

}
