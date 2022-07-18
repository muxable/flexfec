package recover

import(
	"flexfec/bitstring"
	"github.com/pion/rtp"
)

func GetBlockBitstring(packets []rtp.Packet) [][]byte {
	var bitStrings [][]byte

	for _, pkt := range packets {
		bitStrings = append(bitStrings, bitstring.ToBitString(&pkt))
	}

	return bitStrings
}