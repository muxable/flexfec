package main

import (
	fech "flexfec/fec_header"
	"flexfec/util"
	"fmt"
)

func testFecHeaderLD() {
	var LDheader fech.FecHeader = fech.NewFecHeaderLD(false, true, false, false, 11, false, 127, 1234, 23432532, 345, 5, 6)
	buf1 := LDheader.Marshal()

	var resLD fech.FecHeaderLD = fech.FecHeaderLD{}
	resLD.Unmarshal(buf1)

	fmt.Println("LDheader\n", LDheader, "\n")
	util.PrintBytes(buf1)
	fmt.Println(resLD)

	fmt.Println("\n-----------------------------------------")
}

func testFecHeaderFlexibleMask() {
	var maskheader fech.FecHeader = fech.NewFecHeaderFlexibleMask(false, false, false, false, 11, false, 100, 500, 435343, 487, true, 36160, 3229756930, true, 13871700391609118210)
	buf2 := maskheader.Marshal()

	var resmask fech.FecHeaderFlexibleMask = fech.FecHeaderFlexibleMask{}
	resmask.Unmarshal(buf2)

	fmt.Println("Flexible mask header\n", maskheader, "\n")
	util.PrintBytes(buf2)
	fmt.Println(resmask)

	fmt.Println("\n-----------------------------------------")
}

func testFecHeaderRetransmission() {
	var retransheader fech.FecHeader = fech.NewFecHeaderRetransmission(true, false, false, false, 15, false, 34, 2342, 4296729, 4296729)
	buf3 := retransheader.Marshal()

	var resretrans fech.FecHeaderRetransmission = fech.FecHeaderRetransmission{}
	resretrans.Unmarshal(buf3)

	fmt.Println("Retransmission header\n", retransheader, "\n")
	util.PrintBytes(buf3)
	fmt.Println(resretrans)

	fmt.Println("\n-----------------------------------------")
}

func main() {
	testFecHeaderLD()
	testFecHeaderFlexibleMask()
	testFecHeaderRetransmission()
}
