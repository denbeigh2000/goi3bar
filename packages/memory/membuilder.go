package memory

import (
	. "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"time"
)

func (m memoryBuilder) Build(c config.Config) (Producer, error) {
	err := c.ParseConfig(&m)
	if err != nil {
		return nil, err
	}

	interval, err := time.ParseDuration(m.Refresh)
	if err != nil {
		return nil, err
	}

	return &BaseProducer{
		Generator: Memory{},
		Name:      m.Name,
		Interval:  interval,
	}, nil
}

type memoryBuilder struct {
	Name          string
	Refresh       string
	WarnThreshold int
	CritThreshold int
}
