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
	Format string
	// The IANA Timezone database zone name to show the time for
	Location string
	Color    string
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

func Build(options interface{}) (i3.Producer, error) {
	optMap, ok := options.(map[string]interface{})
	if !ok {
		return nil, CorruptedConfigErr
	}

	c := Clock{
		Color:    optMap["color"].(string),
		Format:   optMap["format"].(string),
		Location: optMap["location"].(string),
	}

	item := i3.BaseProducer{
		Generator: c,
		Interval:  1 * time.Second,
		Name:      "time",
	}

	return &item, nil
}
