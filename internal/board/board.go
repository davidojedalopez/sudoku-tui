package board

// Board is a 9x9 Sudoku board.
type Board struct {
	Cells    [9][9]Cell
	Solution [9][9]int
}

// New creates a new empty board.
func New() *Board {
	return &Board{}
}

// Set places a user digit at (row, col). Does nothing if cell is KindGiven.
func (b *Board) Set(row, col, val int) {
	if b.Cells[row][col].Kind == KindGiven {
		return
	}
	b.Cells[row][col].Value = val
	b.Cells[row][col].Kind = KindUser
	b.Cells[row][col].Notes = [9]bool{}
	b.updateConflicts()
}

// Erase clears a user cell at (row, col).
func (b *Board) Erase(row, col int) {
	if b.Cells[row][col].Kind == KindGiven {
		return
	}
	b.Cells[row][col].Value = 0
	b.Cells[row][col].Kind = KindEmpty
	b.Cells[row][col].Notes = [9]bool{}
	b.Cells[row][col].Conflict = false
	b.updateConflicts()
}

// ToggleNote toggles note n (1-9) on cell (row, col).
func (b *Board) ToggleNote(row, col, n int) {
	if b.Cells[row][col].Kind == KindGiven || b.Cells[row][col].Value != 0 {
		return
	}
	if n < 1 || n > 9 {
		return
	}
	b.Cells[row][col].Kind = KindUser
	b.Cells[row][col].Notes[n-1] = !b.Cells[row][col].Notes[n-1]
}

// Peers returns all (row, col) peers of (r, c) — same row, col, or 3x3 box.
func (b *Board) Peers(r, c int) [][2]int {
	seen := map[[2]int]bool{}
	var peers [][2]int

	// Same row
	for col := 0; col < 9; col++ {
		if col != c {
			p := [2]int{r, col}
			if !seen[p] {
				seen[p] = true
				peers = append(peers, p)
			}
		}
	}
	// Same col
	for row := 0; row < 9; row++ {
		if row != r {
			p := [2]int{row, c}
			if !seen[p] {
				seen[p] = true
				peers = append(peers, p)
			}
		}
	}
	// Same box
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

// Candidates returns the set of valid candidates for cell (r, c).
func (b *Board) Candidates(r, c int) map[int]bool {
	if b.Cells[r][c].Value != 0 {
		return nil
	}
	cands := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true, 9: true}
	for _, p := range b.Peers(r, c) {
		v := b.Cells[p[0]][p[1]].Value
		if v != 0 {
			delete(cands, v)
		}
	}
	return cands
}

// AllCandidates returns candidate maps for all empty cells.
func (b *Board) AllCandidates() [9][9]map[int]bool {
	var result [9][9]map[int]bool
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b.Cells[r][c].Value == 0 {
				result[r][c] = b.Candidates(r, c)
			}
		}
	}
	return result
}

// CellsWithValue returns all (row, col) positions with the given value.
func (b *Board) CellsWithValue(val int) [][2]int {
	var result [][2]int
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b.Cells[r][c].Value == val {
				result = append(result, [2]int{r, c})
			}
		}
	}
	return result
}

// RemainingCounts returns how many of each digit (1-9) are still needed.
func (b *Board) RemainingCounts() [9]int {
	placed := [9]int{}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			v := b.Cells[r][c].Value
			if v != 0 {
				placed[v-1]++
			}
		}
	}
	var remaining [9]int
	for i := 0; i < 9; i++ {
		remaining[i] = 9 - placed[i]
	}
	return remaining
}

// IsSolved returns true if all cells are filled correctly.
func (b *Board) IsSolved() bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if b.Cells[r][c].Value == 0 {
				return false
			}
			if b.Cells[r][c].Conflict {
				return false
			}
		}
	}
	return true
}

// LoadGivens loads a puzzle from an 81-character string.
func (b *Board) LoadGivens(s string) {
	for i, ch := range s {
		r, c := i/9, i%9
		if ch >= '1' && ch <= '9' {
			b.Cells[r][c].Value = int(ch - '0')
			b.Cells[r][c].Kind = KindGiven
		} else {
			b.Cells[r][c].Value = 0
			b.Cells[r][c].Kind = KindEmpty
		}
	}
	b.updateConflicts()
}

// Clone returns a deep copy of the board.
func (b *Board) Clone() *Board {
	nb := &Board{}
	nb.Cells = b.Cells
	nb.Solution = b.Solution
	return nb
}

// ToGrid returns the current values as a 9x9 int array.
func (b *Board) ToGrid() [9][9]int {
	var grid [9][9]int
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			grid[r][c] = b.Cells[r][c].Value
		}
	}
	return grid
}

func (b *Board) updateConflicts() {
	// Reset all conflict flags
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			b.Cells[r][c].Conflict = false
		}
	}
	// Check each cell against peers
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			v := b.Cells[r][c].Value
			if v == 0 {
				continue
			}
			for _, p := range b.Peers(r, c) {
				if b.Cells[p[0]][p[1]].Value == v {
					b.Cells[r][c].Conflict = true
					b.Cells[p[0]][p[1]].Conflict = true
				}
			}
		}
	}
}

// UpdateConflicts re-evaluates and marks conflicting cells on the board.
// It is the public equivalent of updateConflicts, for use by external packages.
func (b *Board) UpdateConflicts() {
	b.updateConflicts()
}
