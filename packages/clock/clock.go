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

type multiClock struct {
	clocks        map[string]Clock
	multiProducer i3.MultiProducer
}

//func MultiClock(format string, times map[string]time.Timezone) Generator {
//	clocks := make(map[string]Clock)
//	for abbr, zone := range times {
//	}
//}

//func (m *Multi

// Generate implements i3.Generator
func (c Clock) Generate() ([]i3.Output, error) {
	st := timeFormat.Format(c.format, time.Now())

	o := i3.Output{
		FullText:  st,
		Color:     "#FFFFFF",
		Separator: true,
	}

	return []i3.Output{o}, nil
}
