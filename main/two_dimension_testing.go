package main

import (
	"flexfec/util"
	"flexfec/recover"
	"net"
	"fmt"
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
func sender() {
	serverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", listenPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, serverAddr)
	if err != nil {
		panic(err)
	}

	// generate packets 5X3
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
func receiver() {
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
	srcBlock := []rtp.Packet{}
	
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
			fmt.Println(string(Blue), "Recieved Repair PKt")
			util.PrintPkt(currPkt)
			fmt.Println()
			// if(row){
			// 	// add to repairPacketRows
			// }else{
			// 	// add to repairPacketColumns
			// }
			// repairPacket = currPkt

		} else {
			fmt.Println(string(White), "recieved src pkt")
			util.PrintPkt(currPkt)
			fmt.Println()
			srcBlock = append(srcBlock, currPkt)
		}
	}

}
func main() {
	go sender()
	receiver()
}