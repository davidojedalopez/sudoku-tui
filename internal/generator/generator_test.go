package generator

import (
	"testing"

	"github.com/davidojeda/sudoku-tui/internal/board"
)

func TestGenerateAndSolve(t *testing.T) {
	puzzle, solution := Generate(Easy)
	if len(puzzle) != 81 {
		t.Errorf("expected 81 chars, got %d", len(puzzle))
	}
	if len(solution) != 81 {
		t.Errorf("expected 81 chars solution, got %d", len(solution))
	}

	// Verify solution is valid
	var grid [9][9]int
	for i, ch := range solution {
		if ch >= '1' && ch <= '9' {
			grid[i/9][i%9] = int(ch - '0')
		}
	}
	if !board.IsValidGrid(grid) {
		t.Error("generated solution is not a valid Sudoku")
	}
}

func TestUniqueSolution(t *testing.T) {
	puzzle, _ := Generate(Medium)
	var grid [9][9]int
	for i, ch := range puzzle {
		if ch >= '1' && ch <= '9' {
			grid[i/9][i%9] = int(ch - '0')
		}
	}
	count := CountSolutions(grid, 2)
	if count != 1 {
		t.Errorf("expected unique solution, got %d solutions", count)
	}
}

func TestDifficultyClueCount(t *testing.T) {
	for _, diff := range []Difficulty{Easy, Medium, Hard, Expert} {
		puzzle, _ := Generate(diff)
		clues := 0
		for _, ch := range puzzle {
			if ch >= '1' && ch <= '9' {
				clues++
			}
		}
		t.Logf("Difficulty %s: %d clues", diff, clues)
	}
}

func TestStressGenerate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}
	for i := 0; i < 20; i++ {
		puzzle, solution := Generate(Hard)

		// Solution must be valid
		var solGrid [9][9]int
		for j, ch := range solution {
			if ch >= '1' && ch <= '9' {
				solGrid[j/9][j%9] = int(ch - '0')
			}
		}
		if !board.IsValidGrid(solGrid) {
			t.Errorf("iteration %d: invalid solution", i)
		}

		// Puzzle must have unique solution
		var puzGrid [9][9]int
		for j, ch := range puzzle {
			if ch >= '1' && ch <= '9' {
				puzGrid[j/9][j%9] = int(ch - '0')
			}
		}
		if n := CountSolutions(puzGrid, 2); n != 1 {
			t.Errorf("iteration %d: expected 1 solution, got %d", i, n)
		}
	}
}
