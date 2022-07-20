package main

import (
	"fmt"
	"net"
	"time"
	"flexfec/util"
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/fec_header"
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

// Global variables
var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)
var BUFFER_ROW_REC map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)

var is_2d_row bool  = false
var is_2d bool = false
var col_count uint8  = uint8(0)


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
	stream := util.GenerateRTP(10, 10); 
	util.PadPackets(&stream)

	// test case list
	// variant 0 -> row, 1 -> col, 2 -> 2D
	testCaseList := [][]int{
		{4, 3, 2},
		{4, 3, 0},
		{4, 3, 1},
	}

	index := 0
	for _, item := range testCaseList {
		L := item[0]
		D := item[1]
		variant := item[2]

		srcBlock := stream[index : index + L * D]
		index += L * D

		testCaseMap := util.GetTestCaseMap(variant)

		repairPackets := recover.GenerateRepairLD(&srcBlock, L, D, variant)

		for i:= 0 ; i < len(srcBlock); i++ {
			_, isPresent := testCaseMap[i]
			if isPresent {
				file.WriteString("Sending src block\n")
				file.WriteString(util.PrintPkt(srcBlock[i]))

				// fmt.Println(string(Green), "Sending src block")
				// fmt.Println(util.PrintPkt(srcBlock[i]))
	
				buf, _ := srcBlock[i].Marshal()
				conn.Write(buf)
			} else {
				file.WriteString("Missing Packet at sender end\n")
				file.WriteString(util.PrintPkt(srcBlock[i]))

				// fmt.Println(string(Red), "Missing Packet at sender end")
				// fmt.Println(util.PrintPkt(srcBlock[i]))
			}
			time.Sleep(500 * time.Millisecond)
		}
	
		
		// sending repair packets, row first then column
		// fmt.Println(string(Blue), "*** Sending repair pkts ***")
		file.WriteString("*** Sending repair pkts ***\n")
		for i := 0; i < len(repairPackets); i++ {
			time.Sleep(500 * time.Millisecond)
	
			// fmt.Println(string(Blue), "Sending a repair packet")
			// fmt.Println(util.PrintPkt(repairPackets[i]))

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

	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) // stops reading after 25 seconds
	
	for {
		buf := make([]byte, mtu)
		i, _, err := conn.ReadFrom(buf)

		if err != nil {
			break
		}

		currPkt := rtp.Packet{}
		currPkt.Unmarshal(buf[:i])
	
		if currPkt.SSRC == repairSSRC {
			file.WriteString("Recieved Repair PKt\n")
			file.WriteString(util.PrintPkt(currPkt))

			// fmt.Println(string(Blue), "Recieved Repair PKt")
			// fmt.Println(util.PrintPkt(currPkt))

			// Unmarshal payload to get the values of L and D to seggregate row and column repair packets
			var repairheader fech.FecHeaderLD = fech.FecHeaderLD{}
			repairheader.Unmarshal(currPkt.Payload[:12])

			// check R, F for fec variant
			// condition for 2D
			if repairheader.D == uint8(1) {
				file.WriteString("First round of row recovery\n")

				if is_2d_row == false {
					is_2d_row = true
					is_2d = true
					col_count = 0
				}

				buffer.Update(BUFFER_ROW_REC, currPkt)
			} else {
				is_2d_row = false
				if is_2d {
					col_count++
				}
			}
			
			// Repair using repair packet
			associatedSrcPackets := buffer.Extract(BUFFER, currPkt)
			recoveredPacket, status := recover.RecoverMissingPacket(&associatedSrcPackets, currPkt)
		
			if status == 0 {
				buffer.Update(BUFFER, recoveredPacket) // update recoveredPacket to buffer

				file.WriteString("*** Recovered Packet ***\n")
				file.WriteString(util.PrintPkt(recoveredPacket))

				// fmt.Println(string(Red), "*** Recovered Packet ***")
				// fmt.Println(util.PrintPkt(recoveredPacket))
			}
			

			// fmt.Println(string(White), "Length of associatedSrcPackets:",len(associatedSrcPackets))
			file.WriteString("Length of associatedSrcPackets:" + strconv.Itoa(len(associatedSrcPackets)) + "\n")

			if col_count == repairheader.L{
				// fmt.Println("Second round of row recovery")
				file.WriteString("Second round of row recovery\n")
				// second round row
				// for all pkts in ROW_BUFFER
					// reapir using repair again
				
				for _,repairPacket:=range BUFFER_ROW_REC {
					associatedSrcPackets := buffer.Extract(BUFFER, repairPacket)
					recoveredPacket, status := recover.RecoverMissingPacket(&associatedSrcPackets, repairPacket)

					if status == 0 {
						buffer.Update(BUFFER, recoveredPacket) // update recoveredPacket to buffer
						// fmt.Println(string(Red), "*** Recovered Packet ***")
						// fmt.Println(util.PrintPkt(recoveredPacket))
						file.WriteString("*** Recovered Packet ***\n")
						file.WriteString(util.PrintPkt(recoveredPacket))
					}
					
					// fmt.Println(string(White), "Length of associatedSrcPackets:",len(associatedSrcPackets))
					file.WriteString("Length of associatedSrcPackets:" + strconv.Itoa(len(associatedSrcPackets)) + "\n")
					
				}
				// reset the variables
				is_2d = false
				col_count = 0
			}

		} else {
			// fmt.Println(string(White), "recieved src pkt")
			// fmt.Println(util.PrintPkt(currPkt))

			// file.WriteString("recieved src pkt\n")
			// file.WriteString(util.PrintPkt(currPkt))

			buffer.Update(BUFFER, currPkt)
		}
	}
	
	fmt.Println("Printing Row recovery packets form Buffer:",BUFFER_ROW_REC)
	BUFFER_ROW_REC = make(map[buffer.Key]rtp.Packet)

	fmt.Println("Printing All the Packets form Buffer:", BUFFER)

	// Check if retransmission is required
	// Print or save all the packets
	BUFFER = make(map[buffer.Key]rtp.Packet)
}

func main() {
	go encoder()
	decoder()
}


