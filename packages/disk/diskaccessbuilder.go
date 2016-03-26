package disk

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"time"
)

type diskIOConfig struct {
	WarnThreshold float64 `json:"warn_threshold"`
	CritThreshold float64 `json:"crit_threshold"`

	Interval string `json:"interval"`

	Items []DiskIOItem `json:"items"`
}

type diskIOBuilder struct{}

func (b diskIOBuilder) Build(c config.Config) (p i3.Producer, err error) {
	conf := diskIOConfig{}
	err = c.ParseConfig(&conf)
	if err != nil {
		return
	}

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return
	}

	return &DiskIOGenerator{
		WarnThreshold: conf.WarnThreshold,
		CritThreshold: conf.CritThreshold,
		Interval:      interval,
		Items:         conf.Items,
	}, nil
}

func init() {
	config.Register("disk_access", diskIOBuilder{})
}
