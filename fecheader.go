// This file includes code for generating bitstring with xor operation and creation of FEC header.

package main

import (
	"encoding/binary"
	"fmt"
	)
// Extension RTP Header extension
type Extension struct {
	id      uint8
	payload []byte
}

// Header represents an RTP packet header
type Header struct {
	Version          uint8
	Padding          bool
	Extension        bool
	Marker           bool
	PayloadType      uint8
	SequenceNumber   uint16
	Timestamp        uint32
	SSRC             uint32
	CSRC             []uint32
	ExtensionProfile uint16
	Extensions       []Extension
}

// Packet represents an RTP Packet
type Packet struct {
	Header
	Payload     []byte
	PaddingSize byte
}

//----------------------------------------------------------------------------
func String(p *Packet) (out []byte){

	// 1st byte (4 bits)
	out = append(out,uint8(p.Version))
	out[0] = out[0] << 6
	if p.Padding{
	out[0] = (out[0] | 1 << 5 )
	} 
	if p.Extension{
	out[0] = (out[0] | 1 << 4 )		
	} 

	// 2nd byte 
	out = append(out,uint8(0))	
	if p.Marker{
	out[1] = (out[1] | 1 << 7 )		
	}	
	var payload_type = uint8(p.PayloadType) & 0x7F
	out[1] = out[1] | payload_type
	
	// 3rd & 4th bytes
	// second 16 bits - 2 bytes - length of source packet condition not given
	var paddingSize_int = int(p.PaddingSize)
	var length = len(p.CSRC)+len(p.Extensions)+len(p.Payload)+paddingSize_int
	out = append(out,0,0)
	binary.BigEndian.PutUint16(out[2:4], uint16(length))		
	
	// 5th to 8th bytes (32 bits)
	out = append(out,0,0,0,0)	
	binary.BigEndian.PutUint32(out[4:8], p.Timestamp)

	// 9th byte and the rest
	for _, s := range p.Payload {
		out = append(out,uint8(s))
	}
	out = append(out,uint8(paddingSize_int))
	return out
}
// To be completed to check if fec bitstring output is correct
// func unMarshalString(buf []byte)(fecHeader FecHeaderFlexibleMask,err error){
// }


func printHeader(buf []byte) {
	for index, value := range buf {
		for i := 7; i >= 0; i-- {
			fmt.Print((value >> i) & 1)
		}
		fmt.Print(" ")
		if (index + 1) % 4 == 0 {
			fmt.Println()
		}
	}
}
func fecBitString(buf [][]byte) []byte {
	var xor_out []byte
	xor_out=append(buf[0])

	m:=len(xor_out)
	n:=len(buf)

	for i:=1;i<n;i++{
		for j:=0;j<m;j++{
			xor_out[j] ^= buf[i][j]
		}
	}
	return xor_out
}
func FECHeader(fecbitstring []byte) []byte {
	size:=12
	fecheader := make([]byte, size)
	// 1st byte
	// R - consider R=1
	fecheader[0]= (1<<7)
	// F - consider F=1
	fecheader[0] |= (1<<6)
	// P - from fecbitstring 3rd bit
	fecheader[0] |= ((fecbitstring[0] >> 5)<<5)
	// X - from fecbitstring 4th bit
	fecheader[0] |= ((fecbitstring[0] >> 4)<<4)
	// CC - from fecbitstring 5-8th bit
	fecheader[0] |= uint8((fecbitstring[0] & uint8(0xF)))

	// 2nd byte
	// initial allotment of 0 byte
	fecheader[1] = 0<<7
	// M - from fecbitstring 1 bit
	fecheader[1] |= ((fecbitstring[1]>>7)<<7)
	// PT - from fecbitstring 7 bits
	fecheader[1] |=((fecbitstring[1]<<1)>>1)

	// 3rd & 4th byte
	// same as FEC bit string
	fecheader[2] = fecbitstring[2]
	fecheader[3] = fecbitstring[3]

	// 5th to 8th byte
	// next 32 bits of the FEC bit string are written into the TS recovery field in the FEC header
	// TS recovery field
	fecheader[4]=fecbitstring[4]
	fecheader[5]=fecbitstring[5]
	fecheader[6]=fecbitstring[6]
	fecheader[7]=fecbitstring[7]

	// lowest Sequence Number -- YET TO BE DONE
	return fecheader

}
//---------------------------------------------------------------------------------------
func main() {
	rawPkt := []byte{
		0x90, 0xe0, 0x69, 0x8f, 0xd9, 0xc2, 0x93, 0xda, 0x1c, 0x64,
		0x27, 0x82, 0x00, 0x01, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0x98, 0x36, 0xbe, 0x88, 0x9e,
	}
	parsedPacket1 := &Packet{
		Header: Header{
			Padding:          true,
			Marker:           false,
			Extension:        false,
			ExtensionProfile: 1,
			Extensions: []Extension{
				{0, []byte{
					0xFF, 0xFF, 0xFF, 0xFF,
				}},
			},
			Version:        3,
			PayloadType:    15,
			SequenceNumber: 27023,
			Timestamp:      36,
			SSRC:           476325762,
			CSRC:           []uint32{},
		},
		Payload:     rawPkt[20:],
		PaddingSize: 5,
	}
	parsedPacket2 := &Packet{
		Header: Header{
			Padding:          true,
			Marker:           true,
			Extension:        true,
			ExtensionProfile: 1,
			Extensions: []Extension{
				{0, []byte{
					0xFF, 0xF1, 0xFF, 0xFF,
				}},
			},
			Version:        3,
			PayloadType:    9,
			SequenceNumber: 27023,
			Timestamp:      33,
			SSRC:           476325762,
			CSRC:           []uint32{},
		},
		Payload:     rawPkt[20:],
		PaddingSize: 5,
	}
	bitstring1:=String(parsedPacket1)
	bitstring2:=String(parsedPacket2)

	fmt.Println("packet 1")
	printHeader(bitstring1)
	fmt.Println("\npacket 2")
	printHeader(bitstring2)
	fmt.Println("\nAfter xor operation")

	var buf [][]byte

	buf=append(buf,	bitstring1)
	buf=append(buf,	bitstring2)

	// Input is type [][]byte
	// Output is of type []byte
	fecbitstring:=fecBitString(buf)
	printHeader(fecbitstring)

	fmt.Println("\nFEC header")
	fecheader:=FECHeader(fecbitstring)
	printHeader(fecheader)



}
