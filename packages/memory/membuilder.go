package memory

import (
	. "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"time"
)

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

	return &BaseProducer{
		Generator: Memory{},
		Name:      "memory",
		Interval:  interval,
	}, nil
}

func init() {
	config.Register("memory", memoryBuilder{})
}
