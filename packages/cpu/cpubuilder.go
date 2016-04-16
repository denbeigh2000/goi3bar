package cpu

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"time"
)

type cpuConfig struct {
	Interval      string  `json:"interval"`
	CritThreshold float64 `json:"crit_threshold"`
	WarnThreshold float64 `json:"warn_threshold"`
}

type cpuBuilder struct {
	perc bool
}

func (b cpuBuilder) Build(c config.Config) (p i3.Producer, err error) {
	conf := cpuConfig{}
	err = c.ParseConfig(&conf)
	if err != nil {
		return
	}

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return
	}

	if b.perc {
		p = &CpuPerc{
			Name:          "cpu_util",
			WarnThreshold: conf.WarnThreshold,
			CritThreshold: conf.CritThreshold,
			Interval:      interval,
		}
	} else {
		p = &i3.BaseProducer{
			Generator: &Cpu{
				Name:          "cpu_load",
				WarnThreshold: conf.WarnThreshold,
				CritThreshold: conf.CritThreshold,
			},
			Interval: interval,
		}
	}

	return
}

func init() {
	config.Register("cpu_load", cpuBuilder{perc: false})
	config.Register("cpu_util", cpuBuilder{perc: true})
}
