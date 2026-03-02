package techniques

// DetectPointingPair returns true if pointing pairs exist for the given value.
// A pointing pair occurs when a value in a box is constrained to one row/col,
// allowing elimination from that row/col outside the box.
func DetectPointingPair(cands [9][9]map[int]bool, val int) bool {
	for boxR := 0; boxR < 3; boxR++ {
		for boxC := 0; boxC < 3; boxC++ {
			// Collect rows and cols where val can go in this box.
			rows := map[int]bool{}
			cols := map[int]bool{}
			for r := boxR * 3; r < boxR*3+3; r++ {
				for c := boxC * 3; c < boxC*3+3; c++ {
					if cands[r][c] != nil && cands[r][c][val] {
						rows[r] = true
						cols[c] = true
					}
				}
			}

			// Pointing pair in a row: val is locked to a single row within the box.
			if len(rows) == 1 {
				row := -1
				for r := range rows {
					row = r
				}
				// Check if val also exists in this row outside this box.
				for c := 0; c < 9; c++ {
					if c/3 != boxC && cands[row][c] != nil && cands[row][c][val] {
						return true
					}
				}
			}

			// Pointing pair in a col: val is locked to a single column within the box.
			if len(cols) == 1 {
				col := -1
				for c := range cols {
					col = c
				}
				// Check if val also exists in this col outside this box.
				for r := 0; r < 9; r++ {
					if r/3 != boxR && cands[r][col] != nil && cands[r][col][val] {
						return true
					}
				}
			}
		}
	}
	return false
}
