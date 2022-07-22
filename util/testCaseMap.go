package util


func GetTestCaseMap(variant int) map[int]int {

	if variant == 0 {
		/*
			0 1  2  X
			4 5  X  7
			X 9 10 11
			c1 c2 c3 c4
		*/

		return map[int]int {
			0 : 1,
			1 : 1, 
			2 : 1, 
			// 3 : 1, 
			4 : 1,
			5 : 1, 
			6 : 1,
			7 : 1, 
			// 8 : 1,
			9 : 1, 
			10 : 1, 
			11 : 1, 
		}
	} else if variant == 1 {
		/*
			0 X  2  X
			4 5  6  7
			X 9  X 11
			c1 c2 c3 c4
		*/

		return map[int]int {
			0 : 1,
			// 1 : 1,
			2 : 1,
			// 3 : 1, 
			4 : 1, 
			5 : 1, 
			6 : 1, 
			7 : 1, 
			// 8 : 1,
			9 : 1, 
			// 10 : 1,
			11 : 1, 
		}
	}

	/*
		X X  2  3 |0
		4 X  X  7 |1
		8 9  X  X |2
		3 4  5  6 
	*/

	return map[int]int {
		// 0 : 1,
		// 1 : 1,
		2 : 1,
		3 : 1, 
		4 : 1, 
		// 5 : 1, 
		// 6 : 1, 
		7 : 1, 
		8 : 1, 
		9 : 1, 
		// 10 : 1, 
		// 11 : 1,
	}
	
	
	
}
