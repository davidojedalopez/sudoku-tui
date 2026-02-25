package generator

import (
	"math/rand"
)

// Difficulty levels with target clue counts.
type Difficulty int

const (
	Easy   Difficulty = iota // 36-40 clues
	Medium                   // 30-35 clues
	Hard                     // 25-29 clues
	Expert                   // 17-24 clues
)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	case Expert:
		return "Expert"
	}
	return "Unknown"
}

// TargetClues returns the target number of given clues for a difficulty.
func (d Difficulty) TargetClues() int {
	switch d {
	case Easy:
		return 38
	case Medium:
		return 32
	case Hard:
		return 27
	case Expert:
		return 22
	}
	return 32
}

// Generate creates a new puzzle grid with the given difficulty.
// Returns (puzzle, solution) as 81-char strings.
func Generate(diff Difficulty) (puzzle string, solution string) {
	grid := generateFilledGrid()
	sol := gridToString(grid)

	// Remove clues
	target := diff.TargetClues()
	puzzle81 := removeClues(grid, target)

	return puzzle81, sol
}

// generateFilledGrid creates a valid completed Sudoku grid.
func generateFilledGrid() [9][9]int {
	var grid [9][9]int

	// Fill 3 diagonal boxes randomly
	for box := 0; box < 3; box++ {
		fillBox(&grid, box*3, box*3)
	}

	// Solve the rest
	Solve(&grid)
	return grid
}

// fillBox fills a 3x3 box starting at (startR, startC) with random digits.
func fillBox(grid *[9][9]int, startR, startC int) {
	digits := rand.Perm(9)
	idx := 0
	for r := startR; r < startR+3; r++ {
		for c := startC; c < startC+3; c++ {
			grid[r][c] = digits[idx] + 1
			idx++
		}
	}
}

// removeClues removes cells from a completed grid until target clues remain,
// ensuring the puzzle still has a unique solution.
func removeClues(grid [9][9]int, target int) string {
	puzzle := grid

	// Create list of all positions
	positions := rand.Perm(81)

	clues := 81
	for _, pos := range positions {
		if clues <= target {
			break
		}
		r, c := pos/9, pos%9
		if puzzle[r][c] == 0 {
			continue
		}

		// Try removing this cell
		saved := puzzle[r][c]
		puzzle[r][c] = 0

		// Check unique solution
		if CountSolutions(puzzle, 2) == 1 {
			clues--
		} else {
			// Restore
			puzzle[r][c] = saved
		}
	}

	return gridToString(puzzle)
}

func gridToString(grid [9][9]int) string {
	b := make([]byte, 81)
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			v := grid[r][c]
			if v == 0 {
				b[r*9+c] = '0'
			} else {
				b[r*9+c] = byte('0' + v)
			}
		}
	}
	return string(b)
}
