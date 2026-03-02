package historyscreen

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/davidojeda/sudoku-tui/internal/msgs"
	"github.com/davidojeda/sudoku-tui/internal/history"
	"github.com/davidojeda/sudoku-tui/internal/theme"
)

const pageSize = 8

// Model is the history screen.
type Model struct {
	store   *history.Store
	entries []history.Entry
	cursor  int
	page    int
	width   int
	height  int
	theme   *theme.Theme
}

// New creates a new history screen model.
func New(th *theme.Theme, store *history.Store) *Model {
	entries := []history.Entry{}
	if store != nil {
		entries = store.All()
	}
	return &Model{
		store:   store,
		entries: entries,
		theme:   th,
	}
}

// SetTheme updates the theme on the history screen.
func (m *Model) SetTheme(th *theme.Theme) { m.theme = th }

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return msgs.NavigateMsg{To: msgs.ScreenMenu} }
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.page*pageSize {
					m.page--
				}
			}
		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
				if m.cursor >= (m.page+1)*pageSize {
					m.page++
				}
			}
		case "pgdown":
			m.page++
			maxPage := (len(m.entries) - 1) / pageSize
			if maxPage < 0 {
				maxPage = 0
			}
			if m.page > maxPage {
				m.page = maxPage
			}
			m.cursor = m.page * pageSize
		case "pgup":
			if m.page > 0 {
				m.page--
			}
			m.cursor = m.page * pageSize
		}
	}
	return m, nil
}

func (m *Model) View() string {
	th := m.theme
	header := th.Header.Bar.Width(m.width).Render(" " + th.Header.Title.Render("HISTORY"))

	var stats history.Stats
	if m.store != nil {
		stats = m.store.Stats()
	} else {
		stats = history.Stats{BestTimes: map[string]int64{}}
	}
	statsRow := m.renderStats(stats)

	table := m.renderTable()

	total := len(m.entries)
	pages := (total + pageSize - 1) / pageSize
	if pages == 0 {
		pages = 1
	}
	start := m.page*pageSize + 1
	end := minInt(m.page*pageSize+pageSize, total)
	if total == 0 {
		start = 0
		end = 0
	}
	pagination := th.History.Pagination.Render(
		fmt.Sprintf("  SHOWING %d-%d OF %d   PAGE %d/%d", start, end, total, m.page+1, pages))

	hints := th.Footer.KeyHint.Render("j/k") + th.Footer.KeyLabel.Render(" Scroll") +
		"   " + th.Footer.KeyHint.Render("PgDn/PgUp") + th.Footer.KeyLabel.Render(" Page") +
		"   " + th.Footer.KeyHint.Render("Esc") + th.Footer.KeyLabel.Render(" Back")
	footer := th.Footer.Bar.Width(m.width).Render("  " + hints)

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		statsRow,
		"",
		table,
		"",
		pagination,
	)

	bodyHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)
	body = lipgloss.NewStyle().Height(bodyHeight).Render(body)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m *Model) renderStats(stats history.Stats) string {
	th := m.theme
	boxes := []string{
		th.History.StatBox.Render(
			th.History.StatLabel.Render("TOTAL GAMES") + "\n" +
				th.History.StatValue.Render(fmt.Sprintf("   %d   ", stats.Total))),
		th.History.StatBox.Render(
			th.History.StatLabel.Render(" WIN RATE ") + "\n" +
				th.History.StatValue.Render(fmt.Sprintf("  %d%%  ", stats.WinRate))),
		th.History.StatBox.Render(
			th.History.StatLabel.Render(" BEST (HARD)") + "\n" +
				th.History.StatValue.Render("  "+stats.BestTimeFormatted("Hard")+"  ")),
		th.History.StatBox.Render(
			th.History.StatLabel.Render("  STREAK  ") + "\n" +
				th.History.StatValue.Render(fmt.Sprintf("    %d    ", stats.Streak))),
	}
	return "  " + lipgloss.JoinHorizontal(lipgloss.Top, boxes[0], "  ", boxes[1], "  ", boxes[2], "  ", boxes[3])
}

func (m *Model) renderTable() string {
	th := m.theme
	header := th.History.TableHeader.Render(
		fmt.Sprintf("  %-20s %-12s %-10s %-10s", "DATE", "DIFFICULTY", "TIME", "RESULT"))
	sep := th.History.TableHeader.Render("  " + strings.Repeat("-", 56))

	var lines []string
	lines = append(lines, header, sep)

	start := m.page * pageSize
	end := start + pageSize
	if end > len(m.entries) {
		end = len(m.entries)
	}

	for i := start; i < end; i++ {
		e := m.entries[i]
		dateStr := e.Date
		if len(dateStr) > 16 {
			dateStr = dateStr[:16]
		}
		var resultStyle lipgloss.Style
		resultStr := string(e.Result)
		if e.Result == history.ResultWin {
			resultStyle = th.History.BadgeWin
		} else {
			resultStyle = th.History.BadgeLoss
		}

		line := fmt.Sprintf("  %-20s %-12s %-10s", dateStr, e.Difficulty, e.FormattedTime())
		if i == m.cursor {
			line = th.History.RowActive.Render(">" + line[1:]) + resultStyle.Render(resultStr)
		} else {
			line = th.History.Row.Render(line) + resultStyle.Render(resultStr)
		}
		lines = append(lines, line)
	}

	if len(m.entries) == 0 {
		lines = append(lines, th.History.Row.Render("  No games played yet."))
	}

	return strings.Join(lines, "\n")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
