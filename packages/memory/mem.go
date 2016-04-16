package memory

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"github.com/pivotal-golang/bytefmt"
	"github.com/shirou/gopsutil/mem"

	"fmt"
)

type Memory struct {
	Name          string
	WarnThreshold int
	CritThreshold int
}

const (
	FormatString = "Mem: %v%% (%v/%v)"

	DefaultWarnThreshold = 75
	DefaultCritThreshold = 85
)

func (m Memory) IsWarn(i int) bool {
	switch m.WarnThreshold {
	case 0:
		return i >= DefaultWarnThreshold
	default:
		return i >= m.WarnThreshold
	}
}

func (m Memory) IsCrit(i int) bool {
	switch m.CritThreshold {
	case 0:
		return i >= DefaultCritThreshold
	default:
		return i >= m.CritThreshold
	}
}

func (m Memory) Generate() ([]i3.Output, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	total := mem.Total
	used := total - mem.Buffers - mem.Cached - mem.Free

	percUsed := (100 * used) / total

	var color string
	switch {
	case m.IsCrit(int(percUsed)):
		color = i3.DefaultColors.Crit
	case m.IsWarn(int(percUsed)):
		color = i3.DefaultColors.Warn
	default:
		color = i3.DefaultColors.OK
	}

	out := make([]i3.Output, 1)
	out[0] = i3.Output{
		Name:      m.Name,
		FullText:  fmt.Sprintf(FormatString, percUsed, bytefmt.ByteSize(used), bytefmt.ByteSize(total)),
		Color:     color,
		Separator: true,
	}

	return out, nil
}
