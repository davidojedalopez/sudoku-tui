package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidojeda/sudoku-tui/internal/msgs"
)

// Re-export screen constants and message types from msgs package for convenience.
type Screen = msgs.Screen

const (
	ScreenMenu    = msgs.ScreenMenu
	ScreenGame    = msgs.ScreenGame
	ScreenHistory = msgs.ScreenHistory
	ScreenLibrary = msgs.ScreenLibrary
)

// NavigateMsg requests navigation to a different screen.
type NavigateMsg = msgs.NavigateMsg

// StartGameMsg requests starting a new generated game.
type StartGameMsg = msgs.StartGameMsg

// GameOverMsg is sent when a game ends.
type GameOverMsg = msgs.GameOverMsg

// Navigate returns a command to navigate to a screen.
func Navigate(s Screen) tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{To: s}
	}
}
