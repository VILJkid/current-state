package ui

import (
	"time"

	"github.com/VILJkid/current-state/handlers"
	"github.com/VILJkid/current-state/types"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateMenu(app *tview.Application) tview.Primitive {
	list := tview.NewList()

	listItems := []types.ListItem{
		{
			PrimaryText: "Select an item from below...",
			Shortcut:    'm',
			Err:         nil,
		},
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

	for _, item := range listItems {
		list.AddItem(item.PrimaryText, "", item.Shortcut, item.Action)
	}

	// Start dynamic updates for handler items (indices 1, 2)
	go updateSelectedItem(app, list, listItems)

	list.SetChangedFunc(func(index int, _, _ string, _ rune) {
		// Clear all secondary texts first
		for i, listItem := range listItems {
			list.SetItemText(i, listItem.PrimaryText, "")
		}

		// Set secondary text color to green as default
		list.SetSecondaryTextColor(tcell.ColorGreen)

		// Set secondary text for current item only
		currentListItem := listItems[index]
		if currentListItem.Err != nil {
			// Show error modal
			errorModal := GetOKModal(app, list, currentListItem.Err.Error())
			app.SetRoot(errorModal, false)

			// Set secondary text color to red
			list.SetSecondaryTextColor(tcell.ColorRed)
		}

		list.SetItemText(index, currentListItem.PrimaryText, currentListItem.SecondaryText)
	})

	// Wrap the list in a centered flex container
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false). // Top spacer
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(tview.NewBox(), 0, 1, false). // Left spacer
			AddItem(list, 0, 1, true).            // Your list
			AddItem(tview.NewBox(), 0, 1, false), // Right spacer
							0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false) // Bottom spacer

	return flex
}

// updateSelectedItem refreshes secondary text for ONLY the currently selected item every 5 seconds
func updateSelectedItem(app *tview.Application, list *tview.List, listItems []types.ListItem) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		app.QueueUpdateDraw(func() {
			// Get the currently selected item index
			currentIndex := list.GetCurrentItem()

			// Only update if it's one of the handler items (indices 1, 2)
			if currentIndex >= 1 && currentIndex < 3 {
				var freshItem types.ListItem

				switch currentIndex {
				case 1:
					freshItem = handlers.MemoryHandler()
				case 2:
					freshItem = handlers.DiskHandler()
				}

				// Only update if there's no error
				if freshItem.Err == nil {
					list.SetItemText(currentIndex, freshItem.PrimaryText, freshItem.SecondaryText)
				}
			}
		})
	}
}
