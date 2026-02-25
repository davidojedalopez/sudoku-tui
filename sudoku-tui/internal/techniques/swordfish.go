package techniques

// DetectSwordfish returns true if a Swordfish pattern exists for the given value.
// Swordfish: val appears in 2-3 cells in three different rows,
// and those cells together cover exactly 3 columns.
func DetectSwordfish(cands [9][9]map[int]bool, val int) bool {
	type rowInfo struct {
		row  int
		cols []int
	}
	var rowInfos []rowInfo

	for r := 0; r < 9; r++ {
		var cols []int
		for c := 0; c < 9; c++ {
			if cands[r][c] != nil && cands[r][c][val] {
				cols = append(cols, c)
			}
		}
		if len(cols) >= 2 && len(cols) <= 3 {
			rowInfos = append(rowInfos, rowInfo{r, cols})
		}
	}

	if len(rowInfos) < 3 {
		return false
	}

	// Check all combinations of 3 rows.
	for i := 0; i < len(rowInfos); i++ {
		for j := i + 1; j < len(rowInfos); j++ {
			for k := j + 1; k < len(rowInfos); k++ {
				colSet := map[int]bool{}
				for _, c := range rowInfos[i].cols {
					colSet[c] = true
				}
				for _, c := range rowInfos[j].cols {
					colSet[c] = true
				}
				for _, c := range rowInfos[k].cols {
					colSet[c] = true
				}
				if len(colSet) == 3 {
					return true
				}
			}
		}
	}
	return false
}
