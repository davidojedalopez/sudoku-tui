package game

import (
	"fmt"
	"strings"

	"github.com/davidojeda/sudoku-tui/internal/theme"
)

// renderRemaining returns the "REMAINING" sidebar section as a string.
func renderRemaining(remaining [9]int, cursorVal int, th *theme.Theme) string {
	var sb strings.Builder
	sb.WriteString(th.Sidebar.Title.Render("REMAINING") + "\n")

	for i := 0; i < 9; i++ {
		digit := i + 1
		rem := remaining[i]
		placed := 9 - rem

		barWidth := 9
		filled := placed
		if filled > barWidth {
			filled = barWidth
		}

		digitStr := fmt.Sprintf("%d", digit)

		var line string
		if rem == 0 {
			row := fmt.Sprintf("%s %s DONE", digitStr, strings.Repeat("█", barWidth))
			line = th.Sidebar.RemainingDone.Render(row)
		} else {
			barRendered := th.Sidebar.RemainingBar.Render(strings.Repeat("█", filled)) +
				th.Sidebar.RemainingBarEmpty.Render(strings.Repeat("░", barWidth-filled))
			countStr := th.Sidebar.RemainingCount.Render(fmt.Sprintf("(%d)", rem))

			if digit == cursorVal {
				line = th.Sidebar.RemainingActive.Render(digitStr) + " " + barRendered + " " + countStr
			} else {
				line = th.Sidebar.RemainingDigit.Render(digitStr) + " " + barRendered + " " + countStr
			}
		}

		sb.WriteString(line + "\n")
	}

	return sb.String()
}
