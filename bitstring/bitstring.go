package bitstring

import (
	"encoding/binary"
	"github.com/pion/rtp"
	"fmt"
)

// Converts rtp packet to bitstring as per fec scheme-20
func ToBitString(pkt *rtp.Packet) (out []byte) {
	buf, err := pkt.Marshal()

	if err != nil {
		fmt.Println(err)
	}

	length := uint16(len(buf))

	// replace SN with length
	binary.BigEndian.PutUint16(buf[2:4], length)

	// remove SSRC
	bitstring := make([]byte, length - 4)
	copy(bitstring[:8], buf[:8])
	copy(bitstring[8:], buf[12:])

	return bitstring
}

