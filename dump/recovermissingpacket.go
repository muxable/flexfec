package main

import (
	"fmt"

	"github.com/pion/rtp"
)

func missingPacket(srcBlock *[]rtp.Packet, repairPacket rtp.Packet)(rtp.Header, []byte){
// 5 6 7

// 7 6 5 have to check if it exists in arr

// Task 1
// [SN_base:SN_base+L-1] find missing rtp packet


// Task 2
// header : bitstring of srcblock headers ^ repairpacket header
// payload : bitstring of srcbloack payloads ^ repairpacket payload

// create new rtp packet 

// Calculate bitstring of srcBlock
// missing packet bit string : xor of bitstring with bit string of recoverypacket
}

// 1d 1 row
func recoverMissingPacket(srcBlock *[]rtp.Packet, repairPacket rtp.Packet) (rtp.Header, []byte) {

var fecheader fech.FecHeaderLD = fech.FecHeaderLD{}
fecheader.Unmarshal(repairPacket.Payload[:12])

SN_base:=fecheader.SN_base
L:=fecheader.L

lengthofsrcBlock=len(*srcBlock)
if len(lengthofsrcBlock)!=L{
	if (L-len(lengthofsrcBlock))>1{
		// retransmission
		fmt.Println("retransmission")
	}else{
		// recovery
		return missingPacket(srcBlock,repairPacket);
	}
}
else{
	// successful,  No error
	fmt.Println("All packets transmitted correctly")
}
}


sender struct-->byte array -=-------> array -> struct