// Function to extract relevant packets from buffer for L, D variant
// in case of row and col
package main

import (
	"flexfec/bitstring"
	"flexfec/buffer"
	"flexfec/fec_header"
	"flexfec/recover"
	"flexfec/util"
	"fmt"

	"github.com/pion/rtp"
)

func test1() {
	BUFFER := make(map[buffer.Key]rtp.Packet)
	srcBlock := util.GenerateRTP(4, 3)

	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)
	bitsrings := bitstring.GetBlockBitstring(&srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPacketsRow := recover.GenerateRepairLD(&bitsrings, 4, 3, 0, SN_Base)

	for i := 0; i < len(srcBlock); i++ {
		fmt.Println(util.PrintPkt(srcBlock[i]))
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
	for _, pkt := range repairPacketsRow {
		fmt.Println("len :", len(buffer.Extract(BUFFER, pkt)))
		fmt.Println("missing:", buffer.CountMissing(BUFFER, pkt))
	}

}


func test2() {
	BUFFER := make(map[buffer.Key]rtp.Packet)

	srcBlock := util.GenerateRTP(3, 4)
	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)

	bitsrings := bitstring.GetBlockBitstring(&srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPackets := recover.GenerateRepairLD(&bitsrings, 3, 4, 0, SN_Base)

	for i := 0; i < len(srcBlock); i++ {
		fmt.Println(util.PrintPkt(srcBlock[i]))
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
		fmt.Println("missing:", buffer.CountMissing(BUFFER, pkt))
	}
}

func test3() {
	BUFFER := make(map[buffer.Key]rtp.Packet)

	// sender
	srcBlock := util.GenerateRTP(10, 10)

	mask := uint16(36160)                         // 1000110101000000 16 bit 3,4,6,8
	optionalmask1 := uint32(3229756930)           // 11000000100000100010111000000010 32 bit 15,22,28,32,34,35,36,44,
	optionalmask2 := uint64(13871700391609118210) // 1100000010000010001011100000001011000000100000100010111000000010 64 bit 46,47,54,60,64,66,67,68,76,78,79,86,92,96,98,99
	k1 := true
	k2 := true

	sn_base := srcBlock[0].Header.SequenceNumber
	fmt.Println("sn base :", sn_base)


	// creating dummy mask repair packet for buffer testing
	fecheader := fech.NewFecHeaderFlexibleMask(false, false, false, false, 11, false, 56, 23434, 342334, sn_base, k1, mask, optionalmask1, k2, optionalmask2)
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
	for index, pkt := range associatedPkts {
		fmt.Println("index : ", index)
		fmt.Println(util.PrintPkt(pkt))
	}

}


func main() {
	/*
		Uncomment to test the respective test case
		test1() -> row fec with buffer
		test2() -> col fec with buffer
		test3() -> mask fec with buffer
	*/

	test1()
	// test2()
	// test3()
}
