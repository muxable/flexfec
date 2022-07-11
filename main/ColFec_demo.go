package main

import (
	"flexfec/buffer"
	"flexfec/recover"
	"flexfec/util"
	"fmt"
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

	L = 4
	D = 3
)

var BUFFER map[buffer.Key]rtp.Packet = make(map[buffer.Key]rtp.Packet)

func encoder() {
	serverAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", listenPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, serverAddr)
	if err != nil {
		panic(err)
	}

	srcBlock := util.GenerateRTP(L, D)
	util.PadPackets(&srcBlock)

	repairPacketsCol := recover.GenerateRepairLD(&srcBlock, L, D)

	fmt.Println(string(Green), "Send src block")
	for i := 0; i < len(srcBlock); i++ {
		time.Sleep(1 * time.Second)

		if i != 1 && i != 2 && i != 6 {

			fmt.Println(string(Green), "Sending a src packet")
			util.PrintPkt(srcBlock[i])
			fmt.Println()
			buf, _ := srcBlock[i].Marshal()
			conn.Write(buf)

		} else {
			fmt.Println(string(Red), "missing packet")
			// recoveredPacket, _ := recover.RecoverMissingPacket(&srcBlock, repairPacket)
			util.PrintPkt(srcBlock[i])

		}

		if (i)/L >= D-1 {
			// sending repair pkt
			time.Sleep(1 * time.Second)

			fmt.Println(string(Blue), "Send repair pkt")
			repairPacket := repairPacketsCol[i%L]
			util.PrintPkt(repairPacket)
			fmt.Println()

			repairBuf, _ := repairPacket.Marshal()
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

	conn.SetReadDeadline(time.Now().Add(20 * time.Second))

	// srcBlock := []rtp.Packet{}
	repairPacket := rtp.Packet{}
	repairSSRC := uint32(2868272638)

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
			repairPacket = currPkt

			associatedSrcPackets := buffer.Extract(BUFFER, repairPacket)
			fmt.Println("len : ", len(associatedSrcPackets))
			fmt.Println(string(Red), "Recovered missing packer")
			recoveredPacket, _ := recover.RecoverMissingPacketLD(&associatedSrcPackets, repairPacket)
			util.PrintPkt(recoveredPacket)

		} else {
			fmt.Println(string(White), "recieved src pkt")
			util.PrintPkt(currPkt)
			fmt.Println()
			// srcBlock = append(srcBlock, currPkt)

			buffer.Update(BUFFER, currPkt)
		}

	}

	// fmt.Println(string(Red), "Recovered missing packer")
	// recoveredPacket, _ := recover.RecoverMissingPacket(&srcBlock, repairPacket)
	// util.PrintPkt(recoveredPacket)
}

func main() {
	go encoder()
	decoder()
}
