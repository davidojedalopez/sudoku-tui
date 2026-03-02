package theme

import "github.com/charmbracelet/lipgloss"

// Palette holds the raw color values for a theme.
type Palette struct {
	Bg          lipgloss.Color
	BgAlt       lipgloss.Color
	Fg          lipgloss.Color
	FgDim       lipgloss.Color
	FgMuted     lipgloss.Color
	Primary     lipgloss.Color
	Accent      lipgloss.Color
	Success     lipgloss.Color
	Error       lipgloss.Color
	Border      lipgloss.Color
	BorderHeavy lipgloss.Color
}

// ThemeStrings holds per-theme text strings.
type ThemeStrings struct {
	AppName         string
	MenuSubtitle    string
	DiffTagline     string
	VictoryBadge    string
	VictoryTitle    string
	VictorySubtitle string
}

// BoardStyles holds styles for the board.
type BoardStyles struct {
	Border     lipgloss.Style
	BoxBorder  lipgloss.Style
	CellBorder lipgloss.Style
}

// CellStyles holds styles for individual cells.
type CellStyles struct {
	Given      lipgloss.Style
	User       lipgloss.Style
	Empty      lipgloss.Style
	Conflict   lipgloss.Style
	Highlight  lipgloss.Style
	Cursor     lipgloss.Style
	Notes      lipgloss.Style
	NotesDigit lipgloss.Style
}

// HeaderStyles holds styles for header/footer bars.
type HeaderStyles struct {
	Bar   lipgloss.Style
	Title lipgloss.Style
	Meta  lipgloss.Style
}

// FooterStyles holds footer styles.
type FooterStyles struct {
	Bar      lipgloss.Style
	KeyHint  lipgloss.Style
	KeyLabel lipgloss.Style
}

// MenuStyles holds main menu styles.
type MenuStyles struct {
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Item       lipgloss.Style
	ItemActive lipgloss.Style
	ItemPrefix lipgloss.Style
}

// DifficultyStyles holds difficulty selector styles.
type DifficultyStyles struct {
	ModalBorder lipgloss.Style
	Title       lipgloss.Style
	Tagline     lipgloss.Style
	Option      lipgloss.Style
	Active      lipgloss.Style
	ClueCount   lipgloss.Style
	Desc        lipgloss.Style
	DescIcon    lipgloss.Style
}

// SidebarStyles holds game sidebar styles.
type SidebarStyles struct {
	Border                lipgloss.Style
	Title                 lipgloss.Style
	Timer                 lipgloss.Style
	ModeIndicatorActive   lipgloss.Style
	ModeIndicatorInactive lipgloss.Style
	RemainingDigit        lipgloss.Style
	RemainingBar          lipgloss.Style
	RemainingBarEmpty     lipgloss.Style
	RemainingCount        lipgloss.Style
	RemainingDone         lipgloss.Style
	RemainingActive       lipgloss.Style
}

// BadgeStyles holds difficulty badge styles.
type BadgeStyles struct {
	Easy   lipgloss.Style
	Medium lipgloss.Style
	Hard   lipgloss.Style
	Expert lipgloss.Style
}

// LibraryStyles holds puzzle library styles.
type LibraryStyles struct {
	PanelBorder     lipgloss.Style
	Item            lipgloss.Style
	ItemActive      lipgloss.Style
	DetailTitle     lipgloss.Style
	DetailLabel     lipgloss.Style
	DetailValue     lipgloss.Style
	DetailDesc      lipgloss.Style
	PreviewBorder   lipgloss.Style
	PreviewCell     lipgloss.Style
	LoadButton      lipgloss.Style
	FilterChip      lipgloss.Style
	FilterChipActive lipgloss.Style
}

// HistoryStyles holds history screen styles.
type HistoryStyles struct {
	StatBox     lipgloss.Style
	StatLabel   lipgloss.Style
	StatValue   lipgloss.Style
	TableHeader lipgloss.Style
	Row         lipgloss.Style
	RowActive   lipgloss.Style
	RowAccent   lipgloss.Style
	BadgeWin    lipgloss.Style
	BadgeLoss   lipgloss.Style
	Pagination  lipgloss.Style
}

// VictoryStyles holds victory/celebration screen styles.
type VictoryStyles struct {
	ModalBorder     lipgloss.Style
	Badge           lipgloss.Style
	Title           lipgloss.Style
	Subtitle        lipgloss.Style
	StatBox         lipgloss.Style
	StatLabel       lipgloss.Style
	StatValue       lipgloss.Style
	Button          lipgloss.Style
	ButtonSecondary lipgloss.Style
	Particle        []lipgloss.Color
}

// ToastStyles holds technique toast styles.
type ToastStyles struct {
	Border lipgloss.Style
	Icon   lipgloss.Style
	Text   lipgloss.Style
}

// Theme is the complete theme definition.
type Theme struct {
	Name    string
	Palette Palette
	Board   BoardStyles
	Cell    CellStyles
	Header  HeaderStyles
	Footer  FooterStyles
	Menu    MenuStyles
	Diff    DifficultyStyles
	Sidebar SidebarStyles
	Badges  BadgeStyles
	Library LibraryStyles
	History HistoryStyles
	Victory VictoryStyles
	Toast   ToastStyles
	Strings ThemeStrings
}

// Registry holds all available themes.
var Registry = map[string]*Theme{}

// Default is the active theme.
var Default *Theme

func init() {
	Registry["modern-charm"] = ModernCharmTheme()
	Registry["zen-monolith"] = ZenMonolithTheme()
	Registry["retro-phosphor"] = RetroPhosphorTheme()
	Registry["matrix"] = MatrixTheme()
	Default = Registry["modern-charm"]
}
