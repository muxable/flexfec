package util


func GetTestCaseMap(variant int) map[int]int {

	if variant == 0 {
		/*
			a  b  c  X r1   0 1  2  X
			e  f  X  h r1   4 5  X  7
			X  j  k  l r3   X 9 10 11
			c1 c2 c3 c4
		*/

		return map[int]int {
			0 : 1, 1 : 1, 2 : 1, 4 : 1, 5 : 1, 7 : 1, 9 : 1, 10 : 1, 11 : 1, 
		}
	} else if variant == 1 {
		/*
			a  X  c  X r1   0 X  2  X
			e  f  g  h r1   4 5  6  7
			X  j  X  l r3   X 9  X 11
			c1 c2 c3 c4
		*/

		return map[int]int {
			0 : 1, 2 : 1, 4 : 1, 5 : 1, 6 : 1, 7 : 1, 9 : 1, 11 : 1, 
		}
	}

	/*
		a  X  X  X r1   0 X  X  X
		e  f  X  h r1   X 5  X  7
		X  j  k  l r3   X 9 10 11
		c1 c2 c3 c4
	*/

	return map[int]int {
		0 : 1, 5 : 1, 7 : 1, 9 : 1, 10 : 1, 11 : 1,
	}
	
	
	
}
