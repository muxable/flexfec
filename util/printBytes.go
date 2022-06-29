package util

import (
	"fmt"
)

func PrintBytes(buf []byte) {
	for index, value := range buf {
		for i := 7; i >= 0; i-- {
			fmt.Print((value >> i) & 1)
		}
		fmt.Print(" ")
		if (index + 1) % 4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

