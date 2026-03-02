package game

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/davidojeda/sudoku-tui/internal/board"
	"github.com/davidojeda/sudoku-tui/internal/msgs"
	"github.com/davidojeda/sudoku-tui/internal/session"
	"github.com/davidojeda/sudoku-tui/internal/techniques"
	"github.com/davidojeda/sudoku-tui/internal/theme"
)

// InputMode is digit or note mode.
type InputMode int

const (
	ModDigit InputMode = iota
	ModNote
)

// undoEntry records a cell state for undo.
type undoEntry struct {
	row, col int
	oldValue int
	oldKind  board.CellKind
	oldNotes [9]bool
}

// toastState manages the technique detection toast.
type toastState struct {
	message string
	expires time.Time
	visible bool
}

// Model is the game screen model.
type Model struct {
	board         *board.Board
	cursor        [2]int
	inputMode     InputMode
	undoStack     []undoEntry
	timer         timerState
	toast         toastState
	celebration   celebrationState
	solved        bool
	width         int
	height        int
	difficulty    string
	puzzleID      string
	puzzle        string
	solution      string
	resumeElapsed int64
	theme         *theme.Theme
	snap          techniques.Snapshot
}

// New creates a new game model.
func New(puzzle, solution, difficulty, puzzleID string, th *theme.Theme) *Model {
	b := board.New()
	b.LoadGivens(puzzle)

	for i, ch := range solution {
		if ch >= '1' && ch <= '9' {
			b.Solution[i/9][i%9] = int(ch - '0')
		}
	}

	m := &Model{
		board:       b,
		difficulty:  difficulty,
		puzzleID:    puzzleID,
		puzzle:      puzzle,
		solution:    solution,
		theme:       th,
		celebration: newCelebration(th.Victory.Particle),
	}
	m.snap = techniques.TakeSnapshot(b)
	return m
}

// NewFromSaved restores a game from a previously saved session.
func NewFromSaved(saved *session.SavedGame, th *theme.Theme) *Model {
	b := board.New()
	b.LoadGivens(saved.Puzzle)
	for i, ch := range saved.Solution {
		if ch >= '1' && ch <= '9' {
			b.Solution[i/9][i%9] = int(ch - '0')
		}
	}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			cs := saved.Cells[r][c]
			b.Cells[r][c].Value = cs.Value
			b.Cells[r][c].Kind = board.CellKind(cs.Kind)
			b.Cells[r][c].Notes = cs.Notes
		}
	}
	b.UpdateConflicts()

	m := &Model{
		board:         b,
		difficulty:    saved.Difficulty,
		puzzleID:      saved.PuzzleID,
		puzzle:        saved.Puzzle,
		solution:      saved.Solution,
		resumeElapsed: saved.Elapsed,
		theme:         th,
		celebration:   newCelebration(th.Victory.Particle),
	}
	m.snap = techniques.TakeSnapshot(b)
	return m
}

// Snapshot captures the current game state for persistence.
func (m *Model) Snapshot() *session.SavedGame {
	var cells [9][9]session.CellState
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			cell := m.board.Cells[r][c]
			cells[r][c] = session.CellState{
				Value: cell.Value,
				Kind:  int(cell.Kind),
				Notes: cell.Notes,
			}
		}
	}
	return &session.SavedGame{
		Puzzle:     m.puzzle,
		Solution:   m.solution,
		Difficulty: m.difficulty,
		PuzzleID:   m.puzzleID,
		Elapsed:    m.timer.elapsedSeconds(),
		Cells:      cells,
	}
}

// SetTheme updates the theme on the game screen.
func (m *Model) SetTheme(th *theme.Theme) { m.theme = th }

// Init starts the timer and returns the first tick command.
func (m *Model) Init() tea.Cmd {
	if m.resumeElapsed > 0 {
		return m.timer.startAt(time.Duration(m.resumeElapsed) * time.Second)
	}
	return m.timer.start()
}

// Update handles messages.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		if m.solved {
			m.celebration.update()
			if m.celebration.active {
				return m, tickCmd()
			}
			return m, nil
		}
		return m, m.timer.tick()

	case tea.KeyMsg:
		if m.solved {
			switch msg.String() {
			case "enter", "esc":
				elapsed := m.timer.elapsedSeconds()
				diff := m.difficulty
				return m, func() tea.Msg {
					return msgs.GameOverMsg{Won: true, Elapsed: elapsed, Diff: diff}
				}
			}
			return m, nil
		}
		cmd := m.handleKey(msg)
		return m, cmd

	case tea.FocusMsg:
		return m, m.timer.resume()

	case tea.BlurMsg:
		m.timer.pause()
		return m, nil
	}

	return m, nil
}

// View renders the game screen.
func (m *Model) View() string {
	if m.width < 80 || m.height < 24 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			"Terminal too small.\nPlease resize to at least 80x24.")
	}

	if m.solved {
		return m.renderSolved()
	}

	return m.renderGame()
}

// ElapsedSeconds returns the total elapsed seconds.
func (m *Model) ElapsedSeconds() int64 {
	return m.timer.elapsedSeconds()
}

// Difficulty returns the game difficulty string.
func (m *Model) Difficulty() string {
	return m.difficulty
}

// placeDigit places a digit (or note) at the cursor.
func (m *Model) placeDigit(digit int, forceNote bool) {
	r, c := m.cursor[0], m.cursor[1]
	cell := m.board.Cells[r][c]
	if cell.Kind == board.KindGiven {
		return
	}

	m.undoStack = append(m.undoStack, undoEntry{
		row: r, col: c,
		oldValue: cell.Value,
		oldKind:  cell.Kind,
		oldNotes: cell.Notes,
	})

	isNote := forceNote || m.inputMode == ModNote
	if isNote {
		m.board.ToggleNote(r, c, digit)
	} else {
		m.snap = techniques.TakeSnapshot(m.board)
		m.board.Set(r, c, digit)
	}
}

func (m *Model) eraseCell() {
	r, c := m.cursor[0], m.cursor[1]
	cell := m.board.Cells[r][c]
	if cell.Kind == board.KindGiven {
		return
	}
	m.undoStack = append(m.undoStack, undoEntry{
		row: r, col: c,
		oldValue: cell.Value,
		oldKind:  cell.Kind,
		oldNotes: cell.Notes,
	})
	m.board.Erase(r, c)
}

func (m *Model) undo() {
	if len(m.undoStack) == 0 {
		return
	}
	last := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]

	cell := &m.board.Cells[last.row][last.col]
	cell.Value = last.oldValue
	cell.Kind = last.oldKind
	cell.Notes = last.oldNotes
	m.board.UpdateConflicts()
}

func (m *Model) checkTechnique() tea.Cmd {
	r, c := m.cursor[0], m.cursor[1]
	val := m.board.Cells[r][c].Value
	if val == 0 {
		return nil
	}

	tech := techniques.DetectTechnique(m.board, m.snap, r, c, val)
	if tech != techniques.TechNone {
		m.toast = toastState{
			message: tech.Message(),
			expires: time.Now().Add(3 * time.Second),
			visible: true,
		}
	}

	if m.board.IsSolved() {
		m.solved = true
		m.timer.pause()
		m.celebration.start(m.width, m.height)
		return tickCmd()
	}

	return nil
}

func (m *Model) cursorValue() int {
	return m.board.Cells[m.cursor[0]][m.cursor[1]].Value
}

// renderGame renders the main game layout.
func (m *Model) renderGame() string {
	header := m.renderHeader()
	footer := m.renderFooter()

	bodyHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)

	boardStr := m.renderBoard()
	sidebar := m.renderSidebar()

	toast := ""
	if m.toast.visible && time.Now().Before(m.toast.expires) {
		toast = m.renderToast()
	} else {
		m.toast.visible = false
	}

	body := lipgloss.JoinHorizontal(lipgloss.Top, boardStr, "  ", sidebar)
	body = lipgloss.NewStyle().Height(bodyHeight).Render(body)

	result := lipgloss.JoinVertical(lipgloss.Left, header, body, toast, footer)
	return result
}

// renderSolved renders the victory screen.
func (m *Model) renderSolved() string {
	th := m.theme
	w, h := m.width, m.height

	badge := th.Victory.Badge.Render(th.Strings.VictoryBadge)
	title := th.Victory.Title.Render(th.Strings.VictoryTitle)
	subtitle := th.Victory.Subtitle.Render(th.Strings.VictorySubtitle)

	timeStr := th.Victory.StatValue.Render(m.timer.formatted())
	diffStr := th.Victory.StatValue.Render(strings.ToUpper(m.difficulty))

	stats := th.Victory.StatBox.Render(
		"  " + th.Victory.StatLabel.Render("TIME      ") + timeStr + "\n" +
			"  " + th.Victory.StatLabel.Render("DIFFICULTY") + diffStr,
	)

	btn := th.Victory.Button.Render("[ PRESS ENTER ]")
	secondary := th.Victory.ButtonSecondary.Render("Return to Menu (Esc)")

	modalContent := lipgloss.JoinVertical(lipgloss.Center,
		"",
		badge,
		"",
		title,
		subtitle,
		"",
		stats,
		"",
		btn,
		secondary,
		"",
	)

	modal := th.Victory.ModalBorder.Width(40).Render(modalContent)
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, modal)
}

func (m *Model) renderHeader() string {
	th := m.theme
	title := th.Header.Title.Render("SUDOKU")
	diff := th.Header.Meta.Render(strings.ToUpper(m.difficulty))
	gap := strings.Repeat(" ", maxInt(0, m.width-lipgloss.Width(title)-lipgloss.Width(diff)-4))
	return th.Header.Bar.Width(m.width).Render(" " + title + gap + diff + " ")
}

func (m *Model) renderFooter() string {
	th := m.theme
	hints := []string{
		th.Footer.KeyHint.Render("hjkl") + th.Footer.KeyLabel.Render(" Move"),
		th.Footer.KeyHint.Render("1-9") + th.Footer.KeyLabel.Render(" Digit"),
		th.Footer.KeyHint.Render("n") + th.Footer.KeyLabel.Render(" Mode"),
		th.Footer.KeyHint.Render("x") + th.Footer.KeyLabel.Render(" Erase"),
		th.Footer.KeyHint.Render("u") + th.Footer.KeyLabel.Render(" Undo"),
		th.Footer.KeyHint.Render("Esc") + th.Footer.KeyLabel.Render(" Menu"),
	}
	return th.Footer.Bar.Width(m.width).Render("  " + strings.Join(hints, "   "))
}

func (m *Model) renderSidebar() string {
	th := m.theme
	var sb strings.Builder

	sb.WriteString(th.Sidebar.Title.Render("ELAPSED TIME") + "\n")
	sb.WriteString(th.Sidebar.Timer.Render(m.timer.formatted()) + "\n")
	sb.WriteString("\n")

	sb.WriteString(th.Sidebar.Title.Render("INPUT MODE") + "\n")
	var digitBtn, noteBtn string
	if m.inputMode == ModDigit {
		digitBtn = th.Sidebar.ModeIndicatorActive.Render("● DIGIT")
		noteBtn = th.Sidebar.ModeIndicatorInactive.Render("  NOTE ")
	} else {
		digitBtn = th.Sidebar.ModeIndicatorInactive.Render("  DIGIT")
		noteBtn = th.Sidebar.ModeIndicatorActive.Render("● NOTE ")
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, digitBtn, " ", noteBtn) + "\n")
	sb.WriteString("\n")

	remaining := m.board.RemainingCounts()
	sb.WriteString(renderRemaining(remaining, m.cursorValue(), th))

	return sb.String()
}

func (m *Model) renderToast() string {
	th := m.theme
	msg := th.Toast.Icon.Render("*") + " " + th.Toast.Text.Render(m.toast.message)
	return th.Toast.Border.Render("  " + msg + "  ")
}

// maxInt returns the larger of two ints.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
