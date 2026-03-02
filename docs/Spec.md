# Objective
Build a complete, professional Sudoku game as a terminal UI (TUI), providing a distraction-free Sudoku experience for people who live in their terminals.

# Audience
Developers and terminal-power-users who are waiting for builds, AI prompts, or need a quick mental break during meetings.

# Tech Stack
 - Language: Go
 - TUI Framework: Charm ecosystem (Bubble Tea, Lip Gloss, Bubbles, Huh)
 - Animation: Harmonica (spring physics for celebrations)
 ---

# Screens
## Main Menu
 - Options: New Game, Puzzle Library, History, Quit
 - Selecting "New Game" opens the difficulty selector
## Difficulty Selector
 - Four levels: Easy, Medium, Hard, Expert
 - Selecting a difficulty auto-generates a puzzle and starts the game
## Game Board (Main Gameplay)
 - 9×9 Sudoku grid with visual distinction between 3×3 boxes
 - Cell notes, conflict highlighting, same-number highlighting
 - Timer, mode indicator, remaining numbers, keyboard legend
 - Celebration animation on completion
 - Technique detection toasts
## Puzzle Library
 - Browse pre-curated puzzles from an external JSON file
 - Shows puzzle name, difficulty, and author
 - Selecting a puzzle starts the game
## History View
 - Table of past games: date, difficulty, time, completion status
 - Stored locally as JSON at ~/.config/sudoku-tui/history.json

---

# Core Features

## Puzzle sources
 - Auto-generated: Constraint-propagation + backtracking solver, difficulty graded by required technique sets
 - Pre-curated library: Loaded from puzzles/curated.json (embedded via go:embed, also loadable from external file)
 - ~~Bring your own Sudoku~~ (deferred): Image-to-puzzle import — not in initial build
## Cell Notes
 - Notes placed in correct position within a 3×3 grid inside each cell (1=top-left, 5=center, 9=bottom-right)
 - Toggle between digit mode and note mode with n key, with a clear visual indicator of current mode
 - While in digit mode, use Shift+number to place individual notes without switching modes
 - Notes are cleared when a digit is placed in the cell
## Conflict Highlighting (Mistakes)
 - When a user-placed number conflicts with the same number in its row, column, or 3×3 box, both conflicting cells are highlighted in red + bold
 - We do not track a mistake counter — this is purely visual feedback
## Same-Number Highlighting
 - When the cursor is on a cell containing a digit, all other cells with the same digit are visually highlighted
## Remaining Numbers Viewer
 - Below the board, show digits 1-9
 - Digits that appear 9 times on the board are grayed out (complete)
## Undo
 - Full undo stack for digit placement, note changes, and erasure
 - Triggered with u or Ctrl+Z
 ---

# Advanced Features

## Advanced Technique Detection (must have)
 - Detect when the user applies advanced Sudoku techniques:
   - Basic: Naked Single, Hidden Single
   - Intermediate: Naked Pair/Triple, Pointing Pair
   - Advanced: X-Wing, Y-Wing, Swordfish
 - Detection is reactive: snapshot candidates before each move, compare after placement
 - On detection: display a themed toast message (e.g., "X-Wing detected! Nice!") for ~3 seconds
## Timer
 - Tracks time to solve each puzzle
 - Pauses when terminal loses focus (uses tea.FocusMsg / tea.BlurMsg)
 - Displayed in the status bar below the board
## Celebration Animation
 - ASCII fireworks/confetti using particle physics (Harmonica springs)
 - Triggered on successful board completion
 - Runs for ~3 seconds, then shows completion stats
## History
 - Stored at ~/.config/sudoku-tui/history.json
 - Each entry records: difficulty, start time, duration, completion status, board snapshot (puzzle digits vs user digits vs solution)

  ---

## Theme System
 - Architecture supports multiple themes via a Theme struct with Lip Gloss styles for every visual element
 - Ships with 1 default dark theme
 - Users can override colors via ~/.config/sudoku-tui/theme.json (partial overrides on top of base theme)
 - Designed for easy extension — adding a new theme is just defining a new Theme struct

 ---

# Keyboard Controls

| Key                     | Action                                       |
| ----------------------- | -------------------------------------------- |
| Arrow keys / h, j, k, l | Navigate the board                           |
| 1–9                     | Place digit (digit mode) or note (note mode) |
| Shift + 1–9             | Place note while in digit mode               |
| n                       | Toggle digit/note mode                       |
| Backspace / x           | Erase cell value                             |
| u / Ctrl+Z              | Undo last action                             |
| Esc                     | Return to menu                               |
| Ctrl+C                  | Quit                                         |
A keyboard legend is always visible below the board.

---

# Testing Requirements
 - Unit tests for all packages (board logic, generator correctness, technique detection, history persistence, curated loader)
 - Model tests for Bubble Tea screens (synthetic message → assert state)
 - Table-driven tests for conflict detection and technique patterns
 - Stress tests: generate 100+ puzzles, verify validity and unique solutions
 - Run with race detector: go test -v -race ./...
--- 

# Non-Goals (for initial build)
 - Image-based puzzle import (Bring Your Own Sudoku)
 - Online multiplayer or leaderboards
 - Mobile or web interface
 - Multiple pre-shipped themes (just the extensible system + 1 default)