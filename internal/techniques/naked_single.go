package techniques

// IsNakedSingle returns true if the cell at (r,c) has exactly one candidate.
func IsNakedSingle(cands [9][9]map[int]bool, r, c int) bool {
	return len(cands[r][c]) == 1
}
