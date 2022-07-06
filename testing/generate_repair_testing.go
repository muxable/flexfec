package main

import (
	"flexfec/recover"
	"flexfec/util"
	"fmt"
)

func main() {

	// sender
	packets := util.GenerateRTP(5, 3) // L=5 and D=3
	util.PadPackets(&packets)

	//  L>0, D=0 Row Fec
	/*
		fmt.Println("Source Packets:")
		for i := 0; i < 3; i++ {
			fmt.Println("Row", i+1, "\n")
			for j := 0; j < 5; j++ {
				util.PrintPkt(packets[i*5+j])
			}
		}
		fmt.Println()

		repairPacketsRow := recover.GenerateRepairRowFec(&packets, 5, 0)

		fmt.Println("repair packets:")
		for _, r_packet := range repairPacketsRow {
			util.PrintPkt(r_packet)
		}
	*/

	// ___________________________________________________________________

	//  L>0, D=1 Col Fec
	/*
		fmt.Println("Source Packets:")
		for i := 0; i < 5; i++ {
			fmt.Println("Col", i+1, "\n")
			for j := 0; j < 3; j++ {
				util.PrintPkt(packets[j*3+i])
			}
		}
		fmt.Println()

		repairPacketsCol := recover.GenerateRepairColFec(&packets, 5, 1)

		fmt.Println("repair packets:")
		for _, r_packet := range repairPacketsCol {
			util.PrintPkt(r_packet)
		}
	*/

	// ________________________________________________________________

	// fmt.Println("Source Packets:")

	rowFecPackets, colFecPackets := recover.GenerateRepair2dFec(&packets, 5, 3)

	fmt.Println("Row repair:")
	for _, r_packet := range rowFecPackets {
		util.PrintPkt(r_packet)
	}

	fmt.Println("Col repair:")
	for _, c_packet := range colFecPackets {
		util.PrintPkt(c_packet)
	}
}
