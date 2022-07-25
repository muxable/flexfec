package main

import (
	"flexfec/bitstring"
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/util"
	"fmt"
	"os"
	"sort"

	"github.com/pion/rtp"
)

const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	White   = "\033[37m"
	Blue    = "\033[34m"
	L       = 4
	D       = 3
	variant = 2 // 0,1,2,3: row, col, 2d, flex

	mask          = uint16(36160)                // 1|000110101000000 16 bit 3,4,6,8
	optionalmask1 = uint32(3229756930)           // 1|1000000100000100010111000000010 32 bit 15,22,28,32,34,35,36,44,
	optionalmask2 = uint64(13871700391609117696) // 1100000010000010001011100000001011000000100000100010110000000000 64 bit 46,47,54,60,64,66,67,68,76,78,79,86,92,96,98,99
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

	srcBlock := []rtp.Packet{}
	// Sender
	if variant != 3 {
		srcBlock = util.GenerateRTP(L, D)
	} else if variant == 3 {
		srcBlock = util.GenerateRTP(10, 10)
	} else {
		fmt.Println("incorrect variant did not generate src packets")
	}

	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)

	bitsrings := bitstring.GetBlockBitstring(&srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	var repairPackets []rtp.Packet
	if variant < 3 {
		repairPackets = recover.GenerateRepairLD(&bitsrings, 4, 3, variant, SN_Base)
	} else if variant == 3 {
		repairPackets = recover.GenerateRepairFlex(&bitsrings, mask, optionalmask1, optionalmask2, SN_Base)
	} else {
		fmt.Println("invalid variant")
	}

	var recievedPackets []rtp.Packet

	if variant < 3 {
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

	} else if variant == 3 {
		for i, pkt := range srcBlock {
			if i != 54 {
				file.WriteString(util.PrintPkt(pkt))
				recievedPackets = append(recievedPackets, pkt)
			} else {
				fmt.Fprintln(file, "missing packet")
				file.WriteString(util.PrintPkt(pkt))
			}
		}
	} else {
		fmt.Println("Incorrect variant provided")
	}

	//recevier
	for _, pkt := range recievedPackets {
		buffer.Update(BUFFER, pkt)
	}

	for _, pkt := range repairPackets {
		REPAIR_BUFFER = append(REPAIR_BUFFER, pkt)

		for len(REPAIR_BUFFER) > 0 {

			// not needed for flex
			if variant != 3 {
				sort.Slice(REPAIR_BUFFER, func(i, j int) bool {
					return buffer.CountMissing(BUFFER, REPAIR_BUFFER[i]) < buffer.CountMissing(BUFFER, REPAIR_BUFFER[j])
				})
			}

			// printing buffer
			printBuffer(REPAIR_BUFFER)

			currRecPkt := REPAIR_BUFFER[0]
			REPAIR_BUFFER = REPAIR_BUFFER[1:]

			associatedSrcPackets := []rtp.Packet{}
			recoveredPacket := rtp.Packet{}
			status := 2

			if variant != 3 {
				associatedSrcPackets = buffer.Extract(BUFFER, currRecPkt)
				recoveredPacket, status = recover.RecoverMissingPacket(&associatedSrcPackets, currRecPkt)

			} else {
				associatedSrcPackets = buffer.ExtractMask(BUFFER, currRecPkt)
				recoveredPacket, status = recover.RecoverMissingPacketFlex(&associatedSrcPackets, currRecPkt)
			}

			if status == 1 {
				fmt.Println(string(White), "Repair packet ", currRecPkt.SequenceNumber, " fully recovered")
			} else if status == 0 {
				fmt.Println(string(White), "Using repair packet ", currRecPkt.SequenceNumber, "to recover")
				fmt.Println(string(Red), "Recovered packet")
				fmt.Println(util.PrintPkt(recoveredPacket))
				buffer.Update(BUFFER, recoveredPacket)
			} else if status == -1 {
				fmt.Println(string(White), "Recovery not possible\n")
				REPAIR_BUFFER = append(REPAIR_BUFFER, currRecPkt)
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

func main() {
	testFlexFEC()
}
