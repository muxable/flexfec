package main

import (
	"flexfec/util"
	"flexfec/buffer"
	"flexfec/recover"
	"net"
	"fmt"
	"time"
	fech "flexfec/fec_header"
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


func encoder() {
	serverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", listenPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, serverAddr)
	if err != nil {
		panic(err)
	}

	// generate packets 2x2
	srcBlock := util.GenerateRTP(2, 2);

	// have check if we need to do row and column wise
	util.PadPackets(&srcBlock)

	repairPacketsRow,repairPacketsColumns:=recover.GenerateRepair2dFec(&srcBlock,2,2)

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

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	// --------------------------
	// srcBlock := []rtp.Packet{}
	
	repairPacketRows := rtp.Packet{}
	repairPacketColumns := rtp.Packet{}

	repairSSRC := uint32(2868272638)

	for {
		buffer := make([]byte, mtu)
		i, _, err := conn.ReadFrom(buffer)

		if err != nil {
			break
		}

		currPkt := rtp.Packet{}
		currPkt.Unmarshal(buffer[:i])

		
		if currPkt.SSRC == repairSSRC {

		// if condition for 2D
			fmt.Println(string(Blue), "Recieved Repair PKt")
			util.PrintPkt(currPkt)
			fmt.Println()

			// Unmarshal payload to get the values of L and D to seggregate row and column repair packets
			var repairheader fech.FecHeaderLD = fech.FecHeaderLD{}
			repairheader.Unmarshal(currPkt.Payload[:12])

			// row repair packets
			if(repairheader.D==1){
				buffer.Update(BUFFER_ROW_REC, currPkt)
				
				//Check and Call for packet recovery
				// Requires of creation oif 2d buffer for packets
				associatedSrcPackets := buffer.Extract(BUFFER, currPkt)
				fmt.Println("len : ", len(associatedSrcPackets))
				fmt.Println(string(Red), "Recovered missing packer")
				recoveredPacket, _ := recover.RecoverMissingPacketLD(&associatedSrcPackets, currPkt)
				util.PrintPkt(recoveredPacket)
				
				// Add recoveredPacket to buffer
				buffer.Update(BUFFER, recoveredPacket)
				
			// column repair packets
			}else{
				repairPacketColumns=currPkt
				// Check and Call for packet recovery
			}

		} else {
			fmt.Println(string(White), "recieved src pkt")
			util.PrintPkt(currPkt)
			fmt.Println()

			// without buffer
			// srcBlock = append(srcBlock, currPkt)

			// with buffer
			buffer.Update(BUFFER, currPkt)
		}
	}

}
func main() {
	go encoder()
	decoder()
}