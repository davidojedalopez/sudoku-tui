package techniques

// DetectXWing returns true if an X-Wing pattern exists for the given value.
// X-Wing: val appears in exactly 2 cells in two different rows,
// and those cells share the same two columns.
func DetectXWing(cands [9][9]map[int]bool, val int) bool {
	type rowCols struct {
		row  int
		cols [2]int
	}
	var candidates []rowCols

	for r := 0; r < 9; r++ {
		var cols []int
		for c := 0; c < 9; c++ {
			if cands[r][c] != nil && cands[r][c][val] {
				cols = append(cols, c)
			}
		}
		if len(cols) == 2 {
			candidates = append(candidates, rowCols{r, [2]int{cols[0], cols[1]}})
		}
	}

	// Check if any two rows share the same column pair.
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].cols == candidates[j].cols {
				return true
			}
		}
	}
	return false
}
