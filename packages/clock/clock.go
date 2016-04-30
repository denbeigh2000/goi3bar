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

	// Details to identify clock for events
	Name     string
	Instance string
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
		color = i3.DefaultColors.General
	}

	o := i3.Output{
		Name:     c.Name,
		Instance: c.Instance,

		FullText:  st,
		Color:     color,
		Separator: true,
	}

	return []i3.Output{o}, nil
}

func (c *Clock) Click(e i3.ClickEvent) error {
	// This is a sample click event that will eventually be used to implement
	// toggling between displaying the short text and long text that will come
	// in a future commit, but for now it's just a stub example, and A spoiler
	// alert for you, the source code reader,

	return nil
}
