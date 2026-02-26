package library

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/davidojeda/sudoku-tui/internal/curated"
	"github.com/davidojeda/sudoku-tui/internal/generator"
	"github.com/davidojeda/sudoku-tui/internal/msgs"
	"github.com/davidojeda/sudoku-tui/internal/theme"
)

var filterCycle = []string{"", "easy", "medium", "hard", "expert"}

// Model is the puzzle library screen.
type Model struct {
	puzzles      []curated.Puzzle
	cursor       int
	scrollOffset int // index of the first visible list item
	listHeight   int // visible list rows (set each render, used by Update)
	diffFilter   string // "" = all, or "easy"/"medium"/"hard"/"expert"
	width        int
	height       int
	theme        *theme.Theme
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
		case "tab":
			m.cycleFilter()
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scrollOffset {
					m.scrollOffset = m.cursor
				}
			}
		case "down", "j":
			filtered := m.filteredAndSorted()
			if m.cursor < len(filtered)-1 {
				m.cursor++
				if h := m.visibleHeight(); m.cursor >= m.scrollOffset+h {
					m.scrollOffset = m.cursor - h + 1
				}
			}
		case "pgdown":
			filtered := m.filteredAndSorted()
			h := m.visibleHeight()
			m.cursor = min(m.cursor+h, len(filtered)-1)
			if m.cursor >= m.scrollOffset+h {
				m.scrollOffset = m.cursor - h + 1
			}
		case "pgup":
			h := m.visibleHeight()
			m.cursor = max(m.cursor-h, 0)
			if m.cursor < m.scrollOffset {
				m.scrollOffset = m.cursor
			}
		case "enter":
			filtered := m.filteredAndSorted()
			if m.cursor < len(filtered) {
				puzzle := filtered[m.cursor]
				return m, m.loadPuzzle(puzzle)
			}
		}
	}
	return m, nil
}

func (m *Model) cycleFilter() {
	for i, f := range filterCycle {
		if f == m.diffFilter {
			m.diffFilter = filterCycle[(i+1)%len(filterCycle)]
			m.cursor = 0
			m.scrollOffset = 0
			return
		}
	}
	m.diffFilter = ""
	m.cursor = 0
	m.scrollOffset = 0
}

// normalizeDiff maps raw difficulty strings (including "very hard") to canonical filter keys.
func normalizeDiff(d string) string {
	switch strings.ToLower(d) {
	case "easy":
		return "easy"
	case "medium":
		return "medium"
	case "hard":
		return "hard"
	case "expert", "very hard":
		return "expert"
	}
	return strings.ToLower(d)
}

func diffWeight(d string) int {
	switch normalizeDiff(d) {
	case "easy":
		return 0
	case "medium":
		return 1
	case "hard":
		return 2
	case "expert":
		return 3
	}
	return 4
}

func (m *Model) filteredAndSorted() []curated.Puzzle {
	var result []curated.Puzzle
	for _, p := range m.puzzles {
		if m.diffFilter == "" || normalizeDiff(p.Difficulty) == m.diffFilter {
			result = append(result, p)
		}
	}
	if m.diffFilter == "" {
		// No filter: sort by difficulty first, then alphabetically within each level.
		sort.Slice(result, func(i, j int) bool {
			wi, wj := diffWeight(result[i].Difficulty), diffWeight(result[j].Difficulty)
			if wi != wj {
				return wi < wj
			}
			return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
		})
	} else {
		// Filter active: sort alphabetically.
		sort.Slice(result, func(i, j int) bool {
			return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
		})
	}
	return result
}

// visibleHeight returns the number of list rows currently visible.
// Falls back to a safe default before the first render.
func (m *Model) visibleHeight() int {
	if m.listHeight > 0 {
		return m.listHeight
	}
	return 20
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
	case "expert", "very hard":
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
	filtered := m.filteredAndSorted()

	header := th.Header.Bar.Width(m.width).Render(" " + th.Header.Title.Render("PUZZLE LIBRARY"))

	filterBar := m.renderFilterBar()

	hints := th.Footer.KeyHint.Render("j/k") + th.Footer.KeyLabel.Render(" Scroll") +
		"   " + th.Footer.KeyHint.Render("PgDn/PgUp") + th.Footer.KeyLabel.Render(" Page") +
		"   " + th.Footer.KeyHint.Render("Tab") + th.Footer.KeyLabel.Render(" Filter") +
		"   " + th.Footer.KeyHint.Render("Enter") + th.Footer.KeyLabel.Render(" Load") +
		"   " + th.Footer.KeyHint.Render("Esc") + th.Footer.KeyLabel.Render(" Back")
	footer := th.Footer.Bar.Width(m.width).Render("  " + hints)

	filterBarHeight := lipgloss.Height(filterBar)
	bodyHeight := m.height - lipgloss.Height(header) - filterBarHeight - lipgloss.Height(footer)
	if bodyHeight < 1 {
		bodyHeight = 1
	}

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

	// Compute how many list rows fit inside the panel border (2 lines overhead).
	lh := bodyHeight - 2
	if lh < 1 {
		lh = 1
	}
	m.listHeight = lh

	list := m.renderList(filtered, listWidth)
	detail := m.renderDetail(filtered, detailWidth)

	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		th.Library.PanelBorder.Width(listWidth).Render(list),
		"  ",
		th.Library.PanelBorder.Width(detailWidth).Render(detail),
	)

	body := lipgloss.NewStyle().Height(bodyHeight).Render(panels)

	return lipgloss.JoinVertical(lipgloss.Left, header, filterBar, body, footer)
}

func (m *Model) renderFilterBar() string {
	th := m.theme
	type chip struct {
		key   string
		label string
	}
	chips := []chip{
		{"", "ALL"},
		{"easy", "EASY"},
		{"medium", "MEDIUM"},
		{"hard", "HARD"},
		{"expert", "EXPERT"},
	}

	var parts []string
	parts = append(parts, th.Footer.KeyLabel.Render("  Filter "))
	for _, c := range chips {
		if c.key == m.diffFilter {
			parts = append(parts, th.Library.FilterChipActive.Render(c.label))
		} else {
			parts = append(parts, th.Library.FilterChip.Render(c.label))
		}
		parts = append(parts, " ")
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}

func (m *Model) renderList(puzzles []curated.Puzzle, width int) string {
	th := m.theme
	if len(puzzles) == 0 {
		return th.Library.Item.Render("  No puzzles found.")
	}

	h := m.visibleHeight()
	start := m.scrollOffset
	end := start + h
	if end > len(puzzles) {
		end = len(puzzles)
	}

	var lines []string
	for i := start; i < end; i++ {
		p := puzzles[i]
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
	switch normalizeDiff(diff) {
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
