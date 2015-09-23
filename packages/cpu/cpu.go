package cpu

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Cpu struct {
	// If the CPU loads exceeds these thresholds, they will be rendered in the
	// corresponding state.
	WarnThreshold float64
	CritThreshold float64

	oneLoad     float64
	fiveLoad    float64
	fifteenLoad float64
}

const (
	WarnThreshold = 0.7
	CritThreshold = 0.85

	expMatches  = 3
	rxStr       = "^([0-9]+.[0-9]+) ([0-9]+.[0-9]+) ([0-9]+.[0-9]+)"
	LoadavgPath = "/proc/loadavg"
)

var (
	rx = regexp.MustCompile(rxStr)
)

// Generate implements Generator
func (c *Cpu) Generate() ([]i3.Output, error) {
	f, err := os.Open(LoadavgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	matches := rx.FindStringSubmatch(buf.String())

	// There should be total 4 matches - The whole match, then the 3 groups
	if len(matches) != expMatches+1 {
		return nil,
			fmt.Errorf("Incorrect number of matches for %v, (got %v expected %v)",
				LoadavgPath, len(matches), expMatches+1)
	}

	c.oneLoad, err = strconv.ParseFloat(matches[1], 32)
	if err != nil {
		return nil, err
	}

	c.fiveLoad, err = strconv.ParseFloat(matches[2], 32)
	if err != nil {
		return nil, err
	}

	c.fifteenLoad, err = strconv.ParseFloat(matches[3], 32)
	if err != nil {
		return nil, err
	}

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
