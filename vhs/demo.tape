# Demo: Theme switch, Easy puzzle start, notes, and digit entry
Output vhs/demo.gif

Require bash

Set Shell "bash"
Set FontSize 14
Set Width 1280
Set Height 900
Set Padding 20
Set Theme "Catppuccin Mocha"
Set Framerate 24
Set PlaybackSpeed 1.15
Set TypingSpeed 0ms
Set CursorBlink false

# Clean saved state and start the app
Hide
Type "rm -f ~/.config/sudoku-tui/session.json ~/.config/sudoku-tui/theme.json"
Enter
Sleep 300ms
Type "unset NO_COLOR; TERM=xterm-256color COLORTERM=truecolor ./bin/sudoku-tui"
Enter
Sleep 1s
Show

# Switch once from the default app theme to Retro Phosphor
Sleep 1s
Down 3
Sleep 400ms
Enter
Sleep 700ms
Down 2
Sleep 400ms
Enter
Sleep 1.5s

# Start the first curated Easy Sudoku from the puzzle library
Up 10
Sleep 200ms
Down
Sleep 300ms
Enter
Sleep 1s
Up 10
Sleep 200ms
Enter
Sleep 1.2s

# Add notes to the first empty cell in the top row
Right 3
Sleep 250ms
Type "n"
Sleep 350ms
Type "7"
Sleep 300ms
Type "8"
Sleep 300ms
Type "9"
Sleep 700ms

# Move to the next empty cell and place a digit
Right 4
Sleep 250ms
Type "n"
Sleep 350ms
Type "4"
Sleep 2s

# Quit
Type "q"
Sleep 500ms
