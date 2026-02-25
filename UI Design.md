# Sudoku TUI — UI Design & Theme Specification

> This document is the single source of truth for screen layouts, component specifications, and all four themes. An implementer should be able to build the complete UI solely from this document without referring to the inspiration mockups.

---

## 1. Design Principles

### Philosophy
- **Keyboard-first.** Every action is reachable without a mouse. Keyboard hints are always visible.
- **Dark-mode only.** All four themes use dark backgrounds. No light-mode variant.
- **Monospace-native.** Layouts are specified in fixed-width characters. Every element aligns on a monospace grid.
- **Structural consistency across themes.** Layouts are identical across themes — only colors, border characters, and text casing change.
- **Information density without noise.** Show only what the player needs at that moment. No score, no mistake counter, no redundant labels.

### Universal Layout Rules
1. Every screen has three zones: **Header** (1 line), **Body** (fills remaining height), **Footer** (1–2 lines).
2. Minimum terminal size: **80 × 24** characters. Ideal: **120 × 36**.
3. The game board screen requires at minimum **100 × 30**.
4. When the terminal is too small, show a single centered message: `Terminal too small. Please resize to at least 80×24.`
5. All borders use box-drawing characters (not ASCII `+`, `-`, `|`). Heavy/double borders mark structural divisions; single-line borders mark content cells.

### Typography Conventions per Theme
| Theme | Casing | Number style |
|-------|--------|-------------|
| Modern Charm | Mixed case (Title Case labels, sentence case descriptions) | Proportional feel, normal weight |
| Zen Monolith | ALL CAPS for labels, mixed for descriptions | Bold given cells, normal user cells |
| Retro Phosphor | ALL CAPS everywhere | Block cursor, inverted highlight |
| Matrix | ALL CAPS labels, `//` prefixes on system messages | Glitchy double-print effect in flavor text |

---

## 2. Screen Layouts

### 2.1 Main Menu

**Zones:** Full-height centered layout. Logo top-third, menu mid-third, footer bottom.

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│                                                                              │
│              ███████╗██╗   ██╗██████╗  ██████╗ ██╗  ██╗██╗   ██╗           │
│              ██╔════╝██║   ██║██╔══██╗██╔═══██╗██║ ██╔╝██║   ██║           │
│              ███████╗██║   ██║██║  ██║██║   ██║█████╔╝ ██║   ██║           │
│              ╚════██║██║   ██║██║  ██║██║   ██║██╔═██╗ ██║   ██║           │
│              ███████║╚██████╔╝██████╔╝╚██████╔╝██║  ██╗╚██████╔╝           │
│              ╚══════╝ ╚═════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝ ╚═════╝            │
│                                                                              │
│                           [theme subtitle line]                              │
│                                                                              │
│                                                                              │
│                         ▶  New Game                                          │
│                            Puzzle Library                                    │
│                            History                                           │
│                            Quit                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│  v1.0.0          [k] / [j] Navigate     [Enter] Select     [q] Quit         │
└──────────────────────────────────────────────────────────────────────────────┘
```

**Style token map:**
- Logo block → `MenuTitle`
- Subtitle line → `MenuSubtitle`
- Active menu item → `MenuItemActive` (prefix `▶` + highlight)
- Inactive menu items → `MenuItem`
- Footer bar → `FooterBar`
- Key hints → `KeyHint` (key chip) + `KeyLabel` (action text)

**Per-theme subtitle text:**
| Theme | Subtitle |
|-------|----------|
| Modern Charm | `Modern Charm Edition` |
| Zen Monolith | `GEOMETRIC LOGIC ENGINE` |
| Retro Phosphor | `// SYSTEM READY // WAITING FOR INPUT_` |
| Matrix | `WAKE UP, NEO... THE GRID AWAITS` |

---

### 2.2 Difficulty Selector

Rendered as a **modal overlay** centered in the terminal (not a full-screen replacement). The main menu remains visible but dimmed behind it.

```
                    ┌─────────────────────────────────────┐
                    │                                     │
                    │       [ SELECT DIFFICULTY ]         │  ← DifficultyTitle
                    │          [theme tagline]            │  ← DifficultyTagline
                    │  ───────────────────────────────    │
                    │                                     │
                    │    EASY                             │  ← DifficultyOption
                    │    36–40 clues                      │  ← DifficultyClueCount
                    │                                     │
                    │  ▶ MEDIUM                           │  ← DifficultyActive
                    │    30–35 clues                      │
                    │                                     │
                    │    HARD                             │
                    │    25–29 clues                      │
                    │                                     │
                    │    EXPERT                           │
                    │    17–24 clues                      │
                    │                                     │
                    │  ───────────────────────────────    │
                    │  ⓘ [description of selected level]  │  ← DifficultyDesc
                    │                                     │
                    │  [k] UP  [j] DOWN  [Esc] BACK  [Enter] SELECT │
                    └─────────────────────────────────────┘
```

**Difficulty descriptions:**
| Level | Description |
|-------|-------------|
| Easy | Solvable with naked and hidden singles only. A warm-up. |
| Medium | Requires naked pairs and pointing pairs. Good daily practice. |
| Hard | Advanced patterns required. Focused solving needed. |
| Expert | X-Wing, Y-Wing, and Swordfish. Not for the faint-hearted. |

**Per-theme tagline:**
| Theme | Tagline |
|-------|---------|
| Modern Charm | `Choose your challenge` |
| Zen Monolith | `INITIATE SEQUENCE` |
| Retro Phosphor | `INITIALIZE MATRIX` |
| Matrix | `LOAD PROGRAM` |

---

### 2.3 Game Board

The most complex screen. Two-panel layout: board on the left, sidebar on the right.

```
┌─ SUDOKU ─────────────────────────────── HARD ────────────────────────────────┐
│                                               │  ELAPSED TIME                │
│  ╔═════╤═════╤═════╦═════╤═════╤═════╦═════╤═════╤═════╗  │  00:04:23            │
│  ║ 1 2 │     │     ║  7  │     │     ║     │     │     ║  │                      │
│  ║   5 │  4  │     ║     │     │     ║     │  5  │     ║  │  INPUT MODE          │
│  ║ 7 8 │     │     ║     │     │     ║     │     │     ║  │  [●DIGIT] [ NOTE ]   │
│  ╟─────┼─────┼─────╫─────┼─────┼─────╫─────┼─────┼─────╢  │                      │
│  ║     │     │  8  ║     │     │  3  ║     │     │     ║  │  REMAINING           │
│  ║  6  │     │     ║     │  5  │     ║  2  │     │     ║  │  1 ████░░░░  (2)     │
│  ║     │     │     ║     │     │     ║     │  1  │     ║  │  2 █████████ DONE    │
│  ╟─────┼─────┼─────╫─────┼─────┼─────╫─────┼─────┼─────╢  │  3 ████░░░░  (4)    │
│  ║     │  9  │     ║     │     │     ║     │     │  4  ║  │  4 ██░░░░░░░  (7)    │
│  ║     │     │     ║  8  │     │  3  ║     │  1  │     ║  │  5 ████░░░░░  (3)    │
│  ║     │     │     ║     │     │     ║     │     │     ║  │  6 █░░░░░░░░  (8)    │
│  ╠═════╪═════╪═════╬═════╪═════╪═════╬═════╪═════╪═════╣  │  7 ██░░░░░░░  (7)    │
│  ║     │     │     ║     │  6  │     ║     │     │  3  ║  │  8 █░░░░░░░░  (8)    │
│  ║  7  │     │  2  ║     │     │     ║  6  │     │     ║  │  9 ████░░░░░  (4)    │
│  ║     │     │     ║     │     │     ║     │     │     ║  │                      │
│  ╟─────┼─────┼─────╫─────┼─────┼─────╫─────┼─────┼─────╢  │                      │
│  ║     │  6  │     ║     │     │     ║  2  │  8  │     ║  │                      │
│  ║     │     │     ║  4  │  1  │  9  ║     │     │  5  ║  │                      │
│  ║     │     │     ║     │     │     ║     │     │     ║  │                      │
│  ╠═════╪═════╪═════╬═════╪═════╪═════╬═════╪═════╪═════╣  │                      │
│  ║     │     │     ║     │  8  │     ║  7  │  9  │     ║  │                      │
│  ║     │     │     ║     │     │     ║     │     │     ║  │                      │
│  ║     │     │     ║     │     │     ║     │     │     ║  │                      │
│  ╚═════╧═════╧═════╩═════╧═════╧═════╩═════╧═════╧═════╝  │                      │
├──────────────────────────────────────────────────────────────────────────────────┤
│  [h/j/k/l] Move  [1-9] Digit  [Shift+1-9] Note  [n] Mode  [x] Erase  [u] Undo  [Esc] Menu │
└──────────────────────────────────────────────────────────────────────────────────┘
```

#### Cell Rendering Detail

Each cell is **5 characters wide × 3 lines tall**. The two states:

**Notes mode (pencil marks):**
```
 1 2 3
 4   6     ← note 5 absent means 5 eliminated
 7 8 9
```

**Digit mode (filled cell):**
```

  7

```

**Cell states and their token:**
| State | Description | Token |
|-------|-------------|-------|
| Given | Puzzle clue, immutable | `CellGiven` |
| User digit | Player-placed number | `CellUser` |
| Empty | No value, no notes | `CellEmpty` |
| Conflict | User digit conflicts with row/col/box | `CellConflict` |
| Same-number highlight | Same digit as cursor cell | `CellHighlight` |
| Cursor | Current cursor position | `CellCursor` |
| Notes cell | Has pencil marks but no digit | `CellNotes` (note digits use `CellNotesDigit`) |

#### Grid Border Characters

```
Heavy (3×3 box boundaries):  ╔ ╦ ╗ ╠ ╬ ╣ ╚ ╩ ╝ ═ ║
Light (cell boundaries):      ┼ ┤ ├ ┬ ┴ ─ │ ╟ ╢ ╫ ╤ ╧
```

For **Retro Phosphor** and **Matrix** themes, use double-line ASCII borders for the 3×3 boxes and single-line for cells (same characters, different colors).

#### Sidebar Layout

```
  ELAPSED TIME               ← SidebarTitle
  00:04:23                   ← Timer (large)

  INPUT MODE                 ← SidebarTitle
  [● DIGIT] [ NOTE ]         ← ModeIndicatorActive / ModeIndicatorInactive

  REMAINING                  ← SidebarTitle
  1 ████░░░░  (2)            ← RemainingDigit + progress bar + count
  2 █████████ DONE           ← RemainingDone (grayed out when 9 placed)
  ...
```

The mode indicator shows two segmented buttons side by side. The active mode is filled/highlighted; the inactive one is outlined/dimmed.

#### Technique Toast

When a technique is detected, a **toast** appears above the footer bar for ~3 seconds:

```
  ╭─────────────────────────────────────╮
  │  ✦ X-Wing detected! Nice move.      │   ← ToastTechnique
  ╰─────────────────────────────────────╯
```

Position: bottom-center, above the keyboard legend.

---

### 2.4 Puzzle Library

Two-panel layout: scrollable list on left (60%), detail panel on right (40%).

```
┌─ PUZZLE LIBRARY ─────────────────────────────────────────────────────────────┐
│                                                                              │
│  ┌──────────────────────────────────────┐  ┌───────────────────────────┐   │
│  │  [/] Filter...                       │  │  monolith_001             │   │
│  ├──────────────────────────────────────┤  │  ─────────────────────    │   │
│  │ ▶ monolith_001          [HARD]       │  │  Difficulty:  HARD        │   │
│  │   monolith_002          [MEDIUM]     │  │  Author:      curator     │   │
│  │   classic_symmetry_01   [EASY]       │  │                           │   │
│  │   xwing_practice        [HARD]       │  │  "A classic symmetrical   │   │
│  │   ywing_showcase        [EXPERT]     │  │  puzzle requiring X-Wing  │   │
│  │   beginner_001          [EASY]       │  │  technique to crack."     │   │
│  │   pointing_pairs_demo   [MEDIUM]     │  │                           │   │
│  │   swordfish_master      [EXPERT]     │  │  ┌─────────────────────┐  │   │
│  │                                      │  │  │ · · 5 │ · · · │ · · │  │   │
│  │  (j/k to scroll)                    │  │  │ · · · │ 3 · · │ 7 · │  │   │
│  └──────────────────────────────────────┘  │  │ 2 · · │ · · · │ · · │  │   │
│                                             │  └─────────────────────┘  │   │
│                                             │                           │   │
│                                             │  [Enter] LOAD PUZZLE      │   │
│                                             └───────────────────────────┘   │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│  [j/k] Scroll  [/] Filter  [Enter] Load  [Esc] Back                        │
└──────────────────────────────────────────────────────────────────────────────┘
```

**Difficulty badge colors** (consistent across all themes via semantic tokens):
| Difficulty | Token |
|------------|-------|
| Easy | `BadgeEasy` |
| Medium | `BadgeMedium` |
| Hard | `BadgeHard` |
| Expert | `BadgeExpert` |

---

### 2.5 History

Stats summary row + paginated table.

```
┌─ HISTORY ────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │  TOTAL GAMES │  │   WIN RATE   │  │  BEST (HARD) │  │    STREAK    │   │
│  │     42       │  │    85%       │  │    14:02     │  │      5       │   │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘   │
│                                                                              │
│  DATE              DIFFICULTY    TIME      RESULT                           │
│  ──────────────────────────────────────────────────────                     │
│ ▶ 2023-10-27 14:32  HARD          14:02     WIN                             │
│   2023-10-26 09:15  MEDIUM        08:45     WIN                             │
│   2023-10-25 18:20  EXPERT        --:--     GAVE UP                         │
│   2023-10-24 22:04  HARD          16:10     WIN                             │
│   2023-10-24 11:15  EASY          04:22     WIN                             │
│   2023-10-23 08:38  HARD          15:45     WIN                             │
│   2023-10-22 19:45  EXPERT        22:15     WIN                             │
│   2023-10-21 12:19  MEDIUM        09:12     WIN                             │
│                                                                              │
│  SHOWING 1-8 OF 42   PAGE 1/6                                               │
│                                                                              │
├──────────────────────────────────────────────────────────────────────────────┤
│  [j/k] Scroll  [PgDn/PgUp] Page  [Esc] Back                                │
└──────────────────────────────────────────────────────────────────────────────┘
```

**Result colors** (semantic tokens):
- WIN / SOLVED → `HistoryBadgeWin`
- GAVE UP → `HistoryBadgeLoss`
- Active row → `HistoryRowActive` (left accent bar)

---

### 2.6 Victory / Celebration

A **centered modal** overlaid on the completed (frozen) board. The fireworks/confetti animation renders behind the modal but in front of the board.

```
                    ┌─────────────────────────────────────┐
                    │                                     │
                    │            [theme badge]            │  ← VictoryBadge
                    │                                     │
                    │         PUZZLE SOLVED!              │  ← VictoryTitle
                    │   [theme flavor subtitle]           │  ← VictorySubtitle
                    │                                     │
                    │  ┌─────────────────────────────┐   │
                    │  │  ⏱ TIME        14:02        │   │  ← VictoryStat
                    │  ├─────────────────────────────┤   │
                    │  │  ◈ DIFFICULTY  HARD         │   │
                    │  └─────────────────────────────┘   │
                    │                                     │
                    │  ┌─────────────────────────────┐   │
                    │  │     [ PRESS ENTER ]         │   │  ← VictoryButton
                    │  └─────────────────────────────┘   │
                    │       ← Return to Menu (Esc)        │
                    │                                     │
                    └─────────────────────────────────────┘
```

**Fireworks / Confetti particles** (ASCII, colored with theme palette):
- Characters: `* + · ° ✦ ✧ ⁂ ∘` (use subset that renders well in common terminals)
- Colors: cycle through theme's primary, accent, success, and a bright variant
- Animation: particles emit from a few random points, travel upward with slight horizontal drift, fade out over ~60 frames at 30fps
- Duration: ~3 seconds, then modal becomes the only focus

**Per-theme badge & flavor text:**
| Theme | Badge | Subtitle |
|-------|-------|----------|
| Modern Charm | `🏆` trophy icon | `Excellent work! You've mastered this grid.` |
| Zen Monolith | `◆ S RANK ◆` diamond | `ZEN MONOLITH SYSTEM · PUZZLE CLEARED` |
| Retro Phosphor | `SYSTEM RESTORED` ASCII art | `SEQUENCE COMPLETED SUCCESSFULLY` |
| Matrix | `ENCRYPTION BROKEN` badge | `System anomaly detected. You are the one.` |

---

## 3. Component Token Inventory

This is the complete list of all named style tokens. Every theme must define a value for every token.

### Board & Cells
| Token | Description |
|-------|-------------|
| `BoardBorder` | Outer board border (heavy box-drawing) |
| `BoardBoxBorder` | 3×3 box dividers (heavy) |
| `BoardCellBorder` | Individual cell dividers (light) |
| `CellGiven` | Immutable clue digit — foreground color |
| `CellUser` | Player-placed digit — foreground color |
| `CellEmpty` | Empty cell background |
| `CellConflict` | Conflicting cell — background highlight + bold |
| `CellHighlight` | Same-number highlight — subtle background tint |
| `CellCursor` | Cursor cell — strong background fill |
| `CellNotes` | Note cell background |
| `CellNotesDigit` | Pencil mark digit foreground |

### Header & Footer
| Token | Description |
|-------|-------------|
| `HeaderBar` | Top bar background + foreground |
| `HeaderTitle` | App name / screen title in header |
| `HeaderMeta` | Right-side header info (difficulty, version) |
| `FooterBar` | Bottom bar background + foreground |
| `KeyHint` | Key chip (e.g. `[h]`) — bordered/inverted style |
| `KeyLabel` | Action label next to key chip |

### Main Menu
| Token | Description |
|-------|-------------|
| `MenuTitle` | ASCII logo block foreground |
| `MenuSubtitle` | Tagline below logo |
| `MenuItem` | Inactive menu item |
| `MenuItemActive` | Active menu item (highlighted row + prefix) |
| `MenuItemPrefix` | `▶` chevron on active item |

### Difficulty Selector
| Token | Description |
|-------|-------------|
| `DifficultyModalBorder` | Modal frame border |
| `DifficultyTitle` | `SELECT DIFFICULTY` heading |
| `DifficultyTagline` | Sub-heading below title |
| `DifficultyOption` | Inactive difficulty option |
| `DifficultyActive` | Active/highlighted difficulty row |
| `DifficultyClueCount` | Clue count description (smaller text) |
| `DifficultyDesc` | Dynamic description block at bottom |
| `DifficultyDescIcon` | `ⓘ` info icon |

### Sidebar (Game Screen)
| Token | Description |
|-------|-------------|
| `SidebarBorder` | Vertical divider between board and sidebar |
| `SidebarTitle` | Section heading labels (ELAPSED TIME, etc.) |
| `Timer` | Large timer display |
| `ModeIndicatorActive` | Active mode button (filled) |
| `ModeIndicatorInactive` | Inactive mode button (outline) |
| `RemainingDigit` | Digit label in remaining tracker |
| `RemainingBar` | Filled portion of progress bar |
| `RemainingBarEmpty` | Empty portion of progress bar |
| `RemainingCount` | `(n)` count text |
| `RemainingDone` | Style when digit is fully placed (dim/grayed) |
| `RemainingActive` | Highlighted row when cursor is on that digit |

### Difficulty Badges (Library + History)
| Token | Description |
|-------|-------------|
| `BadgeEasy` | Easy badge fg + bg |
| `BadgeMedium` | Medium badge fg + bg |
| `BadgeHard` | Hard badge fg + bg |
| `BadgeExpert` | Expert badge fg + bg |

### Puzzle Library
| Token | Description |
|-------|-------------|
| `LibraryPanelBorder` | Border around list/detail panels |
| `LibraryItem` | Inactive list item |
| `LibraryItemActive` | Active/selected list item |
| `LibraryDetailTitle` | Puzzle name in detail panel |
| `LibraryDetailLabel` | Field labels (Difficulty:, Author:) |
| `LibraryDetailValue` | Field values |
| `LibraryDetailDesc` | Puzzle description text |
| `LibraryPreviewBorder` | Mini-grid preview border |
| `LibraryPreviewCell` | Mini-grid cell values |
| `LibraryLoadButton` | Load puzzle button |

### History
| Token | Description |
|-------|-------------|
| `HistoryStatBox` | Stats summary box borders |
| `HistoryStatLabel` | Stat box label (TOTAL GAMES, etc.) |
| `HistoryStatValue` | Stat box value (large number) |
| `HistoryTableHeader` | Column header row |
| `HistoryRow` | Normal table row |
| `HistoryRowActive` | Highlighted/selected row |
| `HistoryRowAccent` | Left accent bar on active row |
| `HistoryBadgeWin` | WIN result badge |
| `HistoryBadgeLoss` | GAVE UP result badge |
| `HistoryPagination` | Pagination text |

### Victory Screen
| Token | Description |
|-------|-------------|
| `VictoryModalBorder` | Modal frame border |
| `VictoryBadge` | Theme-specific badge/rank display |
| `VictoryTitle` | `PUZZLE SOLVED!` heading |
| `VictorySubtitle` | Theme flavor subtitle |
| `VictoryStatBox` | Stat box border |
| `VictoryStatLabel` | Stat label (TIME, DIFFICULTY) |
| `VictoryStatValue` | Stat value |
| `VictoryButton` | Primary action button |
| `VictoryButtonSecondary` | Secondary action (return to menu) |
| `CelebrationParticle` | Fireworks particle characters (array of colors) |

### Toast
| Token | Description |
|-------|-------------|
| `ToastBorder` | Toast border |
| `ToastIcon` | Leading `✦` icon |
| `ToastText` | Toast message text |

---

## 4. Theme Architecture

### Theme Struct (Go)

```go
type Theme struct {
    Name string

    // Raw palette (used to derive all tokens)
    Palette Palette

    // All style tokens
    Board    BoardStyles
    Cell     CellStyles
    Header   HeaderStyles
    Footer   FooterStyles
    Menu     MenuStyles
    Diff     DifficultyStyles
    Sidebar  SidebarStyles
    Badges   BadgeStyles
    Library  LibraryStyles
    History  HistoryStyles
    Victory  VictoryStyles
    Toast    ToastStyles

    // Theme-specific text strings
    Strings ThemeStrings
}

type Palette struct {
    Bg       lipgloss.Color
    BgAlt    lipgloss.Color  // slightly lighter bg for elevated panels
    Fg       lipgloss.Color
    FgDim    lipgloss.Color  // muted/secondary text
    FgMuted  lipgloss.Color  // very dim text (done digits, hints)
    Primary  lipgloss.Color  // main accent (cursor, active items)
    Accent   lipgloss.Color  // secondary accent
    Success  lipgloss.Color  // correct, win, done
    Error    lipgloss.Color  // conflict, gave up
    Border   lipgloss.Color  // normal border color
    BorderHeavy lipgloss.Color // thick border color (box dividers)
}

type ThemeStrings struct {
    AppName        string
    MenuSubtitle   string
    DiffTagline    string
    VictoryBadge   string
    VictoryTitle   string
    VictorySubtitle string
    // ... etc
}
```

### Registry & JSON Override

```go
// internal/theme/theme.go
var Registry = map[string]*Theme{
    "modern-charm":   ModernCharmTheme,
    "zen-monolith":   ZenMonolithTheme,
    "retro-phosphor": RetroPhosphorTheme,
    "matrix":         MatrixTheme,
}

var Default = Registry["modern-charm"]
```

JSON override file at `~/.config/sudoku-tui/theme.json`:

```json
{
  "base": "modern-charm",
  "overrides": {
    "palette": {
      "Primary": "#FF6B6B"
    }
  }
}
```

Overrides are merged on top of the base theme at startup. Only `palette` overrides are supported in v1 (all style tokens re-derive from the updated palette).

---

## 5. Theme Definitions

### 5.1 Modern Charm (Default)

**Personality:** Warm, polished, approachable. Catppuccin Mocha base. Feels like a well-designed developer tool.

**Palette:**
```
Bg          #1E1E2E   deep violet-black
BgAlt       #27273A   elevated panel bg
Fg          #CAD3F5   lavender white
FgDim       #8087A2   muted lavender
FgMuted     #51526A   very dim (grayed-out items)
Primary     #F5A97F   peach (cursor, active items, CTAs)
Accent      #8AADF4   sky blue (timer, badges, highlights)
Success     #A6DA95   green
Error       #ED8796   pink-red (conflicts, gave up)
Border      #363654   subtle border
BorderHeavy #494970   box divider border
```

**Text casing:** Mixed case. Normal prose for descriptions.

**Border style:** Single `─ │ ┼` for cells, heavy `═ ║ ╬` for 3×3 boxes.

**Key style tokens:**
```
CellGiven:       Bold, Fg (#CAD3F5)
CellUser:        Accent (#8AADF4)
CellCursor:      Bg=Primary (#F5A97F), Fg=Bg (#1E1E2E), Bold
CellConflict:    Bg=Error (#ED8796), Fg=Bg (#1E1E2E), Bold
CellHighlight:   Bg=BgAlt (#27273A), subtle tint
MenuItemActive:  Fg=Primary (#F5A97F), Bold, chevron prefix
Timer:           Accent (#8AADF4), Bold, large
ModeActive:      Bg=Primary (#F5A97F), Fg=Bg
ToastText:       Fg=Primary, border=Accent
BadgeEasy:       Bg=#2A3829, Fg=#A6DA95
BadgeMedium:     Bg=#2A2D47, Fg=#8AADF4
BadgeHard:       Bg=#3D2E29, Fg=#F5A97F
BadgeExpert:     Bg=#3D2929, Fg=#ED8796
```

---

### 5.2 Zen Monolith

**Personality:** Cold, precise, minimal. Near-black with electric blue. Like a command-line OS from the future.

**Palette:**
```
Bg          #0A0C14   near-black with blue tint
BgAlt       #121726   slightly lighter
Fg          #D0D8F0   cool white
FgDim       #6D7A9E   steel gray
FgMuted     #3D4560   very dim
Primary     #2B7FFF   electric blue (cursor, active, CTAs)
Accent      #5BA4FB   lighter blue (highlights, timer)
Success     #4ECBA0   teal-green
Error       #F04F5F   cold red
Border      #1E2640   very dark border
BorderHeavy #2B3A5C   medium border
```

**Text casing:** ALL CAPS for all labels, headings, and menu items.

**Border style:** Single `─ │ ┼` for cells, heavy `═ ║ ╬` for 3×3 boxes. Borders in `BorderHeavy` color.

**Key style tokens:**
```
CellGiven:       Bold, Fg (#D0D8F0)
CellUser:        Accent (#5BA4FB)
CellCursor:      Bg=Primary (#2B7FFF), Fg=Bg (#0A0C14), Bold
CellConflict:    Bg=Error (#F04F5F), Fg=Bg, Bold
CellHighlight:   Bg=BgAlt, left border in Primary
MenuTitle:       Primary (#2B7FFF), Bold — horizontal scan-line feel
MenuItemActive:  Primary, Underline, `▶ ` prefix
Timer:           Accent, large monospace, `00:00:00` format
ModeActive:      Bg=Primary, Fg=Bg, Uppercase
BadgeEasy:       Bg=dark, Fg=Success
BadgeHard:       Bg=dark, Fg=Primary
BadgeExpert:     Bg=dark, Fg=Error
VictoryBadge:    Diamond `◆ S RANK ◆` in Primary, Bordered
```

---

### 5.3 Retro Phosphor

**Personality:** Amber phosphor CRT terminal. Monochromatic. Everything in 4 shades of amber.

**Palette:**
```
Bg          #0A0600   near-black with amber tint
BgAlt       #150D00   slightly lighter
Fg          #FFD15C   bright amber (primary text, given digits)
FgDim       #C98A10   mid amber (secondary text)
FgMuted     #7A5000   dim amber (done items, borders)
Primary     #FFB000   amber (cursor, active, CTAs)
Accent      #FFD15C   bright amber (same as Fg — it's monochrome)
Success     #FFD15C   bright amber (wins look the same)
Error       #FF6A00   orange-amber (conflicts only color break)
Border      #4A3000   dark amber border
BorderHeavy #7A5000   mid-dim amber
```

**Text casing:** ALL CAPS everywhere without exception.

**Border style:** Double `═ ║ ╬ ╔ ╗ ╚ ╝` for 3×3 boxes, single `─ │ ┼` for cells. All borders in `FgMuted`.

**Special effects (via lipgloss):**
- Active/cursor cells: **inverted** (Bg=Primary, Fg=Bg)
- Menu active item: **inverted** full-width row
- Block cursor character (`█`) prepended to active menu item

**Key style tokens:**
```
CellGiven:       Bold, Fg=Fg (#FFD15C)
CellUser:        Fg=Primary (#FFB000)
CellCursor:      Bg=Primary, Fg=Bg — inverted block
CellConflict:    Fg=Error (#FF6A00), Bold, Underline
CellHighlight:   Fg=Fg, Bg=BgAlt
MenuTitle:       Fg=Fg, Bold — ASCII block letters
MenuItemActive:  Bg=Primary, Fg=Bg — full-width inversion
Timer:           Fg=Fg, Bold — `[ SYSTEM_TIME ] 00:04:21`
ModeActive:      Bg=Primary, Fg=Bg — `[ INSERT ]` style label
ToastBorder:     FgMuted double-border
VictoryTitle:    ASCII art `SYSTEM RESTORED` in Fg
```

---

### 5.4 Matrix

**Personality:** Green-on-black hacker aesthetic. Digital rain. Every label is a system message.

**Palette:**
```
Bg          #050A05   near-black with green tint
BgAlt       #0A180A   slightly lighter
Fg          #25F447   bright green (primary text)
FgDim       #13872A   mid green (secondary text)
FgMuted     #0A4015   dim green (borders, done items)
Primary     #0BDA0B   vivid green (cursor, active, CTAs)
Accent      #13EC5B   bright cyan-green (highlights, timer)
Success     #25F447   bright green (wins — same as Fg)
Error       #FF4444   red — the one non-green color (conflicts)
Border      #0A3010   very dark green border
BorderHeavy #0F5020   medium green border
```

**Text casing:** ALL CAPS. Labels prefixed with `//` or `>` for terminal feel.

**Border style:** Single `─ │ ┼` for cells, heavy `═ ║ ╬` for 3×3 boxes. All green tones.

**Key style tokens:**
```
CellGiven:       Bold, Fg=Fg (#25F447)
CellUser:        Fg=Accent (#13EC5B)
CellCursor:      Bg=Primary (#0BDA0B), Fg=Bg — bright block
CellConflict:    Bg=Error (#FF4444), Fg=Bg — red is intentional contrast
CellHighlight:   Bg=BgAlt, subtle green tint
MenuTitle:       Fg=Primary — glitch double-shadow in Victory
MenuItemActive:  Bg=Primary, Fg=Bg, Bold — full-width button
Timer:           Accent, Bold — `SYSTEM TIME 00:00:45`
ModeActive:      Bg=Primary, Fg=Bg — `> NOTES_MODE`
ToastBorder:     Fg=FgDim, single
ToastText:       Fg=Accent
VictoryTitle:    Fg=Primary, Very Bold — `YOU ARE THE ONE`
```

---

## 6. Keyboard Legend Specification

The footer keyboard legend is rendered as: `[key] ACTION` pairs separated by two spaces. Keys are styled with `KeyHint` (a bordered/inverted chip); action text uses `KeyLabel`.

### Per-Screen Legend

**Main Menu:**
```
[k/j] Navigate    [Enter] Select    [q] Quit
```

**Difficulty Selector:**
```
[k] Up    [j] Down    [Enter] Select    [Esc] Back
```

**Game Board:**
```
[h/j/k/l] Move    [1-9] Digit    [Shift+1-9] Note    [n] Mode    [x] Erase    [u] Undo    [Esc] Menu
```

**Puzzle Library:**
```
[j/k] Scroll    [/] Filter    [Enter] Load    [Esc] Back
```

**History:**
```
[j/k] Scroll    [PgDn/PgUp] Page    [Esc] Back
```

**Victory:**
```
[Enter] New Game    [Esc] Menu
```

---

## 7. Implementation Notes

### Lip Gloss Integration

- All styles are `lipgloss.Style` values stored in the theme struct.
- The board renderer in `internal/screen/game/view.go` reads from the active theme — never hardcodes colors.
- Theme is injected into every screen model at construction time.
- `lipgloss.NewStyle().Background(...).Foreground(...).Bold(...)` chains build each token.

### Theme Loading Order (`internal/theme/config.go`)
1. Load base theme from registry by name (default: `"modern-charm"`)
2. Check `~/.config/sudoku-tui/theme.json`; if present, unmarshal overrides
3. Apply palette overrides → re-derive all style tokens from updated palette
4. Return final theme to `app.go`

### Minimum Terminal Size Check

In `app.go`'s `View()`, check `lipgloss.Width` / `lipgloss.Height` of the window. If below 80×24, render only:

```
┌──────────────────────────────────────┐
│                                      │
│   Terminal too small.                │
│   Please resize to at least 80×24.  │
│                                      │
└──────────────────────────────────────┘
```

### ASCII Logo Generation

The `SUDOKU` block-letter logo for the main menu is a static multi-line string stored in `internal/screen/menu/logo.go`. Each theme's `Strings.AppName` overrides just the text name in the header bar (not the logo itself). The Retro Phosphor theme uses a different compact block-letter style (smaller, chunky).

---

## 8. Verification Checklist

Before considering the UI implementation complete:

- [ ] All 6 screens render correctly at 80×24 (minimum) and 120×36 (ideal)
- [ ] All cell states render visually distinct for all 4 themes
- [ ] Theme can be switched by changing `theme.json` and restarting
- [ ] Keyboard legend updates correctly per screen
- [ ] Technique toast appears and disappears after ~3s without breaking layout
- [ ] Fireworks animation renders without corrupting the board behind it
- [ ] Remaining digits bar shows correct counts and grays out at 9 placements
- [ ] Mode indicator switches correctly on `n` keypress
- [ ] Conflict highlighting appears on both conflicting cells simultaneously
- [ ] Same-number highlighting activates on cursor move, clears on move away
- [ ] History table paginates correctly with j/k scrolling
- [ ] Puzzle library filter works case-insensitively
- [ ] Victory modal shows correct time and difficulty label
