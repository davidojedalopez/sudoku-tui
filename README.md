# sudoku-tui

Terminal Sudoku built in Go with Bubble Tea.

## What it does

- Starts in a full-screen TUI with a main menu.
- Generates new puzzles by difficulty: Easy, Medium, Hard, Expert.
- Lets you play curated puzzles from the built-in library.
- Supports digit mode, note mode, undo, conflict highlighting, a timer, and technique toasts.
- Saves finished games to history and can resume an in-progress game from disk.
- Includes multiple built-in themes.

## Requirements

- Go 1.24+
- A terminal at least `80x24`

## Run

```bash
make run
```

Or build the binary first:

```bash
make build
./bin/sudoku-tui
```

## Controls

### Menus

- `j` / `k` or arrow keys: move
- `Enter` / `Space`: select
- `Esc`: back
- `q` or `Ctrl+C`: quit

### In game

- Arrow keys or `h` `j` `k` `l`: move cursor
- `1`-`9`: place a digit
- `n`: toggle note mode
- `x`, `Backspace`, or `Delete`: erase
- `u` or `Ctrl+Z`: undo
- `Esc`: leave the current game and return to menu

## Data files

The app stores local state under `~/.config/sudoku-tui/`:

- `history.json`: completed or abandoned game history
- `session.json`: resumable in-progress game

## Project layout

- [main.go](/home/perro/dev/sudoku-tui/main.go): application entrypoint
- [internal/app/app.go](/home/perro/dev/sudoku-tui/internal/app/app.go): screen routing
- [internal/screen/menu/menu.go](/home/perro/dev/sudoku-tui/internal/screen/menu/menu.go): main menu, difficulty, themes
- [internal/screen/game/game.go](/home/perro/dev/sudoku-tui/internal/screen/game/game.go): gameplay model
- [internal/screen/library/library.go](/home/perro/dev/sudoku-tui/internal/screen/library/library.go): curated puzzle browser
- [internal/screen/historyscreen/history.go](/home/perro/dev/sudoku-tui/internal/screen/historyscreen/history.go): history view

## Development

```bash
make test
```
