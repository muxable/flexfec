package main

import(
	"fmt"
	"flexfec/util"
	"flexfec/bitstring"
)

func main(){
	packets := util.GenerateRTP(1,3)

	bitstrings := [][]byte{}

	for _, pkt := range packets {
		fmt.Println(util.PrintPkt(pkt))
		
		bitstr := bitstring.ToBitString(&pkt)
		bitstrings = append(bitstrings, bitstr)
	}

	util.PadBitStrings(&bitstrings, -1)

	for i, bitstr := range bitstrings {
		fmt.Println("Bitstring ", i + 1)
		util.PrintBytes(bitstr)
		fmt.Println()
	}

	fecBitstring := bitstring.ToFecBitString(&bitstrings)

	fmt.Println("FECbitstring : ")
	util.PrintBytes(fecBitstring)
}