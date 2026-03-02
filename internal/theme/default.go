package theme

import "github.com/charmbracelet/lipgloss"

// ModernCharmTheme returns the default Modern Charm theme.
func ModernCharmTheme() *Theme {
	p := Palette{
		Bg:          lipgloss.Color("#1E1E2E"),
		BgAlt:       lipgloss.Color("#27273A"),
		Fg:          lipgloss.Color("#CAD3F5"),
		FgDim:       lipgloss.Color("#8087A2"),
		FgMuted:     lipgloss.Color("#51526A"),
		Primary:     lipgloss.Color("#F5A97F"),
		Accent:      lipgloss.Color("#8AADF4"),
		Success:     lipgloss.Color("#A6DA95"),
		Error:       lipgloss.Color("#ED8796"),
		Border:      lipgloss.Color("#363654"),
		BorderHeavy: lipgloss.Color("#494970"),
	}

	return buildTheme("modern-charm", p, ThemeStrings{
		AppName:         "Sudoku",
		MenuSubtitle:    "Modern Charm Edition",
		DiffTagline:     "Choose your challenge",
		VictoryBadge:    "🏆",
		VictoryTitle:    "PUZZLE SOLVED!",
		VictorySubtitle: "Excellent work! You've mastered this grid.",
	})
}

// ZenMonolithTheme returns the Zen Monolith theme.
func ZenMonolithTheme() *Theme {
	p := Palette{
		Bg:          lipgloss.Color("#0A0C14"),
		BgAlt:       lipgloss.Color("#121726"),
		Fg:          lipgloss.Color("#D0D8F0"),
		FgDim:       lipgloss.Color("#6D7A9E"),
		FgMuted:     lipgloss.Color("#3D4560"),
		Primary:     lipgloss.Color("#2B7FFF"),
		Accent:      lipgloss.Color("#5BA4FB"),
		Success:     lipgloss.Color("#4ECBA0"),
		Error:       lipgloss.Color("#F04F5F"),
		Border:      lipgloss.Color("#1E2640"),
		BorderHeavy: lipgloss.Color("#2B3A5C"),
	}

	return buildTheme("zen-monolith", p, ThemeStrings{
		AppName:         "SUDOKU",
		MenuSubtitle:    "GEOMETRIC LOGIC ENGINE",
		DiffTagline:     "INITIATE SEQUENCE",
		VictoryBadge:    "◆ S RANK ◆",
		VictoryTitle:    "PUZZLE CLEARED",
		VictorySubtitle: "ZEN MONOLITH SYSTEM · PUZZLE CLEARED",
	})
}

// RetroPhosphorTheme returns the Retro Phosphor amber CRT theme.
func RetroPhosphorTheme() *Theme {
	p := Palette{
		Bg:          lipgloss.Color("#0A0600"),
		BgAlt:       lipgloss.Color("#150D00"),
		Fg:          lipgloss.Color("#FFD15C"),
		FgDim:       lipgloss.Color("#C98A10"),
		FgMuted:     lipgloss.Color("#7A5000"),
		Primary:     lipgloss.Color("#FFB000"),
		Accent:      lipgloss.Color("#FFD15C"),
		Success:     lipgloss.Color("#FFD15C"),
		Error:       lipgloss.Color("#FF6A00"),
		Border:      lipgloss.Color("#4A3000"),
		BorderHeavy: lipgloss.Color("#7A5000"),
	}

	return buildTheme("retro-phosphor", p, ThemeStrings{
		AppName:         "SUDOKU",
		MenuSubtitle:    "// SYSTEM READY // WAITING FOR INPUT_",
		DiffTagline:     "INITIALIZE MATRIX",
		VictoryBadge:    "SYSTEM RESTORED",
		VictoryTitle:    "SEQUENCE COMPLETED",
		VictorySubtitle: "SEQUENCE COMPLETED SUCCESSFULLY",
	})
}

// MatrixTheme returns the Matrix green-on-black theme.
func MatrixTheme() *Theme {
	p := Palette{
		Bg:          lipgloss.Color("#050A05"),
		BgAlt:       lipgloss.Color("#0A180A"),
		Fg:          lipgloss.Color("#25F447"),
		FgDim:       lipgloss.Color("#13872A"),
		FgMuted:     lipgloss.Color("#0A4015"),
		Primary:     lipgloss.Color("#0BDA0B"),
		Accent:      lipgloss.Color("#13EC5B"),
		Success:     lipgloss.Color("#25F447"),
		Error:       lipgloss.Color("#FF4444"),
		Border:      lipgloss.Color("#0A3010"),
		BorderHeavy: lipgloss.Color("#0F5020"),
	}

	return buildTheme("matrix", p, ThemeStrings{
		AppName:         "SUDOKU",
		MenuSubtitle:    "WAKE UP, NEO... THE GRID AWAITS",
		DiffTagline:     "LOAD PROGRAM",
		VictoryBadge:    "ENCRYPTION BROKEN",
		VictoryTitle:    "YOU ARE THE ONE",
		VictorySubtitle: "System anomaly detected. You are the one.",
	})
}

// buildTheme constructs a complete Theme from a palette.
func buildTheme(name string, p Palette, strings ThemeStrings) *Theme {
	t := &Theme{
		Name:    name,
		Palette: p,
		Strings: strings,
	}

	// Board styles
	t.Board = BoardStyles{
		Border:     lipgloss.NewStyle().Foreground(p.BorderHeavy),
		BoxBorder:  lipgloss.NewStyle().Foreground(p.BorderHeavy),
		CellBorder: lipgloss.NewStyle().Foreground(p.Border),
	}

	// Cell styles
	t.Cell = CellStyles{
		Given:      lipgloss.NewStyle().Foreground(p.Fg).Bold(true),
		User:       lipgloss.NewStyle().Foreground(p.Accent),
		Empty:      lipgloss.NewStyle().Foreground(p.FgMuted),
		Conflict:   lipgloss.NewStyle().Background(p.Error).Foreground(p.Bg).Bold(true),
		Highlight:  lipgloss.NewStyle().Background(p.BgAlt),
		Cursor:     lipgloss.NewStyle().Background(p.Primary).Foreground(p.Bg).Bold(true),
		Notes:      lipgloss.NewStyle().Background(p.Bg),
		NotesDigit: lipgloss.NewStyle().Foreground(p.FgDim),
	}

	// Header styles
	t.Header = HeaderStyles{
		Bar:   lipgloss.NewStyle().Background(p.BgAlt).Foreground(p.Fg),
		Title: lipgloss.NewStyle().Background(p.BgAlt).Foreground(p.Primary).Bold(true),
		Meta:  lipgloss.NewStyle().Background(p.BgAlt).Foreground(p.FgDim),
	}

	// Footer styles
	t.Footer = FooterStyles{
		Bar:      lipgloss.NewStyle().Background(p.BgAlt).Foreground(p.FgDim),
		KeyHint:  lipgloss.NewStyle().Background(p.Primary).Foreground(p.Bg).Bold(true).Padding(0, 1),
		KeyLabel: lipgloss.NewStyle().Foreground(p.FgDim),
	}

	// Menu styles
	t.Menu = MenuStyles{
		Title:      lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		Subtitle:   lipgloss.NewStyle().Foreground(p.FgDim),
		Item:       lipgloss.NewStyle().Foreground(p.Fg),
		ItemActive: lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		ItemPrefix: lipgloss.NewStyle().Foreground(p.Primary),
	}

	// Difficulty styles
	t.Diff = DifficultyStyles{
		ModalBorder: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(p.BorderHeavy),
		Title:       lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		Tagline:     lipgloss.NewStyle().Foreground(p.FgDim),
		Option:      lipgloss.NewStyle().Foreground(p.Fg),
		Active:      lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		ClueCount:   lipgloss.NewStyle().Foreground(p.FgDim),
		Desc:        lipgloss.NewStyle().Foreground(p.Fg),
		DescIcon:    lipgloss.NewStyle().Foreground(p.Accent),
	}

	// Sidebar styles
	t.Sidebar = SidebarStyles{
		Border:                lipgloss.NewStyle().Foreground(p.Border),
		Title:                 lipgloss.NewStyle().Foreground(p.FgDim).Bold(true),
		Timer:                 lipgloss.NewStyle().Foreground(p.Accent).Bold(true),
		ModeIndicatorActive:   lipgloss.NewStyle().Background(p.Primary).Foreground(p.Bg).Bold(true).Padding(0, 1),
		ModeIndicatorInactive: lipgloss.NewStyle().Foreground(p.FgDim).BorderStyle(lipgloss.NormalBorder()).BorderForeground(p.Border).Padding(0, 1),
		RemainingDigit:        lipgloss.NewStyle().Foreground(p.FgDim),
		RemainingBar:          lipgloss.NewStyle().Foreground(p.Accent),
		RemainingBarEmpty:     lipgloss.NewStyle().Foreground(p.FgMuted),
		RemainingCount:        lipgloss.NewStyle().Foreground(p.FgDim),
		RemainingDone:         lipgloss.NewStyle().Foreground(p.FgMuted).Strikethrough(true),
		RemainingActive:       lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
	}

	// Badge styles
	t.Badges = BadgeStyles{
		Easy:   lipgloss.NewStyle().Background(lipgloss.Color("#2A3829")).Foreground(p.Success).Padding(0, 1),
		Medium: lipgloss.NewStyle().Background(lipgloss.Color("#2A2D47")).Foreground(p.Accent).Padding(0, 1),
		Hard:   lipgloss.NewStyle().Background(lipgloss.Color("#3D2E29")).Foreground(p.Primary).Padding(0, 1),
		Expert: lipgloss.NewStyle().Background(lipgloss.Color("#3D2929")).Foreground(p.Error).Padding(0, 1),
	}

	// Library styles
	t.Library = LibraryStyles{
		PanelBorder:   lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(p.Border),
		Item:          lipgloss.NewStyle().Foreground(p.Fg),
		ItemActive:    lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		DetailTitle:   lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		DetailLabel:   lipgloss.NewStyle().Foreground(p.FgDim),
		DetailValue:   lipgloss.NewStyle().Foreground(p.Fg),
		DetailDesc:    lipgloss.NewStyle().Foreground(p.FgDim),
		PreviewBorder:    lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(p.Border),
		PreviewCell:      lipgloss.NewStyle().Foreground(p.Fg),
		LoadButton:       lipgloss.NewStyle().Background(p.Primary).Foreground(p.Bg).Bold(true).Padding(0, 2),
		FilterChip:       lipgloss.NewStyle().Foreground(p.FgDim).Padding(0, 1),
		FilterChipActive: lipgloss.NewStyle().Background(p.Accent).Foreground(p.Bg).Bold(true).Padding(0, 1),
	}

	// History styles
	t.History = HistoryStyles{
		StatBox:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(p.Border).Padding(0, 2),
		StatLabel:   lipgloss.NewStyle().Foreground(p.FgDim),
		StatValue:   lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		TableHeader: lipgloss.NewStyle().Foreground(p.FgDim).Bold(true),
		Row:         lipgloss.NewStyle().Foreground(p.Fg),
		RowActive:   lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		RowAccent:   lipgloss.NewStyle().Foreground(p.Primary),
		BadgeWin:    lipgloss.NewStyle().Foreground(p.Success).Bold(true),
		BadgeLoss:   lipgloss.NewStyle().Foreground(p.FgDim),
		Pagination:  lipgloss.NewStyle().Foreground(p.FgDim),
	}

	// Victory styles
	t.Victory = VictoryStyles{
		ModalBorder:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(p.Primary),
		Badge:           lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		Title:           lipgloss.NewStyle().Foreground(p.Primary).Bold(true),
		Subtitle:        lipgloss.NewStyle().Foreground(p.FgDim),
		StatBox:         lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(p.Border).Padding(0, 1),
		StatLabel:       lipgloss.NewStyle().Foreground(p.FgDim),
		StatValue:       lipgloss.NewStyle().Foreground(p.Fg).Bold(true),
		Button:          lipgloss.NewStyle().Background(p.Primary).Foreground(p.Bg).Bold(true).Padding(0, 3),
		ButtonSecondary: lipgloss.NewStyle().Foreground(p.FgDim),
		Particle:        []lipgloss.Color{p.Primary, p.Accent, p.Success, p.Fg},
	}

	// Toast styles
	t.Toast = ToastStyles{
		Border: lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(p.Accent),
		Icon:   lipgloss.NewStyle().Foreground(p.Primary),
		Text:   lipgloss.NewStyle().Foreground(p.Fg),
	}

	return t
}
