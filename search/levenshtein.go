package search

func LevenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)

	if n > m {
		r1, r2 = r2, r1
		n, m = m, n
	}

	row2 := make([]int, n+1)
	for j := range row2 {
		row2[j] = j
	}

	var row1 []int
	for i := 1; i < m+1; i++ {
		row1, row2 = row2, make([]int, n+1)
		row2[0] = i

		for j := 1; j < n+1; j++ {
			del := row2[j-1] + 1
			ins := row1[j] + 1
			sub := row1[j-1]

			if r1[j-1] != r2[i-1] {
				sub++
			}

			switch {
			case del <= ins && del <= sub:
				row2[j] = del
			case ins <= del && ins <= sub:
				row2[j] = ins
			default:
				row2[j] = sub
			}
		}
	}
	return row2[n]
}
