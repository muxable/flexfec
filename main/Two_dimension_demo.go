package main

import (
	"flexfec/bitstring"
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/util"
	"fmt"
	"net"
	"os"
	"sort"
	"time"

	"github.com/pion/rtp"
)

const (
	repairSSRC = uint32(2868272638)
	listenPort = 6420
	ssrc       = 5000
	mtu        = 200
	Red        = "\033[31m"
	Green      = "\033[32m"
	White      = "\033[37m"
	Blue       = "\033[34m"
)

func printBuffer(repairBuffer []rtp.Packet) {
	fmt.Print(string(Green), "REPAIR BUFFER : [ ")
	for _, pkt := range repairBuffer {
		fmt.Print(pkt.SequenceNumber, " ")
	}
	fmt.Println("]")
}

// Global variables
var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)
var REPAIR_BUFFER []rtp.Packet

func encoder() {
	serverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", listenPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, serverAddr)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("output/sender.txt")

	if err != nil {
		fmt.Println("file error")
	}

	// generate packets
	stream := util.GenerateRTP(10, 10)

	// test case list
	// variant 0 -> row, 1 -> col, 2 -> 2D
	testCaseList := [][]int{
		{4, 3, 2},
		// {4, 3, 0},
		// {4, 3, 1},
	}

	index := 0
	for _, item := range testCaseList {
		L := item[0]
		D := item[1]
		variant := item[2]

		srcBlock := stream[index : index+L*D]
		SN_Base := uint16(srcBlock[0].Header.SequenceNumber)

		bitsrings := bitstring.GetBlockBitstring(&srcBlock)
		util.PadBitStrings(&bitsrings, -1)

		index += L * D

		testCaseMap := util.GetTestCaseMap(variant)

		repairPackets := recover.GenerateRepairLD(&bitsrings, L, D, variant, SN_Base)

		for i := 0; i < len(srcBlock); i++ {
			_, isPresent := testCaseMap[i]
			if isPresent {
				file.WriteString("Sending src block\n")
				file.WriteString(util.PrintPkt(srcBlock[i]))

				buf, _ := srcBlock[i].Marshal()
				conn.Write(buf)
			} else {
				file.WriteString("Missing Packet at sender end\n")
				file.WriteString(util.PrintPkt(srcBlock[i]))
			}
			time.Sleep(500 * time.Millisecond)
		}

		// sending repair packets, row first then column
		file.WriteString("*** Sending repair pkts ***\n")
		for i := 0; i < len(repairPackets); i++ {
			time.Sleep(500 * time.Millisecond)

			file.WriteString("Sending a repair packet\n")
			file.WriteString(util.PrintPkt(repairPackets[i]))

			repairBuf, _ := repairPackets[i].Marshal()
			conn.Write(repairBuf)
		}

		file.WriteString("-----------------------------------------------------------------\n")
	}
}

func decoder() {
	serverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", listenPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp4", serverAddr)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("output/receiver.txt")

	if err != nil {
		fmt.Println("file error")
	}

	conn.SetReadDeadline(time.Now().Add(20 * time.Second)) // stops reading after n seconds

	for {
		buf := make([]byte, mtu)
		i, _, err := conn.ReadFrom(buf)

		if err != nil {
			break
		}

		currPkt := rtp.Packet{}
		currPkt.Unmarshal(buf[:i])

		if currPkt.SSRC == repairSSRC {
			file.WriteString("Recieved Repair Packet\n")
			file.WriteString(util.PrintPkt(currPkt))
			REPAIR_BUFFER = append(REPAIR_BUFFER, currPkt)
		} else {
			buffer.Update(BUFFER, currPkt)
		}

		for len(REPAIR_BUFFER) > 0 {
			sort.Slice(REPAIR_BUFFER, func(i, j int) bool {
				return buffer.CountMissing(BUFFER, REPAIR_BUFFER[i]) < buffer.CountMissing(BUFFER, REPAIR_BUFFER[j])
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
				file.WriteString("Recovered Packet\n")
				file.WriteString(util.PrintPkt(recoveredPacket))
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
	go encoder()
	decoder()
}
