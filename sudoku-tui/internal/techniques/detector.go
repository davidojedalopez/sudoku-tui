package techniques

import (
	"github.com/davidojeda/sudoku-tui/internal/board"
)

// Technique represents a detected solving technique.
type Technique int

const (
	TechNone Technique = iota
	TechNakedSingle
	TechHiddenSingle
	TechNakedPair
	TechPointingPair
	TechXWing
	TechYWing
	TechSwordfish
)

// String returns the display name of the technique.
func (t Technique) String() string {
	switch t {
	case TechNakedSingle:
		return "Naked Single"
	case TechHiddenSingle:
		return "Hidden Single"
	case TechNakedPair:
		return "Naked Pair"
	case TechPointingPair:
		return "Pointing Pair"
	case TechXWing:
		return "X-Wing"
	case TechYWing:
		return "Y-Wing"
	case TechSwordfish:
		return "Swordfish"
	}
	return ""
}

// Message returns the toast message for the technique.
func (t Technique) Message() string {
	switch t {
	case TechNakedSingle:
		return "Naked Single — only one candidate left!"
	case TechHiddenSingle:
		return "Hidden Single — nice deduction."
	case TechNakedPair:
		return "Naked Pair detected! Great pattern."
	case TechPointingPair:
		return "Pointing Pair! Locked candidates found."
	case TechXWing:
		return "✦ X-Wing detected! Nice move."
	case TechYWing:
		return "✦ Y-Wing! Impressive technique."
	case TechSwordfish:
		return "✦ Swordfish! Master-level solve."
	}
	return ""
}

// Snapshot captures candidate sets before a move.
type Snapshot struct {
	Candidates [9][9]map[int]bool
}

// TakeSnapshot captures the current candidate state from a board.
func TakeSnapshot(b *board.Board) Snapshot {
	return Snapshot{Candidates: b.AllCandidates()}
}

// DetectTechnique identifies which technique was used to place val at (r, c).
// It compares the snapshot (before) to the current board state.
func DetectTechnique(b *board.Board, snap Snapshot, r, c, val int) Technique {
	prevCands := snap.Candidates[r][c]

	if prevCands == nil {
		return TechNone
	}

	// Naked single: only one candidate was available.
	if len(prevCands) == 1 {
		return TechNakedSingle
	}

	// Hidden single: val was only possible in one cell in a unit.
	if DetectHiddenSingle(snap.Candidates, r, c, val) {
		return TechHiddenSingle
	}

	// Advanced techniques detected from the snapshot.
	if DetectXWing(snap.Candidates, val) {
		return TechXWing
	}

	if DetectYWing(snap.Candidates) {
		return TechYWing
	}

	if DetectSwordfish(snap.Candidates, val) {
		return TechSwordfish
	}

	if DetectNakedPair(snap.Candidates, r, c) {
		return TechNakedPair
	}

	if DetectPointingPair(snap.Candidates, val) {
		return TechPointingPair
	}

	return TechNone
}
