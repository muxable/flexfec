package bitstring

import(
	"github.com/pion/rtp"
)

func GetBlockBitstring(packets *[]rtp.Packet) [][]byte {
	var bitStrings [][]byte

	for _, pkt := range *packets {
		bitStrings = append(bitStrings, ToBitString(&pkt))
	}

	return bitStrings
}