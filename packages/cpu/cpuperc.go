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
	// Name to identify this widget when receiving events
	Name string

	// If the CPU loads exceeds these thresholds, they will be rendered in the
	// corresponding state.
	WarnThreshold float64
	CritThreshold float64
	Interval      time.Duration

	collecting  bool
	percentages chan float64
}

func (c *CpuPerc) report(kill <-chan struct{}) chan float64 {
	out := make(chan float64)
	go func() {
		defer close(out)

		out <- 0.0

		for {
			select {
			case <-kill:
				return
			default:
				p, err := cpu.Percent(c.Interval, false)
				if err != nil {
					continue
				}
				out <- p[0]
			}
		}
	}()

	return out
}

func (c CpuPerc) format(p float64) []i3.Output {
	var color string
	switch {
	case p >= c.CritThreshold:
		color = i3.DefaultColors.Crit
	case p >= c.WarnThreshold:
		color = i3.DefaultColors.Warn
	default:
		color = i3.DefaultColors.OK
	}

	return []i3.Output{{
		Name:      c.Name,
		FullText:  fmt.Sprintf("CPU: %.2f%%", p),
		Color:     color,
		Separator: true,
	}}
}

func (c *CpuPerc) Produce(kill <-chan struct{}) <-chan []i3.Output {
	out := make(chan []i3.Output)
	c.percentages = make(chan float64)

	go func() {
		defer close(out)

		percs := c.report(kill)

		for {
			select {
			case <-kill:
				return
			case p := <-percs:
				out <- c.format(p)
			}
		}
	}()
	return out
}
