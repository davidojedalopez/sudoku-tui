package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const historyFile = "history.json"

// Store manages the history file.
type Store struct {
	path    string
	entries []Entry
}

// NewStore creates or loads the history store.
func NewStore() (*Store, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "sudoku-tui")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create history dir: %w", err)
	}

	s := &Store{path: filepath.Join(dir, historyFile)}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Add appends a new entry and saves.
func (s *Store) Add(difficulty string, elapsed int64, result Result, puzzleID string) error {
	e := Entry{
		ID:         fmt.Sprintf("%d", time.Now().UnixNano()),
		Date:       time.Now().Format(time.RFC3339),
		Difficulty: difficulty,
		Elapsed:    elapsed,
		Result:     result,
		PuzzleID:   puzzleID,
	}
	s.entries = append(s.entries, e)
	// Sort newest first.
	sort.Slice(s.entries, func(i, j int) bool {
		return s.entries[i].Date > s.entries[j].Date
	})
	return s.save()
}

// All returns all entries.
func (s *Store) All() []Entry {
	return s.entries
}

// Stats computes aggregate stats.
func (s *Store) Stats() Stats {
	stats := Stats{
		BestTimes: map[string]int64{},
	}
	stats.Total = len(s.entries)

	wins := 0
	currentStreak := 0
	maxStreak := 0

	for _, e := range s.entries {
		if e.Result == ResultWin {
			wins++
			if prev, ok := stats.BestTimes[e.Difficulty]; !ok || e.Elapsed < prev {
				stats.BestTimes[e.Difficulty] = e.Elapsed
			}
			currentStreak++
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}
		} else {
			currentStreak = 0
		}
	}

	if stats.Total > 0 {
		stats.WinRate = wins * 100 / stats.Total
	}
	stats.Streak = maxStreak
	return stats
}

// Stats holds aggregate history statistics.
type Stats struct {
	Total     int
	WinRate   int
	Streak    int
	BestTimes map[string]int64
}

// BestTimeFormatted returns the best time for a difficulty as a formatted string.
func (s Stats) BestTimeFormatted(difficulty string) string {
	secs, ok := s.BestTimes[difficulty]
	if !ok {
		return "--:--"
	}
	h := secs / 3600
	m := (secs % 3600) / 60
	sec := secs % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, sec)
	}
	return fmt.Sprintf("%02d:%02d", m, sec)
}
