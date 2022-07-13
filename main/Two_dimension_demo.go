package main

import (
	"flexfec/buffer"
	fech "flexfec/fec_header"
	"flexfec/recover"
	"flexfec/util"
	"fmt"
	"math"
	"net"
	"time"
	"github.com/pion/rtp"
)

const (
	listenPort = 6420
	ssrc       = 5000
	mtu        = 200
	Red        = "\033[31m"
	Green      = "\033[32m"
	White      = "\033[37m"
	Blue       = "\033[34m"
)

var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)
var BUFFER_ROW_REC map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)

func min(a uint16, b uint16) uint16 {
    if a < b {
        return a
    }
    return b
}

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
	srcBlock := util.GenerateRTP(4, 3);

	// have check if we need to do row and column wise
	util.PadPackets(&srcBlock)

	repairPacketsRow,repairPacketsColumns:=recover.GenerateRepair2dFec(&srcBlock,4,3)


	// removing srcBlock[2] in new Block
	var newBlock []rtp.Packet
	newBlock = append(newBlock, srcBlock[:1]...)
	newBlock = append(newBlock, srcBlock[5:6]...)
	newBlock = append(newBlock, srcBlock[7:8]...)
	newBlock = append(newBlock, srcBlock[9:]...)

	fmt.Println(string(Red), "Missing Packet at sender end")
	fmt.Println(newBlock)
	fmt.Println()

	srcBlock=newBlock
	// sending packets
	fmt.Println(string(Green), "Send src block")
	for i := 0; i < len(srcBlock); i++ {
		time.Sleep(1 * time.Second)

		fmt.Println(string(Green), "Sending a src packet")
		util.PrintPkt(srcBlock[i])
		fmt.Println()
		buf, _ := srcBlock[i].Marshal()
		conn.Write(buf)
	}
	
	// sending repair packets, row first then column
	fmt.Println(string(Blue), "*** Sending row repair pkt ***")
	for i := 0; i < len(repairPacketsRow); i++ {
		time.Sleep(1 * time.Second)

		fmt.Println(string(Blue), "Sending a row repair packet")
		util.PrintPkt(repairPacketsRow[i])
		fmt.Println()
		repairBuf, _ := repairPacketsRow[i].Marshal()
		conn.Write(repairBuf)
	}

	// sending repair packets,  column
	fmt.Println(string(Blue), "*** Sending column repair pkt ***")
	for i := 0; i < len(repairPacketsColumns); i++ {
		time.Sleep(1 * time.Second)

		fmt.Println(string(Blue), "Sending a column repair packet")
		util.PrintPkt(repairPacketsColumns[i])
		fmt.Println()
		repairBuf, _ := repairPacketsColumns[i].Marshal()
		conn.Write(repairBuf)
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

	conn.SetReadDeadline(time.Now().Add(25 * time.Second))
	
	// --------------------------	
	repairSSRC := uint32(2868272638)

	for {
		buf := make([]byte, mtu)
		i, _, err := conn.ReadFrom(buf)

		if err != nil {
			break
		}

		currPkt := rtp.Packet{}
		currPkt.Unmarshal(buf[:i])
		// ---------------------------
		// min seq number among 2d row
		var curr_min uint16=math.MaxUint16
		is_2d_row:=false
		// no of 2d row packets
		curr_count:=1
		col_count:=uint8(0)
		
		if currPkt.SSRC == repairSSRC {
			
			fmt.Println(string(Blue), "Recieved Repair PKt")
			util.PrintPkt(currPkt)
			fmt.Println()

			// Unmarshal payload to get the values of L and D to seggregate row and column repair packets
			var repairheader fech.FecHeaderLD = fech.FecHeaderLD{}
			repairheader.Unmarshal(currPkt.Payload[:12])

			// condition for 2D
			fmt.Println("printing D value",repairheader.D)
			if repairheader.D==uint8(1){
				fmt.Println("Entering the 1st phase of recovery")
				if is_2d_row{
					curr_count++
					curr_min=min(curr_min,currPkt.SequenceNumber)

				}else{
					is_2d_row=true
					curr_count=1
					col_count=0
					curr_min=currPkt.SequenceNumber
				}
				buffer.Update(BUFFER_ROW_REC, currPkt)

			}else{
				is_2d_row=false
				col_count++
			}
			
			// Repair using repair packet

			associatedSrcPackets := buffer.Extract(BUFFER, currPkt)
			fmt.Println("Length of associatedSrcPackets:",len(associatedSrcPackets))
			
			recoveredPacket, _ := recover.RecoverMissingPacket(&associatedSrcPackets, currPkt)
			// update recoveredPacket to buffer
			buffer.Update(BUFFER, recoveredPacket)
			
			fmt.Println("col_count:",col_count)
			fmt.Println("repairheader.L",repairheader.L)
			
			if col_count==repairheader.L{
				fmt.Println("Entering Second row recovery phase-------")
				// second round row
				// for all pkts in EXTRACT(CURRMIN to CURRMIN + CURR_COUNT from ROW_BUFFER)
				// reapir using repair again
				// reset the variables

				for _,repairPacket:=range BUFFER_ROW_REC {
					associatedSrcPackets := buffer.Extract(BUFFER, repairPacket)
					fmt.Println("Length of associatedSrcPackets:",len(associatedSrcPackets))
					recoveredPacket, _ := recover.RecoverMissingPacket(&associatedSrcPackets, repairPacket)
					// update recoveredPacket to buffer
					buffer.Update(BUFFER, recoveredPacket)
				}
				// delete buffer
			}

		}else{
			fmt.Println(string(White), "recieved src pkt")
			util.PrintPkt(currPkt)
			fmt.Println()

			buffer.Update(BUFFER, currPkt)
		}
	}
	fmt.Println("Printing Row recovery packets form Buffer:",BUFFER_ROW_REC)
	BUFFER_ROW_REC =make(map[buffer.Key]rtp.Packet)

	fmt.Println("Printing All the Packets form Buffer:",BUFFER)
	// Check if retransmission is required
	// Print or save all the packets
	BUFFER=make(map[buffer.Key]rtp.Packet)
}
func main() {
	go encoder()
	decoder()
}


//  a  X  X  X r1 1
//  X  f  X  h r2 2
//  X  j  k  l r3 3
//  c1 c2 c3 c4
// 1	2	1	2