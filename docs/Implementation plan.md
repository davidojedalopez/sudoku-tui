Context

 Build a complete, professional Sudoku TUI in Go using the Charm ecosystem. The target audience is developers/terminal users who want a distraction-free Sudoku experience. This is a greenfield project — no
 existing code.

 Key decisions:
 - Build all features in one pass (no phasing)
 - "Bring your own Sudoku" via image import is deferred
 - Advanced technique detection (X-Wing, Y-Wing, Swordfish) is must have
 - Pre-curated puzzles loaded from external JSON file (also embedded via go:embed)
 - History stored as JSON at ~/.config/sudoku-tui/history.json
 - Theme system: 1 default theme, architecture must support easy extensions via config
 - Celebration: ASCII fireworks/confetti animation on completion

 Project Structure

 sudoku-tui/
 ├── go.mod / go.sum
 ├── main.go
 ├── Makefile
 ├── puzzles/
 │   └── curated.json
 ├── internal/
 │   ├── app/
 │   │   ├── app.go              # Root Bubble Tea model, screen routing
 │   │   └── messages.go         # Cross-screen navigation messages
 │   ├── board/
 │   │   ├── cell.go             # Cell struct (value, kind, notes, conflict)
 │   │   ├── board.go            # 9x9 Board + methods
 │   │   └── validation.go       # Conflict detection (row/col/box)
 │   ├── generator/
 │   │   ├── solver.go           # Constraint-propagation + backtracking
 │   │   ├── generator.go        # Puzzle generation (fill grid, remove clues)
 │   │   └── difficulty.go       # Difficulty grading by technique requirements
 │   ├── techniques/
 │   │   ├── detector.go         # Snapshot-based detection framework
 │   │   ├── naked_single.go / hidden_single.go
 │   │   ├── naked_pair.go / pointing_pair.go
 │   │   └── xwing.go / ywing.go / swordfish.go
 │   ├── screen/
 │   │   ├── menu/menu.go        # Main menu + difficulty selector
 │   │   ├── game/
 │   │   │   ├── game.go         # Game model
 │   │   │   ├── view.go         # Board rendering (5×3 cells with notes)
 │   │   │   ├── input.go        # Keyboard handling
 │   │   │   ├── timer.go        # Timer with focus/blur pause
 │   │   │   ├── remaining.go    # Remaining numbers display
 │   │   │   └── celebration.go  # ASCII fireworks particle system
 │   │   ├── history/history.go  # History viewer (table)
 │   │   └── library/library.go  # Curated puzzle browser (list)
 │   ├── theme/
 │   │   ├── theme.go            # Theme struct + registry
 │   │   ├── default.go          # Default dark theme
 │   │   └── config.go           # JSON-based theme overrides
 │   ├── history/
 │   │   ├── entry.go            # History entry struct
 │   │   └── store.go            # JSON file persistence
 │   └── curated/
 │       └── loader.go           # Load curated puzzles from JSON

 Dependencies

 ┌────────────────────────────────────┬──────────────────────────────────────────┐
 │               Module               │                 Purpose                  │
 ├────────────────────────────────────┼──────────────────────────────────────────┤
 │ charm.land/bubbletea/v2            │ TUI framework (Elm architecture)         │
 ├────────────────────────────────────┼──────────────────────────────────────────┤
 │ charm.land/lipgloss/v2             │ Styling, layout, borders                 │
 ├────────────────────────────────────┼──────────────────────────────────────────┤
 │ charm.land/bubbles/v2              │ List, table, help components             │
 ├────────────────────────────────────┼──────────────────────────────────────────┤
 │ charm.land/huh/v2                  │ Form-based difficulty selector           │
 ├────────────────────────────────────┼──────────────────────────────────────────┤
 │ github.com/charmbracelet/harmonica │ Spring physics for celebration animation │
 └────────────────────────────────────┴──────────────────────────────────────────┘

 Key Data Structures

 Cell

 - Value int (0=empty, 1-9), Kind (Empty/Given/User), Notes [9]bool, Conflict bool

 Board

 - Cells [9][9]Cell, Solution [9][9]int
 - Methods: Set, Erase, Peers, Candidates, AllCandidates, CellsWithValue, RemainingCounts, IsSolved

 Game Model

 - Board + Cursor + InputMode (digit/note) + NoteHeld flag
 - Undo stack (pos, old value, old kind, old notes)
 - Timer (StartTime, Elapsed, Paused) — pauses on tea.BlurMsg, resumes on tea.FocusMsg
 - Technique detector + celebration state

 Theme

 - Lip Gloss styles for all visual elements (board, cells, modes, timer, etc.)
 - Registry map + JSON config overrides from ~/.config/sudoku-tui/theme.json

 Architecture

 Elm pattern with Bubble Tea v2. Root model (app.go) routes between 4 screens: menu, game, history, library. Screens communicate upward via tea.Cmd returning messages from messages.go.

 main.go options: tea.WithAltScreen(), tea.WithReportFocus() (timer pause), tea.WithKeyboardEnhancements(tea.WithKeyReleases) (shift-hold for notes).

 Sudoku Generation

 1. Generate solved grid: Fill 3 diagonal boxes randomly, solve rest with constraint-propagation + backtracking (MRV heuristic)
 2. Remove clues: Randomly remove cells, verify unique solution after each removal, stop at target clue count (Easy=45, Medium=35, Hard=28, Expert=22)
 3. Grade difficulty: Solve with limited technique sets to confirm label accuracy

 Technique Detection (Reactive)

 4. Before user move: snapshot all candidate sets
 5. After placement: check if simpler techniques could have found this digit
 6. If not: check advanced patterns (X-Wing, Y-Wing, Swordfish)
 7. On detection: show themed toast for ~3 seconds

 Board Rendering

 Each cell: 5 chars wide × 3 lines tall for 3×3 note grid:
 Notes:    Digit:
 1 2 3
 4   6       7
 7 8 9
 Heavy borders between 3×3 boxes, thin borders between cells.

 Keyboard Controls

 ┌───────────────┬───────────────────────────────────────┐
 │      Key      │                Action                 │
 ├───────────────┼───────────────────────────────────────┤
 │ Arrows / hjkl │ Move cursor                           │
 ├───────────────┼───────────────────────────────────────┤
 │ 1-9           │ Place digit or note (depends on mode) │
 ├───────────────┼───────────────────────────────────────┤
 │ Shift+1-9     │ Place note while in digit mode        │
 ├───────────────┼───────────────────────────────────────┤
 │ n             │ Toggle digit/note mode                │
 ├───────────────┼───────────────────────────────────────┤
 │ Backspace/x   │ Erase cell                            │
 ├───────────────┼───────────────────────────────────────┤
 │ u / Ctrl+Z    │ Undo                                  │
 ├───────────────┼───────────────────────────────────────┤
 │ Esc           │ Back to menu                          │
 └───────────────┴───────────────────────────────────────┘

 Build Order

 ┌───────┬────────────────────────────────────────────────────────────┬───────────────┐
 │   #   │                            What                            │   Milestone   │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 1     │ Scaffold: go.mod, main.go, Makefile                        │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 2-4   │ board/ package (cell, board, validation) + tests           │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 5-7   │ generator/ package (solver, generator, difficulty) + tests │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 8     │ theme/ package                                             │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 9-13  │ screen/game/ (model, view, input, timer, remaining)        │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 14-15 │ screen/menu/ + app/ wiring                                 │ Playable game │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 16-17 │ history/ package + screen/history/                         │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 18-19 │ curated/ + screen/library/                                 │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 20    │ techniques/ detection + all technique files + tests        │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 21    │ screen/game/celebration.go (fireworks)                     │               │
 ├───────┼────────────────────────────────────────────────────────────┼───────────────┤
 │ 22    │ Polish: help text, edge cases, theme tuning                │ Complete      │
 └───────┴────────────────────────────────────────────────────────────┴───────────────┘

 Testing Strategy

 - Unit tests for every package (board, generator, techniques, history, curated)
 - Model tests: call Update() with synthetic messages, assert state
 - Table-driven tests for conflict detection and technique patterns
 - Stress test: generate 100 puzzles, all valid + unique solution
 - Run: go test -v -race ./...

 Verification

 1. make build compiles
 2. make test passes with race detector
 3. Play: menu → difficulty → generate → play → complete → fireworks
 4. Notes, conflicts, highlighting, timer pause, undo all work
 5. Browse curated library → play
 6. History saves to ~/.config/sudoku-tui/history.json