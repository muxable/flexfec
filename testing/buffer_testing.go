// Function to extract relevant packets from buffer for L, D variant
// in case of row and col
package main

import (
	"flexfec/buffer"
	"flexfec/fec_header"
	"flexfec/recover"
	"flexfec/util"
	"fmt"

	"github.com/pion/rtp"
)

func test1() {
	BUFFER := make(map[buffer.Key]rtp.Packet)

	srcBlock := util.GenerateRTP(3, 4)
	util.PadPackets(&srcBlock)
	
	repairPackets := recover.GenerateRepairRowFec(&srcBlock, 4)


	for i :=0; i < len(srcBlock) ; i++ {
		util.PrintPkt(srcBlock[i])
		if (i + 1) % 4 == 0 {
			fmt.Println("------------------------------------------")
		}
	}

	/* 
		0 X  X  3
		4 5  X  7
	 	8 9 10 11
	*/


	// Assume packets received
	for i :=0; i < len(srcBlock) ; i++ {
		if i == 1 || i == 2 || i == 6 {
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

func test2() {
	BUFFER := make(map[buffer.Key]rtp.Packet)

	srcBlock := util.GenerateRTP(3, 4)
	util.PadPackets(&srcBlock)

	repairPackets := recover.GenerateRepairColFec(&srcBlock, 4, 1)

	for i := 0; i < len(srcBlock); i++ {
		util.PrintPkt(srcBlock[i])
		if (i+1)%4 == 0 {
			fmt.Println("------------------------------------------")
		}
	}

	/* 
		0 X  X  3
		4 5  X  7
	 	8 9 10 11
	*/
	
	// Assume packets received
	for i := 0; i < len(srcBlock); i++ {
		if i == 1 || i == 2 || i == 6 {
			continue
		}
		buffer.Update(BUFFER, srcBlock[i])
	}

	fmt.Println(BUFFER)

	// repair packets received
	for _, pkt := range repairPackets {
		fmt.Println("len :", len(buffer.Extract(BUFFER, pkt)))
		util.PrintPkt(pkt)
	}
}

func test3() {
	BUFFER := make(map[buffer.Key]rtp.Packet)

	// sender
	srcBlock := util.GenerateRTP(3, 3)
	util.PadPackets(&srcBlock)

	mask := uint16(3392) // 000110101000000
	sn_base := srcBlock[0].Header.SequenceNumber
	fmt.Println("sn base :", sn_base)
	

	// creating dummy mask repair packet for buffer testing
	fecheader := fech.NewFecHeaderFlexibleMask(false, false, false, false, 11, false, 56, 23434, 342334, sn_base, false, mask, 0, false, 0)
	payload := []byte{12, 23, 34, 54}

	repairPacket := rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			Padding:        false,
			Extension:      false,
			Marker:         false,
			PayloadType:    15,
			SequenceNumber: 23433,
			Timestamp:      54243243,
			SSRC:           2343244,
			CSRC:           []uint32{},
		},
		Payload: append(fecheader.Marshal(), payload...),
	}

	// receiver

	for _, pkt := range srcBlock {
		buffer.Update(BUFFER, pkt)
	}

	// Extracting protected packets by looking at the mask
	associatedPkts := buffer.ExtractMask(BUFFER, repairPacket)
	for _, pkt := range associatedPkts {
		util.PrintPkt(pkt)
	}

}

func main() {
	/* 
		Uncomment to test the respective test case
		test1() -> row fec with buffer
		test2() -> col fec with buffer
		test3() -> mask fec with buffer
	*/

	// test1()
	// test2()
	test3()
}
