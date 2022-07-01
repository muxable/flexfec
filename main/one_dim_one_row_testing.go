package main

import (
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

	srcBlock := util.GenerateRTP(5, 1)
	util.PadPackets(&srcBlock)

	fmt.Println("Missing Packet at sender end")
	util.PrintPkt(srcBlock[2])
	fmt.Println()
	// bitStr := bitstring.ToBitString(&srcBlock[2])
	// Y := binary.BigEndian.Uint16(bitStr[2:4])
	// fmt.Println("Y should have been : ", Y)

	repairPacket := recover.GenerateRepair(&srcBlock, 5, 1)

	// removing srcBlock[2] in new Block
	var newBlock []rtp.Packet
	newBlock = append(newBlock, srcBlock[:2]...)
	newBlock = append(newBlock, srcBlock[3:]...)

	// defer conn.Close()

	fmt.Println("Send src block")
	for i := 0; i < len(newBlock); i++ {
		time.Sleep(2 * time.Second)

		fmt.Println("Sending a src packet")
		util.PrintPkt(newBlock[i])
		fmt.Println()
		buf, _ := newBlock[i].Marshal()
		conn.Write(buf)
	}

	// sending repair pkt
	fmt.Println("Send repair pkt")
	util.PrintPkt(repairPacket)
	fmt.Println()

	repairBuf, _ := repairPacket.Marshal()
	conn.Write(repairBuf)
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

	srcBlock := []rtp.Packet{}

	repairPacket := rtp.Packet{}

	repairSSRC := uint32(2868272638)

	start := time.Now()
	for {
		buffer := make([]byte, mtu)

		i, _, err := conn.ReadFrom(buffer)

		if err != nil {
			panic(err)
		}

		// fmt.Println("recieved packet")
		// fmt.Println("buffer : ", buffer[:i], "\n")

		currPkt := rtp.Packet{}
		currPkt.Unmarshal(buffer[:i])

		if currPkt.SSRC == repairSSRC {
			//
			fmt.Println("Recieved Repair PKt")
			util.PrintPkt(currPkt)
			fmt.Println()
			repairPacket = currPkt

		} else {
			fmt.Println("recieved src pkt")
			util.PrintPkt(currPkt)
			fmt.Println()
			srcBlock = append(srcBlock, currPkt)
		}

		if time.Since(start) > 12*time.Second {
			break
		}
	}

	// Not RUNNING
	fmt.Println("Recovered missing packer")
	recoveredPacket, _ := recover.RecoverMissingPacket(&srcBlock, repairPacket)
	util.PrintPkt(recoveredPacket)
}

func main() {
	go sender()
	receiver()
}
