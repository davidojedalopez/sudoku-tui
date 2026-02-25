package history

import "fmt"

// Result represents the game outcome.
type Result string

const (
	ResultWin    Result = "WIN"
	ResultGaveUp Result = "GAVE UP"
)

// Entry represents a single game history record.
type Entry struct {
	ID         string `json:"id"`
	Date       string `json:"date"`
	Difficulty string `json:"difficulty"`
	Elapsed    int64  `json:"elapsed_seconds"`
	Result     Result `json:"result"`
	PuzzleID   string `json:"puzzle_id,omitempty"`
}

// FormattedTime returns the elapsed time as MM:SS or HH:MM:SS.
func (e Entry) FormattedTime() string {
	if e.Result == ResultGaveUp {
		return "--:--"
	}
	secs := e.Elapsed
	h := secs / 3600
	m := (secs % 3600) / 60
	s := secs % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
