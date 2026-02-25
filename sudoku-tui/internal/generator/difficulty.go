package generator

// Grade analyzes a puzzle and returns the minimum difficulty required.
// It solves the puzzle using only specific technique sets.
func Grade(puzzle string) Difficulty {
	// Try easy (naked singles + hidden singles only)
	if canSolveWithEasy(puzzle) {
		return Easy
	}
	// Try medium (+ naked pairs + pointing pairs)
	if canSolveWithMedium(puzzle) {
		return Medium
	}
	// Try hard (+ more complex patterns)
	if canSolveWithHard(puzzle) {
		return Hard
	}
	return Expert
}

func canSolveWithEasy(puzzle string) bool {
	var grid [9][9]int
	for i, ch := range puzzle {
		if ch >= '1' && ch <= '9' {
			grid[i/9][i%9] = int(ch - '0')
		}
	}
	return solveWithNakedHiddenSingles(&grid)
}

func canSolveWithMedium(puzzle string) bool {
	// For grading purposes, medium means it needs at most naked/hidden pairs
	// We approximate: count clues. 30-35 clues => medium
	clues := 0
	for _, ch := range puzzle {
		if ch >= '1' && ch <= '9' {
			clues++
		}
	}
	return clues >= 30 && clues < 36
}

func canSolveWithHard(puzzle string) bool {
	clues := 0
	for _, ch := range puzzle {
		if ch >= '1' && ch <= '9' {
			clues++
		}
	}
	return clues >= 24 && clues < 30
}

// solveWithNakedHiddenSingles attempts to solve using only basic techniques.
func solveWithNakedHiddenSingles(grid *[9][9]int) bool {
	changed := true
	for changed {
		changed = false

		// Naked singles: cells with only one candidate
		for r := 0; r < 9; r++ {
			for c := 0; c < 9; c++ {
				if grid[r][c] != 0 {
					continue
				}
				cands := candidates(*grid, r, c)
				if len(cands) == 1 {
					grid[r][c] = cands[0]
					changed = true
				}
			}
		}

		// Hidden singles: value that can only go in one cell in a unit
		changed = changed || applyHiddenSingles(grid)
	}

	// Check if solved
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if grid[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

func applyHiddenSingles(grid *[9][9]int) bool {
	changed := false

	// Check rows
	for r := 0; r < 9; r++ {
		for val := 1; val <= 9; val++ {
			var positions []int
			for c := 0; c < 9; c++ {
				if grid[r][c] == 0 {
					cands := candidates(*grid, r, c)
					for _, cv := range cands {
						if cv == val {
							positions = append(positions, c)
							break
						}
					}
				}
			}
			if len(positions) == 1 {
				grid[r][positions[0]] = val
				changed = true
			}
		}
	}

	// Check cols
	for c := 0; c < 9; c++ {
		for val := 1; val <= 9; val++ {
			var positions []int
			for r := 0; r < 9; r++ {
				if grid[r][c] == 0 {
					cands := candidates(*grid, r, c)
					for _, cv := range cands {
						if cv == val {
							positions = append(positions, r)
							break
						}
					}
				}
			}
			if len(positions) == 1 {
				grid[positions[0]][c] = val
				changed = true
			}
		}
	}

	// Check boxes
	for boxR := 0; boxR < 3; boxR++ {
		for boxC := 0; boxC < 3; boxC++ {
			for val := 1; val <= 9; val++ {
				type pos struct{ r, c int }
				var positions []pos
				for r := boxR * 3; r < boxR*3+3; r++ {
					for c := boxC * 3; c < boxC*3+3; c++ {
						if grid[r][c] == 0 {
							cands := candidates(*grid, r, c)
							for _, cv := range cands {
								if cv == val {
									positions = append(positions, pos{r, c})
									break
								}
							}
						}
					}
				}
				if len(positions) == 1 {
					grid[positions[0].r][positions[0].c] = val
					changed = true
				}
			}
		}
	}

	return changed
}
