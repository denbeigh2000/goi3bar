package disk

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"time"
)

type diskUsageConfig struct {
	Interval string             `json:"interval"`
	Options  DiskUsageGenerator `json:"options"`
}

type diskUsageBuilder struct{}

func (b diskUsageBuilder) Build(c config.Config) (p i3.Producer, err error) {
	conf := diskUsageConfig{}
	err = c.ParseConfig(&conf)
	if err != nil {
		return
	}

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return
	}

	return &i3.BaseProducer{
		Generator: conf.Options,
		Interval:  interval,
		Name:      "disk usage",
	}, nil
}

func init() {
	config.Register("disk_usage", diskUsageBuilder{})
}
