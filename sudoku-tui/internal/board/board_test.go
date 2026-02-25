package board

import (
	"testing"
)

func TestPeers(t *testing.T) {
	b := New()
	peers := b.Peers(0, 0)
	if len(peers) != 20 {
		t.Errorf("expected 20 peers, got %d", len(peers))
	}
}

func TestConflictDetection(t *testing.T) {
	b := New()
	b.Cells[0][0] = Cell{Value: 5, Kind: KindUser}
	b.Cells[0][1] = Cell{Value: 5, Kind: KindUser}
	b.updateConflicts()
	if !b.Cells[0][0].Conflict {
		t.Error("expected conflict at [0][0]")
	}
	if !b.Cells[0][1].Conflict {
		t.Error("expected conflict at [0][1]")
	}
}

func TestRemainingCounts(t *testing.T) {
	b := New()
	b.Cells[0][0] = Cell{Value: 1, Kind: KindGiven}
	b.Cells[0][1] = Cell{Value: 1, Kind: KindGiven}
	remaining := b.RemainingCounts()
	if remaining[0] != 7 { // 9 - 2 = 7 remaining 1s
		t.Errorf("expected 7 remaining 1s, got %d", remaining[0])
	}
}

func TestLoadGivens(t *testing.T) {
	b := New()
	givens := "530070000600195000098000060800060003400803001700020006060000280000419005000080079"
	b.LoadGivens(givens)
	if b.Cells[0][0].Value != 5 || b.Cells[0][0].Kind != KindGiven {
		t.Error("expected cell [0][0] to be given 5")
	}
	if b.Cells[0][3].Value != 0 || b.Cells[0][3].Kind != KindEmpty {
		t.Error("expected cell [0][3] to be empty")
	}
}

func TestIsValidPlacement(t *testing.T) {
	var grid [9][9]int
	grid[0][0] = 5
	if IsValidPlacement(grid, 0, 1, 5) {
		t.Error("expected invalid: 5 already in row 0")
	}
	if !IsValidPlacement(grid, 0, 1, 3) {
		t.Error("expected valid placement of 3")
	}
}
