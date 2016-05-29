package command

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"time"
)

const Identifier = "command"

type CommandBuilder struct {
	Interval string  `json:"interval"`
	Color    string  `json:"color"`
	Options  Command `json:"options"`
}

func (b CommandBuilder) Build(c config.Config) (p i3.Producer, err error) {
	conf := Command{}

	interval, err := time.ParseDuration(b.Interval)
	if err != nil {
		return
	}

	err = c.ParseConfig(&b)
	if err != nil {
		return
	}

	color, err := i3.ParseColor(b.Color)
	if err != nil {
		return
	}
	conf.Color = color

	p = &i3.BaseProducer{
		Generator: conf,
		Interval:  interval,
		Name:      Identifier,
	}

	return
}

func init() {
	config.Register("command", CommandBuilder{})
}
