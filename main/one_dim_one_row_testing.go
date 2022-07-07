package main

import (
	"flexfec/recover"
	"flexfec/util"
	"flexfec/buffer"
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
)

var buffer map[Key]rtp.Packet

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

	fmt.Println(string(Red), "Missing Packet at sender end")
	util.PrintPkt(srcBlock[2])
	fmt.Println()

	repairPacket := recover.GenerateRepair(&srcBlock, 5, 1)

	// removing srcBlock[2] in new Block
	var newBlock []rtp.Packet
	newBlock = append(newBlock, srcBlock[:2]...)
	newBlock = append(newBlock, srcBlock[3:]...)

	// defer conn.Close()

	fmt.Println(string(Green), "Send src block")
	for i := 0; i < len(newBlock); i++ {
		time.Sleep(1 * time.Second)

		fmt.Println(string(Green), "Sending a src packet")
		util.PrintPkt(newBlock[i])
		fmt.Println()
		buf, _ := newBlock[i].Marshal()
		conn.Write(buf)
	}

	// sending repair pkt
	fmt.Println(string(Blue), "Send repair pkt")
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

	conn.SetReadDeadline(time.Now().Add(7 * time.Second))

	srcBlock := []rtp.Packet{}
	repairPacket := rtp.Packet{}
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
			repairPacket = currPkt

		} else {
			fmt.Println(string(White), "recieved src pkt")
			util.PrintPkt(currPkt)
			fmt.Println()
			srcBlock = append(srcBlock, currPkt)
		}

	}

	fmt.Println(string(Red), "Recovered missing packer")
	recoveredPacket, _ := recover.RecoverMissingPacket(&srcBlock, repairPacket)
	util.PrintPkt(recoveredPacket)
}

func main() {
	go sender()
	receiver()
}
