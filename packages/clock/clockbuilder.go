package clock

import (
	i3 "github.com/denbeigh2000/goi3bar"
	config "github.com/denbeigh2000/goi3bar/config"

	"time"
)

const Identifier = "clock"

type clockBuilder struct{}

func (b clockBuilder) Build(conf config.Config) (i3.Producer, error) {
	c := Clock{}
	err := conf.ParseConfig(&c)
	if err != nil {
		return nil, err
	}

	switch conf.Name {
	case "":
		c.Name = Identifier
	default:
		c.Name = conf.Name
	}

	switch c.Location {
	case "":
		c.Instance = "Local"
	default:
		c.Instance = c.Location
	}

	return &i3.BaseProducerClicker{
		GeneratorClicker: &c,
		Interval:         1 * time.Second,
		Name:             "time",
	}, nil
}

func init() {
	config.Register(Identifier, clockBuilder{})
}
