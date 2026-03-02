package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/davidojeda/sudoku-tui/internal/history"
	"github.com/davidojeda/sudoku-tui/internal/msgs"
	"github.com/davidojeda/sudoku-tui/internal/screen/game"
	"github.com/davidojeda/sudoku-tui/internal/screen/historyscreen"
	"github.com/davidojeda/sudoku-tui/internal/screen/library"
	"github.com/davidojeda/sudoku-tui/internal/screen/menu"
	"github.com/davidojeda/sudoku-tui/internal/session"
	"github.com/davidojeda/sudoku-tui/internal/theme"
)

// App is the root Bubble Tea model.
type App struct {
	screen    msgs.Screen
	menu      *menu.Model
	game      *game.Model
	history   *historyscreen.Model
	library   *library.Model
	store     *history.Store
	theme     *theme.Theme
	savedGame *session.SavedGame
	width     int
	height    int
}

// New creates a new App model.
func New() (*App, error) {
	th := theme.Load()

	store, _ := history.NewStore()

	savedGame, hasSave, _ := session.Load()

	a := &App{
		screen:    msgs.ScreenMenu,
		theme:     th,
		store:     store,
		savedGame: savedGame,
	}
	a.menu = menu.New(th, hasSave)
	a.history = historyscreen.New(th, store)
	a.library = library.New(th)
	return a, nil
}

func (a *App) Init() tea.Cmd {
	return a.menu.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, a.forwardToActive(msg)

	case msgs.NavigateMsg:
		a.screen = msg.To
		switch msg.To {
		case msgs.ScreenHistory:
			a.history = historyscreen.New(a.theme, a.store)
			return a, a.history.Init()
		case msgs.ScreenLibrary:
			return a, tea.Batch(
				a.library.Init(),
				func() tea.Msg { return tea.WindowSizeMsg{Width: a.width, Height: a.height} },
			)
		case msgs.ScreenMenu:
			return a, a.menu.Init()
		}
		return a, nil

	case msgs.StartGameMsg:
		go session.Clear() //nolint:errcheck
		a.savedGame = nil
		a.menu.SetHasSavedGame(false)
		a.screen = msgs.ScreenGame
		diff := msg.Difficulty.String()
		a.game = game.New(msg.Puzzle, msg.Solution, diff, msg.PuzzleID, a.theme)
		return a, tea.Batch(
			a.game.Init(),
			func() tea.Msg { return tea.WindowSizeMsg{Width: a.width, Height: a.height} },
		)

	case msgs.ResumeGameMsg:
		if a.savedGame != nil {
			a.screen = msgs.ScreenGame
			a.game = game.NewFromSaved(a.savedGame, a.theme)
			return a, tea.Batch(
				a.game.Init(),
				func() tea.Msg { return tea.WindowSizeMsg{Width: a.width, Height: a.height} },
			)
		}
		return a, nil

	case msgs.GameOverMsg:
		if msg.Won {
			go session.Clear() //nolint:errcheck
			a.savedGame = nil
			a.menu.SetHasSavedGame(false)
		} else if a.game != nil {
			snap := a.game.Snapshot()
			a.savedGame = snap
			go session.Save(snap) //nolint:errcheck
			a.menu.SetHasSavedGame(true)
		}
		if a.store != nil {
			result := history.ResultGaveUp
			if msg.Won {
				result = history.ResultWin
			}
			_ = a.store.Add(msg.Diff, msg.Elapsed, result, "")
		}
		a.screen = msgs.ScreenMenu
		return a, a.menu.Init()

	case msgs.ChangeThemeMsg:
		th, ok := theme.Registry[msg.ThemeName]
		if !ok {
			return a, nil
		}
		a.theme = th
		a.menu.SetTheme(th)
		a.history.SetTheme(th)
		a.library.SetTheme(th)
		if a.game != nil {
			a.game.SetTheme(th)
		}
		go theme.Save(msg.ThemeName) //nolint:errcheck
		return a, func() tea.Msg { return tea.WindowSizeMsg{Width: a.width, Height: a.height} }
	}

	return a, a.forwardToActive(msg)
}

// forwardToActive forwards a message to the active screen model and updates state.
func (a *App) forwardToActive(msg tea.Msg) tea.Cmd {
	switch a.screen {
	case msgs.ScreenMenu:
		newModel, cmd := a.menu.Update(msg)
		a.menu = newModel.(*menu.Model)
		return cmd
	case msgs.ScreenGame:
		if a.game != nil {
			newModel, cmd := a.game.Update(msg)
			a.game = newModel.(*game.Model)
			return cmd
		}
	case msgs.ScreenHistory:
		newModel, cmd := a.history.Update(msg)
		a.history = newModel.(*historyscreen.Model)
		return cmd
	case msgs.ScreenLibrary:
		newModel, cmd := a.library.Update(msg)
		a.library = newModel.(*library.Model)
		return cmd
	}
	return nil
}

func (a *App) View() string {
	if a.width < 80 || a.height < 24 {
		return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center,
			"Terminal too small.\nPlease resize to at least 80x24.")
	}

	switch a.screen {
	case msgs.ScreenMenu:
		return a.menu.View()
	case msgs.ScreenGame:
		if a.game != nil {
			return a.game.View()
		}
	case msgs.ScreenHistory:
		return a.history.View()
	case msgs.ScreenLibrary:
		return a.library.View()
	}
	return ""
}
