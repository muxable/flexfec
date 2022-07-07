package main

import (
	"flexfec/recover"
	"flexfec/util"
	"fmt"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	White = "\033[37m"
	Blue  = "\033[34m"
)

func main() {

	// sender
	packets := util.GenerateRTP(4, 3) // L=4 and D=3
	util.PadPackets(&packets)

	fmt.Println(string(Green), "source packets")

	for i := 0; i < 4; i++ {
		fmt.Println("Col", i+1)
		for j := 0; j < 3; j++ {
			util.PrintPkt(packets[j*4+i])
		}
	}

	//  L>0, D=0 Row Fec
	repairPackets := recover.GenerateRepairLD(&packets, 4, 3)

	fmt.Println(string(Blue), "repair packets")

	for i := 0; i < len(repairPackets); i++ {
		util.PrintPkt(repairPackets[i])
	}

}
