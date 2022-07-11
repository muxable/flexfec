package recover

/*
func recover_row_packets(srcBlock *[]rtp.Packet, repairPacketRows *[]rtp.Packet, L int, D int) int {

}

func recover_column_packets(srcBlock *[]rtp.Packet, repairPacketColumns *[]rtp.Packet, L int, D int) int {

}


// check if srcBlock is ordered or unordered. Assumed ordered [] for now.
func IterativeDecoding(srcBlock *[]rtp.Packet, repairPacketRows *[]rtp.Packet, repairPacketColumns *[]rtp.Packet, L int, D int) []rtp.Packet {

	num_recovered_until_this_iteration := 0
	num_recovered_so_far := 0

	for {

		num_recovered_so_far += recover_row_packets(&srcBlock, &repairPacketRows, L, D)
		num_recovered_so_far += recover_column_packets(&srcBlock, &repairPacketColumns, L, D)

		if num_recovered_so_far <= num_recovered_until_this_iteration {
			break
		}
		num_recovered_until_this_iteration = num_recovered_so_far

		// Row wise recovery of missing packets
		// num_recovered_so_far++
		// can use MissingPacket function for recover_missing_packet.go

		// column wise recovery of missing packets
		// num_recovered_so_far++

		// if num_recovered_until_this_iteration<num_recovered_so_far
		// then num_recovered_until_this_iteration=num_recovered_so_far
		// Reiterate

		// else terminate / return
	}
	return rtp
}

*/
