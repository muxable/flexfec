package bitstring

// Computes the xor of all the packets in the input array
func ToFecBitString(buf [][]byte) []byte {
	var buf_xor []byte
	buf_xor=append(buf[0])

	m:=len(buf_xor)
	n:=len(buf)

	for i:=1;i<n;i++{
		for j:=0;j<m;j++{
			// xor operation
			buf_xor[j] ^= buf[i][j]
		}
	}
	return buf_xor
}
// ------------------------------------------------------------