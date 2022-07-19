package main

import (
	"flexfec/recover"
	"flexfec/bitstring"
	"flexfec/util"
	"fmt"
)

func testrow() {
	srcBlock := util.GenerateRTP(4, 3)
	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)
	bitsrings := bitstring.GetBlockBitstring(srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPacketsRow := recover.GenerateRepairLD(&bitsrings, 4, 3, 0, SN_Base)

	for i := 0; i < len(srcBlock); i++ {
		fmt.Println(util.PrintPkt(srcBlock[i]))
	}

	fmt.Println("-----------------------------------------")

	for _, rowRepair := range repairPacketsRow {
		fmt.Println(util.PrintPkt(rowRepair))
	}


}

func testcol() {
	srcBlock := util.GenerateRTP(4, 3)
	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)
	bitsrings := bitstring.GetBlockBitstring(srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPacketsCol := recover.GenerateRepairLD(&bitsrings, 4, 3, 1, SN_Base)

	for i := 0; i < len(srcBlock); i++ {
		fmt.Println(util.PrintPkt(srcBlock[i]))
	}

	fmt.Println("-----------------------------------------")

	for _, rowRepair := range repairPacketsCol {
		fmt.Println(util.PrintPkt(rowRepair))
	}


}

func test2D() {
	srcBlock := util.GenerateRTP(4, 3)
	SN_Base := uint16(srcBlock[0].Header.SequenceNumber)
	bitsrings := bitstring.GetBlockBitstring(srcBlock)
	util.PadBitStrings(&bitsrings, -1)

	repairPackets2D := recover.GenerateRepairLD(&bitsrings, 4, 3, 2, SN_Base)

	for i := 0; i < len(srcBlock); i++ {
		fmt.Println(util.PrintPkt(srcBlock[i]))
	}

	fmt.Println("-----------------------------------------")

	for _, rowRepair := range repairPackets2D {
		fmt.Println(util.PrintPkt(rowRepair))
	}

}

func main() {
	// testrow()
	// testcol()
	test2D()
}
