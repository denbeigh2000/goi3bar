package clock

import (
	i3 "github.com/denbeigh2000/goi3bar"
	config "github.com/denbeigh2000/goi3bar/config"

	"time"
)

type clockBuilder struct{}

func (b clockBuilder) Build(conf config.Config) (i3.Producer, error) {
	c := Clock{}
	err := conf.ParseConfig(&c)
	if err != nil {
		return nil, err
	}

	return &i3.BaseProducer{
		Generator: c,
		Interval:  1 * time.Second,
		Name:      "time",
	}, nil
}

func init() {
	config.Register("clock", clockBuilder{})
}
