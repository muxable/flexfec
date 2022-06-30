package main

import (
	"flexfec/bitstring"
	fech "flexfec/fec_header"
	"flexfec/util"
	"fmt"
)

func main() {
	packets := util.GenerateRTP(5, 1)

	util.PadPackets(&packets)

	var bitStrings [][]byte

	for _, pkt := range packets {
		bitStrings = append(bitStrings, bitstring.ToBitString(&pkt))
	}

	fecBitString := bitstring.ToFecBitString(bitStrings)

	fecHeader, _ := fech.ToFecHeader(fecBitString)

	fmt.Println(fecHeader)
	util.PrintBytes(fecHeader.Marshal())
}
