package curated

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed curated.json
var curatedJSON []byte

// Puzzle represents a curated puzzle entry.
type Puzzle struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Difficulty  string `json:"difficulty"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Givens      string `json:"givens"`
}

// Load returns all curated puzzles.
func Load() ([]Puzzle, error) {
	var puzzles []Puzzle
	if err := json.Unmarshal(curatedJSON, &puzzles); err != nil {
		return nil, fmt.Errorf("parse curated puzzles: %w", err)
	}
	return puzzles, nil
}
