package handlers

import (
	"github.com/shirou/gopsutil/v4/mem"

	"fmt"

	"github.com/VILJkid/current-state/pkg/system"
	"github.com/VILJkid/current-state/types"
)

func MemoryHandler() types.ListItem {
	listItem := types.ListItem{
		PrimaryText:   "Get memory usage",
		SecondaryText: "No memory usage information available",
		Shortcut:      'a',
	}

	memUsage, err := mem.VirtualMemory()
	if err != nil {
		listItem.Err = err
		return listItem
	}

	listItem.SecondaryText = fmt.Sprintf(
		"All: %s | Used: %s | Available: %s",
		system.FormatSize(memUsage.Total),
		system.FormatSize(memUsage.Used),
		system.FormatSize(memUsage.Available),
	)
	return listItem
}
