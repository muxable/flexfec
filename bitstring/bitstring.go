package bitstring

import (
	"encoding/binary"

	"github.com/pion/rtp"
)

// Converts rtp packet to bitstring as per fec scheme-20
func ToBitString(pkt *rtp.Packet) (out []byte) {
	buf, _ := pkt.Marshal()
	length := uint16(len(buf))

	// replace SN with length
	binary.BigEndian.PutUint16(buf[2:4], length)

	// remove SSRC
	bitstring := []byte{}
	bitstring = append(bitstring, buf[:8]...)
	bitstring = append(bitstring, buf[12:]...)

	return bitstring
}

