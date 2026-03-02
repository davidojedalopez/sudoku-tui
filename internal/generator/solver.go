package generator

import (
	"github.com/davidojeda/sudoku-tui/internal/board"
)

// Solve solves the grid using constraint propagation + backtracking.
// Returns true if a solution was found. Modifies grid in place.
func Solve(grid *[9][9]int) bool {
	return solve(grid)
}

// CountSolutions counts solutions up to max (for uniqueness check).
func CountSolutions(grid [9][9]int, max int) int {
	count := 0
	countSolutions(&grid, max, &count)
	return count
}

func solve(grid *[9][9]int) bool {
	// Find the cell with minimum remaining values (MRV heuristic)
	bestR, bestC := -1, -1
	bestCount := 10

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if grid[r][c] != 0 {
				continue
			}
			cands := candidates(*grid, r, c)
			if len(cands) == 0 {
				return false // dead end
			}
			if len(cands) < bestCount {
				bestCount = len(cands)
				bestR, bestC = r, c
				if bestCount == 1 {
					goto found
				}
			}
		}
	}

found:
	if bestR == -1 {
		return true // solved
	}

	for _, val := range candidates(*grid, bestR, bestC) {
		grid[bestR][bestC] = val
		if solve(grid) {
			return true
		}
		grid[bestR][bestC] = 0
	}
	return false
}

func countSolutions(grid *[9][9]int, max int, count *int) {
	if *count >= max {
		return
	}

	bestR, bestC := -1, -1
	bestCount := 10

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if grid[r][c] != 0 {
				continue
			}
			cands := candidates(*grid, r, c)
			if len(cands) == 0 {
				return
			}
			if len(cands) < bestCount {
				bestCount = len(cands)
				bestR, bestC = r, c
				if bestCount == 1 {
					goto found
				}
			}
		}
	}

found:
	if bestR == -1 {
		*count++
		return
	}

	for _, val := range candidates(*grid, bestR, bestC) {
		grid[bestR][bestC] = val
		countSolutions(grid, max, count)
		grid[bestR][bestC] = 0
	}
}

func candidates(grid [9][9]int, r, c int) []int {
	used := [10]bool{}
	// Row
	for col := 0; col < 9; col++ {
		used[grid[r][col]] = true
	}
	// Col
	for row := 0; row < 9; row++ {
		used[grid[row][c]] = true
	}
	// Box
	boxR, boxC := (r/3)*3, (c/3)*3
	for row := boxR; row < boxR+3; row++ {
		for col := boxC; col < boxC+3; col++ {
			used[grid[row][col]] = true
		}
	}
	var result []int
	for v := 1; v <= 9; v++ {
		if !used[v] {
			result = append(result, v)
		}
	}
	return result
}

// SolveBoard solves a Board and stores solution.
func SolveBoard(b *board.Board) bool {
	grid := b.ToGrid()
	if !Solve(&grid) {
		return false
	}
	b.Solution = grid
	return true
}
