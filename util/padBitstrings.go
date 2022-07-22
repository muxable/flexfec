package util
// import "fmt"

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
		// fmt.Println("-------------------------max length updated------------------")
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