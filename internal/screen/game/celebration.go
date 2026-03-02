package game

import (
	"math"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	celebrationDuration = 90
	numParticles        = 25
)

type particle struct {
	x, y    float64
	vx, vy  float64
	life    int
	maxLife int
	char    rune
	color   lipgloss.Color
}

type celebrationState struct {
	particles []particle
	frame     int
	active    bool
	width     int
	height    int
	colors    []lipgloss.Color
}

var particleChars = []rune{'*', '+', '.', 'o', '#', '@', '^', '~'}

func newCelebration(colors []lipgloss.Color) celebrationState {
	return celebrationState{colors: colors}
}

func (c *celebrationState) start(width, height int) {
	c.active = true
	c.frame = 0
	c.width = width
	c.height = height
	c.particles = nil
	c.spawnParticles()
}

func (c *celebrationState) spawnParticles() {
	newParticles := make([]particle, numParticles)
	w := c.width
	if w < 1 {
		w = 1
	}
	h := c.height / 2
	if h < 1 {
		h = 1
	}
	for i := range newParticles {
		x := float64(rand.Intn(w))
		y := float64(rand.Intn(h))
		angle := rand.Float64() * 2 * math.Pi
		speed := 0.5 + rand.Float64()*1.5
		maxLife := 20 + rand.Intn(30)
		color := c.colors[rand.Intn(len(c.colors))]
		ch := particleChars[rand.Intn(len(particleChars))]
		newParticles[i] = particle{
			x: x, y: y,
			vx: speed * math.Cos(angle), vy: speed * math.Sin(angle),
			life: maxLife, maxLife: maxLife,
			char: ch, color: color,
		}
	}
	c.particles = append(c.particles, newParticles...)
}

func (c *celebrationState) update() {
	if !c.active {
		return
	}
	c.frame++
	if c.frame >= celebrationDuration {
		c.active = false
		return
	}
	if c.frame%20 == 0 {
		c.spawnParticles()
	}
	alive := c.particles[:0]
	for i := range c.particles {
		p := &c.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += 0.05
		p.life--
		if p.life > 0 {
			alive = append(alive, *p)
		}
	}
	c.particles = alive
}

func (c *celebrationState) render(width, height int) string {
	if !c.active || width == 0 || height == 0 {
		return strings.Repeat("\n", height)
	}
	type cell struct {
		ch    rune
		color lipgloss.Color
	}
	grid := make([][]cell, height)
	for i := range grid {
		grid[i] = make([]cell, width)
	}
	for _, p := range c.particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < width && y >= 0 && y < height {
			grid[y][x] = cell{ch: p.char, color: p.color}
		}
	}
	var sb strings.Builder
	for y, row := range grid {
		for _, cl := range row {
			if cl.ch == 0 {
				sb.WriteByte(' ')
			} else {
				sb.WriteString(lipgloss.NewStyle().Foreground(cl.color).Render(string(cl.ch)))
			}
		}
		if y < height-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
