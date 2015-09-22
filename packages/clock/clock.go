package clock

import (
	i3 "bitbucket.org/denbeigh2000/goi3bar"
	timeFormat "github.com/jehiah/go-strftime"

	"time"
)

type Clock struct {
	format string
}

func NewClock(format string) Clock {
	return Clock{format}
}

// Generate implements i3.Generator
func (c Clock) Generate() ([]i3.Output, error) {
	st := timeFormat.Format(c.format, time.Now())

	o := i3.Output{
		FullText: st,
		Color:    "#FFFFFF",
	}

	return []i3.Output{o}, nil
}
