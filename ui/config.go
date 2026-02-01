package ui

import "time"

// UI layout and timing configuration. Tweak these values to tune the TUI layout.
var (
	// MinTerminalWidth is the minimal comfortable terminal width (columns)
	MinTerminalWidth = 80

	// Column weight ratios for the main content flex (list : details)
	LeftColumnWeight  = 1
	RightColumnWeight = 2

	// Heights (in rows) for title/footer/spacers when using fixed-size AddItem
	TitleHeight     = 1
	FooterHeight    = 1
	TopSpacerHeight = 1

	// RefreshInterval controls how often handler data is refreshed
	RefreshInterval = 5 * time.Second
)
