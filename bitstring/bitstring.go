package bitstring

import (
	"encoding/binary"

	"github.com/pion/rtp"
)

// Converts rtp packet to bitstring as per fec scheme-20
func ToBitString(p *rtp.Packet) (out []byte) {

	// 1st byte (4 bits)
	out = append(out, uint8(p.Version))
	out[0] = out[0] << 6
	if p.Padding {
		out[0] = (out[0] | 1<<5)
	}
	if p.Extension {
		out[0] = (out[0] | 1<<4)
	}

	// 2nd byte
	out = append(out, uint8(0))
	if p.Marker {
		out[1] = (out[1] | 1<<7)
	}
	var payload_type = uint8(p.PayloadType) & 0x7F
	out[1] = out[1] | payload_type

	// 3rd & 4th bytes
	// second 16 bits - 2 bytes - length of source packet condition not given
	var paddingSize_int = int(p.PaddingSize)
	var length = len(p.CSRC) + len(p.Extensions) + len(p.Payload)
	out = append(out, 0, 0)
	binary.BigEndian.PutUint16(out[2:4], uint16(length))

	// 5th to 8th bytes (32 bits)
	out = append(out, 0, 0, 0, 0)
	binary.BigEndian.PutUint32(out[4:8], p.Timestamp)

	// 9th byte and the rest
	for _, s := range p.Payload {
		out = append(out, uint8(s))
	}
	out = append(out, uint8(paddingSize_int))
	return out
}

//---------------------------------------------------------------------------------------
