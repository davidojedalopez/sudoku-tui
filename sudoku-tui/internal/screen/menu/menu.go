package menu

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/davidojeda/sudoku-tui/internal/msgs"
	"github.com/davidojeda/sudoku-tui/internal/generator"
	"github.com/davidojeda/sudoku-tui/internal/theme"
)

// diffItem represents a difficulty option.
type diffItem struct {
	label string
	clues string
	desc  string
	diff  generator.Difficulty
}

var diffItems = []diffItem{
	{"EASY", "36-40 clues", "Solvable with naked and hidden singles only. A warm-up.", generator.Easy},
	{"MEDIUM", "30-35 clues", "Requires naked pairs and pointing pairs. Good daily practice.", generator.Medium},
	{"HARD", "25-29 clues", "Advanced patterns required. Focused solving needed.", generator.Hard},
	{"EXPERT", "17-24 clues", "X-Wing, Y-Wing, and Swordfish. Not for the faint-hearted.", generator.Expert},
}

// menuLabel is a simple label+index pair for menu items.
type menuLabel struct {
	label string
}

var menuLabels = []menuLabel{
	{"New Game"},
	{"Puzzle Library"},
	{"History"},
	{"Themes"},
	{"Quit"},
}

// themeItem represents a selectable theme.
type themeItem struct {
	key   string
	label string
}

var themeItems = []themeItem{
	{"modern-charm", "Modern Charm"},
	{"zen-monolith", "Zen Monolith"},
	{"retro-phosphor", "Retro Phosphor"},
	{"matrix", "Matrix"},
}

// Model is the main menu screen.
type Model struct {
	cursor      int
	width       int
	height      int
	showDiff    bool
	diffCursor  int
	showTheme   bool
	themeCursor int
	activeTheme string
	theme       *theme.Theme
	generating  bool
}

// New creates a new menu model.
func New(th *theme.Theme) *Model {
	themeCursor := 0
	for i, item := range themeItems {
		if item.key == th.Name {
			themeCursor = i
			break
		}
	}
	return &Model{theme: th, activeTheme: th.Name, themeCursor: themeCursor}
}

// SetTheme updates the theme on the menu model.
func (m *Model) SetTheme(th *theme.Theme) {
	m.theme = th
	m.activeTheme = th.Name
	for i, item := range themeItems {
		if item.key == th.Name {
			m.themeCursor = i
			break
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.generating {
			return m, nil
		}
		if m.showDiff {
			return m, m.handleDiffKey(msg)
		}
		if m.showTheme {
			return m, m.handleThemeKey(msg)
		}
		return m, m.handleMenuKey(msg)

	case generateMsg:
		return m, m.startGenerate(msg.diff)
	}

	return m, nil
}

func (m *Model) handleMenuKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(menuLabels)-1 {
			m.cursor++
		}
	case "enter", " ":
		return m.activateItem(m.cursor)
	case "q", "ctrl+c":
		return tea.Quit
	}
	return nil
}

func (m *Model) activateItem(idx int) tea.Cmd {
	switch idx {
	case 0: // New Game
		m.showDiff = true
		return nil
	case 1: // Puzzle Library
		return func() tea.Msg { return msgs.NavigateMsg{To: msgs.ScreenLibrary} }
	case 2: // History
		return func() tea.Msg { return msgs.NavigateMsg{To: msgs.ScreenHistory} }
	case 3: // Themes
		m.showTheme = true
		return nil
	case 4: // Quit
		return tea.Quit
	}
	return nil
}

type generateMsg struct {
	diff generator.Difficulty
}

func (m *Model) handleDiffKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if m.diffCursor > 0 {
			m.diffCursor--
		}
	case "down", "j":
		if m.diffCursor < len(diffItems)-1 {
			m.diffCursor++
		}
	case "enter", " ":
		diff := diffItems[m.diffCursor].diff
		m.generating = true
		m.showDiff = false
		return func() tea.Msg { return generateMsg{diff: diff} }
	case "esc":
		m.showDiff = false
	}
	return nil
}

func (m *Model) handleThemeKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if m.themeCursor > 0 {
			m.themeCursor--
		}
	case "down", "j":
		if m.themeCursor < len(themeItems)-1 {
			m.themeCursor++
		}
	case "enter", " ":
		selected := themeItems[m.themeCursor].key
		m.showTheme = false
		return func() tea.Msg { return msgs.ChangeThemeMsg{ThemeName: selected} }
	case "esc":
		m.showTheme = false
	}
	return nil
}

func (m *Model) startGenerate(diff generator.Difficulty) tea.Cmd {
	return func() tea.Msg {
		puzzle, solution := generator.Generate(diff)
		m.generating = false
		return msgs.StartGameMsg{
			Difficulty: diff,
			Puzzle:     puzzle,
			Solution:   solution,
		}
	}
}

func (m *Model) View() string {
	if m.width < 80 || m.height < 24 {
		return "Terminal too small. Please resize to at least 80x24."
	}

	th := m.theme

	logo := th.Menu.Title.Render(sudokuLogo)
	subtitle := th.Menu.Subtitle.Render(th.Strings.MenuSubtitle)

	var menuLines []string
	for i, item := range menuLabels {
		if i == m.cursor {
			menuLines = append(menuLines, th.Menu.ItemActive.Render("▶  "+item.label))
		} else {
			menuLines = append(menuLines, th.Menu.Item.Render("   "+item.label))
		}
	}
	menuStr := strings.Join(menuLines, "\n")

	if m.generating {
		menuStr = th.Menu.Subtitle.Render("  Generating puzzle...")
	}

	body := lipgloss.JoinVertical(lipgloss.Center, logo, "", subtitle, "", "", menuStr)
	centeredBody := lipgloss.Place(m.width, m.height-2, lipgloss.Center, lipgloss.Center, body)

	hints := th.Footer.KeyHint.Render("k/j") + th.Footer.KeyLabel.Render(" Navigate") +
		"   " + th.Footer.KeyHint.Render("Enter") + th.Footer.KeyLabel.Render(" Select") +
		"   " + th.Footer.KeyHint.Render("q") + th.Footer.KeyLabel.Render(" Quit")
	footer := th.Footer.Bar.Width(m.width).Render("  " + hints)

	result := centeredBody + "\n" + footer

	if m.showDiff {
		modal := m.renderDiffModal()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	if m.showTheme {
		modal := m.renderThemeModal()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
	}

	return result
}

func (m *Model) renderDiffModal() string {
	th := m.theme

	title := th.Diff.Title.Render("[ SELECT DIFFICULTY ]")
	tagline := th.Diff.Tagline.Render(th.Strings.DiffTagline)
	sep := strings.Repeat("─", 35)

	var lines []string
	lines = append(lines, title)
	lines = append(lines, tagline)
	lines = append(lines, sep)
	lines = append(lines, "")

	for i, d := range diffItems {
		var label string
		if i == m.diffCursor {
			label = th.Diff.Active.Render("▶ " + d.label)
		} else {
			label = th.Diff.Option.Render("  " + d.label)
		}
		clues := th.Diff.ClueCount.Render("  " + d.clues)
		lines = append(lines, label)
		lines = append(lines, clues)
		lines = append(lines, "")
	}

	lines = append(lines, sep)
	desc := th.Diff.Desc.Width(35).Render(th.Diff.DescIcon.Render("i") + " " + diffItems[m.diffCursor].desc)
	lines = append(lines, desc)
	lines = append(lines, "")
	hints := th.Footer.KeyHint.Render("k") + " UP  " + th.Footer.KeyHint.Render("j") + " DOWN  " +
		th.Footer.KeyHint.Render("Esc") + " BACK  " + th.Footer.KeyHint.Render("Enter") + " SELECT"
	lines = append(lines, hints)

	content := strings.Join(lines, "\n")
	return th.Diff.ModalBorder.Padding(1, 2).Render(content)
}

func (m *Model) renderThemeModal() string {
	th := m.theme
	title := th.Diff.Title.Render("[ SELECT THEME ]")
	sep := strings.Repeat("─", 30)

	var lines []string
	lines = append(lines, title)
	lines = append(lines, sep)
	lines = append(lines, "")

	for i, t := range themeItems {
		label := t.label
		if t.key == m.activeTheme {
			label += " ✓"
		}
		var rendered string
		if i == m.themeCursor {
			rendered = th.Diff.Active.Render("▶ " + label)
		} else {
			rendered = th.Diff.Option.Render("  " + label)
		}
		lines = append(lines, rendered)
		lines = append(lines, "")
	}

	lines = append(lines, sep)
	hints := th.Footer.KeyHint.Render("k") + " UP  " + th.Footer.KeyHint.Render("j") + " DOWN  " +
		th.Footer.KeyHint.Render("Esc") + " BACK  " + th.Footer.KeyHint.Render("Enter") + " SELECT"
	lines = append(lines, hints)

	content := strings.Join(lines, "\n")
	return th.Diff.ModalBorder.Padding(1, 2).Render(content)
}
