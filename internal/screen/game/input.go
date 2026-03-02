package game

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/davidojeda/sudoku-tui/internal/msgs"
)

// handleKey processes a key message for the game model.
func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	switch key {
	case "esc":
		return func() tea.Msg {
			return msgs.GameOverMsg{
				Won:     false,
				Elapsed: m.timer.elapsedSeconds(),
				Diff:    m.difficulty,
			}
		}
	case "up", "k":
		if m.cursor[0] > 0 {
			m.cursor[0]--
		}
	case "down", "j":
		if m.cursor[0] < 8 {
			m.cursor[0]++
		}
	case "left", "h":
		if m.cursor[1] > 0 {
			m.cursor[1]--
		}
	case "right", "l":
		if m.cursor[1] < 8 {
			m.cursor[1]++
		}
	case "n":
		if m.inputMode == ModDigit {
			m.inputMode = ModNote
		} else {
			m.inputMode = ModDigit
		}
	case "backspace", "x", "delete":
		m.eraseCell()
	case "u", "ctrl+z":
		m.undo()
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		digit := int(key[0] - '0')
		m.placeDigit(digit, false)
		return m.checkTechnique()
	}
	return nil
}
