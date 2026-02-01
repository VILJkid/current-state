package system

import (
	"syscall"

	"github.com/VILJkid/current-state/types"
)

func GetDiskUsage(path string) (*types.DiskStatus, error) {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(path, &fs); err != nil {
		return nil, err
	}

	d := &types.DiskStatus{
		All:  fs.Blocks * uint64(fs.Bsize),
		Free: fs.Bfree * uint64(fs.Bsize),
	}
	d.Used = d.All - d.Free
	return d, nil
}
