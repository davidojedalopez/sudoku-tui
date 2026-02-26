package library

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/davidojeda/sudoku-tui/internal/curated"
	"github.com/davidojeda/sudoku-tui/internal/generator"
	"github.com/davidojeda/sudoku-tui/internal/msgs"
	"github.com/davidojeda/sudoku-tui/internal/theme"
)

// Model is the puzzle library screen.
type Model struct {
	puzzles []curated.Puzzle
	cursor  int
	filter  string
	width   int
	height  int
	theme   *theme.Theme
}

// New creates a new library screen model.
func New(th *theme.Theme) *Model {
	puzzles, _ := curated.Load()
	return &Model{
		puzzles: puzzles,
		theme:   th,
	}
}

// SetTheme updates the theme on the library screen.
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
			}
		case "down", "j":
			filtered := m.filtered()
			if m.cursor < len(filtered)-1 {
				m.cursor++
			}
		case "enter":
			filtered := m.filtered()
			if m.cursor < len(filtered) {
				puzzle := filtered[m.cursor]
				return m, m.loadPuzzle(puzzle)
			}
		}
	}
	return m, nil
}

func (m *Model) filtered() []curated.Puzzle {
	if m.filter == "" {
		return m.puzzles
	}
	f := strings.ToLower(m.filter)
	var result []curated.Puzzle
	for _, p := range m.puzzles {
		if strings.Contains(strings.ToLower(p.Name), f) ||
			strings.Contains(strings.ToLower(p.Difficulty), f) {
			result = append(result, p)
		}
	}
	return result
}

func (m *Model) loadPuzzle(p curated.Puzzle) tea.Cmd {
	return func() tea.Msg {
		var grid [9][9]int
		for i, ch := range p.Givens {
			if ch >= '1' && ch <= '9' {
				grid[i/9][i%9] = int(ch - '0')
			}
		}
		// Solve to get the full solution
		generator.Solve(&grid)
		diff := parseDiff(p.Difficulty)
		return msgs.StartGameMsg{
			Difficulty: diff,
			Puzzle:     p.Givens,
			Solution:   gridToString(grid),
			PuzzleID:   p.ID,
		}
	}
}

func parseDiff(s string) generator.Difficulty {
	switch strings.ToLower(s) {
	case "easy":
		return generator.Easy
	case "medium":
		return generator.Medium
	case "hard":
		return generator.Hard
	case "expert":
		return generator.Expert
	}
	return generator.Medium
}

func gridToString(grid [9][9]int) string {
	b := make([]byte, 81)
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			v := grid[r][c]
			if v == 0 {
				b[r*9+c] = '0'
			} else {
				b[r*9+c] = byte('0' + v)
			}
		}
	}
	return string(b)
}

func (m *Model) View() string {
	th := m.theme
	filtered := m.filtered()

	header := th.Header.Bar.Width(m.width).Render(" " + th.Header.Title.Render("PUZZLE LIBRARY"))

	// Overhead: 2 (list border) + 2 ("  " separator) + 2 (detail border) = 6
	// No body indent – panels extend to the screen edges.
	// Allocate list first (min 29 so 17-char names fit), then give remainder to detail.
	available := m.width - 6
	if available < 1 {
		available = 1
	}
	listWidth := available * 40 / 100
	if listWidth < 29 { // 29 = maxNameLen(17) + badge overhead(12)
		listWidth = 29
	}
	detailWidth := available - listWidth
	if detailWidth < 27 { // board preview is 25 chars wide; need at least 27
		detailWidth = 27
	}

	list := m.renderList(filtered, listWidth)
	detail := m.renderDetail(filtered, detailWidth)

	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		th.Library.PanelBorder.Width(listWidth).Render(list),
		"  ",
		th.Library.PanelBorder.Width(detailWidth).Render(detail),
	)

	hints := th.Footer.KeyHint.Render("j/k") + th.Footer.KeyLabel.Render(" Scroll") +
		"   " + th.Footer.KeyHint.Render("Enter") + th.Footer.KeyLabel.Render(" Load") +
		"   " + th.Footer.KeyHint.Render("Esc") + th.Footer.KeyLabel.Render(" Back")
	footer := th.Footer.Bar.Width(m.width).Render("  " + hints)

	bodyHeight := m.height - 2 - lipgloss.Height(footer)
	body := lipgloss.NewStyle().Height(bodyHeight).Render(panels)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m *Model) renderList(puzzles []curated.Puzzle, width int) string {
	th := m.theme
	var lines []string
	for i, p := range puzzles {
		badge := m.diffBadge(p.Difficulty)
		name := p.Name
		maxNameLen := width - 12
		if maxNameLen < 1 {
			maxNameLen = 1
		}
		if len(name) > maxNameLen {
			name = name[:maxNameLen]
		}
		var line string
		if i == m.cursor {
			line = th.Library.ItemActive.Render("▶ "+name) + "  " + badge
		} else {
			line = th.Library.Item.Render("  "+name) + "  " + badge
		}
		lines = append(lines, line)
	}
	if len(puzzles) == 0 {
		lines = append(lines, th.Library.Item.Render("  No puzzles found."))
	}
	return strings.Join(lines, "\n")
}

func (m *Model) renderDetail(puzzles []curated.Puzzle, width int) string {
	th := m.theme
	if m.cursor >= len(puzzles) || len(puzzles) == 0 {
		return th.Library.DetailTitle.Render("No puzzle selected")
	}

	p := puzzles[m.cursor]

	title := th.Library.DetailTitle.Render(p.Name)
	diffLine := th.Library.DetailLabel.Render("Difficulty: ") + th.Library.DetailValue.Render(strings.ToUpper(p.Difficulty))
	authorLine := th.Library.DetailLabel.Render("Author:     ") + th.Library.DetailValue.Render(p.Author)
	desc := th.Library.DetailDesc.Render(p.Description)

	var refLine string
	if p.Reference != "" {
		ref := p.Reference
		maxRef := width - 6 // account for "URL: " (5 chars) + "…" (1 char)
		if maxRef > 4 && len(ref) > maxRef {
			ref = ref[:maxRef] + "…"
		}
		refLine = th.Library.DetailLabel.Render("URL: ") + th.Library.DetailDesc.Render(ref)
	}

	preview := m.renderPreview(p.Givens)

	loadBtn := th.Library.LoadButton.Render("[ Enter ] LOAD PUZZLE")

	parts := []string{title, "", diffLine, authorLine}
	if refLine != "" {
		parts = append(parts, refLine)
	}
	parts = append(parts, "", desc, "", preview, "", loadBtn)

	_ = width
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m *Model) renderPreview(givens string) string {
	th := m.theme
	if len(givens) < 81 {
		return ""
	}
	var rows []string
	rows = append(rows, "+-------+-------+-------+")
	for r := 0; r < 9; r++ {
		var sb strings.Builder
		sb.WriteString("| ")
		for c := 0; c < 9; c++ {
			ch := givens[r*9+c]
			if ch == '0' {
				sb.WriteString(th.Library.PreviewCell.Render("."))
			} else {
				sb.WriteString(th.Library.PreviewCell.Render(string(ch)))
			}
			if c == 2 || c == 5 {
				sb.WriteString(" | ")
			} else if c < 8 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString(" |")
		rows = append(rows, sb.String())
		if r == 2 || r == 5 {
			rows = append(rows, "+-------+-------+-------+")
		}
	}
	rows = append(rows, "+-------+-------+-------+")
	return strings.Join(rows, "\n")
}

func (m *Model) diffBadge(diff string) string {
	th := m.theme
	switch strings.ToLower(diff) {
	case "easy":
		return th.Badges.Easy.Render("EASY")
	case "medium":
		return th.Badges.Medium.Render("MEDIUM")
	case "hard":
		return th.Badges.Hard.Render("HARD")
	case "expert":
		return th.Badges.Expert.Render("EXPERT")
	}
	return diff
}
