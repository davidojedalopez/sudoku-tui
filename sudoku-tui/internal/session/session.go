// Package session handles saving and restoring in-progress game state.
package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const sessionFile = "session.json"

// CellState is the serializable state of a single board cell.
type CellState struct {
	Value int     `json:"value"`
	Kind  int     `json:"kind"`
	Notes [9]bool `json:"notes"`
}

// SavedGame holds the complete state needed to resume an in-progress game.
type SavedGame struct {
	Puzzle     string           `json:"puzzle"`
	Solution   string           `json:"solution"`
	Difficulty string           `json:"difficulty"`
	PuzzleID   string           `json:"puzzle_id,omitempty"`
	Elapsed    int64            `json:"elapsed_seconds"`
	Cells      [9][9]CellState  `json:"cells"`
}

func configDir() (string, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "sudoku-tui")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create config dir: %w", err)
	}
	return dir, nil
}

// Save persists the current game state to disk.
func Save(g *SavedGame) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, sessionFile), data, 0644)
}

// Load reads the saved game from disk.
// Returns (nil, false, nil) if no save exists.
func Load() (*SavedGame, bool, error) {
	dir, err := configDir()
	if err != nil {
		return nil, false, err
	}
	data, err := os.ReadFile(filepath.Join(dir, sessionFile))
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var g SavedGame
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, false, err
	}
	return &g, true, nil
}

// Clear removes the saved game from disk.
func Clear() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(dir, sessionFile))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
