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

	// generate packets
	stream := util.GenerateRTP(10, 10); 
	util.PadPackets(&stream)

	// test case list
	// variant 0 -> row, 1 -> col, 2 -> 2D
	testCaseList := [][]int{
		// {4, 3, 2},
		// {2, 2, 0},
		{3, 3, 1},
	}

	index := 0
	for _, item := range testCaseList {
		L := item[0]
		D := item[1]
		variant := item[2]

		srcBlock := stream[index : index + L * D]
		index += L * D

		/*
			a  X  X  X r1   0 X  X  X
			a  X  X  X r1   X 5  X  7
			X  j  k  l r3   X 9 10 11
			c1 c2 c3 c4
		*/

		testCaseMap := map[int]int {
			0 : 1, 5 : 1, 7 : 1, 9 : 1, 10 : 1, 11 : 1, 1 : 1, 6 : 1,
		}

		repairPackets := recover.GenerateRepairLD(&srcBlock, L, D, variant)

		for i:= 0 ; i < len(srcBlock); i++ {
			_, isPresent := testCaseMap[i]
			if isPresent {
				fmt.Println(string(Green), "Sending src block")
				util.PrintPkt(srcBlock[i])
				fmt.Println()
	
				buf, _ := srcBlock[i].Marshal()
				conn.Write(buf)
			} else {
				fmt.Println(string(Red), "Missing Packet at sender end")
				util.PrintPkt(srcBlock[i])
				fmt.Println()
			}
			time.Sleep(1 * time.Second)
		}
	
		
		// sending repair packets, row first then column
		fmt.Println(string(Blue), "*** Sending repair pkts ***")
		for i := 0; i < len(repairPackets); i++ {
			time.Sleep(1 * time.Second)
	
			fmt.Println(string(Blue), "Sending a repair packet")
			util.PrintPkt(repairPackets[i])
			fmt.Println()
			repairBuf, _ := repairPackets[i].Marshal()
			conn.Write(repairBuf)
		}
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

	conn.SetReadDeadline(time.Now().Add(25 * time.Second)) // stops reading after 25 seconds
	
	for {
		buf := make([]byte, mtu)
		i, _, err := conn.ReadFrom(buf)

		if err != nil {
			break
		}

		currPkt := rtp.Packet{}
		currPkt.Unmarshal(buf[:i])
	
		if currPkt.SSRC == repairSSRC {
			fmt.Println(string(Blue), "Recieved Repair PKt")
			util.PrintPkt(currPkt)
			fmt.Println()

			// Unmarshal payload to get the values of L and D to seggregate row and column repair packets
			var repairheader fech.FecHeaderLD = fech.FecHeaderLD{}
			repairheader.Unmarshal(currPkt.Payload[:12])

			// check R, F for fec variant
			// condition for 2D
			if repairheader.D == uint8(1) {
				fmt.Println("First round of row recovery")

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
				fmt.Println(string(Red), "*** Recovered Packet ***")
				util.PrintPkt(recoveredPacket)
			}
			

			fmt.Println(string(White), "Length of associatedSrcPackets:",len(associatedSrcPackets))
			fmt.Println("col_count:",col_count)
			fmt.Println("repairheader.L",repairheader.L)
			fmt.Println()

			if col_count == repairheader.L{
				fmt.Println("Second round of row recovery")
				// second round row
				// for all pkts in ROW_BUFFER
					// reapir using repair again
				
				for _,repairPacket:=range BUFFER_ROW_REC {
					associatedSrcPackets := buffer.Extract(BUFFER, repairPacket)
					recoveredPacket, status := recover.RecoverMissingPacket(&associatedSrcPackets, repairPacket)

					if status == 0 {
						buffer.Update(BUFFER, recoveredPacket) // update recoveredPacket to buffer
						fmt.Println(string(Red), "*** Recovered Packet ***")
						util.PrintPkt(recoveredPacket)
					}
					
					fmt.Println(string(White), "Length of associatedSrcPackets:",len(associatedSrcPackets))
					
				}
				// reset the variables
				is_2d = false
				col_count = 0
			}

		} else {
			fmt.Println(string(White), "recieved src pkt")
			util.PrintPkt(currPkt)
			fmt.Println()

			buffer.Update(BUFFER, currPkt)
		}
	}
	
	fmt.Println("Printing Row recovery packets form Buffer:",BUFFER_ROW_REC)
	BUFFER_ROW_REC = make(map[buffer.Key]rtp.Packet)

	fmt.Println("Printing All the Packets form Buffer:",BUFFER)
	// Check if retransmission is required
	// Print or save all the packets
	BUFFER = make(map[buffer.Key]rtp.Packet)
}

func main() {
	go encoder()
	decoder()
}


