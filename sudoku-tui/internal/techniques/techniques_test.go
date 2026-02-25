package techniques

import (
	"testing"

	"github.com/davidojeda/sudoku-tui/internal/board"
)

func TestNakedSingle(t *testing.T) {
	var cands [9][9]map[int]bool
	cands[0][0] = map[int]bool{5: true} // only one candidate
	cands[0][1] = map[int]bool{3: true, 7: true}

	if !IsNakedSingle(cands, 0, 0) {
		t.Error("expected naked single at [0][0]")
	}
	if IsNakedSingle(cands, 0, 1) {
		t.Error("expected no naked single at [0][1]")
	}
}

func TestHiddenSingle(t *testing.T) {
	var cands [9][9]map[int]bool
	// val=5 only appears in col 0 at row 0
	cands[0][0] = map[int]bool{5: true, 3: true}
	for r := 1; r < 9; r++ {
		cands[r][0] = map[int]bool{3: true, 7: true} // 5 not here
	}

	if !DetectHiddenSingle(cands, 0, 0, 5) {
		t.Error("expected hidden single for val=5 at [0][0] in col")
	}
}

func TestXWing(t *testing.T) {
	var cands [9][9]map[int]bool
	// val=5 appears in rows 0 and 3, both in cols 2 and 7 only
	cands[0][2] = map[int]bool{5: true}
	cands[0][7] = map[int]bool{5: true}
	cands[3][2] = map[int]bool{5: true}
	cands[3][7] = map[int]bool{5: true}

	if !DetectXWing(cands, 5) {
		t.Error("expected X-Wing for val=5")
	}
}

func TestSwordfish(t *testing.T) {
	var cands [9][9]map[int]bool
	// val=3 in rows 0, 3, 6 at cols 1, 4, 7
	cands[0][1] = map[int]bool{3: true}
	cands[0][4] = map[int]bool{3: true}
	cands[3][1] = map[int]bool{3: true}
	cands[3][7] = map[int]bool{3: true}
	cands[6][4] = map[int]bool{3: true}
	cands[6][7] = map[int]bool{3: true}

	if !DetectSwordfish(cands, 3) {
		t.Error("expected Swordfish for val=3")
	}
}

func TestTakeSnapshot(t *testing.T) {
	b := board.New()
	b.LoadGivens("530070000600195000098000060800060003400803001700020006060000280000419005000080079")
	snap := TakeSnapshot(b)

	// Cell [0][3] is empty, should have candidates
	if snap.Candidates[0][3] == nil {
		t.Error("expected candidates for empty cell [0][3]")
	}
	// Cell [0][0] is given (5), should have no candidates
	if snap.Candidates[0][0] != nil {
		t.Error("expected no candidates for given cell [0][0]")
	}
}
