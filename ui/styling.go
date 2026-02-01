package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Colors used throughout the UI
var (
	ColorPrimary   = tcell.ColorGreen
	ColorSecondary = tcell.ColorYellow
	ColorText      = tcell.ColorWhite
	ColorError     = tcell.ColorRed
)

// CreateStyledTextView creates a centered TextView with specified colors
func CreateStyledTextView(text string, textColor, bgColor tcell.Color) *tview.TextView {
	view := tview.NewTextView().
		SetText(text).
		SetTextAlign(tview.AlignCenter)

	if textColor != tcell.ColorDefault {
		view.SetTextColor(textColor)
	}
	if bgColor != tcell.ColorDefault {
		view.SetBackgroundColor(bgColor)
	}

	return view
}

// CreateColoredTextView creates a TextView with dynamic color tags
// Allows mixed colors in a single text using [color]text[/color] tags
func CreateColoredTextView(text string) *tview.TextView {
	view := tview.NewTextView().
		SetText(text).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	return view
}

// BuildHelpboxText constructs the help footer text with color-coded shortcuts
// Keys appear in primary color, descriptions in secondary color
func BuildHelpboxText() string {
	shortcuts := []struct {
		key  string
		desc string
	}{
		{"a", "Memory"},
		{"b", "Disk"},
		{"c", "User"},
		{"q", "Quit"},
		{"↑↓", "Navigate"},
	}

	text := "[yellow::b]Shortcuts:[-::-]  "
	for i, s := range shortcuts {
		text += "[green]" + s.key + "[white]=" + s.desc
		if i < len(shortcuts)-1 {
			text += "  "
		}
	}

	// Append current layout weights for easy tuning visibility
	text += fmt.Sprintf("  [white]| Layout:%d:%d", LeftColumnWeight, RightColumnWeight)
	return text
}

// BuildTitleText constructs the title with color differentiation
func BuildTitleText() string {
	return "[green::b]⚡ current-state[-::-][white] - System Information Dashboard"
}
