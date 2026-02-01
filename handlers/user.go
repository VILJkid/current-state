package handlers

import (
	"fmt"

	"github.com/VILJkid/current-state/pkg/system"
	"github.com/VILJkid/current-state/types"
)

func UserHandler() types.ListItem {
	user, err := system.GetCurrentUser()
	return types.ListItem{
		PrimaryText:   "Get the current logged in user",
		SecondaryText: fmt.Sprintf("Current user: %s", user),
		Shortcut:      'c',
		Err:           err,
	}
}
