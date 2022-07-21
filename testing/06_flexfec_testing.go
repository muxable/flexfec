package main

import (
	"fmt"
	"sort"
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/bitstring"
	"flexfec/util"
	"github.com/pion/rtp"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	White = "\033[37m"
	Blue  = "\033[34m"
	L     = 4
	D     = 3
)

func testrow() {
	var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)

	srcBlock := util.GenerateRTP(4, 3)
	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)
	bitsrings := bitstring.GetBlockBitstring(&srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPacketsRow := recover.GenerateRepairLD(&bitsrings, 4, 3, 0, SN_Base)

	var recievedPackets []rtp.Packet
	testcaseMap := util.GetTestCaseMap(0)

	for i := 0; i < len(srcBlock); i++ {
		_, isPresent := testcaseMap[i]
		if isPresent{
			fmt.Println(string(Green), "Sending a src packet")
			fmt.Println(util.PrintPkt(srcBlock[i]))
			recievedPackets = append(recievedPackets, srcBlock[i])
		} else {
			fmt.Println(string(Red), "missing packet")
			fmt.Println(util.PrintPkt(srcBlock[i]))
		}
	}

	//receiver
	for _, pkt := range recievedPackets {
		buffer.Update(BUFFER, pkt)
	}

	for i := 0; i < len(repairPacketsRow); i++ {
		associatedSrcPackets := buffer.Extract(BUFFER, repairPacketsRow[i])
		fmt.Println(string(White), "recovery")

		recoveredPacket, _ := recover.RecoverMissingPacket(&associatedSrcPackets, repairPacketsRow[i])
		fmt.Println(util.PrintPkt(recoveredPacket))
	}
}

func testcol() {
	var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)

	// Sender
	srcBlock := util.GenerateRTP(4, 3)
	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)

	bitsrings := bitstring.GetBlockBitstring(&srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPacketsCol := recover.GenerateRepairLD(&bitsrings, 4, 3, 1, SN_Base)

	var recievedPackets []rtp.Packet
	testcaseMap := util.GetTestCaseMap(1)

	for i := 0; i < len(srcBlock); i++ {
		_, isPresent := testcaseMap[i]
		if isPresent {
			fmt.Println(string(Green), "Sending a src packet")
			fmt.Println(util.PrintPkt(srcBlock[i]))
			recievedPackets = append(recievedPackets, srcBlock[i])
		} else {
			fmt.Println(string(Red), "missing packet")
			fmt.Println(util.PrintPkt(srcBlock[i]))
		}
	}

	//receiver
	for _, pkt := range recievedPackets {
		buffer.Update(BUFFER, pkt)
	}

	for _, repairPacket := range repairPacketsCol {
		associatedSrcPackets := buffer.Extract(BUFFER, repairPacket)
		recoveredPacket, _ := recover.RecoverMissingPacket(&associatedSrcPackets, repairPacket)
		fmt.Println(string(White), "recovered packets")
		fmt.Println(util.PrintPkt(recoveredPacket))
	}

}

func test2d() {
	var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)
	var REPAIR_BUFFER []rtp.Packet

	// Sender
	srcBlock := util.GenerateRTP(4, 3)
	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)

	bitsrings := bitstring.GetBlockBitstring(&srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPackets2d := recover.GenerateRepairLD(&bitsrings, 4, 3, 2, SN_Base)

	var recievedPackets []rtp.Packet
	testcaseMap := util.GetTestCaseMap(2)

	for i := 0; i < len(srcBlock); i++ {
		_, isPresent := testcaseMap[i]
		if isPresent {
			fmt.Println(string(Green), "Sending a src packet")
			fmt.Println(util.PrintPkt(srcBlock[i]))
			recievedPackets = append(recievedPackets, srcBlock[i])
		} else {
			fmt.Println(string(Red), "missing packet")
			fmt.Println(util.PrintPkt(srcBlock[i]))
		}
	}

	//recevier
	for _, pkt := range recievedPackets {
		buffer.Update(BUFFER, pkt)
	}

	for _, pkt := range repairPackets2d {
		REPAIR_BUFFER = append(REPAIR_BUFFER, pkt)

		for len(REPAIR_BUFFER) > 0 {
			sort.Slice(REPAIR_BUFFER, func(i, j int) bool {
				return buffer.CountMissing(BUFFER,REPAIR_BUFFER[i]) < buffer.CountMissing(BUFFER,REPAIR_BUFFER[j])
			})

			fmt.Println(REPAIR_BUFFER)
	
			currRecPkt := REPAIR_BUFFER[0]
			REPAIR_BUFFER = REPAIR_BUFFER[1:]

			associatedSrcPackets := buffer.Extract(BUFFER, currRecPkt)
			recoveredPacket, status := recover.RecoverMissingPacket(&associatedSrcPackets, currRecPkt)
			
			if status==0{
				fmt.Println(string(White),"Recovered packet")
				fmt.Println(util.PrintPkt(recoveredPacket))
				buffer.Update(BUFFER, recoveredPacket)
			}else if status==-1{
				fmt.Println("Recovery not possible\n")
				REPAIR_BUFFER=append(REPAIR_BUFFER,currRecPkt)
				break
			}
		}
	}

	fmt.Println("Printing All the Packets form Buffer:", BUFFER)

}


func main(){
	// testrow()
	// testcol()
	test2d()
}