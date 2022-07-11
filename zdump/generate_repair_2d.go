package recover

import "github.com/pion/rtp"

// L>0 and D>1 (in fecheader)
func GenerateRepair2dFec(srcBlock *[]rtp.Packet, L, D int) ([]rtp.Packet, []rtp.Packet) {
	rowFecPackets := GenerateRepairRowFec(srcBlock, L, D) // call with D>1( actual number of rows)
	colFecPackets := GenerateRepairColFec(srcBlock, L, D) // call with D>1( actual number of rows)

	return rowFecPackets, colFecPackets
}
