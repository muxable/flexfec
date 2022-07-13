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
	repairSSRC := uint32(2868272638)

	for {
		buffer := make([]byte, mtu)
		i, _, err := conn.ReadFrom(buffer)

		if err != nil {
			break
		}

		currPkt := rtp.Packet{}
		currPkt.Unmarshal(buffer[:i])

		// ---------------------------
		var curr_min uint16=math.MaxUint16
		is_2d_row:=false
		curr_count:=1
		col_count:=0
		
		if currPkt.SSRC == repairSSRC {
			
			// Unmarshal payload to get the values of L and D to seggregate row and column repair packets
			var repairheader fech.FecHeaderLD = fech.FecHeaderLD{}
			repairheader.Unmarshal(currPkt.Payload[:12])

			// condition for 2D
			if repairheader.L>0 && repairheader.D==1{

				if is_2d_row{
					curr_count++
					// check for type compatability
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
			recoveredPacket, _ := recover.RecoverMissingPacketLD(&associatedSrcPackets, currPkt)
			// update recoveredPacket to buffer
			buffer.Update(BUFFER, recoveredPacket)

			if col_count==repairheader.L{

				// second round row
				// for all pkts in EXTRACT(CURRMIN to CURRMIN + CURR_COUNT from ROW_BUFFER)
				// reapir using repair again
				// reset the variables
			}
		}else{
			buffer.Update(BUFFER, currPkt)
		}
	}
}
func main() {
	go encoder()
	decoder()
}