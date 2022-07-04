package fech

import (
	"encoding/binary"
	"errors"
)
func ToFecHeaderLD(buf []byte)(FecHeaderLD, []byte){
	
/*
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |0|1|P|X|  CC   |M| PT recovery |         length recovery       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                          TS recovery                          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           SN base_i           |  L (columns)  |    D (rows)   |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |    ... next SN base and L/D for CSRC_i in CSRC list ...       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   :               Repair "Payload" follows FEC Header             :
   :                                                               :
*/
	
	var fecheader FecHeaderLD
	// first 2 bits are neglected in FEC bit string and replaced by R and F
	fecheader.R = false
	fecheader.F = false
	fecheader.P = (buf[0] >> 5 & 0x1) > 0
	fecheader.X = (buf[0] >> 4 & 0x1) > 0
	fecheader.CC = uint8((buf[0] & uint8(0xF)))

	fecheader.M = (buf[1] >> 7 & 0x1) > 0
	fecheader.PTRecovery = buf[1] & 0x7F

	fecheader.LengthRecovery = binary.BigEndian.Uint16(buf[2:4])

	fecheader.TimestampRecovery = binary.BigEndian.Uint32(buf[4:8])

	// Check: SN_base, L, D
	return fecheader, buf[8:]
}
// -------------------------------------------------------------------------
func ToFecHeaderFlexibleMask(buf []byte)(FecHeaderFlexibleMask, []byte){
	
/*
 	0                   1                   2                   3
      0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |0|0|P|X|  CC   |M| PT recovery |        length recovery        |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |                          TS recovery                          |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |           SN base_i           |k|          Mask [0-14]        |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |k|                   Mask [15-45] (optional)                   |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |                     Mask [46-109] (optional)                  |
     |                                                               |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     |   ... next SN base and Mask for CSRC_i in CSRC list ...       |
     +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
     :               Repair "Payload" follows FEC Header             :
     :                                                               :
*/

}
// -----------------------------------------------------------------------------
func ToFecHeaderRetransmission(buf []byte)(FecHeaderRetransmission, []byte){
	
/*
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |1|0|P|X|  CC   |M| Payload Type|        Sequence Number        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                           Timestamp                           |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                              SSRC                             |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   :          Retransmission "Payload" follows FEC Header          :
   :                                                               :
*/
	var fecheader FecHeaderRetransmission
	fecheader.R = true
	fecheader.F = false
	fecheader.P = (buf[0] >> 5 & 0x1) > 0
	fecheader.X = (buf[0] >> 4 & 0x1) > 0
	fecheader.CC = uint8((buf[0] & uint8(0xF)))

	fecheader.M = (buf[1] >> 7 & 0x1) > 0
	fecheader.PayloadType=buf[1] & 0x7F

	fecheader.SeqNumber=binary.BigEndian.Uint16(buf[2:4])

	fecheader.TimeStamp=binary.BigEndian.Uint32(buf[4:8])

	return fecheader,buf[8:]
}
// ----------------------------------------------------------------------------------------
// function to convert the FEC bit string (type []byte) to FEC header (type FecHeaderLD)
func ToFecHeader(buf []byte, fecvarient string) (FecHeader, []byte, err) {

	// using fecvarient to check 
		// L D
		// flexible mask
		// retransmission
	if(fecvarient=="LD"){
		header,body:=ToFecHeaderLD(buf)
		return header,body,nil;
	}
	if(fecvarient=="flexible mask"){
		header,body:=ToFecHeaderFlexibleMask(buf)
		return header,body,nil;
	}
	if(fecvarient=="retransmission"){
		header,body:=ToFecHeaderRetransmission(buf)
		return header,body,nil;
	}
	return nil,nil,errors.New("Fec varient is not defined correctly.")

}

// --------------------------------------------------------------------
