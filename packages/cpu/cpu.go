package cpu

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"strconv"

	"github.com/shirou/gopsutil/load"
)

// Cpu is a small CPU load monitor. It scrapes /proc/loadavg to display your
// average # waiting threads over 1, 5 and 15 minute averages.
type Cpu struct {
	// If the CPU loads exceeds these thresholds, they will be rendered in the
	// corresponding state.
	WarnThreshold float64 `json:"warn_threshold"`
	CritThreshold float64 `json:"crit_threshold"`

	oneLoad     float64
	fiveLoad    float64
	fifteenLoad float64
}

const (
	expMatches = 3
)

func (c *Cpu) populateValues() error {
	load, err := load.LoadAvg()
	if err != nil {
		return err
	}

	c.oneLoad = load.Load1
	c.fiveLoad = load.Load5
	c.fifteenLoad = load.Load15
	return nil
}

// Generate implements Generator
func (c *Cpu) Generate() ([]i3.Output, error) {
	c.populateValues()

	outputs := make([]i3.Output, expMatches)
	for i, l := range []float64{c.oneLoad, c.fiveLoad, c.fifteenLoad} {
		var color string
		switch {
		case l >= c.CritThreshold:
			color = "#FF0000"
		case l >= c.WarnThreshold:
			color = "#FFA500"
		default:
			color = "#00FF00"
		}

		outputs[i] = i3.Output{
			FullText:  strconv.FormatFloat(l, 'f', 2, 64),
			Color:     color,
			Separator: i == expMatches-1,
		}
	}

	return outputs, nil
}
