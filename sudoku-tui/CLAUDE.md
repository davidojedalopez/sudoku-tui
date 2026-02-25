# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run
go run .
make run

# Build binary to bin/sudoku-tui
make build

# Test (with race detector)
make test
go test -v -race ./...

# Run a single package's tests
go test -v -race ./internal/board/...
go test -v -race ./internal/techniques/...
go test -v -race ./internal/generator/...

# Tidy dependencies
make deps
```

Minimum terminal size: 80×24. The app uses alt-screen mode.

## Architecture

**Entry point:** `main.go` → `internal/app/app.go`

The app is a [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI with four screens routed by a top-level `App` model in `internal/app/`. Screen transitions happen via message passing — screens emit `msgs.NavigateMsg`, `msgs.StartGameMsg`, or `msgs.GameOverMsg` which `App.Update()` intercepts before forwarding to the active screen.

### Message flow (`internal/msgs/msgs.go`)

`msgs` is a dedicated package that exists solely to break import cycles. All three cross-screen message types live here:

- `NavigateMsg{To: Screen}` — navigate to any screen
- `StartGameMsg{Puzzle, Solution, Difficulty, PuzzleID}` — start a game (sent by menu and library)
- `GameOverMsg{Won, Elapsed, Diff}` — game ended; `App` saves history and returns to menu

### Screens

| Screen | Package | Notes |
|--------|---------|-------|
| Menu | `internal/screen/menu/` | Difficulty modal triggers puzzle generation then emits `StartGameMsg` |
| Game | `internal/screen/game/` | Split into `game.go`, `input.go`, `view.go`, `timer.go`, `celebration.go`, `remaining.go` |
| History | `internal/screen/historyscreen/` | Paginated table, reads from `history.Store` |
| Library | `internal/screen/library/` | Curated puzzle browser; solves puzzle before emitting `StartGameMsg` |

### Core packages

- **`internal/board/`** — `Cell` (value, kind, notes), `Board` (9×9 cells + solution), validation and conflict detection
- **`internal/generator/`** — MRV backtracking solver (`solver.go`), puzzle generator (`generator.go`), difficulty grader (`difficulty.go`)
- **`internal/techniques/`** — Snapshot-based technique detection after each digit placement. `TakeSnapshot()` captures candidates before a move; `DetectTechnique()` compares after. Detects: Naked Single/Pair, Hidden Single, Pointing Pair, X-Wing, Y-Wing, Swordfish.
- **`internal/history/`** — `Store` persists game results as JSON at `~/.config/sudoku-tui/history.json`
- **`internal/curated/`** — Uses `//go:embed` to bundle `puzzles/curated.json` into the binary

### Theme system (`internal/theme/`)

`Theme` is a flat struct of typed `lipgloss.Style` groups (Board, Cell, Header, Footer, Menu, Sidebar, etc.) plus a `Palette` and `ThemeStrings`. Four built-in themes are registered in `init()`: `modern-charm` (default), `zen-monolith`, `retro-phosphor`, `matrix`.

`theme.Load()` reads `~/.config/sudoku-tui/theme.json` if present. The config specifies a `base` theme name and optional `overrides.palette` key/value pairs (e.g. `"Primary": "#ff6b6b"`). `buildTheme()` reconstructs all styles from the palette — adding a new theme means implementing a `buildTheme(name, palette, strings)` call.

### Game screen internals

- **Input modes:** `ModDigit` / `ModNote` toggled with `n`; vim-style movement (`hjkl`/arrows)
- **Undo stack:** each mutation pushes `undoEntry{row, col, oldValue, oldKind, oldNotes}`
- **Technique toast:** snapshot taken before each digit placement; technique detected after; 3-second toast shown
- **Timer:** `tea.Tick`-driven; pauses on `tea.BlurMsg`, resumes on `tea.FocusMsg`
- **Celebration:** particle animation on win; `tickMsg` drives both timer and celebration frames
