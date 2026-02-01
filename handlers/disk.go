package handlers

import (
	"fmt"

	"github.com/VILJkid/current-state/pkg/system"
	"github.com/VILJkid/current-state/types"
)

func DiskHandler() types.ListItem {
	listItem := types.ListItem{
		PrimaryText:   "Get disk usage",
		SecondaryText: "No disk usage information available",
		Shortcut:      'b',
	}

	diskUsage, err := system.GetDiskUsage("/")
	if err != nil {
		listItem.Err = system.SanitizeError(err)
		return listItem
	}

	listItem.SecondaryText = fmt.Sprintf(
		"All: %s | Used: %s | Free: %s",
		system.FormatSize(diskUsage.All),
		system.FormatSize(diskUsage.Used),
		system.FormatSize(diskUsage.Free),
	)
	return listItem
}
