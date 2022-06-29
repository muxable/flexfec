package main

import(
	"fmt"
	"flexfec/util"
	"flexfec/bitstring"
	"flexfec/fec_header"
)


func main() {
	packets := util.GenerateRTP(5, 1)

	util.PadPackets(&packets)

	var bitStrings [][]byte

	for _, pkt := range packets {
		bitStrings = append(bitStrings, bitstring.ToBitString(&pkt))
	}

	fecBitString := bitstring.ToFecBitString(bitStrings)

	fecHeader := fech.UnmarshalFec(fecBitString)

	fmt.Println(fecHeader)
	util.PrintBytes(fecHeader.Marshal())
}