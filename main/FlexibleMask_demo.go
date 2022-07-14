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

	L = 10
	D = 10
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

	mask := uint16(36160)                         // 1|000110101000000 16 bit 3,4,6,8
	optionalmask1 := uint32(3229756930)           // 1|1000000100000100010111000000010 32 bit 15,22,28,32,34,35,36,44,
	optionalmask2 := uint64(13871700391609117696) // 1100000010000010001011100000001011000000100000100010110000000000 64 bit 46,47,54,60,64,66,67,68,76,78,79,86,92,96,98,99

	// need to check if mask bits above source block length is set, then signal for resending a correct mask or manually recorrect those bits here to 0.
	flexRepairPacket := recover.GenerateRepairFlex(&srcBlock, mask, optionalmask1, optionalmask2)

	fmt.Println(string(Green), "Send src block")

	for i := 0; i < len(srcBlock); i++ {
		if i != 3 {

			time.Sleep(1 * time.Second)

			fmt.Println(string(Green), "Sending a src packet")
			util.PrintPkt(srcBlock[i])
			fmt.Println()
			buf, _ := srcBlock[i].Marshal()
			conn.Write(buf)
		} else {
			time.Sleep(1 * time.Second)

			fmt.Println(string(Red), "missing packet")
			// recoveredPacket, _ := recover.RecoverMissingPacket(&srcBlock, repairPacket)
			util.PrintPkt(srcBlock[i])
		}

		if i == len(srcBlock)-1 {
			// time.Sleep(1 * time.Second)

			repairBuf, _ := flexRepairPacket.Marshal()
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

	conn.SetReadDeadline(time.Now().Add(120 * time.Second))

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

			associatedSrcPackets := buffer.ExtractMask(BUFFER, repairPacket)
			fmt.Println("len : ", len(associatedSrcPackets))
			fmt.Println(string(Red), "Recovered missing packer")
			recoveredPacket, _ := recover.RecoverMissingPacketFlex(&associatedSrcPackets, repairPacket)
			util.PrintPkt(recoveredPacket)
		} else {
			fmt.Println(string(White), "recieved src pkt")
			util.PrintPkt(currPkt)
			fmt.Println()
			// srcBlock = append(srcBlock, currPkt)

			buffer.Update(BUFFER, currPkt)
		}
	}
}

func main() {
	go encoder()
	decoder()
}
