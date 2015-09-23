package clock

import (
	i3 "github.com/denbeigh2000/goi3bar"
	timeFormat "github.com/jehiah/go-strftime"

	"time"
)

type Clock struct {
	Format   string
	Location string
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
	if c.Location == "" {
		c.Location = "Local"
	}

	l, err := time.LoadLocation(c.Location)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	st := timeFormat.Format(c.Format, t.In(l))

	o := i3.Output{
		FullText:  st,
		Color:     "#FFFFFF",
		Separator: true,
	}

	return []i3.Output{o}, nil
}
