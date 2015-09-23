package cpu

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Cpu struct{}

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
func (c Cpu) Generate() ([]i3.Output, error) {
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

	values := matches[1:4]
	outputs := make([]i3.Output, expMatches)

	for i, v := range values {
		valueF, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil, err
		}

		var color string
		switch {
		case valueF >= CritThreshold:
			color = "#FF0000"
		case valueF >= WarnThreshold:
			color = "#FFA500"
		default:
			color = "#00FF00"
		}

		outputs[i] = i3.Output{
			FullText:  v,
			Color:     color,
			Separator: i == expMatches-1,
		}
	}

	return outputs, nil
}
