package disk

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"github.com/pivotal-golang/bytefmt"

	"fmt"
	"syscall"
)

const (
	freeFormat = "%v: %v free"
)

type DiskUsageItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type DiskUsageGenerator struct {
	CritThreshold int `json:"crit_threshold"`
	WarnThreshold int `json:"warn_threshold"`

	Items []DiskUsageItem `json:"items"`

	// Identifier for receiving click events
	Name string
}

type diskUsageInfo struct {
	Free     uint64
	Total    uint64
	UsedPerc int
}

func getUsage(path string) (info diskUsageInfo, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return
	}

	info.Free = (uint64(stat.Bavail) * uint64(stat.Bsize))
	info.Total = (uint64(stat.Blocks) * uint64(stat.Bsize))

	info.UsedPerc = int(uint64(info.Total-info.Free) * 100 / uint64(info.Total))

	return
}

func (g DiskUsageGenerator) Generate() ([]i3.Output, error) {
	items := make([]i3.Output, len(g.Items))

	for i, item := range g.Items {
		usage, err := getUsage(item.Path)
		if err != nil {
			return nil, err
		}

		free := bytefmt.ByteSize(usage.Free * bytefmt.BYTE)

		items[i].FullText = fmt.Sprintf(freeFormat, item.Name, free)

		freePercent := 100 - int(usage.UsedPerc)

		var color string
		switch {
		case freePercent < g.CritThreshold:
			color = i3.DefaultColors.Crit
		case freePercent < g.WarnThreshold:
			color = i3.DefaultColors.Warn
		default:
			color = i3.DefaultColors.OK
		}

		items[i].Color = color

		items[i].Name = g.Name
		items[i].Instance = item.Path
	}

	items[len(items)-1].Separator = true

	return items, nil
}
