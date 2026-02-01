package handlers

import (
	"fmt"

	"github.com/VILJkid/current-state/pkg/system"
	"github.com/VILJkid/current-state/types"
)

func UserHandler() types.ListItem {
	user, err := system.GetCurrentUser()
	sanitizedErr := system.SanitizeError(err)
	return types.ListItem{
		PrimaryText:   "User",
		SecondaryText: fmt.Sprintf("Current user: %s", user),
		Shortcut:      'c',
		Err:           sanitizedErr,
	}
}
