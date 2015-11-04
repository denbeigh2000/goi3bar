package memory

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"github.com/pivotal-golang/bytefmt"
	"github.com/shirou/gopsutil/mem"

	"fmt"
)

type Memory struct{}

const (
	FormatString = "Mem: %v%% (%v/%v)"

	WarnThreshold = 75
	CritThreshold = 85
)

func (m Memory) IsWarn(i int) bool {
	return i >= WarnThreshold
}

func (m Memory) IsCrit(i int) bool {
	return i >= CritThreshold
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
		color = "#FF0000"
	case m.IsWarn(int(percUsed)):
		color = "#FFA500"
	default:
		color = "#00FF00"
	}

	out := make([]i3.Output, 1)
	out[0] = i3.Output{
		FullText:  fmt.Sprintf(FormatString, percUsed, bytefmt.ByteSize(used), bytefmt.ByteSize(total)),
		Color:     color,
		Separator: true,
	}

	return out, nil
}
