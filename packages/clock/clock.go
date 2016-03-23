package clock

import (
	i3 "github.com/denbeigh2000/goi3bar"
	timeFormat "github.com/jehiah/go-strftime"

	"errors"
	"time"
)

var CorruptedConfigErr = errors.New("Could not parse config options")

type Clock struct {
	// How the time should be formatted. See http://strftime.org/ for reference.
	Format string `json:"format"`
	// The IANA Timezone database zone name to show the time for
	Location string `json:"location"`
	Color    string `json:"color"`
}

// Generate implements i3.Generator
func (c Clock) Generate() ([]i3.Output, error) {
	if c.Location == "" {
		c.Location = "Local"
	}

	l, err := time.LoadLocation(c.Location)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	st := timeFormat.Format(c.Format, t.In(l))

	color := c.Color
	if color == "" {
		color = "#FFFFFF"
	}

	o := i3.Output{
		FullText:  st,
		Color:     color,
		Separator: true,
	}

	return []i3.Output{o}, nil
}
