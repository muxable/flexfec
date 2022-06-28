package main

import (
	"fmt"

	"flexfec/bitstring"

	"github.com/pion/rtp"
)

// Extension RTP Header extension
// type Extension struct {
// 	id      uint8
// 	payload []byte
// }

// Header represents an RTP packet header
// type Header struct {
// 	Version          uint8
// 	Padding          bool
// 	Extension        bool
// 	Marker           bool
// 	PayloadType      uint8
// 	SequenceNumber   uint16
// 	Timestamp        uint32
// 	SSRC             uint32
// 	CSRC             []uint32
// 	ExtensionProfile uint16
// 	Extensions       []Extension
// }

// // Packet represents an RTP Packet
// type Packet struct {
// 	Header
// 	Payload     []byte
// 	PaddingSize byte
// }

//----------------------------------------------------------------------------
// func String(p *rtp.Packet) (out []byte) {

// 	// 1st byte (4 bits)
// 	out = append(out, uint8(p.Version))
// 	out[0] = out[0] << 6
// 	if p.Padding {
// 		out[0] = (out[0] | 1<<5)
// 	}
// 	if p.Extension {
// 		out[0] = (out[0] | 1<<4)
// 	}

// 	// 2nd byte
// 	out = append(out, uint8(0))
// 	if p.Marker {
// 		out[1] = (out[1] | 1<<7)
// 	}
// 	var payload_type = uint8(p.PayloadType) & 0x7F
// 	out[1] = out[1] | payload_type

// 	// 3rd & 4th bytes
// 	// second 16 bits - 2 bytes - length of source packet condition not given
// 	var paddingSize_int = int(p.PaddingSize)
// 	var length = len(p.CSRC) + len(p.Extensions) + len(p.Payload) + paddingSize_int
// 	out = append(out, 0, 0)
// 	binary.BigEndian.PutUint16(out[2:4], uint16(length))

// 	// 5th to 8th bytes (32 bits)
// 	out = append(out, 0, 0, 0, 0)
// 	binary.BigEndian.PutUint32(out[4:8], p.Timestamp)

// 	// 9th byte and the rest
// 	for _, s := range p.Payload {
// 		out = append(out, uint8(s))
// 	}
// 	out = append(out, uint8(paddingSize_int))
// 	return out
// }

func printHeader(buf []byte) {
	for index, value := range buf {
		for i := 7; i >= 0; i-- {
			fmt.Print((value >> i) & 1)
		}
		fmt.Print(" ")
		if (index+1)%4 == 0 {
			fmt.Println()
		}
	}
}

//---------------------------------------------------------------------------------------
func main() {
	rawPkt := []byte{
		0x90, 0xe0, 0x69, 0x8f, 0xd9, 0xc2, 0x93, 0xda, 0x1c, 0x64,
		0x27, 0x82, 0x00, 0x01, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0x98, 0x36, 0xbe, 0x88, 0x9e,
	}
	parsedPacket := &rtp.Packet{
		Header: rtp.Header{
			Padding:   true,
			Marker:    false,
			Extension: false,
			// ExtensionProfile: 1,
			// Extensions: []Extension{
			// 	{0, []byte{
			// 		0xFF, 0xFF, 0xFF, 0xFF,
			// 	}},
			// },
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
	printHeader(bitstring.ToBitString(parsedPacket))
	// fmt.Println()
	// printHeader(rawPkt[20:])
}
