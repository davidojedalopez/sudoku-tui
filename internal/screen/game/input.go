package game

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/davidojeda/sudoku-tui/internal/msgs"
)

// handleKey processes a key message for the game model.
func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	key := msg.String()

	shiftDigit := map[string]int{
		"!": 1, "@": 2, "#": 3, "$": 4, "%": 5,
		"^": 6, "&": 7, "*": 8, "(": 9,
	}

	switch {
	case key == "esc" || key == "q":
		return func() tea.Msg {
			return msgs.GameOverMsg{
				Won:     false,
				Elapsed: m.timer.elapsedSeconds(),
				Diff:    m.difficulty,
			}
		}
	case key == "up" || key == "w":
		if m.cursor[0] > 0 {
			m.cursor[0]--
		}
	case key == "down" || key == "s":
		if m.cursor[0] < 8 {
			m.cursor[0]++
		}
	case key == "left" || key == "a":
		if m.cursor[1] > 0 {
			m.cursor[1]--
		}
	case key == "right" || key == "d":
		if m.cursor[1] < 8 {
			m.cursor[1]++
		}
	case key == "n":
		if m.inputMode == ModDigit {
			m.inputMode = ModNote
		} else {
			m.inputMode = ModDigit
		}
	case key == "backspace" || key == "x" || key == "delete":
		m.eraseCell()
	case key == "u" || key == "ctrl+z":
		m.undo()
	case key >= "1" && key <= "9":
		digit := int(key[0] - '0')
		m.placeDigit(digit, false)
		return m.checkTechnique()
	case shiftDigit[key] != 0:
		m.placeDigit(shiftDigit[key], true)
	}
	return nil
}
