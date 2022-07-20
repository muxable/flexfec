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

type FecHeader interface {
	Marshal() []byte
	Unmarshal(buf []byte)
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
type FecHeaderLD struct {
	R                 bool
	F                 bool
	P                 bool
	X                 bool
	CC                uint8
	M                 bool
	PTRecovery        uint8
	LengthRecovery    uint16
	TimestampRecovery uint32
	SN_base           uint16
	L                 uint8
	D                 uint8
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

// function to convert the FEC bit string (type []byte) to FEC header (type FecHeaderLD) 
func ToFecHeader(buf []byte)(FecHeaderLD){
	var fecheader FecHeaderLD
	fecheader.R=false
	fecheader.F=false
	fecheader.P=(buf[0] >> 5 & 0x1) > 0
	fecheader.X=(buf[0] >> 4 & 0x1) > 0
	fecheader.CC = uint8((buf[0] & uint8(0xF)))
	
	fecheader.M = (buf[1] >> 7 & 0x1) > 0
	fecheader.PTRecovery = buf[1] & 0x7F

	fecheader.LengthRecovery = binary.BigEndian.Uint16(buf[2:4])
	fecheader.TimestampRecovery = binary.BigEndian.Uint32(buf[4:8])
	
	return fecheader
}
//---------------------------------------------------------------------------------------
func NewFecHeaderLD(R bool, F bool, P bool, X bool, CC uint8, M bool, PTRecovery uint8, LengthRecovery uint16, TimestampRecovery uint32, SN_base uint16, L uint8, D uint8) FecHeader {
	return &FecHeaderLD{
		R:                 R,
		F:                 F,
		P:                 P,
		X:                 X,
		M:                 M,
		L:                 L,
		D:                 D,
		CC:                CC,
		SN_base:           SN_base,
		PTRecovery:        PTRecovery,
		LengthRecovery:    LengthRecovery,
		TimestampRecovery: TimestampRecovery,
	}
}

func (fh *FecHeaderLD) Marshal() []byte {
	// size to be detarmined later , for now for L and D variant size = 12 bytes
	size := 12

	buf := make([]byte, size)

	if fh.R {
		buf[0] = (1 << 7)
	}

	if fh.F {
		buf[0] |= (1 << 6)
	}

	if fh.P {
		buf[0] |= (1 << 5)
	}

	if fh.X {
		buf[0] |= (1 << 4)
	}

	buf[0] |= fh.CC

	if fh.M {
		buf[1] = (1 << 7)
	}

	buf[1] |= (fh.PTRecovery & 0x7F)

	binary.BigEndian.PutUint16(buf[2:4], fh.LengthRecovery)
	binary.BigEndian.PutUint32(buf[4:8], fh.TimestampRecovery)

	binary.BigEndian.PutUint16(buf[8:10], fh.SN_base)

	buf[10] = fh.L
	buf[11] = fh.D

	return buf
}

func (fh *FecHeaderLD) Unmarshal (buf []byte) {
	fh.R = (buf[0] >> 7 & 0x1) > 0
	fh.F = (buf[0] >> 6 & 0x1) > 0
	fh.P = (buf[0] >> 5 & 0x1) > 0
	fh.X = (buf[0] >> 4 & 0x1) > 0
	fh.CC = uint8((buf[0] & uint8(0xF)))
	fh.M = (buf[0] >> 7 & 0x1) > 0

	fh.PTRecovery = buf[1] & 0x7F

	fh.LengthRecovery = binary.BigEndian.Uint16(buf[2:4])
	fh.TimestampRecovery = binary.BigEndian.Uint32(buf[4:8])

	fh.SN_base = binary.BigEndian.Uint16(buf[8:10])
	fh.L = buf[10]
	fh.D = buf[11]

}
// --------------------------------------------------------------------------------------
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
	fecheader:=ToFecHeader(fecbitstring)
	fmt.Println(fecheader)

	fecheaderBits := fecheader.Marshal()

	// resfh, _ := Unmarshal(fecheaderBits)
	resfh:=Unmarshal(fecheaderBits)

	for _, BYTE := range fecheaderBits {
		printBits(BYTE)
	}

	fmt.Println(resfh)


}

func Unmarshal(fecheaderBits []byte) {
	panic("unimplemented")
}


