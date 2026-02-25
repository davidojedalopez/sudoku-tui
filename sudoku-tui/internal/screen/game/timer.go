package game

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg is sent every second to update the timer.
type tickMsg struct{}

// timerState holds the game timer state.
type timerState struct {
	startTime time.Time
	elapsed   time.Duration
	paused    bool
	started   bool
}

func newTimer() timerState {
	return timerState{}
}

func (t *timerState) start() tea.Cmd {
	t.startTime = time.Now()
	t.started = true
	t.paused = false
	return tickCmd()
}

func (t *timerState) startAt(offset time.Duration) tea.Cmd {
	t.elapsed = offset
	t.startTime = time.Now()
	t.started = true
	t.paused = false
	return tickCmd()
}

func (t *timerState) pause() {
	if !t.paused && t.started {
		t.elapsed += time.Since(t.startTime)
		t.paused = true
	}
}

func (t *timerState) resume() tea.Cmd {
	if t.paused && t.started {
		t.startTime = time.Now()
		t.paused = false
		return tickCmd()
	}
	return nil
}

func (t *timerState) tick() tea.Cmd {
	if t.paused || !t.started {
		return nil
	}
	return tickCmd()
}

func (t *timerState) totalElapsed() time.Duration {
	if !t.started {
		return 0
	}
	if t.paused {
		return t.elapsed
	}
	return t.elapsed + time.Since(t.startTime)
}

func (t *timerState) elapsedSeconds() int64 {
	return int64(t.totalElapsed().Seconds())
}

func (t *timerState) formatted() string {
	total := t.totalElapsed()
	h := int(total.Hours())
	m := int(total.Minutes()) % 60
	s := int(total.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}
