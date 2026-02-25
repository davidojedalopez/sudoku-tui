package board

// IsValidPlacement returns true if placing val at (r, c) doesn't conflict.
func IsValidPlacement(grid [9][9]int, r, c, val int) bool {
	// Check row
	for col := 0; col < 9; col++ {
		if col != c && grid[r][col] == val {
			return false
		}
	}
	// Check col
	for row := 0; row < 9; row++ {
		if row != r && grid[row][c] == val {
			return false
		}
	}
	// Check 3x3 box
	boxR, boxC := (r/3)*3, (c/3)*3
	for row := boxR; row < boxR+3; row++ {
		for col := boxC; col < boxC+3; col++ {
			if (row != r || col != c) && grid[row][col] == val {
				return false
			}
		}
	}
	return true
}

// IsValidGrid returns true if a completed grid has no conflicts.
func IsValidGrid(grid [9][9]int) bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			v := grid[r][c]
			if v == 0 {
				return false
			}
			grid[r][c] = 0
			if !IsValidPlacement(grid, r, c, v) {
				return false
			}
			grid[r][c] = v
		}
	}
	return true
}
