package main

import (
	fech "flexfec/fec_header"
	"fmt"
)

func printHeader(buf []byte) {
	for index, value := range buf {
		for i := 7; i >= 0; i-- {
			fmt.Print((value >> i) & 1)
		}
		fmt.Print(" ")
		if (index+1)%4 == 0 {
			fmt.Println()
		}
	}
}

func main() {
	var LDheader fech.FecHeader = fech.NewFecHeaderLD(false, true, false, false, 11, false, 127, 1234, 23432532, 345, 5, 6)
	buf1 := LDheader.Marshal()

	var resLD fech.FecHeaderLD = fech.FecHeaderLD{}
	resLD.Unmarshal(buf1)

	fmt.Println("LDheader\n", LDheader)
	printHeader(buf1)
	fmt.Println(resLD)

	fmt.Println()

	// -------------------------------------------------------------------------------------------------------------------

	var maskheader fech.FecHeader = fech.NewFecHeaderFlexibleMask(false, false, false, false, 11, false, 100, 500, 435343, 487, true, 127, true, [3]uint32{213123123, 213123123, 31231231})
	buf2 := maskheader.Marshal()

	var resmask fech.FecHeaderFlexibleMask = fech.FecHeaderFlexibleMask{}
	resmask.Unmarshal(buf2)

	fmt.Println("Flexible mask header\n", maskheader)
	printHeader(buf2)
	fmt.Println(resmask)

	fmt.Println()

	// -------------------------------------------------------------------------------------------------------------------

	var retransheader fech.FecHeader = fech.NewFecHeaderRetransmission(true, false, false, false, 15, false, 34, 2342, 4296729, 4296729)
	buf3 := retransheader.Marshal()

	var resretrans fech.FecHeaderRetransmission = fech.FecHeaderRetransmission{}
	resretrans.Unmarshal(buf3)

	fmt.Println("Retransmission header\n", retransheader)
	printHeader(buf3)
	fmt.Println(resretrans)
}
