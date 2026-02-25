// Package msgs contains shared message types used across screen packages
// and the app package to avoid import cycles.
package msgs

import "github.com/davidojeda/sudoku-tui/internal/generator"

// Screen identifies which screen is active.
type Screen int

const (
	ScreenMenu Screen = iota
	ScreenGame
	ScreenHistory
	ScreenLibrary
)

// NavigateMsg requests navigation to a different screen.
type NavigateMsg struct {
	To Screen
}

// StartGameMsg requests starting a new generated game.
type StartGameMsg struct {
	Difficulty generator.Difficulty
	Puzzle     string
	Solution   string
	PuzzleID   string
}

// GameOverMsg is sent when a game ends.
type GameOverMsg struct {
	Won     bool
	Elapsed int64
	Diff    string
}

// ChangeThemeMsg is sent when the user selects a new theme.
type ChangeThemeMsg struct {
	ThemeName string
}

// ResumeGameMsg is sent when the user selects "Resume Game" from the menu.
type ResumeGameMsg struct{}
