// No need of this fuunction as algorithm was changed to decrease latency. Can be used in case if required.





// package recover

// import(
// 	"github.com/pion/rtp"
// )

// // TODO --------------------------------------------
// func recover_row_packets(srcBlock *[]rtp.Packet, repairPacketRows *[]rtp.Packet, L int,D int) int {
// 	num_recovered_so_far:=0
// 	// Row wise recovery of missing packets
// 	// num_recovered_so_far++
// 	// can use MissingPacket function for recover_missing_packet.go
// 	// return num_recovered_so_far

// 	// Check on how to seperate each row from srcBlock
// 	return num_recovered_so_far
// }

// func recover_column_packets(srcBlock *[]rtp.Packet, repairPacketColumns *[]rtp.Packet, L int,D int) int{
// 	num_recovered_so_far:=0

// 	// column wise recovery of missing packets
// 	// num_recovered_so_far++
// 	// return num_recovered_so_far++

// 	// Check on how to seperate each column from srcBlock
// 	return num_recovered_so_far

// }
// // ----------------------------------------------------

// // check if srcBlock is ordered or unordered. Assumed ordered [] for now.
// func IterativeDecoding(srcBlock *[]rtp.Packet, repairPacketRows *[]rtp.Packet, repairPacketColumns *[]rtp.Packet, L int, D int) []rtp.Packet {

// 	num_recovered_until_this_iteration := 0
// 	num_recovered_so_far := 0
	
// 	for true{
// 		num_recovered_so_far+=recover_row_packets(srcBlock, repairPacketRows,L,D)
// 		num_recovered_so_far+=recover_column_packets(srcBlock, repairPacketColumns,L,D)

// 		if num_recovered_so_far <= num_recovered_until_this_iteration {
// 			break
// 		}
// 		num_recovered_until_this_iteration = num_recovered_so_far

// 	}
// 	return 
// }
