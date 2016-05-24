package command

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"
	"time"
)

const Identifier = "command"

type CommandConfig struct {
	Interval string             `json:"interval"`
	Options  CommandGenerator   `json:"options"`
}

type CommandBuilder struct {
}

func (b CommandBuilder) Build(c config.Config) (p i3.Producer, err error) {
	conf := CommandConfig{}
	err = c.ParseConfig(&conf)
	if err != nil {
		return
	}

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return
	}

	conf.Options.Name = Identifier

	return &i3.BaseProducer{
		Generator: conf.Options,
		Interval:  interval,
		Name:      "cmd",
	}, nil
}

func init() {
	config.Register("command", CommandBuilder{})
}
