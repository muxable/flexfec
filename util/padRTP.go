package util

func PadBitStrings(bitstrings *[][]byte, length int){
	maxSize := -1
	n := len(*bitstrings)

	for _, bitstring := range *bitstrings {
		currSize := len(bitstring)
		if maxSize < currSize {
			maxSize = currSize
		}
	}

	if maxSize < length {
		maxSize = length
	}

	for i := 0; i < n; i++ {
		size := len((*bitstrings)[i])

		if size < maxSize {
			paddingSize := maxSize - size
			padding := make([]byte, paddingSize)
			(*bitstrings)[i] = append((*bitstrings)[i], padding...)
		}

	}

}

// func PadPackets(srcBlock *[]rtp.Packet) {
// 	maxSize := -1
// 	n := len(*srcBlock)

// 	for i := 0; i < n; i++ {
// 		currSize := len((*srcBlock)[i].Payload)
// 		if maxSize < currSize {
// 			maxSize = currSize
// 		}
// 	}

// 	for i := 0; i < n; i++ {
// 		size := len((*srcBlock)[i].Payload)

// 		if size < maxSize {
// 			leftOverPadBytes := uint8(maxSize - size)

// 			// (*srcBlock)[i].PaddingSize = uint8(leftOverPadBytes + 1)
// 			(*srcBlock)[i].PaddingSize = leftOverPadBytes

// 			// work like immutable entity, so replace with new slice
// 			payload := make([]byte, maxSize)
// 			copy(payload, (*srcBlock)[i].Payload)

// 			(*srcBlock)[i].Payload = payload
// 		} else {
// 			(*srcBlock)[i].Padding = false
// 		}
// 	}

// }



