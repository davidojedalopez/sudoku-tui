package theme

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// themeConfig is the JSON structure for theme overrides.
type themeConfig struct {
	Base      string          `json:"base"`
	Overrides paletteOverride `json:"overrides"`
}

type paletteOverride struct {
	Palette map[string]string `json:"palette"`
}

// Save persists the selected theme name to the config file.
func Save(themeName string) error {
	cfgDir := filepath.Join(os.Getenv("HOME"), ".config", "sudoku-tui")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		return err
	}
	cfgPath := filepath.Join(cfgDir, "theme.json")
	cfg := themeConfig{Base: themeName}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, data, 0644)
}

// Load loads the active theme, applying any user overrides.
func Load() *Theme {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".config", "sudoku-tui", "theme.json")

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		// No config file — return default.
		return Default
	}

	var cfg themeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default
	}

	base, ok := Registry[cfg.Base]
	if !ok {
		base = Default
	}

	// Apply palette overrides.
	if len(cfg.Overrides.Palette) == 0 {
		return base
	}

	p := base.Palette
	for key, val := range cfg.Overrides.Palette {
		color := lipgloss.Color(val)
		switch key {
		case "Bg":
			p.Bg = color
		case "BgAlt":
			p.BgAlt = color
		case "Fg":
			p.Fg = color
		case "FgDim":
			p.FgDim = color
		case "FgMuted":
			p.FgMuted = color
		case "Primary":
			p.Primary = color
		case "Accent":
			p.Accent = color
		case "Success":
			p.Success = color
		case "Error":
			p.Error = color
		case "Border":
			p.Border = color
		case "BorderHeavy":
			p.BorderHeavy = color
		}
	}
	return buildTheme(base.Name, p, base.Strings)
}
