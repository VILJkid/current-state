package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/VILJkid/current-state/handlers"
	"github.com/VILJkid/current-state/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateMenu(app *tview.Application) tview.Primitive {
	list := tview.NewList()
	// Default secondary text color for unselected items (cyan)
	list.SetSecondaryTextColor(tcell.ColorDarkCyan)
	// Selected item text color (black)
	list.SetSelectedTextColor(tcell.ColorBlack)

	// (wrapping handled by TextView's word wrap)

	listItems := []types.ListItem{
		handlers.MemoryHandler(),
		handlers.DiskHandler(),
		handlers.UserHandler(),
		{
			PrimaryText:   "Quit",
			SecondaryText: "Press to exit the application",
			Shortcut:      'q',
			Action: func() {
				app.Stop()
			},
			Err: nil,
		},
	}

	// Persistent descriptions for each menu item (shown when not selected)
	descs := []string{
		"View memory usage",
		"Check disk space",
		"Show current user",
		"Quit application",
	}

	for i, item := range listItems {
		sec := ""
		if i < len(descs) {
			sec = descs[i]
		}
		list.AddItem(item.PrimaryText, sec, item.Shortcut, item.Action)
	}

	// Ensure consistent colors for list items
	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSecondaryTextColor(tcell.ColorDarkCyan) // cyan for unselected secondary
	list.SetSelectedTextColor(tcell.ColorBlack)     // black for selected primary
	// Try to set selected secondary color (may not exist on all tview versions)
	if setter, ok := any(list).(interface{ SetSelectedSecondaryTextColor(tcell.Color) *tview.List }); ok {
		setter.SetSelectedSecondaryTextColor(tcell.ColorDarkMagenta)
	}

	// Auto-select Memory item on startup (index 0) and display data immediately
	list.SetCurrentItem(0)

	// Details panel (right side) for multiline output
	detailsView := tview.NewTextView()
	detailsView.SetWrap(true)
	detailsView.SetWordWrap(true)
	detailsView.SetDynamicColors(true)
	detailsView.SetBorder(true)
	detailsView.SetTitle("Details")

	// Manually trigger data fetch for the selected item and populate details
	memoryItem := handlers.MemoryHandler()
	if memoryItem.Err == nil {
		list.SetItemText(0, memoryItem.PrimaryText, memoryItem.SecondaryText)
		detailsView.SetText(memoryItem.SecondaryText)
	}

	// Create cancellation context for goroutine lifecycle management
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel // Store for cleanup (prevents unused variable warning)

	// Start dynamic updates for handler items (indices 0, 1)
	go updateSelectedItem(ctx, app, list, listItems, detailsView)

	// Track previous selected index so we can restore its persistent description
	prevIndex := 0

	list.SetChangedFunc(func(index int, _, _ string, _ rune) {
		// Restore previous item's persistent description
		if prevIndex >= 0 && prevIndex < len(listItems) {
			list.SetItemText(prevIndex, listItems[prevIndex].PrimaryText, descs[prevIndex])
		}

		// Ensure unselected items keep the cyan secondary color
		list.SetSecondaryTextColor(tcell.ColorDarkCyan)

		// For handler items (memory, disk, user), call handlers fresh to get latest data
		currentListItem := listItems[index]
		switch index {
		case 0:
			currentListItem = handlers.MemoryHandler()
		case 1:
			currentListItem = handlers.DiskHandler()
		case 2:
			currentListItem = handlers.UserHandler()
		}

		if currentListItem.Err != nil {
			// Show error modal
			errorModal := GetOKModal(app, list, currentListItem.Err.Error())
			app.SetRoot(errorModal, false)

			// On error, set secondary text to red so user notices
			list.SetSecondaryTextColor(tcell.ColorRed)

			// Update previous index and return early
			prevIndex = index
			return
		}

		// Update the selected item's list entry and details panel
		list.SetItemText(index, currentListItem.PrimaryText, currentListItem.SecondaryText)
		detailsView.SetText(currentListItem.SecondaryText)

		// Update previous index
		prevIndex = index
	})

	// Create title text box with dynamic colors
	titleBox := CreateColoredTextView(BuildTitleText())

	// Create help footer text box with dynamic colors
	helpBox := CreateColoredTextView(BuildHelpboxText())

	// Wrap the list and details in a centered flex container with title and footer
	cols := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(list, 0, LeftColumnWeight, true).
		AddItem(detailsView, 0, RightColumnWeight, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), TopSpacerHeight, 0, false). // Top spacer
		AddItem(titleBox, TitleHeight, 0, false).           // Title
		AddItem(tview.NewBox(), 0, 1, false).               // Middle spacer (proportional)
		AddItem(cols, 0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false).    // Bottom spacer
		AddItem(helpBox, FooterHeight, 0, false) // Footer

	// Live resize check: show warning modal if terminal too narrow
	modalShown := false
	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		w, _ := screen.Size()
		if w < MinTerminalWidth && !modalShown {
			warning := fmt.Sprintf("Terminal width %d too small to display full UI (need >= %d columns). Resize terminal or press OK.", w, MinTerminalWidth)
			modal := GetOKModal(app, flex, warning)
			// Override modal done func so we can reset the shown flag and return to flex
			modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.SetRoot(flex, true)
				modalShown = false
			})
			app.QueueUpdateDraw(func() {
				app.SetRoot(modal, true)
			})
			modalShown = true
			return
		}
		// If terminal is wide enough and modal was previously shown, restore root
		if w >= MinTerminalWidth && modalShown {
			app.QueueUpdateDraw(func() {
				app.SetRoot(flex, true)
			})
			modalShown = false
		}
	})

	return flex
}

// updateSelectedItem refreshes secondary text for ONLY the currently selected item every 5 seconds.
// Respects context cancellation for graceful shutdown.
func updateSelectedItem(ctx context.Context, app *tview.Application, list *tview.List, listItems []types.ListItem, detailsView *tview.TextView) {
	ticker := time.NewTicker(RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Graceful shutdown - exit goroutine
			return
		case <-ticker.C:
			app.QueueUpdateDraw(func() {
				// Get the currently selected item index
				currentIndex := list.GetCurrentItem()

				// Only update if it's one of the handler items (indices 0, 1)
				if currentIndex >= 0 && currentIndex < 2 {
					var freshItem types.ListItem

					switch currentIndex {
					case 0:
						freshItem = handlers.MemoryHandler()
					case 1:
						freshItem = handlers.DiskHandler()
					}

					// Only update if there's no error
					if freshItem.Err == nil {
						list.SetItemText(currentIndex, freshItem.PrimaryText, freshItem.SecondaryText)
						// update details panel as well
						detailsView.SetText(freshItem.SecondaryText)
					}
				}
			})
		}
	}
}
