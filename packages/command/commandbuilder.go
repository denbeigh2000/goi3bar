package command

import (
	i3 "github.com/denbeigh2000/goi3bar"
	config "github.com/denbeigh2000/goi3bar/config"

	"time"
)

const Identifier = "command"

type commandBuilder struct{}

func (b commandBuilder) Build(c config.Config) (i3.Producer, error) {
	conf := CommandConfig{}
	if err := c.ParseConfig(&conf); err != nil {
		return nil, err
	}

	if len(conf.Name) == 0 {
		switch c.Name {
		case "":
			conf.Name = Identifier
		default:
			conf.Name = c.Name
		}
	}

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return nil, err
	}

	return &i3.BaseProducer{
		Generator: &conf,
		Interval:  interval,
		Name:      Identifier,
	}, nil
}

func init() {
	config.Register(Identifier, commandBuilder{})
}
