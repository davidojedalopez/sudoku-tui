package board

// CellKind represents the kind of a cell.
type CellKind int

const (
	KindEmpty CellKind = iota
	KindGiven
	KindUser
)

// Cell represents a single cell on the Sudoku board.
type Cell struct {
	Value    int
	Kind     CellKind
	Notes    [9]bool
	Conflict bool
}

// IsEmpty returns true if the cell has no value.
func (c Cell) IsEmpty() bool {
	return c.Value == 0
}

// HasNote returns true if note n (1-9) is set.
func (c Cell) HasNote(n int) bool {
	if n < 1 || n > 9 {
		return false
	}
	return c.Notes[n-1]
}

// HasAnyNote returns true if any note is set.
func (c Cell) HasAnyNote() bool {
	for _, n := range c.Notes {
		if n {
			return true
		}
	}
	return false
}
