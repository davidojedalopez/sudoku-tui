package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/davidojeda/sudoku-tui/internal/board"
)

// renderBoard renders the 9x9 Sudoku grid.
func (m *Model) renderBoard() string {
	cursorVal := m.cursorValue()

	var rows []string

	rows = append(rows, buildTopBorder())

	for r := 0; r < 9; r++ {
		for line := 0; line < 3; line++ {
			var rowStr strings.Builder
			rowStr.WriteString("║")
			for c := 0; c < 9; c++ {
				cell := m.board.Cells[r][c]
				rowStr.WriteString(m.renderCellLine(r, c, line, cell, cursorVal))
				if (c+1)%3 == 0 {
					rowStr.WriteString("║")
				} else {
					rowStr.WriteString("│")
				}
			}
			rows = append(rows, rowStr.String())
		}

		if r < 8 {
			if (r+1)%3 == 0 {
				rows = append(rows, buildBoxRowSep())
			} else {
				rows = append(rows, buildCellRowSep())
			}
		}
	}

	rows = append(rows, buildBottomBorder())

	return strings.Join(rows, "\n")
}

// renderCellLine renders one of the 3 lines of a cell.
func (m *Model) renderCellLine(r, c, line int, cell board.Cell, cursorVal int) string {
	th := m.theme
	isCursor := r == m.cursor[0] && c == m.cursor[1]
	isSameNum := cursorVal != 0 && cell.Value == cursorVal && !isCursor

	var style lipgloss.Style
	switch {
	case isCursor:
		style = th.Cell.Cursor
	case cell.Conflict:
		style = th.Cell.Conflict
	case isSameNum:
		style = th.Cell.Highlight
	case cell.Kind == board.KindGiven:
		style = th.Cell.Given
	case cell.Kind == board.KindUser && cell.Value != 0:
		style = th.Cell.User
	default:
		style = th.Cell.Empty
	}

	content := cellContent(cell, line)
	return style.Width(5).Align(lipgloss.Center).Render(content)
}

// cellContent returns the content string for line (0,1,2) of a cell.
func cellContent(cell board.Cell, line int) string {
	if cell.Value != 0 {
		if line == 1 {
			return fmt.Sprintf(" %d ", cell.Value)
		}
		return "   "
	}

	if !cell.HasAnyNote() {
		return "   "
	}

	var sb strings.Builder
	base := line * 3
	for i := 0; i < 3; i++ {
		n := base + i + 1
		if cell.HasNote(n) {
			sb.WriteString(fmt.Sprintf("%d", n))
		} else {
			sb.WriteString(" ")
		}
		if i < 2 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

func buildTopBorder() string {
	return "╔" + cellBlock("═", "╤", "╦") + "╗"
}

func buildBottomBorder() string {
	return "╚" + cellBlock("═", "╧", "╩") + "╝"
}

func buildBoxRowSep() string {
	return "╠" + cellBlock("═", "╪", "╬") + "╣"
}

func buildCellRowSep() string {
	return "╟" + cellBlock("─", "┼", "╫") + "╢"
}

func cellBlock(fill, inner, box string) string {
	cell := strings.Repeat(fill, 5)
	var sb strings.Builder
	for c := 0; c < 9; c++ {
		sb.WriteString(cell)
		if c < 8 {
			if (c+1)%3 == 0 {
				sb.WriteString(box)
			} else {
				sb.WriteString(inner)
			}
		}
	}
	return sb.String()
}
