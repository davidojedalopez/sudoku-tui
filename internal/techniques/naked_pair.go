package techniques

// DetectNakedPair returns true if a naked pair exists in the same unit as (r,c).
func DetectNakedPair(cands [9][9]map[int]bool, r, c int) bool {
	mySet := cands[r][c]
	if len(mySet) != 2 {
		return false
	}

	// Check row.
	for col := 0; col < 9; col++ {
		if col == c {
			continue
		}
		if mapsEqual(cands[r][col], mySet) {
			return true
		}
	}

	// Check column.
	for row := 0; row < 9; row++ {
		if row == r {
			continue
		}
		if mapsEqual(cands[row][c], mySet) {
			return true
		}
	}

	// Check box.
	boxR, boxC := (r/3)*3, (c/3)*3
	for row := boxR; row < boxR+3; row++ {
		for col := boxC; col < boxC+3; col++ {
			if row == r && col == c {
				continue
			}
			if mapsEqual(cands[row][col], mySet) {
				return true
			}
		}
	}
	return false
}

func mapsEqual(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}
