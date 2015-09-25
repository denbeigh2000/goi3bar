package cpu

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// Cpu is a small CPU load monitor. It scrapes /proc/loadavg to display your
// average # waiting threads over 1, 5 and 15 minute averages.
type CpuPerc struct {
	// If the CPU loads exceeds these thresholds, they will be rendered in the
	// corresponding state.
	WarnThreshold float64
	CritThreshold float64

	collecting bool
	percentage float64
}

func (c *CpuPerc) collect() {
	for {
		percentages, err := cpu.CPUPercent(5*time.Second, false)
		if err != nil {
			continue
		}
		c.percentage = percentages[0]
	}
}

// Generate implements Generator
func (c *CpuPerc) Generate() ([]i3.Output, error) {
	if !c.collecting {
		c.collecting = true
		go c.collect()
	}

	var color string
	switch {
	case c.percentage >= c.CritThreshold:
		color = "#FF0000"
	case c.percentage >= c.WarnThreshold:
		color = "#FFA500"
	default:
		color = "#00FF00"
	}

	output := i3.Output{
		FullText:  fmt.Sprintf("CPU: %.2f%%", c.percentage),
		Color:     color,
		Separator: true,
	}

	return []i3.Output{output}, nil
}
