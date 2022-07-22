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
		// fmt.Fprintln(file, bitsrings[i], len(bitsrings))
		// fmt.Println("bit str ", bitsrings[i])
		_, isPresent := testcaseMap[i]
		if isPresent {
			// fmt.Println(string(Green), "Sending a src packet")
			// fmt.Println(util.PrintPkt(srcBlock[i]))
			file.WriteString(util.PrintPkt(srcBlock[i]))
			recievedPackets = append(recievedPackets, srcBlock[i])
		} else {
			// fmt.Println(string(Red), "missing packet")
			// fmt.Println(util.PrintPkt(srcBlock[i]))
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
	testFlexFEC()
}