// Function to extract relevant packets from buffer for L, D variant
// in case of row and col
package main

import(
	"fmt"
	"flexfec/util"
	"flexfec/recover"
	"flexfec/buffer"
	"github.com/pion/rtp"
)

func main() {
	BUFFER := make(map[buffer.Key]rtp.Packet)

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
		buffer.Update(BUFFER, srcBlock[i])
	}

	fmt.Println(BUFFER)

	// repair packets received
	for _, pkt := range repairPackets {
		fmt.Println("len :",len(buffer.Extract(BUFFER, pkt)))
		util.PrintPkt(pkt)
	}

}

