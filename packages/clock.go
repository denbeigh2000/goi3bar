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

// Update implements i3.Updater
func (c Clock) Update(o *i3.Output) error {
	st := timeFormat.Format(c.format, time.Now())
	o.FullText = st

	return nil
}
