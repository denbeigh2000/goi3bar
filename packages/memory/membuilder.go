package memory

import (
	. "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"time"
)

const Identifier = "memory"

// MemoryConfig represents the configuration for a Memory plugin in
// JSON format
type MemoryConfig struct {
	Interval      string `json:"interval"`
	WarnThreshold int    `json:"warn_threshold"`
	CritThreshold int    `json:"crit_threshold"`
}

type memoryBuilder struct{}

func (m memoryBuilder) Build(c config.Config) (Producer, error) {
	conf := MemoryConfig{}
	err := c.ParseConfig(&conf)
	if err != nil {
		return nil, err
	}

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return nil, err
	}

	g := Memory{
		Name:          Identifier,
		WarnThreshold: conf.WarnThreshold,
		CritThreshold: conf.CritThreshold,
	}

	return &BaseProducer{
		Generator: g,
		Name:      Identifier,
		Interval:  interval,
	}, nil
}

func init() {
	config.Register(Identifier, memoryBuilder{})
}
