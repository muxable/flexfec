package main

import (
	"os"
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
	variant = 2
)

func printBuffer(repairBuffer []rtp.Packet) {
	fmt.Print(string(Green), "REPAIR BUFFER : [ ")
	for _, pkt := range repairBuffer {
		fmt.Print(pkt.SequenceNumber, " ")
	}
	fmt.Println("]")
}

func testFlexFEC() {
	file, err := os.Create("debug.txt")

	if err != nil {
		fmt.Println(err)
	}

	var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)
	var REPAIR_BUFFER []rtp.Packet

	// Sender
	srcBlock := util.GenerateRTP(4, 3)
	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)

	bitsrings := bitstring.GetBlockBitstring(&srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPackets2d := recover.GenerateRepairLD(&bitsrings, 4, 3, variant, SN_Base)

	var recievedPackets []rtp.Packet
	testcaseMap := util.GetTestCaseMap(variant)

	for i := 0; i < len(srcBlock); i++ {
		_, isPresent := testcaseMap[i]
		if isPresent {
			file.WriteString(util.PrintPkt(srcBlock[i]))
			recievedPackets = append(recievedPackets, srcBlock[i])
		} else {
			fmt.Fprintln(file, "missing packet")
			file.WriteString(util.PrintPkt(srcBlock[i]))
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

			// printing buffer
			printBuffer(REPAIR_BUFFER)

			currRecPkt := REPAIR_BUFFER[0]
			REPAIR_BUFFER = REPAIR_BUFFER[1:]

			associatedSrcPackets := buffer.Extract(BUFFER, currRecPkt)
			recoveredPacket, status := recover.RecoverMissingPacket(&associatedSrcPackets, currRecPkt)
			
			if status == 1 {
				fmt.Println(string(White), "Repair packet ", currRecPkt.SequenceNumber, " fully recovered")
			} else if status == 0 {
				fmt.Println(string(White), "Using repair packet ", currRecPkt.SequenceNumber, "to recover")
				fmt.Println(string(Red), "Recovered packet")
				fmt.Println(util.PrintPkt(recoveredPacket))
				buffer.Update(BUFFER, recoveredPacket)
			} else if status == -1 {
				fmt.Println(string(White), "Recovery not possible\n")
				REPAIR_BUFFER=append(REPAIR_BUFFER,currRecPkt)
				break
			}
		}
	}

	// printing BUFFER
	fmt.Println("\nPrinting All the Packets form Buffer:")
	fmt.Print("REPAIR BUFFER : [ ")
	for _, pkt := range BUFFER {
		fmt.Print(pkt.SequenceNumber, " ")
	}
	fmt.Println("]\n")

}


func main(){
	testFlexFEC()
}