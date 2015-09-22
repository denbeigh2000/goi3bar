package memory

import (
	i3 "bitbucket.org/denbeigh2000/goi3bar"

	"github.com/pivotal-golang/bytefmt"

	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Memory struct{}

var (
	rx = regexp.MustCompile("MemTotal:\\s+([0-9]+) kB\nMemFree:\\s+([0-9]+) kB\nMemAvailable:\\s+([0-9]+) kB")
)

const (
	FormatString = "Mem: %v%% (%v/%v)"
	MeminfoPath  = "/proc/meminfo"

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
	f, err := os.Open(MeminfoPath)
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
	if len(matches) != 4 {
		return nil, fmt.Errorf("Parsing error in %v", MeminfoPath)
	}

	totali, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return nil, err
	}

	freei, err := strconv.ParseUint(matches[3], 10, 64)
	if err != nil {
		return nil, err
	}

	total := totali * bytefmt.KILOBYTE
	free := freei * bytefmt.KILOBYTE

	used := total - free

	//percUsed := int32((float32(used)*100)/float32(total) + .5)
	percUsed := (used * 100) / total

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
