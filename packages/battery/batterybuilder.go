package battery

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"fmt"
	"time"
)

type batteryConfig struct {
	Interval      string `json:"interval"`
	Name          string `json:"name"`
	Identifier    string `json:"identifier"`
	WarnThreshold int    `json:"warn_threshold"`
	CritThreshold int    `json:"crit_threshold"`
}

type batteryBuilder struct{}

func (b batteryBuilder) Build(c config.Config) (i3.Producer, error) {
	conf := batteryConfig{}
	err := c.ParseConfig(&conf)
	if err != nil {
		return nil, err
	}

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return nil, err
	}

	if ok := validateThreshold(conf.WarnThreshold); !ok {
		return nil, fmt.Errorf(
			"WarnThreshold for %v (%v) is outside acceptable range (0, 100)",
			conf.Identifier, conf.WarnThreshold,
		)
	}

	if ok := validateThreshold(conf.CritThreshold); !ok {
		return nil, fmt.Errorf(
			"CritThreshold for %v (%v) is outside acceptable range (0, 100)",
			conf.Identifier, conf.CritThreshold,
		)
	}

	bat := Battery{
		Name:          conf.Name,
		Identifier:    conf.Identifier,
		WarnThreshold: conf.WarnThreshold,
		CritThreshold: conf.CritThreshold,
	}

	return &i3.BaseProducer{
		Generator: &bat,
		Interval:  interval,
		Name:      conf.Identifier + "_bat",
	}, nil
}

func validateThreshold(v int) (ok bool) {
	ok = v > 0 || v < 100

	return
}

func init() {
	config.Register("battery", batteryBuilder{})
}
