package techniques

// DetectHiddenSingle returns true if val is the only candidate for (r,c)
// in its row, column, or box.
func DetectHiddenSingle(cands [9][9]map[int]bool, r, c, val int) bool {
	// Check row.
	count := 0
	for col := 0; col < 9; col++ {
		if cands[r][col] != nil && cands[r][col][val] {
			count++
		}
	}
	if count == 1 {
		return true
	}

	// Check column.
	count = 0
	for row := 0; row < 9; row++ {
		if cands[row][c] != nil && cands[row][c][val] {
			count++
		}
	}
	if count == 1 {
		return true
	}

	// Check box.
	count = 0
	boxR, boxC := (r/3)*3, (c/3)*3
	for row := boxR; row < boxR+3; row++ {
		for col := boxC; col < boxC+3; col++ {
			if cands[row][col] != nil && cands[row][col][val] {
				count++
			}
		}
	}
	return count == 1
}
