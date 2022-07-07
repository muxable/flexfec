package main

import (
	"encoding/binary"
	"flexfec/recover"
	"flexfec/util"
	"fmt"
)

func main() {

	// sender
	packets := util.GenerateRTP(4, 3) // L=4 and D=3
	util.PadPackets(&packets)

	//  L>0, D=0 Row Fec

	fmt.Println("Source Packets:")
	for i := 0; i < 3; i++ {
		fmt.Println("Row", i+1, "\n")
		for j := 0; j < 4; j++ {
			util.PrintPkt(packets[i*4+j])
		}
	}
	fmt.Println()

	repairPacketsRow := recover.GenerateRepairRowFec(&packets, 4, 0)

	fmt.Println("repair packets:")
	for _, r_packet := range repairPacketsRow {
		SN_base := binary.BigEndian.Uint16(r_packet.Payload[8:10])
		fmt.Println("SNbase,L,D :", SN_base)

		fmt.Println()
		util.PrintPkt(r_packet)

	}

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
	/*
		rowFecPackets, colFecPackets := recover.GenerateRepair2dFec(&packets, 5, 3)

		fmt.Println("Row repair:")
		for _, r_packet := range rowFecPackets {
			util.PrintPkt(r_packet)
		}

		fmt.Println("Col repair:")
		for _, c_packet := range colFecPackets {
			util.PrintPkt(c_packet)
		}
	*/
}
