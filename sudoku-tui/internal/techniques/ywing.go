package techniques

// DetectYWing returns true if a Y-Wing (XY-Wing) pattern exists.
// Y-Wing: pivot cell with 2 candidates (A,B), pincer1 with (A,C), pincer2 with (B,C),
// where both pincers share a unit with the pivot.
func DetectYWing(cands [9][9]map[int]bool) bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			pivot := cands[r][c]
			if len(pivot) != 2 {
				continue
			}
			pivotVals := mapKeys(pivot)
			a, b := pivotVals[0], pivotVals[1]

			// Find pincers: cells with 2 candidates sharing one value with pivot.
			pincers := findPincers(cands, r, c, a, b)
			if len(pincers) >= 2 {
				return true
			}
		}
	}
	return false
}

func mapKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func findPincers(cands [9][9]map[int]bool, pivotR, pivotC, a, b int) [][3]int {
	// Returns pincers as [row, col, sharedVal].
	var pincers [][3]int

	seen := getPeers(pivotR, pivotC)

	for _, peer := range seen {
		r, c := peer[0], peer[1]
		if cands[r][c] == nil || len(cands[r][c]) != 2 {
			continue
		}
		// Pincer shares exactly one value with pivot.
		if cands[r][c][a] && !cands[r][c][b] {
			for v := range cands[r][c] {
				if v != a {
					pincers = append(pincers, [3]int{r, c, v})
				}
			}
		}
		if cands[r][c][b] && !cands[r][c][a] {
			for v := range cands[r][c] {
				if v != b {
					pincers = append(pincers, [3]int{r, c, v})
				}
			}
		}
	}
	return pincers
}

func getPeers(r, c int) [][2]int {
	seen := map[[2]int]bool{}
	var peers [][2]int

	for col := 0; col < 9; col++ {
		if col != c {
			p := [2]int{r, col}
			if !seen[p] {
				seen[p] = true
				peers = append(peers, p)
			}
		}
	}
	for row := 0; row < 9; row++ {
		if row != r {
			p := [2]int{row, c}
			if !seen[p] {
				seen[p] = true
				peers = append(peers, p)
			}
		}
	}
	boxR, boxC := (r/3)*3, (c/3)*3
	for row := boxR; row < boxR+3; row++ {
		for col := boxC; col < boxC+3; col++ {
			if row != r || col != c {
				p := [2]int{row, col}
				if !seen[p] {
					seen[p] = true
					peers = append(peers, p)
				}
			}
		}
	}
	return peers
}
